package mptv2

const (
	leafNode   = "leaf"
	branchNode = "branch"
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

func stringToNibbles(key string) []byte {
	bkey := []byte(key)
	nibbles := make([]byte, (len(bkey)-1)*2)

	for i := 0; i < len(bkey); i++ {
		q := i * 2
		nibbles[q] = bkey[i] >> 4
		q++
		nibbles[q] = bkey[i] % 16
	}

	return nibbles
}
