package blc

import (
	// "fmt"
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
func (bc *BlockChain) AddBlockToBlockChain(data string) error {
	// 从数据库获取最后一个区块的hash值
	tx, err := bc.DB.Begin(true)
	if err != nil {
		log.Panic(err)
		return err
	}
	var lastHash []byte
	err = tx.Get(BLOCK_BUCKET, "l", &lastHash)
	if err != nil {
		log.Panic(err)
		return err
	}
	// fmt.Println("lastHash", lastHash)
	// // 获取最后一个区块
	var lastBlock Block
	err = tx.One("Hash", lastHash, &lastBlock)
	if err != nil {
		log.Panic(err)
		return err
	}
	// fmt.Println("lastBlock", lastBlock)
	// 创建新的区块
	newBlock := NewBlock(lastBlock.Height+1, lastBlock.Hash, data)
	// 将新的区块保存到数据库
	err = tx.Save(newBlock)
	if err != nil {
		log.Panic(err)
		tx.Rollback() // 回滚
		return err

	}
	// 更新区块链的最后一个区块的hash值
	err = tx.Set(BLOCK_BUCKET, "l", newBlock.Hash)
	if err != nil {
		log.Panic(err)
		tx.Rollback() // 回滚
		return err
	}
	return tx.Commit() // 提交
}

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
