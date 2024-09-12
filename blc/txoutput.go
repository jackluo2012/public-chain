package blc

import "bytes"

type TXOutput struct {
	// 转账金额
	Value int64
	// 锁定脚本
	Ripemd160Hash []byte //公钥
}

// 解锁
func (txOutput *TXOutput) UnLockScriptPubKeyWithAddress(address string) bool {
	//将地址转换成字节数组
	addressByte := Base58Encode([]byte(address)) //公钥
	// 1 version，20 字节 4 checksum
	// 中字就是 reipemd160Hash
	ripemd160Hash := addressByte[1 : len(addressByte)-4]
	return bytes.Compare(txOutput.Ripemd160Hash, ripemd160Hash) == 0
}

// 上 锁
/**
 * 参考 ，wallet.go 中的 GetAddress() 方法
 */
func (txOutput *TXOutput) Lock(address string) {
	//将地址转换成字节数组
	addressByte := Base58Encode([]byte(address)) //公钥
	// 1 version，20 字节 4 checksum
	// 中字就是 reipemd160Hash
	txOutput.Ripemd160Hash = addressByte[1 : len(addressByte)-4]

}

/**
 *	新建交易输出
 */
func NewTXOutput(value int64, address string) *TXOutput {
	txOutput := &TXOutput{value, nil}
	// 设置锁定脚本 Ripemd160Hash
	txOutput.Lock(address)
	return txOutput
}
