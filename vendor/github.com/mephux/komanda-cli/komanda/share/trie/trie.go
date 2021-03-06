package trie

import (
	"sort"
)

// Node struct
type Node struct {
	val      rune
	term     bool
	depth    int
	meta     interface{}
	mask     uint64
	parent   *Node
	children map[rune]*Node
}

// Trie struct
type Trie struct {
	root *Node
	size int
}

// ByKeys type cast for string array
type ByKeys []string

// Len of keys
func (a ByKeys) Len() int { return len(a) }

// Swap keys to reorder
func (a ByKeys) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Less checks if key i is less than j to determine position
func (a ByKeys) Less(i, j int) bool { return len(a[i]) < len(a[j]) }

const nul = 0x0

// New Trie with an initialized root Node.
func New() *Trie {
	return &Trie{
		root: &Node{children: make(map[rune]*Node), depth: 0},
		size: 0,
	}
}

// Root returns the root node for the Trie.
func (t *Trie) Root() *Node {
	return t.root
}

// Add the key to the Trie, including meta data. Meta data
// is stored as `interface{}` and must be type cast by
// the caller.
func (t *Trie) Add(key string, meta interface{}) *Node {
	t.size++
	runes := []rune(key)
	bitmask := maskruneslice(runes)
	node := t.root
	node.mask |= bitmask
	for i := range runes {
		r := runes[i]
		bitmask = maskruneslice(runes[i:])
		if n, ok := node.children[r]; ok {
			node = n
			node.mask |= bitmask
		} else {
			node = node.NewChild(r, bitmask, nil, false)
		}
	}
	node = node.NewChild(nul, 0, meta, true)
	return node
}

// Find and returns meta data associated
// with `key`.
func (t *Trie) Find(key string) (*Node, bool) {
	node := findNode(t.Root(), []rune(key))
	if node == nil {
		return nil, false
	}
	node = node.Children()[nul]

	if !node.term {
		return nil, false
	}

	return node, true
}

// Remove a key from the trie, ensuring that
// all bitmasks up to root are appropriately recalculated.
func (t *Trie) Remove(key string) {
	var (
		i    int
		rs   = []rune(key)
		node = findNode(t.Root(), []rune(key))
	)

	t.size--
	for n := node.Parent(); n != nil; n = n.Parent() {
		i++
		if len(n.Children()) > 1 {
			r := rs[len(rs)-i]
			n.RemoveChild(r)
			break
		}
	}
}

// Keys will return all the keys currently stored in the trie.
func (t *Trie) Keys() []string {
	return t.PrefixSearch("")
}

// FuzzySearch performs a fuzzy search against the keys in the trie.
func (t Trie) FuzzySearch(pre string) []string {
	keys := fuzzycollect(t.Root(), []rune(pre))
	sort.Sort(ByKeys(keys))
	return keys
}

// PrefixSearch performs a prefix search against the keys in the trie.
func (t Trie) PrefixSearch(pre string) []string {
	node := findNode(t.Root(), []rune(pre))
	if node == nil {
		return nil
	}

	return collect(node)
}

// NewChild creates and returns a pointer to a new child for the node.
func (n *Node) NewChild(val rune, bitmask uint64, meta interface{}, term bool) *Node {
	node := &Node{
		val:      val,
		mask:     bitmask,
		term:     term,
		meta:     meta,
		parent:   n,
		children: make(map[rune]*Node),
		depth:    n.depth + 1,
	}
	n.children[val] = node
	n.mask |= bitmask
	return node
}

// RemoveChild from list
func (n *Node) RemoveChild(r rune) {
	delete(n.children, r)
	for nd := n.parent; nd != nil; nd = nd.parent {
		nd.mask ^= nd.mask
		nd.mask |= uint64(1) << uint64(nd.val-'a')
		for _, c := range nd.children {
			nd.mask |= c.mask
		}
	}
}

// Parent returns the parent of this node.
func (n Node) Parent() *Node {
	return n.parent
}

// Meta returns the meta information of this node.
func (n Node) Meta() interface{} {
	return n.meta
}

// Children returns the children of this node.
func (n Node) Children() map[rune]*Node {
	return n.children
}

// Val returns the value of the node as a rune
func (n Node) Val() rune {
	return n.val
}

// Mask returns a uint64 representing the current
// mask of this node.
func (n Node) Mask() uint64 {
	return n.mask
}

func findNode(node *Node, runes []rune) *Node {
	if node == nil {
		return nil
	}

	if len(runes) == 0 {
		return node
	}

	n, ok := node.Children()[runes[0]]
	if !ok {
		return nil
	}

	var nrunes []rune
	if len(runes) > 1 {
		nrunes = runes[1:]
	} else {
		nrunes = runes[0:0]
	}

	return findNode(n, nrunes)
}

func maskruneslice(rs []rune) uint64 {
	var m uint64
	for _, r := range rs {
		m |= uint64(1) << uint64(r-'a')
	}
	return m
}

func collect(node *Node) []string {
	var (
		keys []string
		n    *Node
		i    int
	)
	nodes := []*Node{node}
	for l := len(nodes); l != 0; l = len(nodes) {
		i = l - 1
		n = nodes[i]
		nodes = nodes[:i]
		for _, c := range n.children {
			nodes = append(nodes, c)
		}
		if n.term {
			word := ""
			for p := n.parent; p.depth != 0; p = p.parent {
				word = string(p.val) + word
			}
			keys = append(keys, word)
		}
	}
	return keys
}

type potentialSubtree struct {
	idx  int
	node *Node
}

func fuzzycollect(node *Node, partial []rune) []string {
	var (
		m    uint64
		i    int
		p    potentialSubtree
		keys []string
	)

	potential := []potentialSubtree{{node: node, idx: 0}}
	for l := len(potential); l > 0; l = len(potential) {
		i = l - 1
		p = potential[i]
		potential = potential[:i]
		m = maskruneslice(partial[p.idx:])
		if (p.node.mask & m) != m {
			continue
		}

		if p.node.val == partial[p.idx] {
			p.idx++
			if p.idx == len(partial) {
				keys = append(keys, collect(p.node)...)
				continue
			}
		}

		for _, c := range p.node.children {
			potential = append(potential, potentialSubtree{node: c, idx: p.idx})
		}
	}
	return keys
}
