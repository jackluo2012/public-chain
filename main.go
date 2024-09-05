package main

import (
	"fmt"
	"public-chain/blc"
)

func main() {

	genesisBlock := blc.NewGenesisBlock("Genesis Block")

	fmt.Println(genesisBlock)
}
