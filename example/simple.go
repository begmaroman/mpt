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
		node, _ := mpt.InitNode([]byte(fmt.Sprintf("testKey%d", i)))
		trie.Add(node)
	}

	root := trie.Root()

	sum, _ := root.Checksum()
	leftSum, _ := root.Left.Checksum()
	rightSum, _ := root.Right.Checksum()

	hash := sha1.New()
	hash.Write(leftSum)
	hash.Write(rightSum)

	fmt.Printf("parent   sum: %v\n", sum)
	fmt.Printf("left     sum: %v\n", leftSum)
	fmt.Printf("right    sum: %v\n", rightSum)
	fmt.Printf("expected sum: %v\n", hash.Sum(nil))

	dumpNode(trie.Root())
}

func dumpNode(node *mpt.Node) {
	if node == nil {
		return
	}

	sum, _ := node.Checksum()

	if node.IsLeaf() {
		//fmt.Printf("value: %s, checksum: %v, is leaf \n", node.Value, sum)
		return
	}

	leftSum, _ := node.Left.Checksum()
	rightSum, _ := node.Right.Checksum()

	hash := sha1.New()
	hash.Write(leftSum)
	hash.Write(rightSum)

	fmt.Printf("value: %v, valid checksum: %t \n", node.Value, reflect.DeepEqual(hash.Sum(nil), sum))

	dumpNode(node.Left)
	dumpNode(node.Right)
}
