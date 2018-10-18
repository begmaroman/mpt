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
    
6. Get hash of the trie:
    
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

// try to get hash of the trie
h, err := trie.Hash()
if err != nil {
	// error while hash calculation
}

fmt.Println("hash of the trie:", h)

```

## Benchmarks

|Name|Iterations|Time|Description|
|----|----------|----|-----------|
|BenchmarkTrie_Get/FromFullTrie100000-4|1000000|3463 ns/op|Get data from the tree which have 100000 key/value pairs|
|BenchmarkTrie_Get/FromEmptyTrie-4|1000000|1262 ns/op|Get data from the empty tree|
|BenchmarkTrie_Put/ToFullTrie100000-4|3000000|513 ns/op|Put data to the tree which have 100000 key/value pairs|
|BenchmarkTrie_Put/ToEmptyTrie-4|3000000|524 ns/op|Put data to empty tree|
|BenchmarkTrie_Update/FullTrie100000-4|1000000|3083 ns/op|Update data in the tree which have 100000 key/value pairs|
|BenchmarkTrie_Update/EmptyTrie-4|2000000|674 ns/op|Update data in empty tree|
