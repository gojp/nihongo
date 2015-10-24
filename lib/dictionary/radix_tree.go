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
	label  []rune
}

func (r RadixEdge) String() string {
	return string(r.label)
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

func (r *RadixTree) findLastMatchingNode(key []rune) (n *RadixNode, ef int) {
	n = r.Root
	ef = 0

	for n != nil && ef < len(key) {
		suffix := key[ef:]
		var nextEdge *RadixEdge
		for i := range n.edges {
			if strings.HasPrefix(string(suffix), string(n.edges[i].label)) {
				nextEdge = &n.edges[i]
				break
			}
		}

		if nextEdge == nil {
			break
		}

		n = nextEdge.target

		ef += len(nextEdge.label)
	}

	return n, ef
}

func compareRunes(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i, r := range a {
		if b[i] != r {
			return false
		}
	}
	return true
}

func (r *RadixTree) Insert(key []rune, id EntryID) {
	n, ef := r.findLastMatchingNode(key)
	if n == nil || ef > len(key) {
		return
	}
	if ef == len(key) {
		// key already exists, so add id to this node
		if n.ids == nil {
			n.ids = []EntryID{}
		}
		n.ids = append(n.ids, id)
		return
	}
	// check if an outgoing edge shares a prefix with us
	suffix := key[ef:]
	prefix := []rune{}
	sharedEdge := -1

	for i := range n.edges {
	inner:
		for u := 0; u < len(n.edges[i].label) && u < len(suffix); u++ {
			if n.edges[i].label[u] == suffix[u] {
				prefix = append(prefix, suffix[u:u+1]...)
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
		return
	}

	oldEdge := n.edges[sharedEdge]

	node := RadixNode{edges: []RadixEdge{}, ids: []EntryID{id}}
	var left RadixEdge
	if !compareRunes(prefix, suffix) {
		left = RadixEdge{target: &node, label: prefix}
	} else {
		left = RadixEdge{target: &node, label: suffix[len(prefix):]}
	}
	right := RadixEdge{target: oldEdge.target, label: oldEdge.label[len(prefix):]}
	child := RadixNode{}

	child.edges = []RadixEdge{left, right}
	n.edges[sharedEdge] = RadixEdge{target: &child, label: prefix}
}

func (r *RadixTree) Get(key []rune) []EntryID {
	n, ef := r.findLastMatchingNode(key)

	// A match is found if we arrive at a leaf node and have used up exactly len(key) elements
	if n != nil && n.IsLeaf() && ef == len(key) {
		return n.Value()
	}

	return nil
}
