package mpt

import (
	"bytes"
)

type Trie struct {
	node node
}

func NewTrie(node node) *Trie {
	return &Trie{
		node: node,
	}
}

// Put inserts key/value pair into tree
func (t *Trie) Put(key, value []byte) bool {
	node, ok := t.put(t.node, bytesToHex(key), NewLeafNode(value))
	if !ok {
		return false
	}

	t.node = node
	return true
}

// Get returns value for incoming key
func (t *Trie) Get(key []byte) ([]byte, bool) {
	val, node, resolved := t.get(t.node, bytesToHex(key))
	if resolved {
		t.node = node
	}

	return val, resolved
}

// Delete remove key/value par from chain
func (t *Trie) Delete(key []byte) bool {
	n, ok := t.delete(t.node, bytesToHex(key))
	if !ok {
		return false
	}

	t.node = n
	return true
}

func (t *Trie) put(n node, key []byte, value node) (node, bool) {
	if len(key) == 0 {
		if val, ok := n.(LeafNode); ok {
			return value, !bytes.Equal(val, value.(LeafNode))
		}
		return value, true
	}

	if n == nil {
		// initialize new ExtensionNode and set key/value pair
		return NewExtensionNode(key, value), true
	}

	return n.put(key, value)
}

func (t *Trie) get(n node, key []byte) ([]byte, node, bool) {
	if n == nil {
		return nil, nil, false
	}

	return n.find(key)
}

func (t *Trie) delete(n node, key []byte) (node, bool) {
	if n == nil {
		return nil, false
	}

	return n.delete(key)
}
