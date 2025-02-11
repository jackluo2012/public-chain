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
# windows 下需要切换到 bash 下面执行   -----start-----
$ ./main send -f '["jackluo","tom","lily"]' -t '["mm","li","one"]' -a '["2","3","4"]'
$ ./main send -f '["jackluo"]' -t '["luo"]' -a '["2"]'
$ ./main send -f '["1G1hfVcYFpiDRoMruJeM82SDunJk8yEbmY"]' -t '["19pQLVBinRNE1ZTFmgp3Lj8DXxzK3nQD8W"]' -a '["2"]'
$ ./main send -f '["144nDLvWQU4pcX3vteN46bWARRQFFgDsQU"]' -t '["1Br9cUuXZ9xfevMUMBmkUiP4oWD6f3XdMV"]' -a '["2"]'
# windows 下需要切换到 bash 下面执行   -----end-----
# 查询交易
$ go run main.go gettransaction -txid "0x1234567890abcdef"

# 查询未确认交易
$ go run main.go getunconfirmedtransactions

# 查询未确认交易数量
$ go run main.go getunconfirmedtransactionscount

# 查询区块奖励
```


### 创建一个钱包地址：
- 1、生成一对公钥和私钥
- 2、想要使用这个地址，需要将公钥进行Base58 哈希，得到一个地址
- 3、想要别人给我转账，把地址给别人，别人将反编码变成公钥，将公钥和数据进行签名
- 4、通过私钥进行解密，解密是单方向的，只有用私钥才能解密

### 1、创建钱包
    1）、私 钥
    2）、公钥
    2、先将公钥进行一次256Hash,再进行一次160hash
    // 20 字节的[]byte
    // version {0} + hash160 -> pubkey

    // 256hash pubkey 几次
    // 256 64
    // 最后的四个字节 取出来
    // 4个字节 + 20个字节 = 24个字节
    // version {0} + hash160 +4个字节 -> 25个字节
    // base58 编码

### 未花费输出
- 存储到别一张表中，存储到数据库中


###  主要内容包括
- pow  工作量证明
- 数据库存储
- 命令行工具
- utxo模型
- 公钥 私钥
- 交易
- 签名
- 钱包
- 多笔转账