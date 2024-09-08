package blc

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
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
	return nil
}

// 子命令 打印区块链
type PrintChainCommand struct {
	Print bool `short:"p" long:"printchain" description:"Print the blockchain"`
}

func (b *PrintChainCommand) Execute(args []string) error {
	return nil
}

type CLI struct {
	Blc *BlockChain
}

func (cli *CLI) Run() {
	// 实例化顶层命令
	var opts BlockCommand
	_, err := flags.Parse(&opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if opts.AddBlock.Data != "" {
		cli.Blc.AddBlockToBlockChain(opts.AddBlock.Data)
	}
	if opts.PrintChain.Print {
		cli.Blc.PrintChain()
	}
	// fmt.Println("====", opts.AddBlock.Data)
}
