package blc

type TXOutput struct {
	// 转账金额
	Value int64
	// 锁定脚本
	ScriptPubKey string //用户名
}

func (txOutput *TXOutput) UnLockWithAddress(address string) bool {
	return txOutput.ScriptPubKey == address
}
