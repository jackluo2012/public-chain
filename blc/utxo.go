package blc

// 未花费的输出
type UTXO struct {
	TxHash    []byte //交易哈希
	Index     int64  //索引
	Output    *TXOutput
	BlockHash []byte //区块哈希
}
