package blc

import (
	"errors"
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
}

// 子命令 添加区块
type SendBlockCommand struct {
	From   string `short:"f" long:"from" description:"Source address"`
	To     string `short:"t" long:"to" description:"Destination address"`
	Amount string `short:"a" long:"amount" description:"Amount to send"`
}

func (s *SendBlockCommand) Execute(args []string) error {
	if s.From == "" || s.To == "" || s.Amount == "" {
		fmt.Println("Invalid command, please check your command")
		return errors.New("Invalid command, please check your command")
	}
	fmt.Printf("sendblock - from: %s - to: %s - amount: %s \n", s.From, s.To, s.Amount)
	return nil
}

type CreateBlockChainWithGenesisBlockCommand struct {
	Address string `short:"a" long:"address" description:"address for the genesis block"`
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

type GetBalanceCommand struct {
	Address string `short:"a" long:"address" description:"address to get balance"`
}

type CLI struct {
}

// 添加区块
func (cli *CLI) AddBlock(data string) {
	blockChain := GetBlockChainObject()
	defer blockChain.DB.Close()
	blockChain.AddBlockToBlockChain([]*Transaction{})
}
func (cli *CLI) CreateBlockChainWithGenesisBlock(address string) {

	blockChain := CreateBlockChainWithGenesisBlock(address)
	defer blockChain.DB.Close()
}

// 打印区块链
func (cli *CLI) PrintChain() {
	blockChain := GetBlockChainObject()
	defer blockChain.DB.Close()
	blockChain.PrintChain()
}

// 转账
func (cli *CLI) Send(form, to, amount []string) {
	blockChain := GetBlockChainObject()
	defer blockChain.DB.Close()
	// 转账
	blockChain.MineNewBlock(form, to, amount)
}

// 获取余额
func (cli *CLI) GetBalance(address string) {
	blockChain := GetBlockChainObject()
	defer blockChain.DB.Close()
	amount := blockChain.GetBalance(address)
	fmt.Printf("address %s balance is %d\n", address, amount)

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
	// fmt.Println("====", opts.AddBlock.Data)
}
