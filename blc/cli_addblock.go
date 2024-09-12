package blc

type AddBlockCommand struct {
	Data string `short:"d" long:"data" description:"Data for the block"`
}

func (b *AddBlockCommand) Execute(args []string) error {
	return nil
}

// 添加区块
func (cli *CLI) AddBlock(data string) {
	blockChain := GetBlockChainObject()
	defer blockChain.DB.Close()
	blockChain.AddBlockToBlockChain([]*Transaction{})
}
