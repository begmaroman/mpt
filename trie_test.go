package mpt

import "testing"

func BenchmarkTrie_AddString(b *testing.B) {
	t := NewTrie()

	for i := 0; i < b.N; i++ {
		t.Add("testString", "testValue")
	}
}
