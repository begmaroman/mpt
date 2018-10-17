package mptv2

type Trie struct {
	db   Database
	root []byte
}

func NewTrie(db Database, root []byte) *Trie {
	return &Trie{
		db:   db,
		root: root,
	}
}

func (t *Trie) Get(key string) []byte {
	return nil
}

func (t *Trie) Put(key string, value []byte) error {
	if value == nil || len(value) == 0 {
		return t.Del(key)
	}

	if t.root == nil || len(t.root) == 0 {
		return t.createInitialNode(key, value)
	}

	return nil
}

func (t *Trie) Del(key string) error {
	return nil
}

// createInitialNode creates the initial node from an empty tree
func (t *Trie) createInitialNode(key string, value []byte) error {
	newNode := NewTrieNode(leafNode, key, value)
	t.root = newNode.hash()
	return t.putNode(newNode)
}

// putNode writes a single node to db
func (t *Trie) putNode(node *TrieNode) error {
	hash := node.hash()
	serialized := node.serialize()
	_, err := t.db.Put(hash, serialized)
	return err
}

func (t *Trie) findPath(key string) {

}

func (t *Trie) lookUpNode(node *TrieNode) {

}
