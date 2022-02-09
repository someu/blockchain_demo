package merkletree

import (
	"crypto/sha256"
)

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

type MerkleTree struct {
	RootNode *MerkleNode
}

func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []MerkleNode

	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}
	for _, datum := range data {
		nodes = append(nodes, *NewMerkleNode(nil, nil, datum))
	}
	for i := 0; i < len(data)/2; i++ {
		var newNodes []MerkleNode
		for j := 0; j < len(nodes); j += 2 {
			newNodes = append(newNodes, *NewMerkleNode(&nodes[j], &nodes[j+1], nil))
		}
		nodes = newNodes
	}

	tree := MerkleTree{
		RootNode: &nodes[0],
	}
	return &tree
}

func NewMerkleNode(left *MerkleNode, right *MerkleNode, data []byte) *MerkleNode {
	node := MerkleNode{
		Left:  left,
		Right: right,
	}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.Data = hash[:]
	} else {
		hash := sha256.Sum256(append(left.Data, right.Data...))
		node.Data = hash[:]
	}

	return &node
}
