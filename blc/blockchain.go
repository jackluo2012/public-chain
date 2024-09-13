package blc

import (
	// "fmt"
	"bytes"
	"crypto/ecdsa"
	"errors"
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
	// 在建立新区块之前，需要先验证交易
	for _, tx := range txs {
		if !bc.VerifyTransaction(tx) {
			// log.Panic("交易验证失败")
			// return errors.New("交易验证失败")
		}
	}

	//2. 创建新的区块
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
				fmt.Printf("%s\n", BytesToStr(in.Signature))
			}
			fmt.Println("Vout:")
			for _, out := range tx.Vouts {
				fmt.Printf("%d\n", out.Value)
				fmt.Printf("%s\n", out.Ripemd160Hash)
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
		tx := NewSimpleTransaction(from[i], to[i], StrToInt64(amount[i]), blockChain, txs)
		txs = append(txs, tx)
	}
	// 2、创建新的区块
	blockChain.AddBlockToBlockChain(txs)
	// 3、更新区块链

	return nil
}

// 如果一个地址的UTXO数量大于0，则返回这些UTXO
func (blockChain *BlockChain) UnUTXOs(address string, txs []*Transaction) []*Transaction {
	// var utxos []*UTXO
	var spentTXs map[string][]int64 = make(map[string][]int64)
	var unspentTXs []*Transaction
	// {hash:[0,1,2]}
	//将未打包的交易传入就行计算未花费的输出
	if len(txs) > 0 {
		//倒序遍历交易
		for i := len(txs) - 1; i >= 0; i-- {
			//将未打包的交易传入
			unspentTXs = append(unspentTXs, txs[i])
		}

	}

	// 获取所有的区块
	it := blockChain.Iterator()
	for {
		block := it.Next()
		if block == nil {
			break
		}
		fmt.Printf("正在遍历第 %d 个区块\n", block.Height)
		for _, tx := range block.Txs {
			// txHash
			txHash := BytesToStr(tx.TxHash)
			// Vouts 获取所有的输出
		Outputs:
			for outid, out := range tx.Vouts {
				if spentTXs[txHash] != nil {
					for _, spentOut := range spentTXs[txHash] {
						if spentOut == int64(outid) {
							continue Outputs
						}
					}
				}
				if out.UnLockScriptPubKeyWithAddress(address) {
					unspentTXs = append(unspentTXs, tx)
				}
			}
			if tx.IsCoinbaseTransaction() == false {
				for _, in := range tx.Vins {
					if in.UnLockScript(address) {
						spentTXs[txHash] = append(spentTXs[txHash], in.Vout)
					}
				}
			}
		}
		if len(block.PrevHash) == 0 {
			break
		}
		fmt.Printf("%v\n", unspentTXs)
	}
	return unspentTXs
}

// 转账时查找可用的utxo
func (blockChain *BlockChain) FindSpendableUTXOs(from string, amount int64, txs []*Transaction) (int64, map[string][]int64) {
	// 1、获取所有的UTXO 未花费的交易输出
	unspentTXs := blockChain.UnUTXOs(from, txs)
	unspentOutputs := make(map[string][]int64)

	var accumulated int64
	if amount == -1 {
		accumulated = -2
	}
Work:
	for _, tx := range unspentTXs {
		txid := BytesToStr(tx.TxHash)

		for outid, out := range tx.Vouts {
			if out.UnLockScriptPubKeyWithAddress(from) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txid] = append(unspentOutputs[txid], int64(outid))

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

// 查询余额
func (blockChain *BlockChain) GetBalance(address string) int64 {
	utxos := blockChain.UnUTXOs(address, []*Transaction{})
	var amount int64
	for _, out := range utxos {
		fmt.Printf("未消费的：%v\n", out)
		for _, vout := range out.Vouts {
			if vout.UnLockScriptPubKeyWithAddress(address) {
				amount += vout.Value
			}
		}
	}
	return amount
}

// 对数据区块链进行数字签名
func (blockChain *BlockChain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	// 创世区块不需要签名
	if tx.IsCoinbaseTransaction() {
		return
	}
	prevTXs := make(map[string]*Transaction)
	for _, vin := range tx.Vins {
		// 查找vin中的交易 对应的output hash
		prevTX, err := blockChain.FindTransaction(vin.TxHash)
		if err != nil {
			log.Panic(err)
		} else {
			prevTXs[BytesToStr(prevTX.TxHash)] = prevTX
		}
	}
	tx.Sign(privKey, prevTXs)

}

// 查找花费了的 tx
func (blockChain *BlockChain) FindTransaction(txHash []byte) (*Transaction, error) {
	it := blockChain.Iterator()
	for {
		block := it.Next()
		for _, tx := range block.Txs {
			if bytes.Compare(tx.TxHash, txHash) == 0 {
				return tx, nil
			}
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return nil, errors.New("交易不存在")
}

// VerifyTransaction 验证交易
func (blockChain *BlockChain) VerifyTransaction(tx *Transaction) bool {
	prevTXs := make(map[string]*Transaction)
	for _, vin := range tx.Vins {
		// 查找vin中的交易 对应的output hash
		prevTX, err := blockChain.FindTransaction(vin.TxHash)
		if err != nil {
			log.Panic(err)
		} else {
			prevTXs[BytesToStr(prevTX.TxHash)] = prevTX
		}
	}
	return tx.Verify(prevTXs)
}
