package mpt_test

import (
	"bytes"
	"github.com/begmaroman/mpt"
	"math/rand"
	"testing"

	"github.com/begmaroman/mpt/node"
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

	trie := mpt.NewTrie(nil)

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

	trie := mpt.NewTrie(nil)

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

	trie := mpt.NewTrie(nil)

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

	trie := mpt.NewTrie(nil)

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

	trie := mpt.NewTrie(nil)

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

func TestTrie_CustomStartNode_BranchNode_Empty(t *testing.T) {
	key, val := []byte("test"), []byte("test value")

	trie := mpt.NewTrie(node.NewBranchNode())
	trie.Put(key, val)

	rVal, ok := trie.Get([]byte("test"))
	if !ok {
		t.Errorf("lookup failed for key %s", key)
	}

	if !bytes.Equal(rVal, val) {
		t.Errorf("check added: invalid value: expected %s got %s", val, rVal)
	}
}

func TestTrie_CustomStartNode_ExtensionNode_Empty(t *testing.T) {
	key, val := []byte("test"), []byte("test value")

	trie := mpt.NewTrie(node.NewExtensionNode(nil, nil))
	trie.Put(key, val)

	rVal, ok := trie.Get([]byte("test"))
	if !ok {
		t.Errorf("lookup failed for key %s", key)
	}

	if !bytes.Equal(rVal, val) {
		t.Errorf("check added: invalid value: expected %s got %s", val, rVal)
	}
}

// Benchmarks
func BenchmarkTrie_Get(b *testing.B) {
	key, val := []byte("tolongwseb;kjnwkjlevnoknbkomnrwobnwotrh34y6onyo"), []byte("lwebniopwjgnipwuhgpv wrhtpovm rhiopwerhgpowuihnopwrntgmopkrtmvport uiwthgipurthg puti")

	fullTrie := mpt.NewTrie(nil)
	fullTrie.Put(key, val)
	for j := 0; j < 100000; j++ {
		tokenKey, tokenValue := make([]byte, 56), make([]byte, 556)
		rand.Read(tokenKey)
		rand.Read(tokenValue)
		fullTrie.Put(tokenKey, tokenValue)
	}
	b.Run("FromFullTrie100000", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fullTrie.Get(key)
		}
	})

	emptyTrie := mpt.NewTrie(nil)
	emptyTrie.Put(key, val)
	b.Run("FromEmptyTrie", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			emptyTrie.Get(key)
		}
	})
}

func BenchmarkTrie_Put(b *testing.B) {
	key, val := []byte("tolongwseb;kjnwkjlevnoknbkomnrwobnwotrh34y6onyo"), []byte("lwebniopwjgnipwuhgpv wrhtpovm rhiopwerhgpowuihnopwrntgmopkrtmvport uiwthgipurthg puti")

	fullTrie := mpt.NewTrie(nil)
	for j := 0; j < 100000; j++ {
		tokenKey, tokenValue := make([]byte, 56), make([]byte, 556)
		rand.Read(tokenKey)
		rand.Read(tokenValue)
		fullTrie.Put(tokenKey, tokenValue)
	}
	b.Run("ToFullTrie100000", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fullTrie.Put(key, val)
		}
	})

	emptyTrie := mpt.NewTrie(nil)
	emptyTrie.Put(key, val)
	b.Run("ToEmptyTrie", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			emptyTrie.Put(key, val)
		}
	})

	preparedBranchTrie := mpt.NewTrie(node.NewBranchNode())
	preparedBranchTrie.Put(key, val)
	b.Run("ToPreparedBranchNodeEmptyTrie", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			preparedBranchTrie.Put(key, val)
		}
	})

	preparedExtTrie := mpt.NewTrie(node.NewExtensionNode(nil, nil))
	preparedExtTrie.Put(key, val)
	b.Run("ToPreparedBranchNodeEmptyTrie", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			preparedExtTrie.Put(key, val)
		}
	})
}
