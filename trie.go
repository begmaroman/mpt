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
func (t *Trie) Put(key, value []byte) bool {
	nd, ok := t.put(t.node, enc.BytesToHex(key), node.NewLeafNode(value), nil)
	if !ok {
		return false
	}

	t.node = nd
	return true
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
func (t *Trie) Delete(key []byte) bool {
	n, ok := t.delete(t.node, enc.BytesToHex(key), nil)
	if !ok {
		return false
	}

	t.node = n
	return true
}

// Update update value by key
func (t *Trie) Update(key, value []byte) bool {
	var n node.Node
	var ok bool
	kHex := enc.BytesToHex(key)

	n, ok = t.delete(t.node, kHex, nil)
	if !ok {
		return false
	}

	if len(value) != 0 {
		n, ok = t.put(n, kHex, node.NewLeafNode(value), nil)
		if !ok {
			return false
		}
	}

	t.node = n
	return true
}

func (t *Trie) put(n node.Node, key []byte, value node.Node, prefix []byte) (node.Node, bool) {
	if len(key) == 0 {
		if val, ok := n.(node.LeafNode); ok {
			return value, !bytes.Equal(val, value.(node.LeafNode))
		}
		return value, true
	}

	if n == nil {
		// initialize new ExtensionNode and set key/value pair
		return node.NewExtensionNode(key, value, nil), true
	}

	return n.Put(key, value)
}

func (t *Trie) get(n node.Node, key []byte) ([]byte, node.Node, bool) {
	if n == nil {
		return nil, nil, false
	}

	return n.Find(key)
}

func (t *Trie) delete(n node.Node, key, prefix []byte) (node.Node, bool) {
	if n == nil {
		return nil, false
	}

	return n.Delete(key)
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
