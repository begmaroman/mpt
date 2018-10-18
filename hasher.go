package mpt

import (
	"hash"
	"sync"

	"github.com/begmaroman/mpt/enc"
	"github.com/begmaroman/mpt/node"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	hPool = sync.Pool{
		New: func() interface{} {
			return &Encryptor{
				rlpBuf: make(enc.SliceBuffer, 0, 550),
				sha:    sha3.NewKeccak256().(keccakState),
			}
		},
	}
)

type keccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}

type Encryptor struct {
	rlpBuf enc.SliceBuffer
	sha    keccakState
}

func NewEncryptor() *Encryptor {
	h := hPool.Get().(*Encryptor)
	return h
}

// hash recursive hash the trie based on keccak256 and sha3
func (sh *Encryptor) hash(n node.Node, db *Database, force bool) (node.Node, node.Node, error) {
	// If we're not storing the node, just hashing, use available cached data
	if cHach, dirty := n.Cache(); cHach != nil && !dirty {
		return node.NewHashNode(cHach), n, nil
	}

	// Trie not processed yet or needs storage, walk the children
	collapsed, cached, err := sh.getChildHash(n, db)
	if err != nil {
		return node.NewHashNode(nil), n, err
	}

	encryptedNode, err := sh.encrypt(collapsed, db, force)
	if err != nil {
		return node.NewHashNode(nil), n, err
	}

	if hashData, ok := encryptedNode.(node.HashNode); ok && db != nil {
		// We are pooling the trie nodes into an intermediate memory cache
		db.Insert(enc.BytesToHash(hashData), sh.rlpBuf, n)
	}

	cachedHash, _ := encryptedNode.(node.HashNode)
	switch cn := cached.(type) {
	case *node.ExtensionNode:
		cn.Hash = cachedHash
		if db != nil {
			cn.Dirty = false
		}
	case *node.BranchNode:
		cn.Hash = cachedHash
		if db != nil {
			cn.Dirty = false
		}
	}

	return encryptedNode, cached, nil
}

// getChildHash hash children nodes
func (sh *Encryptor) getChildHash(original node.Node, db *Database) (node.Node, node.Node, error) {
	var err error

	switch n := original.(type) {
	case *node.ExtensionNode:
		// Hash the short node's child, caching the newly hashed subtree
		collapsed, cached := n.Copy(), n.Copy()
		collapsed.Key = enc.HexToCompact(n.Key)
		cached.Key = enc.CopyBytes(n.Key)

		switch n.Value.(type) {
		case node.LeafNode:
		default:
			if collapsed.Value, cached.Value, err = sh.hash(n.Value, db, false); err != nil {
				return original, original, err
			}
		}

		return collapsed, cached, nil

	case *node.BranchNode:
		// Hash the full node's children, caching the newly hashed subtrees
		collapsed, cached := n.Copy(), n.Copy()

		for i := 0; i < 16; i++ {
			if n.Children[i] == nil {
				continue
			}

			if collapsed.Children[i], cached.Children[i], err = sh.hash(n.Children[i], db, false); err != nil {
				return original, original, err
			}
		}

		cached.Children[16] = n.Children[16]

		return collapsed, cached, nil

	default:
		// Value and hash nodes don't have children so they're left as were
		return n, original, nil
	}
}

// encrypt returns HashNode or incomin node
func (sh *Encryptor) encrypt(n node.Node, db *Database, force bool) (node.Node, error) {
	// Don't store hashes or empty nodes.
	if _, isHash := n.(node.HashNode); n == nil || isHash {
		return n, nil
	}

	// Generate the RLP encoding of the node
	sh.rlpBuf.Reset()
	if err := rlp.Encode(&sh.rlpBuf, n); err != nil {
		return n, err
	}

	// Nodes smaller than 32 bytes are stored inside their parent
	if len(sh.rlpBuf) < 32 && !force {
		return n, nil
	}

	// Larger nodes are replaced by their hash and stored in the database.
	newHash, dirty := n.Cache()
	if newHash == nil || dirty {
		newHash = sh.makeHashNode(sh.rlpBuf)
	}

	return node.NewHashNode(newHash), nil
}

// makeHashNode wrap data on sha3
func (sh *Encryptor) makeHashNode(data []byte) node.HashNode {
	n := make(node.HashNode, sh.sha.Size())
	sh.sha.Reset()
	sh.sha.Write(data)
	sh.sha.Read(n)
	return n
}
