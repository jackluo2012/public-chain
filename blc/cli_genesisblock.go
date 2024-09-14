package blc

import (
	"errors"
	"fmt"
)

type CreateBlockChainWithGenesisBlockCommand struct {
	Address string `short:"a" long:"address" description:"address for the genesis block"`
}

func (b *CreateBlockChainWithGenesisBlockCommand) Execute(args []string) error {
	if IsValidAddress(b.Address) == false {
		fmt.Println("地址无效")
		return errors.New("Invalid address")
	}
	return nil
}

// 创建 创世区块
func (cli *CLI) CreateBlockChainWithGenesisBlock(address string) {

	blockChain := CreateBlockChainWithGenesisBlock(address)
	defer blockChain.DB.Close()
	utxoSet := &UTXOSet{blockChain}
	utxoSet.ResetUTXOSet()
}

//blocks
//utxoTable
