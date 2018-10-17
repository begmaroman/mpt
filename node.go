package mpt

import "bytes"

type node interface {
	find(key []byte) ([]byte, node, bool)
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

func (e *ExtensionNode) find(key []byte) ([]byte, node, bool) {
	if len(key) < len(e.Key) || !bytes.Equal(e.Key, key[:len(e.Key)]) {
		// record not found in tree
		return nil, e, false
	}

	var nd *ExtensionNode
	val, newNode, resolved := e.Value.find(key[len(e.Key):])
	if resolved {
		nd = e.copy()
		nd.Value = newNode
	}

	return val, nd, resolved
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

func (b *BranchNode) find(key []byte) ([]byte, node, bool) {
	var nd *BranchNode
	val, newNode, resolved := b.Children[key[0]].find(key[1:])
	if resolved {
		nd = b.copy()
		nd.Children[key[0]] = newNode
	}

	return val, nd, resolved
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

func (b LeafNode) find(key []byte) ([]byte, node, bool) {
	return b, b, true
}
