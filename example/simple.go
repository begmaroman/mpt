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

	found, err := trie.Get([]byte("key_new"))
	fmt.Println("found value:", string(found))
	fmt.Println("error:", err)
}
