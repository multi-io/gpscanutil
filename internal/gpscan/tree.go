// Package gpscan reads GrandPerspective scan dumps (.gpscan files) into a
// compact in-memory tree and matches up the nodes of two such trees.
package gpscan

import (
	"bytes"
	"slices"
	"strings"
)

// None is the node index used for "no node".
const None int32 = -1

// Info holds the scan metadata from the file header.
type Info struct {
	AppVersion    string
	FormatVersion string
	VolumePath    string
	VolumeSize    int64
	FreeSpace     int64
	ScanTime      string
	SizeMeasure   string // "physical" or "logical"
}

// Node is a file or folder. Folder sizes are the sum of their contents.
// Nodes are kept small (32 bytes) because a scan may contain ~10M of them.
type Node struct {
	Size    int64
	nameOff uint32
	nameLen uint32
	Parent  int32
	first   int32 // first child
	next    int32 // next sibling
	Dir     bool
}

// Tree is a parsed scan. Node 0 is the root folder.
type Tree struct {
	Info     Info
	Nodes    []Node
	NumFiles int
	NumDirs  int
	names    []byte
}

// Root returns the root node index.
func (t *Tree) Root() int32 {
	if len(t.Nodes) == 0 {
		return None
	}
	return 0
}

// NameBytes returns the node's name without copying.
func (t *Tree) NameBytes(i int32) []byte {
	n := &t.Nodes[i]
	return t.names[n.nameOff : n.nameOff+n.nameLen]
}

// Name returns the node's name.
func (t *Tree) Name(i int32) string { return string(t.NameBytes(i)) }

// Children appends the indices of i's children to dst.
func (t *Tree) Children(i int32, dst []int32) []int32 {
	for c := t.Nodes[i].first; c != None; c = t.Nodes[c].next {
		dst = append(dst, c)
	}
	return dst
}

// FirstChild and NextSibling allow allocation-free iteration over children.
func (t *Tree) FirstChild(i int32) int32  { return t.Nodes[i].first }
func (t *Tree) NextSibling(i int32) int32 { return t.Nodes[i].next }

func (t *Tree) rootPath() string {
	if t.Info.VolumePath != "" {
		return t.Info.VolumePath
	}
	return "/"
}

// Path returns the absolute path of node i. The root maps to the volume path.
func (t *Tree) Path(i int32) string {
	var parts []string
	for ; i != None && t.Nodes[i].Parent != None; i = t.Nodes[i].Parent {
		parts = append(parts, t.Name(i))
	}
	slices.Reverse(parts)
	root := t.rootPath()
	if len(parts) == 0 {
		return root
	}
	return strings.TrimSuffix(root, "/") + "/" + strings.Join(parts, "/")
}

// Lookup finds the node with the given absolute path (or a path relative to
// the scan root). It returns None if there is no such node.
func (t *Tree) Lookup(path string) int32 {
	i := t.Root()
	if i == None {
		return None
	}
	root := strings.TrimSuffix(t.rootPath(), "/")
	if path == root || path == root+"/" {
		return i
	}
	if root != "" && strings.HasPrefix(path, root+"/") {
		path = path[len(root):]
	}
	for _, part := range strings.Split(path, "/") {
		if part == "" || part == "." {
			continue
		}
		next := None
		for c := t.Nodes[i].first; c != None; c = t.Nodes[c].next {
			if string(t.NameBytes(c)) == part {
				next = c
				break
			}
		}
		if next == None {
			return None
		}
		i = next
	}
	return i
}

// Pair is a node of tree A matched with the node of the same name and type in
// tree B. One side is None if the node exists in only one tree.
type Pair struct{ A, B int32 }

// Matcher matches up children of nodes in two trees by name and type.
type Matcher struct {
	A, B   *Tree
	ca, cb []int32
}

// ChildPairs returns the matched children of a (in m.A) and b (in m.B).
// Either a or b may be None, in which case all children of the other are
// returned unmatched.
func (m *Matcher) ChildPairs(a, b int32) []Pair {
	m.ca = m.ca[:0]
	m.cb = m.cb[:0]
	if a != None {
		m.ca = m.A.Children(a, m.ca)
	}
	if b != None {
		m.cb = m.B.Children(b, m.cb)
	}
	pairs := make([]Pair, 0, max(len(m.ca), len(m.cb)))
	if len(m.ca) == 0 {
		for _, y := range m.cb {
			pairs = append(pairs, Pair{None, y})
		}
		return pairs
	}
	if len(m.cb) == 0 {
		for _, x := range m.ca {
			pairs = append(pairs, Pair{x, None})
		}
		return pairs
	}
	slices.SortFunc(m.ca, func(x, y int32) int { return compareNodes(m.A, x, m.A, y) })
	slices.SortFunc(m.cb, func(x, y int32) int { return compareNodes(m.B, x, m.B, y) })
	i, j := 0, 0
	for i < len(m.ca) || j < len(m.cb) {
		var c int
		switch {
		case i == len(m.ca):
			c = 1
		case j == len(m.cb):
			c = -1
		default:
			c = compareNodes(m.A, m.ca[i], m.B, m.cb[j])
		}
		switch {
		case c < 0:
			pairs = append(pairs, Pair{m.ca[i], None})
			i++
		case c > 0:
			pairs = append(pairs, Pair{None, m.cb[j]})
			j++
		default:
			pairs = append(pairs, Pair{m.ca[i], m.cb[j]})
			i++
			j++
		}
	}
	return pairs
}

// compareNodes orders by name, then files before folders.
func compareNodes(ta *Tree, x int32, tb *Tree, y int32) int {
	if c := bytes.Compare(ta.NameBytes(x), tb.NameBytes(y)); c != 0 {
		return c
	}
	dx, dy := ta.Nodes[x].Dir, tb.Nodes[y].Dir
	switch {
	case dx == dy:
		return 0
	case !dx:
		return -1
	default:
		return 1
	}
}
