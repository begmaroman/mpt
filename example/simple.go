package main

import (
	"fmt"
	"github.com/begmaroman/mpt"
)

func main() {
	trie := mpt.NewTrie(nil)

	// Created new ExtensionNode and set key and value
	trie.Put([]byte("key"), []byte("val"))

	// Created BranchNode and insert into ExtensionNode
	trie.Put([]byte("key_new"), []byte("val_new"))

	// try to find value
	found, ok := trie.Get([]byte("key_new"))
	fmt.Println("value:", string(found))
	fmt.Println("found:", ok)
	fmt.Println()

	// try to update value
	ok = trie.Update([]byte("key_new"), []byte("val_new_updated"))
	fmt.Println("updated:", ok)
	fmt.Println()

	// try to find value
	found, ok = trie.Get([]byte("key_new"))
	fmt.Println("value:", string(found))
	fmt.Println("found:", ok)
	fmt.Println()

	ok = trie.Delete([]byte("key_new"))
	fmt.Println("deleted:", ok)
	found, ok = trie.Get([]byte("key_new"))
	fmt.Println("value:", string(found))
	fmt.Println("found:", ok)
}
