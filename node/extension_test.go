package node

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/begmaroman/mpt/enc"
)

func TestExtensionNode_Find(t *testing.T) {
	key, value := enc.BytesToHex([]byte("key")), []byte("value")

	nd := NewExtensionNode(key, NewLeafNode(value))

	rVal, rNd, ok := nd.Find(key)
	if !ok {
		t.Error("added data not found")
		return
	}

	if !bytes.Equal(rVal, value) {
		t.Errorf("invalid value: expected %s got %s", value, rVal)
	}

	if !reflect.DeepEqual(nd, rNd) {
		t.Errorf("invalid node: expected %#v got %#v", nd, rNd)
	}
}

func TestExtensionNode_PutAndCheck(t *testing.T) {
	key, value := enc.BytesToHex([]byte("key")), []byte("value")

	nd := NewExtensionNode(enc.BytesToHex([]byte("k")), NewLeafNode([]byte("kat value")))
	newNode, ok := nd.Put(key, NewLeafNode(value))
	if !ok {
		t.Error("put failed")
		return
	}

	rVal, rNd, ok := newNode.Find(key)
	if !ok {
		t.Error("added data not found")
		return
	}

	if !bytes.Equal(rVal, value) {
		t.Errorf("invalid value: expected %s got %s", value, rVal)
	}

	if !reflect.DeepEqual(newNode, rNd) {
		t.Errorf("invalid node: expected %#v got %#v", newNode, rNd)
	}
}

func TestExtensionNode_Cache(t *testing.T) {
	nd := NewExtensionNode(nil, nil)
	nd.Hash = []byte("test")
	nd.Dirty = true

	hash, dirty := nd.Cache()
	if !bytes.Equal(hash, nd.Hash) {
		t.Errorf("invalid hash: expected %s got %s", nd.Hash, hash)
	}

	if nd.Dirty != dirty {
		t.Errorf("invalid dirty: expected %t got %t", nd.Dirty, dirty)
	}
}
