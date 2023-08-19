package gomerkletree

import (
	"math/big"
	"testing"
)

var TestZeroValue *big.Int

func TestMain(m *testing.M) {
	TestZeroValue = EncodeToBigInt([]byte("Meetup Golang SP"))
	m.Run()
}

func TestMerkletreeCreation(t *testing.T) {
	mktree := NewMerkletree(16, TestZeroValue, false)
	t.Log(mktree.Root.String())
	if mktree.Root != nil {
		t.Fatal("Initial merkletree root must be nil")
	}
}

func TestMerkletreeOneLeaf(t *testing.T) {
	mktree := NewMerkletree(16, TestZeroValue, false)
	mktree.InsertLeaf([]byte("Tractor"))
	t.Log(mktree.Root.String())
	if mktree.Root == nil {
		t.Fatal("Merkletree root must not be nil")
	}
}

func TestMerkletreeTwoLeafs(t *testing.T) {
	mktree := NewMerkletree(16, TestZeroValue, false)
	mktree.InsertLeaf([]byte("Tractor"))
	mktree.InsertLeaf(HugeData)
	t.Log(mktree.Root.String())
	if mktree.Root == nil {
		t.Fatal("Merkletree root must not be nil")
	}
}

func TestMerkletreeTwoLeafsMerklePath(t *testing.T) {
	mktree := NewMerkletree(16, TestZeroValue, true)
	mktree.InsertLeaf([]byte("Tractor"))
	mktree.InsertLeaf(HugeData)
	mktree.InsertLeaf([]byte("Golang é uma excelente linguagem para Merkletrees"))
	mktree.InsertLeaf([]byte("lero lero"))
	path, siblings, _, err := mktree.GenerateProof(2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Path:", path)
	t.Log("Siblings:", siblings)
	t.Log(mktree.Root.String())
	if mktree.Root == nil {
		t.Fatal("Merkletree root must not be nil")
	}
	t.Log(mktree.ToString())
}

func TestMerkletreeTwoLeafsMerklePathVerify(t *testing.T) {
	mktree := NewMerkletree(16, TestZeroValue, false)
	path := []uint8{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	siblings := make([]*big.Int, 0)
	siblings = append(siblings, BigIntFromString("8599140134662068052"))
	siblings = append(siblings, BigIntFromString("10067819521047319766852107309646420610919743370087407102768430642652567736822"))
	siblings = append(siblings, BigIntFromString("21708460629667792308055071268939199783370044695388328472701844306535917242192"))
	siblings = append(siblings, BigIntFromString("7485975476463180800053455070211586807979161959627774676249618770478605908096"))
	siblings = append(siblings, BigIntFromString("10347299320411585321646885566934198631104667823051726121058086165915524494935"))
	siblings = append(siblings, BigIntFromString("12903629696625722715888061132032514286322843009182752206466509843250085574704"))
	siblings = append(siblings, BigIntFromString("17391634270714250350967844113061856138585218141161452494485545226128515036033"))
	siblings = append(siblings, BigIntFromString("19804812636803295170207486566715688151323543322527045762096677951497925547180"))
	siblings = append(siblings, BigIntFromString("20695643990884438052769918171177021320753307359239243223324821679432115018141"))
	siblings = append(siblings, BigIntFromString("1358043722075868607704918630151829383256730068737849085451098215893882831100"))
	siblings = append(siblings, BigIntFromString("10002118680962692276545303924198289411957668467416850483925943186426743806747"))
	siblings = append(siblings, BigIntFromString("5647106588072366421599540744240250139957225395205168907265079184478252405279"))
	siblings = append(siblings, BigIntFromString("8802197429240712177227894133008453481233743896581295851359149606031793016065"))
	siblings = append(siblings, BigIntFromString("18092005561719949609558268831494762178361868295037069669555353380447067688942"))
	siblings = append(siblings, BigIntFromString("18254997291417335721341738471922127175603134628737146017357468362805363690980"))
	siblings = append(siblings, BigIntFromString("18890063479545491309347457422119934745757786060816509663123049187886686757457"))
	root := BigIntFromString("10226085413213173741969767641564963431942522917014458725296558191721694814372")
	verified, oldRoot, newRoot := mktree.VerifyProof([]byte("Golang é uma excelente linguagem para Merkletrees"), path, siblings, root)

	if !verified {
		t.Log("Old Root:", oldRoot.String())
		t.Log("New Root:", newRoot.String())
		t.Fatal("Merkle proof is not valid")
	}
}
