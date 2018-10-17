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
func (t *Trie) Put(key, value []byte) error {
	_, node, err := t.put(t.node, bytesToHex(key), NewLeafNode(value))
	if err != nil {
		return err
	}

	t.node = node
	return nil
}

// Get returns value for incoming key
func (t *Trie) Get(key []byte) ([]byte, error) {
	val, node, resolved, err := t.get(t.node, bytesToHex(key))
	if err == nil && resolved {
		t.node = node
	}

	return val, err
}

// Delete remove transaction based on key
func (t *Trie) Delete(key []byte) error {
	return nil
}

func (t *Trie) put(n node, key []byte, value node) (bool, node, error) {
	if len(key) == 0 {
		if val, ok := n.(LeafNode); ok {
			return !bytes.Equal(val, value.(LeafNode)), value, nil
		}
		return true, value, nil
	}

	switch n := n.(type) {
	case *ExtensionNode:
		matchKey := prefixLen(key, n.Key)

		// check if key of current node is compare with key
		if matchKey == len(n.Key) {
			ok, nd, err := t.put(n.Value, key[matchKey:], value)
			if !ok || err != nil {
				return false, n, err
			}
			return true, NewExtensionNode(n.Key, nd), nil
		}

		var err error
		branchNode := NewBranchNode()
		_, branchNode.Children[n.Key[matchKey]], err = t.put(
			nil,
			n.Key[matchKey+1:],
			n.Value,
		)
		if err != nil {
			return false, nil, err
		}

		_, branchNode.Children[key[matchKey]], err = t.put(
			nil,
			key[matchKey+1:],
			value,
		)
		if err != nil {
			return false, nil, err
		}

		if matchKey == 0 {
			return true, branchNode, nil
		}

		return true, NewExtensionNode(key[:matchKey], branchNode), nil
	case *BranchNode:
		ok, nd, err := t.put(n.Children[key[0]], key[1:], value)
		if !ok || err != nil {
			return false, n, err
		}

		newNode := n.copy()
		newNode.Children[key[0]] = nd
		return true, newNode, nil
	case nil:
		// initialize new ExtensionNode and set key/value pair
		return true, NewExtensionNode(key, value), nil
	}

	return false, nil, ErrUndefinedType
}

func (t *Trie) get(n node, key []byte) ([]byte, node, bool, error) {
	if n == nil {
		return nil, nil, false, nil
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
