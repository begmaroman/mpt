package mpt

import "bytes"

type node interface {
	find(key []byte) ([]byte, node, bool, error)
}

// ExtensionNode
type ExtensionNode struct {
	Key   []byte
	Value node
}

func NewExtensionNode(key []byte, val node) *ExtensionNode {
	return &ExtensionNode{
		Key:   key,
		Value: val,
	}
}

func (e *ExtensionNode) find(key []byte) ([]byte, node, bool, error) {
	if len(key) < len(e.Key) || !bytes.Equal(e.Key, key[:len(e.Key)]) {
		// record not found in tree
		return nil, e, false, nil
	}

	var nd *ExtensionNode
	val, newNode, resolved, err := e.Value.find(key[len(e.Key):])
	if err == nil && resolved {
		nd = e.copy()
		nd.Value = newNode
	}

	return val, nd, resolved, err
}

func (e *ExtensionNode) copy() *ExtensionNode {
	newNode := *e
	return &newNode
}

// BranchNode
type BranchNode struct {
	Children [17]node
}

func NewBranchNode() *BranchNode {
	return &BranchNode{}
}

func (b *BranchNode) find(key []byte) ([]byte, node, bool, error) {
	var nd *BranchNode
	val, newNode, resolved, err := b.Children[key[0]].find(key[1:])
	if err == nil && resolved {
		nd = b.copy()
		nd.Children[key[0]] = newNode
	}

	return val, nd, resolved, err
}

func (b *BranchNode) copy() *BranchNode {
	newNode := *b
	return &newNode
}

// LeafNode
type LeafNode []byte

func NewLeafNode(value []byte) LeafNode {
	return LeafNode(value)
}

func (b LeafNode) find(key []byte) ([]byte, node, bool, error) {
	return b, b, false, nil
}
