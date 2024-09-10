package blc

type CreateBlockChainWithGenesisBlockCommand struct {
	Address string `short:"a" long:"address" description:"address for the genesis block"`
}

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

func (cli *CLI) CreateBlockChainWithGenesisBlock(address string) {

	blockChain := CreateBlockChainWithGenesisBlock(address)
	defer blockChain.DB.Close()
}
