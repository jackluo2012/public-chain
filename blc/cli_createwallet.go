package blc

import (
	"fmt"
)

// 子命令 创建钱包
type CreateWalletCommand struct {
	// Create string `short:"c" long:"create" description:"create wallet"`
	Create bool `short:"c" long:"wallet" description:"create wallet"`
}

func (c *CreateWalletCommand) Execute(args []string) error {
	// fmt.Printf("Create wallet: %s\n", c.Create)
	return nil
}

// 创建钱包
func (cli *CLI) CreateWallet() {
	wallets, _ := NewWallets()

	address := wallets.CreateWallet()

	fmt.Printf("Your address is: %s\n", address)
	fmt.Println(len(wallets.WalletsMap))
}
