package blc

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
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
type TXInput struct {
	// 交易hash
	TxHash []byte
	// 索引
	Index int64
	//解锁脚本
	ScriptSig string
}

type TXOutput struct {
	// 转账金额
	Value int64
	// 锁定脚本
	ScriptPubKey string //用户名
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
