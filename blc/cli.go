package blc

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

// 命令行参数
type BlockCommand struct {
	CreateBlockChainWithGenesisBlock CreateBlockChainWithGenesisBlockCommand `command:"createblockchain" description:"Create blockchain with genesis block"`
	AddBlock                         AddBlockCommand                         `command:"addblock" description:"Add a block to the blockchain"`
	PrintChain                       PrintChainCommand                       `command:"printchain" description:"Print all the blocks in the blockchain"`
}

// 子命令 添加区块
type CreateBlockChainWithGenesisBlockCommand struct {
	Data string `short:"d" long:"data" description:"Data for the genesis block"`
}

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
}

// 添加区块
func (cli *CLI) AddBlock(data string) {
	blockChain := GetBlockChainObject()
	defer blockChain.DB.Close()
	blockChain.AddBlockToBlockChain(data)
}

// 打印区块链
func (cli *CLI) PrintChain() {
	blockChain := GetBlockChainObject()
	defer blockChain.DB.Close()
	blockChain.PrintChain()
}
func (cli *CLI) Run() {
	// 实例化顶层命令
	var opts BlockCommand
	_, err := flags.Parse(&opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if opts.CreateBlockChainWithGenesisBlock.Data != "" {
		CreateBlockChainWithGenesisBlock(opts.CreateBlockChainWithGenesisBlock.Data)
	}
	if opts.AddBlock.Data != "" {
		cli.AddBlock(opts.AddBlock.Data)
	}
	if opts.PrintChain.Print {
		cli.PrintChain()
	}
	// fmt.Println("====", opts.AddBlock.Data)
}
