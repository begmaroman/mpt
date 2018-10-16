package mpt

import "fmt"

type ErrChecksum struct {
	node *Node
}

func NewErrChecksum(node *Node) *ErrChecksum {
	return &ErrChecksum{
		node: node,
	}
}

func (err *ErrChecksum) Error() string {
	return fmt.Sprintf("checksum detecting error: %#v", *err.node)
}
