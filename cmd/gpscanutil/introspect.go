package main

import (
	"fmt"
	"slices"

	"github.com/multi-io/gpscanutil/internal/gpscan"
)

func runIntrospect(args []string) error {
	o, files, err := parseArgs("introspect", args, 1)
	if err != nil {
		return err
	}
	t, err := load(files[0], o.quiet)
	if err != nil {
		return err
	}
	root := resolveRoot(t, o.root)
	if root == gpscan.None {
		return fmt.Errorf("%s not found in %s", o.root, files[0])
	}

	printScanHeader(files[0], t)
	fmt.Println()
	if o.tree {
		printSizeTree(t, root, o)
	} else {
		printSizeList(t, root, o)
	}
	return nil
}

func printScanHeader(file string, t *gpscan.Tree) {
	i := t.Info
	fmt.Printf("Scan:         %s  (%s, GrandPerspective %s, %s sizes)\n", file, i.ScanTime, i.AppVersion, i.SizeMeasure)
	fmt.Printf("Volume:       %s  size %s, used %s, free %s\n", i.VolumePath,
		formatSize(i.VolumeSize, false), formatSize(usedSpace(t), false), formatSize(i.FreeSpace, false))
	fmt.Printf("Scanned:      %s in %s files, %s folders\n",
		formatSize(t.Nodes[t.Root()].Size, false), formatCount(t.NumFiles), formatCount(t.NumDirs))
	fmt.Printf("Not in scan:  %s  (used space not accounted for by scanned files)\n", formatSize(notInScan(t), false))
}

// usedSpace returns the used space of the volume. On APFS this includes
// everything in the container: other volumes (VM, Preboot, ...), snapshots,
// files deleted while still open, and folders the scan couldn't read.
func usedSpace(t *gpscan.Tree) int64 { return t.Info.VolumeSize - t.Info.FreeSpace }

// notInScan returns the used space that the scanned files don't account for.
// It can be negative, e.g. when clones or hard links are counted repeatedly.
func notInScan(t *gpscan.Tree) int64 { return usedSpace(t) - t.Nodes[t.Root()].Size }

type sized struct {
	size int64
	idx  int32
}

func biggerFirst(a, b sized) bool {
	if a.size != b.size {
		return a.size > b.size
	}
	return a.idx < b.idx
}

func displayPath(t *gpscan.Tree, i int32) string {
	p := t.Path(i)
	if t.Nodes[i].Dir && p != "/" {
		p += "/"
	}
	return p
}

func printSizeList(t *gpscan.Tree, root int32, o *options) {
	top := newTopK(o.top, biggerFirst)
	min := int64(o.minSize)
	var walk func(i int32, depth int)
	walk = func(i int32, depth int) {
		n := &t.Nodes[i]
		if n.Size < min {
			return // descendants can't be bigger
		}
		if o.include(n.Dir) {
			top.add(sized{n.Size, i})
		}
		if !n.Dir || !o.withinDepth(depth+1) {
			return
		}
		for c := t.FirstChild(i); c != gpscan.None; c = t.NextSibling(c) {
			walk(c, depth+1)
		}
	}
	walk(root, 0)
	for _, e := range top.sorted() {
		fmt.Printf("%12s  %s\n", formatSize(e.size, o.bytes), displayPath(t, e.idx))
	}
}

func printSizeTree(t *gpscan.Tree, root int32, o *options) {
	fmt.Printf("%12s  %s\n", formatSize(t.Nodes[root].Size, o.bytes), displayPath(t, root))
	var walk func(i int32, depth int, prefix string)
	walk = func(i int32, depth int, prefix string) {
		var kids []sized
		for c := t.FirstChild(i); c != gpscan.None; c = t.NextSibling(c) {
			if o.include(t.Nodes[c].Dir) {
				kids = append(kids, sized{t.Nodes[c].Size, c})
			}
		}
		slices.SortFunc(kids, func(a, b sized) int {
			if biggerFirst(a, b) {
				return -1
			}
			return 1
		})
		shown := 0
		for shown < len(kids) && (o.top <= 0 || shown < o.top) && kids[shown].size >= int64(o.minSize) {
			shown++
		}
		var restSize int64
		for _, k := range kids[shown:] {
			restSize += k.size
		}
		for n, k := range kids[:shown] {
			last := n == shown-1 && shown == len(kids)
			branch, indent := "├── ", "│   "
			if last {
				branch, indent = "└── ", "    "
			}
			name := t.Name(k.idx)
			if t.Nodes[k.idx].Dir {
				name += "/"
			}
			fmt.Printf("%12s  %s%s%s\n", formatSize(k.size, o.bytes), prefix, branch, name)
			if t.Nodes[k.idx].Dir && o.withinDepth(depth+1) {
				walk(k.idx, depth+1, prefix+indent)
			}
		}
		if rest := len(kids) - shown; rest > 0 {
			fmt.Printf("%12s  %s└── … %s more\n", formatSize(restSize, o.bytes), prefix, formatCount(rest))
		}
	}
	if t.Nodes[root].Dir && o.withinDepth(1) {
		walk(root, 1, "")
	}
}
