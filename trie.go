package mpt

type Trie struct {
	root *Node

	nodes []*Node
}

func NewTrie() *Trie {
	return &Trie{}
}

func (t *Trie) Root() *Node {
	return t.root
}

// Add add new Node to trie.
func (t *Trie) Add(value []byte) error {
	node := NewNode()
	node.Value = value
	if err := node.LoadChecksum(); err != nil {
		return err
	}

	nodes := t.root.leaf(nil)
	nodes = append(nodes, node)

	return t.build(nodes)
}

// build build Merkle Patricia Tries based nodes
func (t *Trie) build(nodes []*Node) error {
	if nodes == nil {
		return nil
	}

	var err error
	newNodes := nodes
	for {
		newNodes, err = buildNewLevel(newNodes)
		if err != nil {
			return err
		}

		if len(newNodes) == 1 {
			break
		}
	}
	t.root = newNodes[0]
	return nil
}

// buildNewLevel build new level based on nodes list
func buildNewLevel(nodes []*Node) ([]*Node, error) {
	newNodes := make([]*Node, 0, len(nodes)/2)
	last := len(nodes) - 1

	for i := range nodes {
		if i%2 == 0 {
			if i == last {
				newNodes = append(newNodes, nodes[i])
				continue
			}

			n := NewNode()
			n.Left = nodes[i]
			n.Left.Parent = n
			newNodes = append(newNodes, n)
		} else {
			n := newNodes[len(newNodes)-1]
			n.Right = nodes[i]
			n.Right.Parent = n
			if err := n.LoadChecksum(); err != nil {
				return nil, err
			}
		}
	}

	return newNodes, nil
}
