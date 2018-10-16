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
	// value of current node, is nil if Node is not leaf
	Value []byte
	// checksum of current node
	checksum []byte
}

func NewNode() *Node {
	return &Node{}
}

// IsLeaf returns true is current Node is leaf
func (n *Node) IsLeaf() bool {
	return len(n.checksum) != 0 && (n.Left == nil && n.Right == nil)
}

// Checksum returns checksum of node.
func (n *Node) Checksum() []byte {
	return n.checksum
}

// LoadChecksum calculate and set checksum for current node
func (n *Node) LoadChecksum() error {
	checksum, err := n.calculateChecksum()
	if err != nil {
		return err
	}

	n.checksum = checksum
	return nil
}

// calculateChecksum returns checksum which calculate based on Nodes data
func (n *Node) calculateChecksum() ([]byte, error) {
	hash := sha1.New()
	if n.Left != nil && n.Right != nil {
		leftSumChan, rightSumChan := make(chan *checkSumResponse), make(chan *checkSumResponse)
		go func() {
			leftSumChan <- newCheckSumResponse(n.Left.calculateChecksum())
		}()

		go func() {
			rightSumChan <- newCheckSumResponse(n.Right.calculateChecksum())
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
