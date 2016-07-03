package dictionary

import (
	"fmt"
	"strings"
)

type RadixTree struct {
	Root *RadixNode
}

func (r RadixTree) String() string {
	return fmt.Sprintf("%v", r.Root.edges)
}

func NewRadixTree() *RadixTree {
	root := RadixNode{
		edges: []RadixEdge{},
		ids:   nil,
	}
	return &RadixTree{Root: &root}
}

type RadixEdge struct {
	target *RadixNode
	label  string
}

func (r RadixEdge) String() string {
	return r.label
}

type RadixNode struct {
	edges []RadixEdge
	ids   []EntryID
}

func (n RadixNode) IsLeaf() bool {
	return len(n.ids) > 0
}

func (n RadixNode) Value() []EntryID {
	return n.ids
}

func (n RadixNode) FindPrefixedEntries(max int) (entries []EntryID) {
	entries = []EntryID{}

	stack := []*RadixNode{&n}
	added := map[EntryID]bool{}

	var node *RadixNode
	for len(stack) > 0 {
		node, stack = stack[len(stack)-1], stack[:len(stack)-1]
		if node.IsLeaf() {
			for _, v := range node.Value() {
				if _, ok := added[v]; !ok {
					entries = append(entries, v)
				}
				added[v] = true
			}
		}
		if len(entries) >= max {
			return
		}
		for i := range node.edges {
			stack = append(stack, node.edges[i].target)
		}
	}
	return
}

func (r *RadixTree) findLastMatchingNode(key string) (n *RadixNode, elementsFound int) {
	n = r.Root
	elementsFound = 0

	for n != nil && elementsFound < len(key) {
		// Get the next edge to explore based on the elements not yet found in key
		// select edge from n.edges where edge.label is a prefix of key.suffix(elementsFound)
		suffix := key[elementsFound:]
		var nextEdge *RadixEdge
		for i := range n.edges {
			if strings.HasPrefix(suffix, n.edges[i].label) {
				nextEdge = &n.edges[i]
				break
			}
		}

		if nextEdge == nil {
			// terminate loop
			break
		}

		// Was an edge found?
		// Set the next node to explore
		n = nextEdge.target

		// Increment elements found based on the label stored at the edge
		elementsFound += len(nextEdge.label)
	}

	return n, elementsFound
}

func (r *RadixTree) Insert(key string, id EntryID) {
	n, elementsFound := r.findLastMatchingNode(key)
	if n == nil || elementsFound > len(key) {
		return
	}
	if elementsFound == len(key) {
		// key already exists, so add id to this node
		if n.ids == nil {
			n.ids = []EntryID{}
		}
		n.ids = append(n.ids, id)
	} else {
		// check if an outgoing edge shares a prefix with us
		suffix := key[elementsFound:]
		prefix := ""
		sharedEdge := -1

		for i := range n.edges {
		inner:
			for u := 0; u < len(n.edges[i].label) && u < len(suffix); u++ {
				if n.edges[i].label[u] == suffix[u] {
					prefix += suffix[u : u+1]
				} else {
					break inner
				}
			}
			// there can be at most one outgoing edge that shares a prefix
			if len(prefix) > 0 {
				sharedEdge = i
				break
			}
		}

		if sharedEdge == -1 {
			// create a new edge and node
			n.edges = append(n.edges, RadixEdge{target: &RadixNode{edges: []RadixEdge{}, ids: []EntryID{id}}, label: suffix})
		} else {
			oldEdge := n.edges[sharedEdge]

			child := RadixNode{ids: []EntryID{}}
			n.edges[sharedEdge] = RadixEdge{target: &child, label: prefix}

			node := RadixNode{edges: []RadixEdge{}, ids: []EntryID{id}}
			var left RadixEdge
			if prefix == suffix {
				left = RadixEdge{target: &node, label: prefix}
			} else {
				left = RadixEdge{target: &node, label: suffix[len(prefix):]}
			}
			right := RadixEdge{target: oldEdge.target, label: oldEdge.label[len(prefix):]}
			child.edges = []RadixEdge{left, right}
		}
	}
}

func (r *RadixTree) Get(key string) []EntryID {
	n, elementsFound := r.findLastMatchingNode(key)

	// A match is found if we arrive at a leaf node and have used up exactly len(key) elements
	if n != nil && elementsFound == len(key) {
		return n.Value()
	}

	return nil
}

func (r *RadixTree) FindWordsWithPrefix(key string, max int) []EntryID {
	words := []EntryID{}

	n, elementsFound := r.findLastMatchingNode(key)
	if n != nil {
		if elementsFound == len(key) {
			children := n.FindPrefixedEntries(max - len(words))
			words = append(words, children...)
		} else {
			// check if an outgoing edge shares a prefix with us
			suffix := key[elementsFound:]
			prefix := ""
			sharedEdge := -1

			for i := range n.edges {
			inner:
				for u := 0; u < len(n.edges[i].label) && u < len(suffix); u++ {
					if n.edges[i].label[u] == suffix[u] {
						prefix += suffix[u : u+1]
					} else {
						break inner
					}
				}
				// there can be at most one outgoing edge that shares a prefix
				if len(prefix) > 0 {
					sharedEdge = i
					break
				}
			}

			if sharedEdge >= 0 {
				children := n.edges[sharedEdge].target.FindPrefixedEntries(max - len(words))
				words = append(words, children...)
			}
		}
	}

	return words
}
