package gomerkletree

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"

	"github.com/iden3/go-iden3-crypto/poseidon"
)

type MerkleStructure struct {
	Depth         int
	Root          *big.Int
	LastLeafIndex uint64 // index of the leaf in its original array or list before being
	LastLeaf      *big.Int
	Zero          *big.Int
	Siblings      []*big.Int
	// FilledNodes: Last value filled in a layer of the tree
	// see: https://github.com/privacy-scaling-explorations/zk-kit/blob/f9b209d6006fd454fc464bbd7e1ad3cfab6389d0/packages/incremental-merkle-tree.sol/contracts/IncrementalBinaryTree.sol#L14
	FilledNodes []*big.Int
	EmptyNodes  []*big.Int
	Leaves      map[uint64]*big.Int
	PrintLog    bool
}

func NewMerkletree(depth int, zero *big.Int, printLog bool) (tree MerkleStructure) {
	tree.Depth = depth
	tree.Zero = zero
	tree.LastLeafIndex = 0
	tree.PrintLog = printLog
	tree.Siblings = make([]*big.Int, depth)
	tree.FilledNodes = make([]*big.Int, depth)
	tree.Leaves = make(map[uint64]*big.Int)
	tree.ComputeEmptyNodes()
	log.Println("NewMerkletree - Zero: ", zero.String())
	return
}

func (tree *MerkleStructure) ComputeEmptyNodes() {
	tree.EmptyNodes = append(tree.EmptyNodes, tree.Zero)
	for layer := 1; layer <= tree.Depth; layer++ {
		prevLayerHash := tree.EmptyNodes[layer-1]
		tree.EmptyNodes = append(tree.EmptyNodes, HashCalculator(prevLayerHash, prevLayerHash))
	}
}

func (tree *MerkleStructure) getEmptyNode(layer int) *big.Int {
	return tree.EmptyNodes[layer]
}

func (tree *MerkleStructure) InsertLeaf(leafData []byte) (leafDb DatabaseLeaf) {
	tree.LastLeaf = EncodeToBigInt(leafData)
	tree.Root = tree.LastLeaf
	if tree.PrintLog {
		log.Println("InsertLeaf - Initial state - Root: ", tree.Root)
	}
	idx := tree.LastLeafIndex
	for layer := 0; layer < tree.Depth; layer++ {
		isRightNode := idx%2 == 1 // is right node in this level?
		previousRoot := tree.Root.String()
		if isRightNode {
			tree.Siblings[layer] = tree.FilledNodes[layer]
			tree.Root = HashCalculator(tree.Siblings[layer], tree.Root)
		} else {
			tree.FilledNodes[layer] = tree.Root
			tree.Siblings[layer] = tree.getEmptyNode(layer)
			tree.Root = HashCalculator(tree.Root, tree.Siblings[layer])
		}
		if tree.PrintLog {
			log.Println("layer: ", layer, " - path: ", idx%2, " - Previous Root: ", previousRoot, " - Sibling: ", tree.Siblings[layer].String(), " - Root result: ", tree.Root.String())
		}
		idx = idx / 2
	}

	leafDb.LeafData = leafData
	leafDb.LeafIndex = tree.LastLeafIndex
	leafDb.LeafEncoded = tree.LastLeaf

	tree.Leaves[tree.LastLeafIndex] = tree.LastLeaf
	tree.LastLeafIndex++

	if tree.PrintLog {
		log.Println("InsertLeaf - Final state - Root: ", tree.Root)
	}

	return
}

/*
TODO: Implement UpdateLeaf
https://github.com/privacy-scaling-explorations/zk-kit/blob/f9b209d6006fd454fc464bbd7e1ad3cfab6389d0/packages/incremental-merkle-tree.sol/contracts/IncrementalBinaryTree.sol#L171
*/
func (tree *MerkleStructure) UpdateLeaf(leafData []byte, index uint64, path []uint8, siblings []*big.Int) (leafDb DatabaseLeaf) {
	return
}

func (tree *MerkleStructure) GenerateProof(index uint64) (path []uint8, siblings []*big.Int, root *big.Int, err error) {
	if index > tree.LastLeafIndex {
		err = fmt.Errorf("index higher than elements inserted into tree. Index: %d", index)
		return
	}
	root = tree.Leaves[index]
	if root == nil {
		err = fmt.Errorf("element not found for this index: %d", index)
		return
	}
	if tree.PrintLog {
		log.Println("GenerateProof - Initial state - Root: ", root)
	}
	siblings = make([]*big.Int, tree.Depth)
	filledNodes := tree.FilledNodes
	for layer := 0; layer < tree.Depth; layer++ {
		isRightNode := index%2 == 1 // is right node in this level?
		if isRightNode {
			siblings[layer] = filledNodes[layer]
			root = HashCalculator(siblings[layer], root)
			path = append(path, 1)
		} else {
			filledNodes[layer] = root
			siblings[layer] = tree.getEmptyNode(layer)
			root = HashCalculator(root, siblings[layer])
			path = append(path, 0)
		}
		if tree.PrintLog {
			log.Println("G - layer: ", layer, " - path: ", path, " - Sibling: ", siblings[layer].String(), " - Root result: ", root)
		}
		index = index / 2
	}
	if tree.PrintLog {
		log.Println("GenerateProof - Final state - Root: ", root)
	}
	return
}

func (tree *MerkleStructure) VerifyProof(leafData []byte, path []uint8, siblings []*big.Int, root *big.Int) (bool, *big.Int, *big.Int) {
	lastLeaf := EncodeToBigInt(leafData)
	newRoot := lastLeaf
	if tree.PrintLog {
		log.Println("VerifyProof - Initial state - Root: ", newRoot.String())
	}
	for layer := 0; layer < tree.Depth; layer++ {
		isRightNode := path[layer] == 1 // is right node in this level?
		if isRightNode {
			newRoot = HashCalculator(siblings[layer], newRoot)
		} else {
			newRoot = HashCalculator(newRoot, siblings[layer])
		}
		if tree.PrintLog {
			log.Println("V - layer: ", layer, " - path: ", path[layer], " - Sibling: ", siblings[layer].String(), " - Root result: ", newRoot.String())
		}
	}
	if tree.PrintLog {
		log.Println("VerifyProof - Final state - Root: ", newRoot.String())
	}
	return root.Cmp(newRoot) == 0, root, newRoot
}

func (tree *MerkleStructure) ToString() (stats string) {
	var sb strings.Builder
	sb.WriteString("Depth: " + strconv.Itoa(tree.Depth))
	sb.WriteString("\n")
	sb.WriteString("Root: " + tree.Root.String())
	sb.WriteString("\n")
	sb.WriteString("LastLeafIndex: " + strconv.Itoa(int(tree.LastLeafIndex)))
	sb.WriteString("\n")
	sb.WriteString("LastLeaf: " + tree.LastLeaf.String())
	sb.WriteString("\n")
	sb.WriteString("Zero: " + tree.Zero.String())
	sb.WriteString("\n")
	sb.WriteString("PrintLog: " + fmt.Sprintf("%t", tree.PrintLog))
	sb.WriteString("\n")
	i := 0
	for ; i < len(tree.Siblings); i++ {
		sb.WriteString("Sibling - index: " + strconv.Itoa(i) + " - value: " + tree.Siblings[i].String())
		sb.WriteString("\n")
	}
	i = 0
	for ; i < len(tree.FilledNodes); i++ {
		sb.WriteString("FilledNodes - index: " + strconv.Itoa(i) + " - value: " + tree.FilledNodes[i].String())
		sb.WriteString("\n")
	}
	i = 0
	for ; i < len(tree.EmptyNodes); i++ {
		sb.WriteString("EmptyNodes - index: " + strconv.Itoa(i) + " - value: " + tree.EmptyNodes[i].String())
		sb.WriteString("\n")
	}
	i = 0
	for ; i < len(tree.Leaves); i++ {
		sb.WriteString("Leaves - index: " + strconv.Itoa(i) + " - value: " + tree.Leaves[uint64(i)].String())
		sb.WriteString("\n")
	}
	stats = sb.String()
	return
}

// HashCalculator optmistc poseidon hash calculator.
func HashCalculator(a, b *big.Int) (hash *big.Int) {
	var hashInputs []*big.Int
	hashInputs = append(hashInputs, a, b)
	hash, err := poseidon.Hash(hashInputs)
	if err != nil {
		log.Println("HashCalculator - Error: ", err.Error())
	}
	return
}

func IsRootMatch(str, treeRoot string) bool {
	str, _ = strings.CutPrefix(str, "0x")
	for len(treeRoot) < 64 {
		treeRoot = "0" + treeRoot
	}
	// log.Printf("isRootMatch - tree.Root and batch.BusTreeNewRoot [%s] - [%s]", treeRoot, str)
	return treeRoot == str
}

func EncodeToBigInt(data []byte) (encodedData *big.Int) {
	encoded64Data := base64.StdEncoding.EncodeToString(data)
	encoded64Num := binary.LittleEndian.Uint64([]byte(encoded64Data))
	encodedData = big.NewInt(0).SetUint64(encoded64Num)
	return
}

func BigIntFromString(str string) (newint *big.Int) {
	newint, _ = big.NewInt(0).SetString(str, 10)
	return
}
