package blc

import (
	"errors"
	"fmt"
)

type GetBalanceCommand struct {
	Address string `short:"a" long:"address" description:"address to get balance"`
}

func (b *GetBalanceCommand) Execute(args []string) error {
	if IsValidAddress(b.Address) == false {
		fmt.Println("地址无效")
		return errors.New("Invalid address")
	}
	return nil
}

// 获取余额
func (cli *CLI) GetBalance(address string) {
	blockChain := GetBlockChainObject()
	defer blockChain.DB.Close()
	amount := blockChain.GetBalance(address)
	fmt.Printf("address %s balance is %d\n", address, amount)

}
