package mpt

import (
	"crypto/sha1"
)

type checkSumResponse struct {
	checksum []byte
	err      error
}

func newCheckSumResponse(checksum []byte, err error) *checkSumResponse {
	return &checkSumResponse{
		checksum: checksum,
		err:      err,
	}
}

type Node struct {
	Left, Right, Parent *Node
	Value               []byte

	checksum []byte
}

func NewNode(value []byte) *Node {
	return &Node{
		Value: value,
	}
}

func NewEmptyNode() *Node {
	return &Node{}
}

func InitNode(value []byte) (*Node, error) {
	node := NewNode(value)
	sum, err := node.Checksum()
	if err != nil {
		return nil, err
	}
	node.checksum = sum
	return node, nil
}

func (n *Node) IsLeaf() bool {
	return len(n.checksum) != 0 && (n.Left == nil && n.Right == nil)
}

func (n *Node) Checksum() ([]byte, error) {
	if n.checksum != nil {
		return n.checksum, nil
	}

	hash := sha1.New()
	if n.Left != nil && n.Right != nil {
		leftSumChan, rightSumChan := make(chan *checkSumResponse), make(chan *checkSumResponse)
		go func() {
			leftSumChan <- newCheckSumResponse(n.Left.Checksum())
		}()

		go func() {
			rightSumChan <- newCheckSumResponse(n.Right.Checksum())
		}()

		leftSum := <-leftSumChan
		if leftSum.err != nil {
			return nil, leftSum.err
		}
		if _, err := hash.Write(leftSum.checksum); err != nil {
			return nil, err
		}

		rightSum := <-rightSumChan
		if rightSum.err != nil {
			return nil, rightSum.err
		}
		if _, err := hash.Write(rightSum.checksum); err != nil {
			return nil, err
		}
	} else {
		if _, err := hash.Write(n.Value); err != nil {
			return nil, err
		}
	}

	return hash.Sum(nil), nil
}

// leaf collect all leaf of a trie
func (n *Node) leaf(nodes []*Node) []*Node {
	if n == nil {
		return nil
	}

	if nodes == nil {
		nodes = make([]*Node, 0, 2)
	}

	if n.IsLeaf() {
		nodes = append(nodes, n)
		return nodes
	}

	nodes = n.Left.leaf(nodes)
	nodes = n.Right.leaf(nodes)

	return nodes
}
