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
}
