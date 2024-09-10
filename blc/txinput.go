package blc

type TXInput struct {
	// 交易hash
	TxHash []byte
	// 索引
	Vout int64
	//解锁脚本
	ScriptSig string
}

// 判断解锁脚本是否匹配 消费的钱是 谁的
func (txInput *TXInput) UnLockWithAddress(address string) bool {
	return txInput.ScriptSig == address
}
