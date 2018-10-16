package mpt

import "testing"

func BenchmarkTrie_Add(b *testing.B) {
	t := NewTrie()
	str := []byte("test string")
	for i := 0; i < b.N; i++ {
		t.Add(str)
	}
}
