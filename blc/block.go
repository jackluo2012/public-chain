package blc

import (
	"bytes"
	"encoding/gob"
	"time"
)

type Block struct {
	// 区块的高度
	Height int64
	// 上一个区块的哈希
	PrevHash []byte
	// 交易数据
	Txs []*Transaction
	// 时间戳
	Timestamp int64
	// 当前区块的哈希
	Hash []byte `storm:"id"` // primary key

	// Nonce 添加工作量证明的难度
	Nonce int64
}

// 区块 序列化
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	// 创建编码器 打包
	encoder := gob.NewEncoder(&result)
	// 编码
	err := encoder.Encode(b)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}

// 区块 反序列化
func (b *Block) Deserialize(data []byte) *Block {
	var block Block
	// 创建解码器 解包
	decoder := gob.NewDecoder(bytes.NewReader(data))
	// 解码
	err := decoder.Decode(&block)
	if err != nil {
		panic(err)
	}
	return &block //
}

// 将Txs转换成[]byte
func (b *Block) HashTransactions() []byte {
	// var txHashes [][]byte
	// var txHash [32]byte

	// for _, tx := range b.Txs {
	// 	txHashes = append(txHashes, tx.TxHash)
	// }
	// txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	// return txHash[:]
	var transaction [][]byte
	for _, tx := range b.Txs {
		transaction = append(transaction, tx.Serialize())
	}
	mTree := NewMerkleTree(transaction)
	return mTree.RootNode.Data
}

func NewBlock(height int64, prevHash []byte, txs []*Transaction) *Block {
	block := &Block{
		Height:    height,
		PrevHash:  prevHash,
		Txs:       txs,
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
func NewGenesisBlock(txs []*Transaction) *Block {
	return NewBlock(1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, txs)
}
