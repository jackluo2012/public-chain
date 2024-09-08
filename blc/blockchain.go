package blc

import (
	// "fmt"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

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

// 迭代器
func (bc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{bc.Tip, bc.DB}
}

// 添加区块到区块链
func (bc *BlockChain) AddBlockToBlockChain(data string) error {
	// 从数据库获取最后一个区块的hash值
	tx, err := bc.DB.Begin(true)
	if err != nil {
		log.Panic(err)
		return err
	}
	// // 获取最后一个区块
	var lastBlock Block
	err = tx.One("Hash", bc.Tip, &lastBlock)
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
	bc.Tip = newBlock.Hash
	return tx.Commit() // 提交
}

// 判断数据库是否存在
func IsDBExist() bool {
	if _, err := os.Stat(BLOCKCHAIN_DB); os.IsNotExist(err) {
		return false
	}
	return true
}

// 创建带有创世区块的区块链
func CreateBlockChainWithGenesisBlock(data string) {
	var (
		tip []byte
		db  *storm.DB
		err error
	)
	if IsDBExist() {
		fmt.Println("创世区块已经存在，无需再次创建")
		os.Exit(1)
	}
	//创建 或打开 一个数据库
	db, err = storm.Open(BLOCKCHAIN_DB)
	if err != nil {
		log.Panic(err)
	}
	// defer db.Close() // 关闭数据库
	// 创建创世区块
	genesisBlock := NewGenesisBlock(data)
	tx, err := db.Begin(true)
	if err != nil {
		log.Panic(err)
	}
	// 将创世区块存储到数据库中
	err = tx.Save(genesisBlock)
	if err != nil {
		log.Panic(err)
	}
	tip = genesisBlock.Hash
	// 存储最新区块的hash值
	err = tx.Set(BLOCK_BUCKET, "l", tip)
	if err != nil {
		log.Panic(err)
	}
	tx.Commit()

}

// 输出所有区块链
func (bc *BlockChain) PrintChain() {

	//简单的遍历区块链 1
	// var blocks []*Block
	// bc.DB.All(&blocks)
	// for _, block := range blocks {
	// 	fmt.Println(block.Hash)
	// 	fmt.Printf("%s\n", block.Data)
	// }
	// 2
	// currentHash := bc.Tip
	// for {
	// 	var block Block
	// 	err := bc.DB.One("Hash", currentHash, &block)
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}
	// 	fmt.Printf("Height:%d\n", block.Height)
	// 	fmt.Printf("PrevHash:%x\n", block.PrevHash)
	// 	fmt.Printf("Data:%s\n", block.Data)

	// 	fmt.Printf("Timestamp:%s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 15:04:05"))
	// 	fmt.Printf("Hash:%x\n", block.Hash)
	// 	fmt.Printf("Nonce:%d\n", block.Nonce)

	// 	var hashInt big.Int
	// 	hashInt.SetBytes(block.PrevHash)
	// 	if big.NewInt(0).Cmp(&hashInt) == 0 {
	// 		break
	// 	}
	// 	currentHash = block.PrevHash
	// }
	// 3
	blockChainIterator := bc.Iterator()

	for {
		block := blockChainIterator.Next()
		fmt.Println()
		fmt.Printf("Height:%d\n", block.Height)
		fmt.Printf("PrevHash:%x\n", block.PrevHash)
		fmt.Printf("Data:%s\n", block.Data)

		fmt.Printf("Timestamp:%s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 15:04:05"))
		fmt.Printf("Hash:%x\n", block.Hash)
		fmt.Printf("Nonce:%d\n", block.Nonce)
		fmt.Println()

		var hashInt big.Int
		hashInt.SetBytes(block.PrevHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}

}

// 返回 blockchain对象
func GetBlockChainObject() *BlockChain {
	var (
		tip []byte
		db  *storm.DB
		err error
	)
	//创建 或打开 一个数据库
	db, err = storm.Open(BLOCKCHAIN_DB)
	if err != nil {
		log.Panic(err)
		return nil
	}

	// 获取最后一个区块的hash值
	err = db.Get(BLOCK_BUCKET, "l", &tip)
	if err != nil {
		log.Panic(err)
	}
	// 创建区块链
	return &BlockChain{
		Tip: tip,
		DB:  db,
	}
}
