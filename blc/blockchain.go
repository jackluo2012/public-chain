package blc

type BlockChain struct {
	Blocks []*Block // 区块链由多个区块组成
}

// 添加区块到区块链
func (bc *BlockChain) AddBlockToBlockChain(data string) {
	// 获取最后一个区块
	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	// 创建新的区块
	newBlock := NewBlock(lastBlock.Height+1, lastBlock.Hash, data)
	// 将新的区块添加到区块链
	bc.Blocks = append(bc.Blocks, newBlock)
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
