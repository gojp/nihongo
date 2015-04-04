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
	return n.ids != nil
}

func (n RadixNode) Value() []EntryID {
	return n.ids
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

		// Was an edge found?
		if nextEdge != nil {

			// Set the next node to explore
			n = nextEdge.target

			// Increment elements found based on the label stored at the edge
			elementsFound += len(nextEdge.label)
		} else {

			// terminate loop
			break
		}
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

			newChild := RadixNode{}
			n.edges[sharedEdge] = RadixEdge{target: &newChild, label: prefix}

			newNode := RadixNode{edges: []RadixEdge{}, ids: []EntryID{id}}
			newEdgeLeft := RadixEdge{target: &newNode, label: suffix[len(prefix):]}
			newEdgeRight := RadixEdge{target: oldEdge.target, label: oldEdge.label[len(prefix):]}
			newChild.edges = []RadixEdge{newEdgeLeft, newEdgeRight}
		}
	}
}

func (r *RadixTree) Get(key string) []EntryID {
	n, elementsFound := r.findLastMatchingNode(key)

	// A match is found if we arrive at a leaf node and have used up exactly len(key) elements
	if n != nil && n.IsLeaf() && elementsFound == len(key) {
		return n.Value()
	}

	return nil
}
