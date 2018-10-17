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

	found, ok := trie.Get([]byte("key_new"))
	fmt.Println("value:", string(found))
	fmt.Println("found:", ok)

	ok = trie.Delete([]byte("key_new"))
	fmt.Println("deleted:", ok)
	found, ok = trie.Get([]byte("key_new"))
	fmt.Println("value:", string(found))
	fmt.Println("found:", ok)
}
