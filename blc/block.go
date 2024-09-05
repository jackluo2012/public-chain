package blc

import (
	"bytes"
	"crypto/sha256"
	"time"
)

type Block struct {
	// 区块的高度
	Height int64
	// 上一个区块的哈希
	PrevHash []byte
	// 交易数据
	Data []byte
	// 时间戳
	Timestamp int64
	// 当前区块的哈希
	Hash []byte
}

func NewBlock(height int64, prevHash []byte, data string) *Block {
	block := &Block{
		Height:    height,
		PrevHash:  prevHash,
		Data:      []byte(data),
		Timestamp: time.Now().Unix(),
		Hash:      nil,
	}
	// 调用生成哈希的方法
	block.SetHash()
	return block
}
func (block *Block) SetHash() {
	// 1. 将Height、PrevHash、Data、Timestamp拼接成字节数组
	info := bytes.Join([][]byte{
		IntToHex(block.Height),
		block.PrevHash,
		block.Data,
		IntToHex(block.Timestamp),
	}, []byte{})
	// 2. 对拼接好的字节数组进行哈希运算
	hash := sha256.Sum256(info)
	// 3. 将哈希值赋值给block的Hash属性
	block.Hash = hash[:]
}
