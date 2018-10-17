package mpt

type node interface {
}

type ExtensionNode struct {
	Key   []byte
	Value node
}

func NewExtensionNode(key []byte, val node) *ExtensionNode {
	return &ExtensionNode{
		Key:   key,
		Value: val,
	}
}

func (e *ExtensionNode) copy() *ExtensionNode {
	newNode := *e
	return &newNode
}

type BranchNode struct {
	Children [17]node
}

func NewBranchNode() *BranchNode {
	return &BranchNode{}
}

func (b *BranchNode) copy() *BranchNode {
	newNode := *b
	return &newNode
}

type LeafNode []byte

func NewLeafNode(value []byte) LeafNode {
	return LeafNode(value)
}
