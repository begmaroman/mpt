// TODO: need to implement database functional
package mpt

import (
	"github.com/begmaroman/mpt/enc"
	"github.com/begmaroman/mpt/node"
	"sync"
)

type Database struct {
	sync.Mutex
}

func NewDatabase() *Database {
	return &Database{}
}

func (db *Database) GetNode(hash enc.Hash) node.Node {
	return nil
}

func (db *Database) Insert(hash enc.Hash, blob []byte, node node.Node) {
	db.Lock()
	defer db.Unlock()
	// TODO: need to implement storing logic
}
