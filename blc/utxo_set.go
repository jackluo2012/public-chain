package blc

import "log"

// 1. 遍历整个数据库 UTXO集合管理全部的未花费输出，然后将所有的UTXO存储到数据库中
// reset the UTXO set
// 2. 每当有新的交易产生时，遍历整个数据库，将所有与该交易地址相关的UTXO进行更新
// reset
// 去遍历数据时，
// []*TxOutputs

//表名
const utxoTableName = "UTXO"

type UTXOSet struct {
	Blockchain *BlockChain
}

//重置UTXO集合
// 重置数据库表
func (utxoSet *UTXOSet) ResetUTXOSet() {
	//删除 表
	utxoSet.Blockchain.DB.Drop(utxoTableName)

	// 查找TXOutputs [string]*TxOutputs
	txOutputsMap := utxoSet.Blockchain.FindUTXOMap()
	for txID, outputs := range txOutputsMap {
		err := utxoSet.Blockchain.DB.Set(utxoTableName, txID, outputs)
		if err != nil {
			log.Panic(err)
		}
	}
}

// FindUTXOMap
// 查找UTXO集合
func (utxoSet *UTXOSet) FindUTXOMapForAddress(address string) []*UTXO {
	utxoMap := make([]*UTXO, 0)
	// 遍历数据库
	err := utxoSet.Blockchain.DB.Get(utxoTableName, address, &utxoMap)
	if err != nil {
		log.Panic(err)
	}
	return utxoMap
}

// getBalance
// 查询余额
func (utxoSet *UTXOSet) GetBalance(address string) int64 {

	utxoMap := utxoSet.FindUTXOMapForAddress(address)
	var balance int64

	for _, utxo := range utxoMap {
		balance += utxo.Output.Value
	}

	return balance
}
