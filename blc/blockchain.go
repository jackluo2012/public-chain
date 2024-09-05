package blc

type BlockChain struct {
	Blocks []*Block // 区块链由多个区块组成
}

// 创建带有创世区块的区块链
func CreateBlockChainWithGenesisBlock() *BlockChain {
	// 创建创世区块
	genesisBlock := NewGenesisBlock("Genesis Block")
	// 创建区块链
	return &BlockChain{
		Blocks: []*Block{genesisBlock},
	}
}
