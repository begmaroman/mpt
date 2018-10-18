package mpt

import (
	"bytes"
	"github.com/begmaroman/mpt/enc"
	"github.com/begmaroman/mpt/node"
)

var (
	// emptyRoot is the known root hash of an empty trie.
	emptyRoot = enc.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
)

type Trie struct {
	db   *Database
	node node.Node

	cGen, cLimit uint16
}

func NewTrie(node node.Node) *Trie {
	return &Trie{
		node: node,
	}
}

// Hash calculate and return hash of the tree
func (t *Trie) Hash() (enc.Hash, error) {
	bHash, cached, err := t.hashRoot(nil)
	hash := enc.BytesToHash(bHash.(node.HashNode))

	if err == nil {
		t.node = cached
	}

	return hash, err
}

// Put inserts key/value pair into tree
func (t *Trie) Put(key, value []byte) (bool, error) {
	nd, ok, err := t.put(t.node, enc.BytesToHex(key), node.NewLeafNode(value), nil)
	if !ok || err != nil {
		return false, err
	}

	t.node = nd
	return true, nil
}

// Get returns value for incoming key
func (t *Trie) Get(key []byte) ([]byte, bool) {
	val, nd, resolved := t.get(t.node, enc.BytesToHex(key))
	if resolved {
		t.node = nd
	}

	return val, resolved
}

// Delete remove key/value par from chain
func (t *Trie) Delete(key []byte) (bool, error) {
	n, ok, err := t.delete(t.node, enc.BytesToHex(key), nil)
	if !ok || err != nil {
		return false, err
	}

	t.node = n
	return true, nil
}

// Update update value by key
func (t *Trie) Update(key, value []byte) (bool, error) {
	var n node.Node
	var ok bool
	var err error
	kHex := enc.BytesToHex(key)

	n, ok, err = t.delete(t.node, kHex, nil)
	if !ok || err != nil {
		return false, err
	}

	if len(value) != 0 {
		n, ok, err = t.put(n, kHex, node.NewLeafNode(value), nil)
		if !ok || err != nil {
			return false, err
		}
	}

	t.node = n
	return true, nil
}

func (t *Trie) put(n node.Node, key []byte, value node.Node, prefix []byte) (node.Node, bool, error) {
	if len(key) == 0 {
		if val, ok := n.(node.LeafNode); ok {
			return value, !bytes.Equal(val, value.(node.LeafNode)), nil
		}
		return value, true, nil
	}

	if n == nil {
		// initialize new ExtensionNode and set key/value pair
		return node.NewExtensionNode(key, value, nil), true, nil
	}

	if n, ok := n.(node.HashNode); ok {
		rn, err := t.resolveHash(n, prefix)
		if err != nil {
			return nil, false, err
		}

		nn, ok, err := t.put(rn, key, value, prefix)
		if !ok || err != nil {
			return rn, false, err
		}

		return nn, true, nil
	}

	nd, ok := n.Put(key, value)
	return nd, ok, nil
}

func (t *Trie) get(n node.Node, key []byte) ([]byte, node.Node, bool) {
	if n == nil {
		return nil, nil, false
	}

	return n.Find(key)
}

func (t *Trie) delete(n node.Node, key, prefix []byte) (node.Node, bool, error) {
	if n == nil {
		return nil, false, nil
	}

	if n, ok := n.(node.HashNode); ok {
		rn, err := t.resolveHash(n, prefix)
		if err != nil {
			return nil, false, err
		}

		nn, ok, err := t.delete(rn, key, prefix)
		if !ok || err != nil {
			return rn, false, err
		}

		return nn, true, nil
	}

	nn, ok := n.Delete(key)
	return nn, ok, nil
}

// hashToot do hash for root of tree
func (t *Trie) hashRoot(db *Database) (node.Node, node.Node, error) {
	if t.node == nil {
		return node.NewHashNode(emptyRoot.Bytes()), nil, nil
	}
	h := NewEncryptor()
	defer hPool.Put(h)
	return h.hash(t.node, db, true)
}

// resolve returns results of resolveHash if incoming node.Node is HashNode, or return incoming node.Node
func (t *Trie) resolve(n node.Node, prefix []byte) (node.Node, error) {
	if n, ok := n.(node.HashNode); ok {
		return t.resolveHash(n, prefix)
	}
	return n, nil
}

// resolveHash return node.Node from DB or error if not found
func (t *Trie) resolveHash(n node.HashNode, prefix []byte) (node.Node, error) {
	hash := n.Hash() // get hash of HashNode

	if nd := t.db.GetNode(hash, t.cGen); nd != nil { // lookup node.Node in database
		return nd, nil // return node.Node from database if exist
	}

	return nil, NewErrNodeNotFound(prefix, hash) // return error if not found
}
