package blc

// 子命令 打印区块链
type PrintChainCommand struct {
	Print bool `short:"p" long:"printchain" description:"Print the blockchain"`
}

func (b *PrintChainCommand) Execute(args []string) error {
	return nil
}

// 打印区块链
func (cli *CLI) PrintChain() {
	blockChain := GetBlockChainObject()
	defer blockChain.DB.Close()
	blockChain.PrintChain()
}
