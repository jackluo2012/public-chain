package main

import (
	// "fmt"
	// "log"

	"public-chain/blc"
)

func main() {
	// 创世区块
	// genesisBlockChain := blc.CreateBlockChainWithGenesisBlock()
	// //关闭数据库
	// defer genesisBlockChain.DB.Close()
	// cli := blc.CLI{Blc: genesisBlockChain}
	// cli.Run()
	cli := blc.CLI{}
	cli.Run()

	// 新增区块
	// genesisBlockChain.AddBlockToBlockChain("Send 100 RMB to zhangsan")
	// genesisBlockChain.AddBlockToBlockChain("Send 200 RMB to lisi")
	// genesisBlockChain.PrintChain()
	// fmt.Println(genesisBlockChain.Blocks)

	// block := blc.NewBlock(1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, "Send 100 RMB to zhangsan")
	// fmt.Printf("block: %d\n", block.Nonce)
	// fmt.Printf("block: %x\n", block.Hash)

	// bytes := block.Serialize()
	// fmt.Printf("block: %x\n", bytes)
	// block2 := block.Deserialize(bytes)
	// fmt.Printf("block2: %x\n", block2.Hash)
	// fmt.Printf("block2: %d\n", block2.Nonce)

	// proofOfWork := blc.NewProofOfWork(block)
	// fmt.Printf("pow: %v\n", proofOfWork.IsValid())
}
