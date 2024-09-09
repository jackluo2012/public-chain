package blc

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000
// 16进制 64位
// 2^256 / 12 = 2^64
// 16 个0
const targetBits = 12

type ProofOfWork struct {
	// 区块 当前要验证的区块
	Block *Block
	// 目标值 大数据存储，难度值
	target *big.Int
}

// //1.将block的字段拼接成字节数组
func (pow *ProofOfWork) prepareData(nonce int) []byte {

	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.HashTransactions(),
			IntToHex(pow.Block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
			IntToHex(int64(pow.Block.Height)),
		},
		[]byte{},
	)
	return data
}

// 验证挖矿结果
func (pow *ProofOfWork) IsValid() bool {
	var hashInt big.Int
	hashInt.SetBytes(pow.Block.Hash)
	isValid := hashInt.Cmp(pow.target) == -1
	return isValid
}

// 开始挖矿
func (pow *ProofOfWork) Run() (int64, []byte) {
	var hashInt big.Int //存储我们新生成的hash
	var hash [32]byte
	nonce := 0
	//1.将block的字段拼接成字节数组
	// 2、生成hash
	// 3、判断hash是否满足条件，满足则退出循环

	for {
		// 1.将block的字段拼接成字节数组
		data := pow.prepareData(nonce)
		// 2、生成hash
		hash = sha256.Sum256(data)
		// \r 换行，直接 覆盖
		fmt.Printf("\r%x", hash)
		//将hash 转化为 big.Int
		hashInt.SetBytes(hash[:])
		//  3、判断hash是否满足条件，满足则退出循环
		// 判断hashInt 是否小于Block里面的target
		// Cmp  compare x and y and return:
		// -1 if x < y
		// 0 if x == y
		// +1 if x > y
		if hashInt.Cmp(pow.target) == -1 {
			break
		}
		nonce++
	}
	return int64(nonce), hash[:]
}

// 创建新的工作量证明对象
func NewProofOfWork(block *Block) *ProofOfWork {

	// 1.big.Int 对象 1
	// 2.
	// 0000 0001
	// 8 - 2 = 6
	// 1 << 6 = 64
	// 0100 0000

	// 1.创建一个初始值 为1的 target
	target := big.NewInt(1)
	// 2.左移256 - targetBits 位 得到目标值
	target = target.Lsh(target, uint(256-targetBits))
	return &ProofOfWork{
		Block:  block,
		target: target,
	}
}
