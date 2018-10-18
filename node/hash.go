package node

import (
	"github.com/begmaroman/mpt/enc"
)

type HashNode []byte

func NewHashNode(value []byte) HashNode {
	return HashNode(value)
}

func (h HashNode) Hash() enc.Hash {
	return enc.BytesToHash(h)
}

func (h HashNode) Find(key []byte) ([]byte, Node, bool) {
	// TODO: need to implement lookup logic
	return h, h, true
}

func (h HashNode) Put(key []byte, value Node) (Node, bool) {
	// TODO: need to implement insertion logic
	return nil, false
}

func (h HashNode) Delete(key []byte) (Node, bool) {
	return nil, true
}

func (h HashNode) Cache() ([]byte, bool) {
	return nil, true
}
