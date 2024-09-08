package test

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
	"public-chain/blc"
)

func test() {
	// 实例化顶层命令
	var opts blc.BlockCommand
	_, err := flags.Parse(&opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(opts.AddBlock.Data)
}
