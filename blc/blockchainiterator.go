package blc

import (
	"log"

	"github.com/asdine/storm/v3"
)

type BlockChainIterator struct {
	CurrentHash []byte // 当前正在遍历区块的hash值
	DB          *storm.DB
}

// 迭代器遍历下一个区块
func (it *BlockChainIterator) Next() *Block {
	var block Block
	err := it.DB.One("Hash", it.CurrentHash, &block)
	if err != nil {
		log.Panic(err)
	}
	it.CurrentHash = block.PrevHash
	return &block
}
