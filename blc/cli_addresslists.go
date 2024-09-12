package blc

import "fmt"

type AddressListsCommand struct {
	All bool `short:"a" long:"all" description:"list all address lists"`
}

func (b *AddressListsCommand) Execute(args []string) error {

	return nil
}

func (cli *CLI) listAddressLists() []string {
	// TODO
	fmt.Println("打印所有地址列表")
	wallets, _ := NewWallets()
	for address, _ := range wallets.WalletsMap {
		fmt.Printf("address: %s\n", address)
	}
	return nil
}
