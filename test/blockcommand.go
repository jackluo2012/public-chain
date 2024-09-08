package test

import (
	"fmt"
)

// 命令行参数
type BlockCommand struct {
	AddBlock   AddBlockCommand   `command:"addblock" description:"Add a block to the blockchain"`
	PrintChain PrintChainCommand `command:"printchain" description:"Print all the blocks in the blockchain"`
}

// 子命令 添加区块
type AddBlockCommand struct {
	Data string `short:"d" long:"data" description:"Data for the block"`
}

func (b *AddBlockCommand) Execute(args []string) error {
	fmt.Printf("%s\n", b.Data)
	return nil
}

// 子命令 打印区块链
type PrintChainCommand struct {
	Print bool `short:"p" long:"printchain" description:"Print the blockchain"`
}

func (b *PrintChainCommand) Execute(args []string) error {
	if b.Print {
		// 打印区块链
		fmt.Println("Print the blockchain")
	}
	return nil
}
