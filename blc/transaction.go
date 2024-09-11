package blc

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

/**
 * 创建分两种 情况
 * 1.创世区块 Transaction
 * 2.普通区块 转账时 Transaction
 */
// UTXO
type Transaction struct {
	// 交易hash
	TxHash []byte
	// 输入
	Vins []*TXInput
	// 输出
	Vouts []*TXOutput
}

// 是否是coinbase交易
func (tx *Transaction) IsCoinbaseTransaction() bool {
	fmt.Printf("len(tx.Vins):%v\n", tx.Vins)
	return len(tx.Vins[0].TxHash) == 0 && tx.Vins[0].Vout == -1
}

// Transaction 序列化
func (tx *Transaction) HashTransaction() {
	var result bytes.Buffer
	// 创建编码器 打包
	encoder := gob.NewEncoder(&result)
	// 编码
	err := encoder.Encode(tx)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(result.Bytes())
	tx.TxHash = hash[:]
}

// 1.创世区块 Transaction
func NewCoinbaseTx(address string) *Transaction {
	// 1.输入
	txInput := &TXInput{[]byte{}, -1, "Genesis Data"}
	// 2.输出
	txOutput := &TXOutput{10, address}
	// 3.交易hash
	tx := &Transaction{[]byte{}, []*TXInput{txInput}, []*TXOutput{txOutput}}
	// 4.序列化 设置hash 值
	tx.HashTransaction()
	return tx
}

// 1.有一个函数，返回from 这个人所有的未花费的输出 所对应的Transaction
// 2.
// 2.普通区块 转账时 Transaction
func NewSimpleTransaction(from, to string, amount int64, bc *BlockChain, txs []*Transaction) *Transaction {
	// 1.获取所有未花费的输出
	// 2.创建交易输入
	var txInputs []*TXInput
	var txOutputs []*TXOutput

	// 通unSpentTxs 返回 from这个人所有的未花费的输出
	money, spendableUTXOs := bc.FindSpendableUTXOs(from, amount, txs)

	for txHash, indexArray := range spendableUTXOs {
		for _, index := range indexArray {
			// 代表消费
			txInput := &TXInput{StrToBytes(txHash), index, from}
			txInputs = append(txInputs, txInput)
		}
	}
	fmt.Printf("txInputs:%v\n", txInputs)
	fmt.Println(len(txInputs))

	// 转账
	txOutput := &TXOutput{amount, to}
	txOutputs = append(txOutputs, txOutput)
	// 3.找零
	txOutput = &TXOutput{money - amount, from}
	txOutputs = append(txOutputs, txOutput)
	// 4.创建交易
	tx := &Transaction{[]byte{}, txInputs, txOutputs}
	// 5.设置hash
	tx.HashTransaction()
	return tx
}
