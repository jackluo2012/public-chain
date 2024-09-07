package blc

import (
	"log"

	"github.com/asdine/storm/v3"
)

const (
	BLOCKCHAIN_DB = "blockchain.db" // 区块链数据库文件名
	BLOCK_BUCKET  = "blocks"        // 区块桶名
)

type BlockChain struct {
	Tip []byte    // 区块链的最后一个区块的hash值
	DB  *storm.DB // 数据库
}

// 添加区块到区块链
// func (bc *BlockChain) AddBlockToBlockChain(data string) {
// 	// 获取最后一个区块
// 	lastBlock := bc.Blocks[len(bc.Blocks)-1]
// 	// 创建新的区块
// 	newBlock := NewBlock(lastBlock.Height+1, lastBlock.Hash, data)
// 	// 将新的区块添加到区块链
// 	bc.Blocks = append(bc.Blocks, newBlock)
// }

// 创建带有创世区块的区块链
func CreateBlockChainWithGenesisBlock() *BlockChain {
	//创建 或打开 一个数据库
	db, err := storm.Open(BLOCKCHAIN_DB)
	if err != nil {
		log.Panic(err)
	}
	// defer db.Close() // 关闭数据库
	// 创建创世区块
	genesisBlock := NewGenesisBlock("Genesis Block")
	tx, err := db.Begin(true)
	if err != nil {
		log.Panic(err)
	}
	// 将创世区块存储到数据库中
	err = tx.Save(genesisBlock)
	if err != nil {
		log.Panic(err)
	}
	// 存储最新区块的hash值
	err = tx.Set(BLOCK_BUCKET, "l", genesisBlock.Hash)
	if err != nil {
		log.Panic(err)
	}
	tx.Commit()
	// 创建区块链
	return &BlockChain{
		Tip: genesisBlock.Hash,
		DB:  db,
	}
}
