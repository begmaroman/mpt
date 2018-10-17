package mptv2

const (
	leafNode = "leaf"
)

type TrieNode struct {
	Type string
	Raw  [uint8(2)]*TrieNode
}

func NewTrieNode(t string, key string, value []byte) *TrieNode {
	return &TrieNode{
		Type: t,
	}
}

func (t *TrieNode) hash() []byte {
	return []byte{}
}

func (t *TrieNode) serialize() []byte {
	return []byte{}
}
