package blc

import (
	"errors"
	"fmt"
)

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
	// 判断 地址是否合法
	from := JsonToArr(s.From)
	to := JsonToArr(s.To)
	for index, v := range from {
		//检查地址是否合法
		if !IsValidAddress(v) || IsValidAddress(to[index]) {
			fmt.Println("Invalid address")
			return errors.New("Invalid address")
		}
	}

	fmt.Printf("sendblock - from: %s - to: %s - amount: %s \n", s.From, s.To, s.Amount)
	return nil
}

// 转账
func (cli *CLI) Send(form, to, amount []string) {
	blockChain := GetBlockChainObject()
	defer blockChain.DB.Close()
	// 转账
	blockChain.MineNewBlock(form, to, amount)
}
