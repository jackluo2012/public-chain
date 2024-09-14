package blc

import "fmt"

type AtestCommand struct {
	Print bool `short:"p" long:"print" description:"Print test message"` // 添加一个布尔类型的 Print 字段
}

func (t *AtestCommand) Execute(args []string) error {
	return nil
}

// Test
func (c *CLI) Test() {
	blockchain := GetBlockChainObject()
	defer blockchain.DB.Close()
	utxoMap := blockchain.FindUTXOMap()
	fmt.Println(utxoMap)

}
