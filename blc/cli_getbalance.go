package blc

import "fmt"

type GetBalanceCommand struct {
	Address string `short:"a" long:"address" description:"address to get balance"`
}

// 获取余额
func (cli *CLI) GetBalance(address string) {
	blockChain := GetBlockChainObject()
	defer blockChain.DB.Close()
	amount := blockChain.GetBalance(address)
	fmt.Printf("address %s balance is %d\n", address, amount)

}
