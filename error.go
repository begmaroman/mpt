package mpt

import (
	"fmt"

	"github.com/begmaroman/mpt/enc"
)

type ErrNodeNotFound struct {
	Key  []byte
	Hash enc.Hash
}

func NewErrNodeNotFound(key []byte, hash enc.Hash) *ErrNodeNotFound {
	return &ErrNodeNotFound{
		Key:  key,
		Hash: hash,
	}
}

func (e *ErrNodeNotFound) Error() string {
	return fmt.Sprintf("node not found: hash %v key %v", e.Hash, e.Key)
}
