package node

type LeafNode []byte

func NewLeafNode(value []byte) LeafNode {
	return LeafNode(value)
}

func (l LeafNode) Find(key []byte) ([]byte, Node, bool) {
	return l, l, true
}

func (l LeafNode) Put(key []byte, value Node) (Node, bool) {
	return nil, false
}

func (l LeafNode) Delete(key []byte) (Node, bool) {
	return nil, true
}

func (l LeafNode) Cache() ([]byte, bool) {
	return nil, true
}

func (l LeafNode) CanUpload(gen, limit uint16) bool {
	return false
}
