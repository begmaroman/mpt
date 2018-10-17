package mpt

import "bytes"

// behavior of node
type node interface {
	find([]byte) ([]byte, node, bool)
	put([]byte, node) (node, bool)
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

func (e *ExtensionNode) copy() *ExtensionNode {
	newNode := *e
	return &newNode
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

func (e *ExtensionNode) put(key []byte, value node) (node, bool) {
	matchKey := prefixLen(key, e.Key)

	// check if key of current node is compare with key
	if matchKey == len(e.Key) {
		nd, ok := e.Value.put(key[matchKey:], value)
		if !ok {
			return e, false
		}

		return NewExtensionNode(e.Key, nd), true
	}

	branchNode := NewBranchNode()
	branchNode.Children[e.Key[matchKey]] = NewExtensionNode(e.Key[matchKey+1:], e.Value)
	branchNode.Children[key[matchKey]] = NewExtensionNode(key[matchKey+1:], value)
	if matchKey == 0 {
		return branchNode, true
	}

	return NewExtensionNode(key[:matchKey], branchNode), true
}

// BranchNode
type BranchNode struct {
	Children [17]node
}

func NewBranchNode() *BranchNode {
	return &BranchNode{}
}

func (b *BranchNode) copy() *BranchNode {
	newNode := *b
	return &newNode
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

func (b *BranchNode) put(key []byte, value node) (node, bool) {
	nd, ok := b.Children[key[0]].put(key[1:], value)
	if !ok {
		return b, false
	}

	newNode := b.copy()
	newNode.Children[key[0]] = nd
	return newNode, true
}

// LeafNode
type LeafNode []byte

func NewLeafNode(value []byte) LeafNode {
	return LeafNode(value)
}

func (b LeafNode) find(key []byte) ([]byte, node, bool) {
	return b, b, true
}

func (b LeafNode) put(key []byte, value node) (node, bool) {
	return nil, false
}
