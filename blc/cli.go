package blc

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

// 命令行参数
type BlockCommand struct {
	SendBlock                        SendBlockCommand                        `command:"send" description:"send - from - to - amount"`
	CreateBlockChainWithGenesisBlock CreateBlockChainWithGenesisBlockCommand `command:"createblockchain" description:"Create blockchain with genesis block"`
	AddBlock                         AddBlockCommand                         `command:"addblock" description:"Add a block to the blockchain"`
	PrintChain                       PrintChainCommand                       `command:"printchain" description:"Print all the blocks in the blockchain"`
	GetBalance                       GetBalanceCommand                       `command:"getbalance" description:"Get balance of the address"`
	CreateWallet                     CreateWalletCommand                     `command:"createwallet" description:"Create wallet"`
}

type CLI struct {
}

func (cli *CLI) Run() {
	// 实例化顶层命令
	var opts BlockCommand
	_, err := flags.Parse(&opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// 转账
	if opts.SendBlock.Amount != "" && opts.SendBlock.From != "" && opts.SendBlock.To != "" {
		from := JsonToArr(opts.SendBlock.From)
		to := JsonToArr(opts.SendBlock.To)
		amount := JsonToArr(opts.SendBlock.Amount)
		cli.Send(from, to, amount)
	}
	//创世区块
	if opts.CreateBlockChainWithGenesisBlock.Address != "" {
		cli.CreateBlockChainWithGenesisBlock(opts.CreateBlockChainWithGenesisBlock.Address)
	}
	// 添加区块
	if opts.AddBlock.Data != "" {
		cli.AddBlock(opts.AddBlock.Data)
	}
	// 打印区块链
	if opts.PrintChain.Print {
		cli.PrintChain()
	}
	// 获取余额
	if opts.GetBalance.Address != "" {
		cli.GetBalance(opts.GetBalance.Address)
	}
	// 创建 钱包
	if opts.CreateWallet.Create {
		cli.CreateWallet()
	}
	// fmt.Println("====", opts.AddBlock.Data)
}
