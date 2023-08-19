package gomerkletree

import (
	"math/big"

	"gorm.io/gorm"
)

type DatabaseLeaf struct {
	gorm.Model
	LeafIndex   uint64 `gorm:"index:idx_leaf_index,unique"`
	LeafData    []byte
	LeafEncoded *big.Int
}

func NewDatabaseLeaf(index uint64, leafData []byte, leafEncoded *big.Int) (obj DatabaseLeaf) {
	obj.LeafIndex = index
	obj.LeafData = leafData
	obj.LeafEncoded = leafEncoded
	return
}
