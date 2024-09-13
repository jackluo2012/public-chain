package blc

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"math/big"
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
	txInput := &TXInput{[]byte{}, -1, nil, []byte{}}
	// 2.输出
	txOutput := NewTXOutput(10, address)
	// 3.交易hash
	tx := &Transaction{[]byte{}, []*TXInput{txInput}, []*TXOutput{txOutput}}
	// 4.序列化 设置hash 值
	tx.HashTransaction()

	//进行数字签名

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
	wallets, _ := NewWallets()
	wallet := wallets.WalletsMap[from]

	// 通unSpentTxs 返回 from这个人所有的未花费的输出
	money, spendableUTXOs := bc.FindSpendableUTXOs(from, amount, txs)

	for txHash, indexArray := range spendableUTXOs {
		for _, index := range indexArray {
			// 代表消费
			txInput := &TXInput{StrToBytes(txHash), index, wallet.PublicKey, nil}
			txInputs = append(txInputs, txInput)
		}
	}
	fmt.Printf("txInputs:%v\n", txInputs)
	fmt.Println(len(txInputs))

	// 转账
	txOutput := NewTXOutput(amount, to)
	txOutputs = append(txOutputs, txOutput)
	// 3.找零
	txOutput = NewTXOutput(money-amount, from)
	txOutputs = append(txOutputs, txOutput)
	// 4.创建交易
	tx := &Transaction{[]byte{}, txInputs, txOutputs}
	// 5.设置hash
	tx.HashTransaction()
	// 6.进行数字签名
	bc.SignTransaction(tx, wallet.PrivateKey)
	// 7.返回
	return tx
}

// Sign
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]*Transaction) {
	// 1.判断是否是coinbase交易
	if tx.IsCoinbaseTransaction() {
		return
	}
	// 2.获取交易输入
	for _, vin := range tx.Vins {
		//当前的txinput 没有找到 对应的交易
		if prevTXs[BytesToStr(vin.TxHash)] == nil {
			panic("prevTXs[hex.EncodeToString(vin.TxHash)] == nil")
		}
	}
	// 备份 tx
	txCopy := tx.TrimmedCopy()
	// 3.遍历txCopy 的输入
	for inID, vin := range txCopy.Vins {
		// 4.找到vin 对应的交易
		prevTx := prevTXs[BytesToStr(vin.TxHash)]
		// 5.设置txCopy 的输入的签名为nil
		txCopy.Vins[inID].Signature = nil
		// 6.设置txCopy 的输入的公钥为对应交易的公钥
		txCopy.Vins[inID].PublicKey = prevTx.Vouts[vin.Vout].Ripemd160Hash
		// 7.计算txCopy 的hash
		txCopy.TxHash = txCopy.Hash()
		// 8.对txCopy 进行数字签名
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.TxHash)
		if err != nil {
			panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)
		fmt.Printf("signature:%x\n", signature)
		// 9.设置当前tx 的输入的签名
		tx.Vins[inID].Signature = signature
	}
}

// 2.创建txCopy 用于签名
func (tx *Transaction) TrimmedCopy() *Transaction {
	var inputs []*TXInput
	var outputs []*TXOutput
	txCopy := &Transaction{tx.TxHash, inputs, outputs}
	for _, vin := range tx.Vins {
		txCopy.Vins = append(txCopy.Vins, &TXInput{vin.TxHash, vin.Vout, nil, nil})
	}
	for _, vout := range tx.Vouts {
		txCopy.Vouts = append(txCopy.Vouts, &TXOutput{vout.Value, vout.Ripemd160Hash})
	}
	return txCopy
}

// 3.生成hash
func (tx *Transaction) Hash() []byte {

	txCopy := tx
	txCopy.TxHash = []byte{}
	// 1.序列化
	return Sha256Hash(txCopy.Serialize())
}

// tx Serialize
func (tx *Transaction) Serialize() []byte {
	var result bytes.Buffer
	// 创建编码器 打包
	encoder := gob.NewEncoder(&result)
	// 编码
	err := encoder.Encode(tx)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}

// Verify Transaction -- 数字签名 验证
func (tx *Transaction) Verify(prevTXs map[string]*Transaction) bool {
	// 1.判断是否是coinbase交易
	if tx.IsCoinbaseTransaction() {
		return true
	}
	// 2.获取交易输入
	for _, vin := range tx.Vins {
		//当前的txinput 没有找到 对应的交易
		if prevTXs[BytesToStr(vin.TxHash)] == nil {
			panic("prevTXs[hex.EncodeToString(vin.TxHash)] == nil")
		}
	}
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	// 3.遍历txCopy 的输入
	for inID, vin := range txCopy.Vins {
		// 4.找到vin 对应的交易
		prevTx := prevTXs[BytesToStr(vin.TxHash)]
		// 5.设置txCopy 的输入的签名为nil
		txCopy.Vins[inID].Signature = nil
		// 6.设置txCopy 的输入的公钥为对应交易的公钥
		txCopy.Vins[inID].PublicKey = prevTx.Vouts[vin.Vout].Ripemd160Hash
		// 7.计算txCopy 的hash
		txCopy.TxHash = txCopy.Hash()
		txCopy.Vins[inID].PublicKey = nil

		//  私钥 ID
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])
		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PublicKey)
		x.SetBytes(vin.PublicKey[:(keyLen / 2)])
		y.SetBytes(vin.PublicKey[(keyLen / 2):])
		// 7.生成公钥
		publicKey := ecdsa.PublicKey{curve, &x, &y}
		// 8.验证数字签名
		if !ecdsa.Verify(&publicKey, txCopy.TxHash, &r, &s) {
			return false
		}
	}
	return true
}
