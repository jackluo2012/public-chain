package main

import (
	"fmt"
	"public-chain/blc"
)

func main() {
	// 创世区块
	genesisBlockChain := blc.CreateBlockChainWithGenesisBlock()

	// 新增区块
	genesisBlockChain.AddBlockToBlockChain("Send 100 RMB to zhangsan")
	genesisBlockChain.AddBlockToBlockChain("Send 200 RMB to lisi")
	fmt.Println(genesisBlockChain.Blocks)

	block := blc.NewBlock(1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, "Send 100 RMB to zhangsan")
	fmt.Printf("block: %d\n", block.Nonce)
	fmt.Printf("block: %x\n", block.Hash)

	proofOfWork := blc.NewProofOfWork(block)
	fmt.Printf("pow: %v\n", proofOfWork.IsValid())
}
