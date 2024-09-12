package blc

import "bytes"

type TXInput struct {
	// 交易hash
	TxHash []byte
	// 索引
	Vout int64
	//解锁脚本
	// ScriptSig string
	// 数字 签名
	Signature []byte
	// 公钥,钱包里面的，2次256hash后的结果
	PublicKey []byte //原生的公钥
}

// 判断解锁脚本是否匹配 消费的钱是 谁的
func (txInput *TXInput) UnLockRipemd160Hash(ripemd160Hash []byte) bool {
	publicKey := Ripemd160Hash(txInput.PublicKey)

	return bytes.Compare(publicKey, ripemd160Hash) == 0
}

/**
*	直接判断公钥是否匹配
 */
func (txInput *TXInput) UnLockScript(address string) bool {
	addressByte := Base58Encode([]byte(address)) //公钥
	// 1 version，20 字节 4 checksum
	// 中字就是 reipemd160Hash
	ripemd160Hash := addressByte[1 : len(addressByte)-4]
	publicKey := Ripemd160Hash(txInput.PublicKey)
	return bytes.Compare(publicKey, ripemd160Hash) == 0
}
