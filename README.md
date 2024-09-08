# public-chain


## 项目介绍

本项目是区块链项目，通过区块链技术实现数字货币的发行和流通，实现去中心化的交易。

## 项目结构

```
public-chain
├─blc

└─cli

```

## 技术栈

- Go语言
- Golang区块链开发实战
- 区块链
- 数字货币
- 区块链技术
## cli 模块

cli模块是区块链项目的命令行接口，用于与区块链进行交互。它提供了以下功能：

- 添加区块
- 查询区块链
- 查询余额
- 发送交易
- 查询交易
- 查询未确认交易
- 查询未确认交易数量
- 查询区块奖励
- 查询区块奖励数量
``` bash
$ go run main.go addblock -data "hello world"
$ go run main.go getblockchain
$ go run main.go getbalance -address "0x1234567890abcdef"
$ go run main.go send -from "0x1234567890abcdef" -to "0xabcdef1234567890" -amount 10
```