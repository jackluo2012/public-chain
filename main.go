package main

import (
	"fmt"
	"public-chain/blc"
)

func main() {
	block := blc.NewBlock(1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, "Genesis Block")
	fmt.Println(block)
}
