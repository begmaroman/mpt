# MPT - Modified Merkle Patricia Tries

Simple implementation of Modified Merkle Patricia Tries on GoLang.

## Usage

1. Create the Trie:

    - With empty root node:
    
    ```go
    trie := mpt.NewEmptyTrie()
    ```

    - With branch root node:
    
    ```go
    nd := node.NewBranchNode()
    trie := mpt.NewTrie(nd)
    ```

    - With extension root node:
    
    ```go
    nd := node.NewExtensionNode(nil, nil)
    trie := mpt.NewTrie(nd)
    ```
    
2. Add key/value pair:
    
```go
key := []byte("key")
val := []byte("val")

// create new trie and try add key/value pair
trie := mpt.NewEmptyTrie()
if ok := trie.Put(key, val); !ok {
	// key/pair not added
}

```

3. Get data by key:
    
```go
key := []byte("key")
val := []byte("val")

// create new trie and try add key/value pair
trie := mpt.NewEmptyTrie()
if ok := trie.Put(key, val); !ok {
	// key/pair not added
}

// some logic...

// try to get data by key
getVal, ok := trie.Get(key)
if !ok {
	// key/pair not added
}

```
    
4. Update data by key:
    
```go
// initialize data
key := []byte("key")
val := []byte("val")

// create new trie and try add key/value pair
trie := mpt.NewEmptyTrie()
if ok := trie.Put(key, val); !ok {
	// key/pair not added
}

// some logic...

// try to update value by key
newVal := []byte("updated value")
if ok := trie.Update(key, newVal); !ok {
	// data not updated
}

```
    
5. Delete data by key:
    
```go
// initialize data
key := []byte("key")
val := []byte("val")

// create new trie and try add key/value pair
trie := mpt.NewEmptyTrie()
if ok := trie.Put(key, val); !ok {
	// key/pair not added
}

// some logic...

// try to delete data by key
if ok := trie.Delete(key, newVal); !ok {
	// data not deleted
}

```