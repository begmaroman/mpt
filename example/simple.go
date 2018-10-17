package main

import (
	"crypto/sha1"
	"fmt"
	"reflect"

	"github.com/begmaroman/mpt"
)

func main() {
	trie := mpt.NewTrie()

	for i := 0; i < 10; i++ {
		trie.Add([]byte(fmt.Sprintf("testKey%d", i)))
	}

	root := trie.Root()

	dumpNode(root)
}

func dumpNode(node *mpt.Node) {
	if node == nil {
		return
	}

	sum := node.Checksum()

	if node.IsLeaf() {
		//fmt.Printf("value: %s, checksum: %v, is leaf \n", node.Value, sum)
		return
	}

	node.Left.LoadChecksum()
	node.Right.LoadChecksum()

	hash := sha1.New()
	hash.Write(node.Left.Checksum())
	hash.Write(node.Right.Checksum())

	fmt.Printf("valid checksum: %t \n", reflect.DeepEqual(hash.Sum(nil), sum))

	dumpNode(node.Left)
	dumpNode(node.Right)
}
