package blc

type MerkleTree struct {
	RootNode *MerkleNode
}
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

/**
* data 是所有transaction 序列化后的二维数组
 */
func NewMerkleTree(data [][]byte) *MerkleTree {
	//[tx1,tx2,tx3]
	//存放节点的数组
	var nodes []MerkleNode
	// 如果 是 奇数个节点，则复制最后一个节点
	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
		//[tx1,tx2,tx3,tx3]
	}
	// 创建叶子节点
	for _, d := range data {
		node := NewMerkleNode(nil, nil, d)
		nodes = append(nodes, *node)
	}
	// MerkelNode{nil,nil, tx1Bytes} node2
	// MerkelNode{nil,nil, tx2Bytes} node3
	// MerkelNode{nil,nil, tx3Bytes} node4
	// MerkelNode{nil,nil, tx3Bytes} node5

	// 循环两次
	for i := 0; i < len(data)/2; i++ {
		var newLevel []MerkleNode
		// 每次循环，将两个节点合并成一个节点
		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			//2节点合并成一个节点
			newLevel = append(newLevel, *node)
			// MerkleNode{MerkleNode{nil,nil,tx1BYytes},MerkleNode{nil,nil,tx2BYytes} }
			// MerkleNode{MerkleNode{nil,nil,tx3BYytes},MerkleNode{nil,nil,tx3BYytes} }
		}
		if len(newLevel)%2 != 0 {
			newLevel = append(newLevel, newLevel[len(newLevel)-1])
		}
		nodes = newLevel
	}
	// MerkleNode:
	// Left: MerkleNode:
	return &MerkleTree{&nodes[0]}
}

// 创建merkle一个节点
func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	mNode := &MerkleNode{}
	// 创建叶子节点
	if left == nil && right == nil {
		mNode.Data = Sha256Hash(data)
	} else {
		// 创建非叶子节点
		//把左右节点的hash值拼接在一起
		prevHashes := append(left.Data, right.Data...)
		mNode.Data = Sha256Hash(prevHashes)
	}
	mNode.Left = left
	mNode.Right = right

	return mNode
}

//一个区块中有多个transaction
// Block [tx1 tx2 tx3 tx4] node1
// MerkleNode{nil,nil, tx1} node2
// MerkleNode{nil,nil,tx2} node3
// MerkleNode{nil,nil,tx3} node4
// MerkleNode{nil,nil,tx4} node5
