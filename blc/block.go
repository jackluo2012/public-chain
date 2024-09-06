package blc

import (
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

	// Nonce 添加工作量证明的难度
	Nonce int64
}

func NewBlock(height int64, prevHash []byte, data string) *Block {
	block := &Block{
		Height:    height,
		PrevHash:  prevHash,
		Data:      []byte(data),
		Timestamp: time.Now().Unix(),
		Hash:      nil,
		Nonce:     0,
	}
	// 调用工作量证明，并且 返回有效的Hash和Nonce值
	pow := NewProofOfWork(block)
	// 调用工作量证明，并且 返回有效的Hash和Nonce值
	block.Nonce, block.Hash = pow.Run()
	//
	return block
}

// 生成创世区块
func NewGenesisBlock(data string) *Block {
	return NewBlock(1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, data)
}
