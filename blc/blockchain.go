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
func (bc *BlockChain) AddBlockToBlockChain(txs []*Transaction) error {
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
	newBlock := NewBlock(lastBlock.Height+1, lastBlock.Hash, txs)
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
func CreateBlockChainWithGenesisBlock(address string) *BlockChain {
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
	txCoinbase := NewCoinbaseTx(address)
	// 创建创世区块
	genesisBlock := NewGenesisBlock([]*Transaction{txCoinbase})
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
	return &BlockChain{Tip: tip, DB: db}
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
		fmt.Printf("Timestamp:%s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 15:04:05"))
		fmt.Printf("Hash:%x\n", block.Hash)
		fmt.Printf("Nonce:%d\n", block.Nonce)
		fmt.Println("Txs:")
		for _, tx := range block.Txs {
			fmt.Printf("%x\n", tx.TxHash)
			fmt.Println("Vin:")
			for _, in := range tx.Vins {
				fmt.Printf("%x\n", in.TxHash)
				fmt.Printf("%d\n", in.Vout)
				fmt.Printf("%s\n", in.ScriptSig)
			}
			fmt.Println("Vout:")
			for _, out := range tx.Vouts {
				fmt.Printf("%d\n", out.Value)
				fmt.Printf("%s\n", out.ScriptPubKey)
			}
		}
		fmt.Println("--------------------------------------------------\n")

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

// 拿到最的区块
func (blockChain *BlockChain) GetLastBlock() *Block {
	var lastBlock Block
	err := blockChain.DB.Get(BLOCK_BUCKET, blockChain.Tip, &lastBlock)
	if err != nil {
		log.Panic(err)
	}
	return &lastBlock
}

// 获取 未花费的输出

// 如果一个地地对应的TXOut未花费，那么这个Transaction就应该添加到数据中返回
func (blockChain *BlockChain) UnSpentTransationsWithAddress(address string) []*TXOutput {
	// var unSpentTxs []*Transaction
	var unUTXOs []*TXOutput
	spentTXOutputs := make(map[string][]int64) // key是txid，value是vout的索引
	// {hash:[0,1,2]}
	// 获取所有的区块
	it := blockChain.Iterator()
	for {
		block := it.Next()
		fmt.Printf("正在遍历第%v个区块\n", block)

		for _, tx := range block.Txs {
			// txHash
			// Vins 获取所有的输入,消费的UTXO
			if !tx.IsCoinbaseTransaction() {
				for _, in := range tx.Vins {
					// 是否能解锁
					if in.UnLockWithAddress(address) {
						// unSpentTxs = append(unSpentTxs, tx)
						spentTXOutputs[BytesToStr(in.TxHash)] = append(spentTXOutputs[BytesToStr(in.TxHash)], in.Vout) //
					}
				}
			}
			fmt.Printf("消费的：%v\n", spentTXOutputs)
			// Vouts 获取所有的输出
			for index, out := range tx.Vouts {
				if out.UnLockWithAddress(address) {
					fmt.Printf("未消费的：%v\n", out)
					fmt.Printf("len：%v\n", len(spentTXOutputs))
					if spentTXOutputs != nil {
						if len(spentTXOutputs) != 0 {
							for txHash, indexArray := range spentTXOutputs {
								if txHash == BytesToStr(tx.TxHash) {
									for _, i := range indexArray {
										if index == int(i) && txHash == BytesToStr(tx.TxHash) {
											continue
										} else {
											unUTXOs = append(unUTXOs, out)
										}
									}
								}
							}
						} else {
							unUTXOs = append(unUTXOs, out)
						}
					}
				}
			}
		}

		var hashInt big.Int
		hashInt.SetBytes(block.PrevHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
	return unUTXOs
}

// 挖掘新的区块
func (blockChain *BlockChain) MineNewBlock(from []string, to []string, amount []string) *BlockChain {
	fmt.Println("正在挖矿...")
	fmt.Println(from)
	fmt.Println(to)
	fmt.Println(amount)

	// 1、通过相关参数创建一个新的交易 Transaction 数组

	// 1、通过相关参数创建一个新的交易 Transaction 数组
	var txs []*Transaction
	for i := 0; i < len(from); i++ {
		// 1、创建新的交易
		tx := NewSimpleTransaction(from[i], to[i], StrToInt64(amount[i]))
		txs = append(txs, tx)
	}
	// 2、创建新的区块
	blockChain.AddBlockToBlockChain(txs)
	// 3、更新区块链

	return nil
}

// 如果一个地址的UTXO数量大于0，则返回这些UTXO
func (blockChain *BlockChain) UnUTXOs(address string) []*UTXO {
	// var utxos []*UTXO
	var unUTXOs []*UTXO
	spentTXOutputs := make(map[string][]int64) // key是txid，value是vout的索引
	// {hash:[0,1,2]}
	// 获取所有的区块
	it := blockChain.Iterator()
	for {
		block := it.Next()
		fmt.Printf("正在遍历第%v个区块\n", block)

		for _, tx := range block.Txs {
			// txHash
			// Vins 获取所有的输入,消费的UTXO
			if !tx.IsCoinbaseTransaction() {
				for _, in := range tx.Vins {
					// 是否能解锁
					if in.UnLockWithAddress(address) {
						// unSpentTxs = append(unSpentTxs, tx)
						spentTXOutputs[BytesToStr(in.TxHash)] = append(spentTXOutputs[BytesToStr(in.TxHash)], in.Vout) //
					}
				}
			}
			// Vouts 获取所有的输出
			for index, out := range tx.Vouts {
				if out.UnLockWithAddress(address) {
					fmt.Printf("len：%v\n", len(spentTXOutputs))
					if spentTXOutputs != nil {
						if len(spentTXOutputs) != 0 {

							for txHash, indexArray := range spentTXOutputs {
								for _, i := range indexArray {
									if index == int(i) && txHash == BytesToStr(tx.TxHash) {
										continue
									} else {
										utxo := &UTXO{tx.TxHash, int64(index), out, block.Hash}
										unUTXOs = append(unUTXOs, utxo)
									}
								}
							}
						} else {
							utxo := &UTXO{tx.TxHash, int64(index), out, block.Hash}
							unUTXOs = append(unUTXOs, utxo)
						}
					}
				}
			}
		}

		var hashInt big.Int
		hashInt.SetBytes(block.PrevHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
	return unUTXOs
}

// 查询余额
func (blockChain *BlockChain) GetBalance(address string) int64 {
	utxos := blockChain.UnUTXOs(address)
	var amount int64
	for _, out := range utxos {
		fmt.Printf("未消费的：%v\n", out)
		amount += out.Output.Value
	}
	return amount
}
