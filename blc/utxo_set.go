package blc

import (
	"bytes"
	"fmt"
	"log"
)

// 1. 遍历整个数据库 UTXO集合管理全部的未花费输出，然后将所有的UTXO存储到数据库中
// reset the UTXO set
// 2. 每当有新的交易产生时，遍历整个数据库，将所有与该交易地址相关的UTXO进行更新
// reset
// 去遍历数据时，
// []*TxOutputs

// 表名
const utxoTableName = "UTXO"

type UTXOSet struct {
	Blockchain *BlockChain
}

// 重置UTXO集合
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

// 查找未打包的交易的utxo
func (utxoSet *UTXOSet) FindUnpackageSpendableUTXOs(address string, txs []*Transaction) []*UTXO {
	var unUTXOs []*UTXO
	spentTXOutputs := make(map[string][]int64)

	// 1. 查找未打包的交易
	for _, tx := range txs {
		if tx.IsCoinbaseTransaction() == false {
			for _, in := range tx.Vins {
				if in.UnLockScript(address) {
					key := BytesToStr(in.TxHash)
					spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
				}
			}
		}
	}
	for _, tx := range txs {
	Work1:
		for index, out := range tx.Vouts {
			if out.UnLockScriptPubKeyWithAddress(address) {
				if len(spentTXOutputs[BytesToStr(tx.TxHash)]) > 0 {
					fmt.Println(address)
					if len(spentTXOutputs[BytesToStr(tx.TxHash)]) == 0 {
						utxo := &UTXO{TxHash: tx.TxHash, Index: int64(index), Output: out}
						unUTXOs = append(unUTXOs, utxo)
					} else {
						for hash, indexArray := range spentTXOutputs {
							txHashStr := BytesToStr(tx.TxHash)
							if hash == txHashStr {
								var isUnSpentUTXO bool
								for _, outIndex := range indexArray {
									if index == int(outIndex) {
										isUnSpentUTXO = true
										continue Work1
									}
									if isUnSpentUTXO == false {
										utxo := &UTXO{TxHash: tx.TxHash, Index: int64(index), Output: out}
										unUTXOs = append(unUTXOs, utxo)
									}
								}

							} else {
								//处理未打包的交易
								utxo := &UTXO{TxHash: tx.TxHash, Index: int64(index), Output: out}
								unUTXOs = append(unUTXOs, utxo)
							}

						}
					}
				}
			}
		}
	}
	return unUTXOs
}

//FindSpendableUTXOs
// 查找可花费的UTXO
/**
* address: 查询的地址
* amount: 需要花费的金额
* txs: 未打包的交易
*  * 返回值：要凑的金额，可花费的UTXO
 */
func (utxoSet *UTXOSet) FindSpendableUTXOs(address string, amount int64, txs []*Transaction) (int64, map[string][]int64) {
	unPackageUTXOS := utxoSet.FindUnpackageSpendableUTXOs(address, txs)
	spentableUTXO := make(map[string][]int64)
	var money int64
	//如果未打包的交易金额大于要花费的金额，则直接返回
	for _, utxo := range unPackageUTXOS {
		money += utxo.Output.Value
		txHash := BytesToStr(utxo.TxHash)
		spentableUTXO[txHash] = append(spentableUTXO[txHash], utxo.Index)
		if money >= amount {
			break
		}
	}
	// 如果未打包的交易金额小于要花费的金额，则从数据库中查找
	if money < amount {
		utxos := utxoSet.FindUTXOMapForAddress(address)
		for _, utxo := range utxos {
			money += utxo.Output.Value
			txHash := BytesToStr(utxo.TxHash)
			spentableUTXO[txHash] = append(spentableUTXO[txHash], utxo.Index)
			if money >= amount {
				break
			}
		}
	}

	return money, spentableUTXO
}

// Update
// 更新UTXO集合
func (utxoSet *UTXOSet) Update(block *Block) {
	// 1. 遍历区块中的交易
	// blocks -> 最新 的区块
	bock := utxoSet.Blockchain.Iterator().Next()
	// spentUTXO := make(map[string][]int64)
	ins := []*TXInput{}
	// outsMap := make(map[string][]*TXOutput)
	// 找到所有我要删除的数据
	for _, tx := range bock.Txs {
		for _, in := range tx.Vins {
			// 2. 将输入的UTXO标记为已花费
			// spentUTXO[BytesToStr(in.TxHash)] = append(spentUTXO[BytesToStr(in.TxHash)], in.Vout)
			ins = append(ins, in)
			txID := BytesToStr(in.TxHash)
			utxo := []*UTXO{} //make([]*UTXO,0)
			utxos := []*UTXO{}
			err := utxoSet.Blockchain.DB.Get(utxoTableName, txID, &utxo)
			if err != nil {
				log.Panic(err)
			}
			// 判断是否需要
			isNeedDelete := false
			for _, utxo := range utxo {
				//删除使用过的UTXO
				if in.Vout == utxo.Index && bytes.Compare(utxo.Output.Ripemd160Hash, Ripemd160Hash(in.PublicKey)) == 0 {
					// utxoSet.Blockchain.DB.Delete(utxoTableName, txID)
					isNeedDelete = true
				} else {
					utxos = append(utxos, utxo)
				}
			}
			if isNeedDelete {
				utxoSet.Blockchain.DB.Delete(utxoTableName, txID)
				if len(utxos) > 0 {
					utxoSet.Blockchain.DB.Set(utxoTableName, txID, utxos)
				}
			}
		}
		utxos := []*UTXO{}
		// 未花费的UTXO
		for index, out := range tx.Vouts {
			isSpent := false
			for _, in := range ins {
				if in.Vout == int64(index) && bytes.Compare(tx.TxHash, in.TxHash) == 0 && bytes.Compare(out.Ripemd160Hash, Ripemd160Hash(in.PublicKey)) == 0 {
					isSpent = true
					continue
				}
			}
			if !isSpent {
				utxos = append(utxos, &UTXO{TxHash: tx.TxHash, Index: int64(index), Output: out})
			}
		}
		// 3. 将输出的UTXO添加到UTXO集合中
		if len(utxos) > 0 {
			utxoSet.Blockchain.DB.Set(utxoTableName, BytesToStr(tx.TxHash), utxos)
		}
	}

}
