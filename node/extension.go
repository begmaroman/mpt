package node

import (
	"bytes"

	"github.com/begmaroman/mpt/enc"
)

type ExtensionNode struct {
	Key   []byte
	Value Node

	Hash  []byte
	Dirty bool
}

func NewExtensionNode(key []byte, val Node) *ExtensionNode {
	return &ExtensionNode{
		Key:   key,
		Value: val,
	}
}

func (e *ExtensionNode) Copy() *ExtensionNode {
	newNode := *e
	return &newNode
}

func (e *ExtensionNode) Find(key []byte) ([]byte, Node, bool) {
	if len(key) < len(e.Key) || !bytes.Equal(e.Key, key[:len(e.Key)]) {
		// record not found in tree
		return nil, e, false
	}

	var nd *ExtensionNode
	val, newNode, resolved := e.Value.Find(key[len(e.Key):])
	if resolved {
		nd = e.Copy()
		nd.Value = newNode
	}

	return val, nd, resolved
}

func (e *ExtensionNode) Put(key []byte, value Node) (Node, bool) {
	if e.Key == nil {
		return NewExtensionNode(key, value), true
	}

	matchKey := enc.PrefixLen(key, e.Key)

	// check if key of current node is compare with key
	if matchKey == len(e.Key) {
		nd, ok := e.Value.Put(key[matchKey:], value)
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

func (e *ExtensionNode) Delete(key []byte) (Node, bool) {
	matchKey := enc.PrefixLen(key, e.Key)

	// if key not fully compare with node's key
	if matchKey < len(e.Key) {
		return e, false
	}

	if matchKey == len(key) {
		return nil, true
	}

	childNode, ok := e.Value.Delete(key[len(e.Key):])
	if !ok {
		return e, false
	}

	switch childNode := childNode.(type) {
	case *ExtensionNode:
		r := make([]byte, len(e.Key)+len(childNode.Key))
		copy(r, e.Key)
		copy(r[len(e.Key):], childNode.Key)
		return NewExtensionNode(r, childNode.Value), true
	default:
		return NewExtensionNode(e.Key, childNode), true
	}
}

func (e *ExtensionNode) Cache() ([]byte, bool) {
	return e.Hash, e.Dirty
}
