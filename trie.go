package mpt

import (
	"bytes"
	"errors"
)

var (
	ErrUndefinedType = errors.New("undefined node type")
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

// Delete remove transaction based on key
func (t *Trie) Delete(key []byte) error {
	return nil
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

func (t *Trie) delete(n node, key []byte) (node, bool, error) {
	switch n := n.(type) {
	case *ExtensionNode:
		matchKey := prefixLen(key, n.Key)

		// if key not fully compare with node's key
		if matchKey < len(n.Key) {
			return n, false, nil
		}

		if matchKey == len(key) {
			return nil, true, nil
		}

		childNode, ok, err := t.delete(n.Value, key[len(n.Key):])
		if !ok || err != nil {
			return n, false, err
		}

		switch childNode := childNode.(type) {
		case *ExtensionNode:
			return NewExtensionNode(concat(n.Key, childNode.Key...), childNode.Value), true, nil
		default:
			return NewExtensionNode(n.Key, childNode), true, nil
		}
	case *BranchNode:
		nd, ok, err := t.delete(n.Children[key[0]], key[1:])
		if !ok || err != nil {
			return n, false, err
		}

		n = n.copy()
		n.Children[key[0]] = nd

		position := -1
		for i, child := range &n.Children {
			if child == nil {
				continue
			}

			if position == -1 {
				position = i
				continue
			}

			position = -2
			break
		}

		// TODO: realize logic
		/*if position >= 0 {
			if position != 16 {
				cNode, err
			}
		}*/
	case LeafNode:
		return nil, true, nil
	case nil:
		return nil, false, nil
	}

	return n, false, ErrUndefinedType
}
