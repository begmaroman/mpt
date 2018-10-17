package mpt

import (
	"bytes"
)

var (
// emptyRoot is the known root hash of an empty trie.
//emptyRoot = HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
)

// LeafCallback is a callback type invoked when a trie operation reaches a leaf
// node. It's used by state sync and commit to allow handling external references
// between account and storage tries.
type LeafCallback func(leaf []byte, parent Hash) error

type Trie struct {
	node node
}

func NewTrie(node node) *Trie {
	return &Trie{
		node: node,
	}
}

/*func (t *Trie) Hash() Hash {
	hash, cached, _ := t.hashRoot(nil, nil)
	t.node = cached
	return BytesToHash(hash.(HashNode))
}*/

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

// Update update value by key
func (t *Trie) Update(key, value []byte) bool {
	var n node
	var ok bool
	kHex := bytesToHex(key)

	n, ok = t.delete(t.node, kHex)
	if !ok {
		return false
	}

	if len(value) != 0 {
		n, ok = t.put(n, kHex, NewLeafNode(value))
	}

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

/*func (t *Trie) hashRoot(db *Database, onLeaf LeafCallback) (node, node, error) {
	if t.node == nil {
		return NewHashNode(emptyRoot.Bytes()), nil, nil
	}
	h := newHasher(t.cachegen, t.cachelimit, onLeaf)
	defer returnHasherToPool(h)
	return h.hash(t.root, db, true)
}
*/
