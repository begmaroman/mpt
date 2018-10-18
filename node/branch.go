package node

type BranchNode struct {
	Children [17]Node

	Hash  []byte
	Dirty bool
}

func NewBranchNode() *BranchNode {
	return &BranchNode{}
}

func (b *BranchNode) Copy() *BranchNode {
	newNode := *b
	return &newNode
}

func (b *BranchNode) Find(key []byte) ([]byte, Node, bool) {
	var nd *BranchNode
	val, newNode, resolved := b.Children[key[0]].Find(key[1:])
	if resolved {
		nd = b.Copy()
		nd.Children[key[0]] = newNode
	}

	return val, nd, resolved
}

func (b *BranchNode) Put(key []byte, value Node) (Node, bool) {
	var n Node
	var ok bool

	if b.Children[key[0]] == nil {
		n = NewExtensionNode(key[1:], value)
	} else {
		n, ok = b.Children[key[0]].Put(key[1:], value)
		if !ok {
			return b, false
		}
	}

	newNode := b.Copy()
	newNode.Children[key[0]] = n
	return newNode, true
}

func (b *BranchNode) Delete(key []byte) (Node, bool) {
	nd, ok := b.Children[key[0]].Delete(key[1:])
	if !ok {
		return b, false
	}

	newNode := b.Copy()
	newNode.Children[key[0]] = nd

	position := -1
	for i, child := range &newNode.Children {
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

	if position >= 0 {
		if position != 16 {
			if cNode, ok := newNode.Children[position].(*ExtensionNode); ok {
				k := append([]byte{byte(position)}, cNode.Key...)
				return NewExtensionNode(k, cNode.Value), true
			}
		}

		return NewExtensionNode([]byte{byte(position)}, newNode.Children[position]), true
	}

	return newNode, true
}

func (b *BranchNode) Cache() ([]byte, bool) {
	return b.Hash, b.Dirty
}
