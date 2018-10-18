package mpt

import (
	"bytes"
	"github.com/begmaroman/mpt/node"
	"testing"
)

func TestTrie_Put(t *testing.T) {
	testData := []*struct {
		key, value []byte
	}{
		{
			key:   []byte("test"),
			value: []byte("test value with words"),
		},
		{
			key:   []byte("test1"),
			value: []byte("test value with words 1"),
		},
		{
			key:   []byte("test2"),
			value: []byte("test value with words 2"),
		},
	}

	trie := NewTrie(nil)

	for _, data := range testData {
		if !trie.Put(data.key, data.value) {
			t.Errorf("key/value not added to tree: %s: %s", data.key, data.value)
		}
	}
}

func TestTrie_PutAndCheck(t *testing.T) {
	testData := []*struct {
		key, value, upValue []byte
	}{
		{
			key:     []byte("test"),
			value:   []byte("test value with words"),
			upValue: []byte("test updated value"),
		},
		{
			key:     []byte("test1"),
			value:   []byte("test value with words 1"),
			upValue: []byte("test updated value 1"),
		},
		{
			key:     []byte("test2"),
			value:   []byte("test value with words 2"),
			upValue: []byte("test updated value 2"),
		},
	}

	trie := NewTrie(nil)

	for _, data := range testData {
		if !trie.Put(data.key, data.value) {
			t.Errorf("key/value not added to tree: %s: %s", data.key, data.value)
		}
	}

	for _, data := range testData {
		rVal, ok := trie.Get(data.key)
		if !ok {
			t.Errorf("lookup failed for key %s", data.key)
		}

		if !bytes.Equal(rVal, []byte(data.value)) {
			t.Errorf("invalid value: expected %s got %s", data.value, rVal)
		}
	}
}

func TestTrie_CheckHash(t *testing.T) {
	testData := []*struct {
		key, value, upValue []byte
	}{
		{
			key:     []byte("test"),
			value:   []byte("test value with words"),
			upValue: []byte("test updated value"),
		},
		{
			key:     []byte("test1"),
			value:   []byte("test value with words 1"),
			upValue: []byte("test updated value 1"),
		},
		{
			key:     []byte("test2"),
			value:   []byte("test value with words 2"),
			upValue: []byte("test updated value 2"),
		},
	}

	trie := NewTrie(nil)

	for _, data := range testData {
		if !trie.Put(data.key, data.value) {
			t.Errorf("key/value not added to tree: %s: %s", data.key, data.value)
		}
	}

	hash, err := trie.Hash()
	if err != nil {
		t.Error(err)
		return
	}

	if len(hash) != 32 {
		t.Errorf("invalid length of hash: expected 32 got %d", len(hash))
	}

	trie.Update(testData[0].key, []byte("new value for key 0"))

	newhash, err := trie.Hash()
	if err != nil {
		t.Error(err)
		return
	}

	if len(hash) != 32 {
		t.Errorf("invalid length of hash after update: expected 32 got %d", len(hash))
	}

	if newhash == hash {
		t.Error("hash not rebuilding after update")
	}
}

func TestTrie_PutUpdateAndCheck(t *testing.T) {
	testData := []*struct {
		key, value, upValue []byte
	}{
		{
			key:     []byte("test"),
			value:   []byte("test value with words"),
			upValue: []byte("test updated value"),
		},
		{
			key:     []byte("test1"),
			value:   []byte("test value with words 1"),
			upValue: []byte("test updated value 1"),
		},
		{
			key:     []byte("test2"),
			value:   []byte("test value with words 2"),
			upValue: []byte("test updated value 2"),
		},
	}

	trie := NewTrie(nil)

	for _, data := range testData {
		if !trie.Put(data.key, data.value) {
			t.Errorf("key/value not added: %s: %s", data.key, data.value)
		}
	}

	for _, data := range testData {
		rVal, ok := trie.Get(data.key)
		if !ok {
			t.Errorf("lookup failed for key %s", data.key)
		}

		if !bytes.Equal(rVal, []byte(data.value)) {
			t.Errorf("invalid value: expected %s got %s", data.value, rVal)
		}
	}

	for _, data := range testData {
		if !trie.Update(data.key, data.upValue) {
			t.Errorf("value not updated: %s: %s", data.key, data.upValue)
		}
	}

	for _, data := range testData {
		rVal, ok := trie.Get(data.key)
		if !ok {
			t.Errorf("lookup failed for key %s", data.key)
		}

		if !bytes.Equal(rVal, []byte(data.upValue)) {
			t.Errorf("invalid value: expected %s got %s", data.upValue, rVal)
		}
	}
}

func TestTrie_PutDeleteAndCheck(t *testing.T) {
	testData := []*struct {
		key, value []byte
	}{
		{
			key:   []byte("test"),
			value: []byte("test value with words"),
		},
		{
			key:   []byte("test1"),
			value: []byte("test value with words 1"),
		},
		{
			key:   []byte("test2"),
			value: []byte("test value with words 2"),
		},
	}

	trie := NewTrie(nil)

	// Add data
	for _, data := range testData {
		if !trie.Put(data.key, data.value) {
			t.Errorf("key/value not added: %s: %s", data.key, data.value)
		}
	}

	// Check added data
	for _, data := range testData {
		rVal, ok := trie.Get(data.key)
		if !ok {
			t.Errorf("get added: lookup failed for key %s", data.key)
		}

		if !bytes.Equal(rVal, []byte(data.value)) {
			t.Errorf("check added: invalid value: expected %s got %s", data.value, rVal)
		}
	}

	// Delete Data
	for _, data := range testData {
		if !trie.Delete(data.key) {
			t.Errorf("delete: value not deleted for key %s", data.key)
		}
	}

	// Check deleted data
	for _, data := range testData {
		_, ok := trie.Get(data.key)
		if ok {
			t.Errorf("data not deleted for key %s", data.key)
		}
	}
}

func TestTrie_CustomStartNode(t *testing.T) {
	key, val := []byte("test"), []byte("test value")

	trie := NewTrie(node.NewBranchNode())
	trie.Put(key, val)

	rVal, ok := trie.Get([]byte("test"))
	if !ok {
		t.Errorf("lookup failed for key %s", key)
	}

	if !bytes.Equal(rVal, val) {
		t.Errorf("check added: invalid value: expected %s got %s", val, rVal)
	}
}
