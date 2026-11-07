package main

import (
	"fmt"
	"os"
	"slices"
	"sync"

	"github.com/oklischat/gpscanutil/internal/gpscan"
)

func runDiff(args []string) error {
	o, files, err := parseArgs("diff", args, 2)
	if err != nil {
		return err
	}
	var trees [2]*gpscan.Tree
	var errs [2]error
	var wg sync.WaitGroup
	for i := range files {
		wg.Go(func() { trees[i], errs[i] = load(files[i], o.quiet) })
	}
	wg.Wait()
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	d := &differ{m: gpscan.Matcher{A: trees[0], B: trees[1]}, o: o}
	ra, rb := resolveRoot(trees[0], o.root), resolveRoot(trees[1], o.root)
	if ra == gpscan.None && rb == gpscan.None {
		return fmt.Errorf("%s not found in either scan", o.root)
	}

	printDiffHeader(files, trees)
	fmt.Println()
	root := gpscan.Pair{A: ra, B: rb}
	if o.tree {
		d.printTree(root)
	} else {
		d.printList(root)
	}
	return nil
}

func printDiffHeader(files []string, trees [2]*gpscan.Tree) {
	a, b := trees[0], trees[1]
	for n, t := range trees {
		fmt.Printf("Scan %d:  %s  (%s, %s sizes)\n", n+1, files[n], t.Info.ScanTime, t.Info.SizeMeasure)
	}
	if a.Info.SizeMeasure != b.Info.SizeMeasure {
		fmt.Fprintf(os.Stderr, "warning: scans use different size measures (%s vs %s); sizes are not comparable\n",
			a.Info.SizeMeasure, b.Info.SizeMeasure)
	}
	if a.Info.VolumePath != b.Info.VolumePath {
		fmt.Fprintf(os.Stderr, "warning: scans are of different volumes (%s vs %s)\n", a.Info.VolumePath, b.Info.VolumePath)
	}
	fmt.Println()
	fmt.Printf("%-12s %16s %16s %16s\n", "", "scan 1", "scan 2", "change")
	row := func(label string, x, y int64) {
		fmt.Printf("%-12s %16s %16s %16s\n", label, formatSize(x, false), formatSize(y, false), formatDelta(y-x, false))
	}
	countRow := func(label string, x, y int) {
		fmt.Printf("%-12s %16s %16s %16s\n", label, formatCount(x), formatCount(y), formatCountDelta(y-x))
	}
	row("Free space", a.Info.FreeSpace, b.Info.FreeSpace)
	row("Scanned", a.Nodes[a.Root()].Size, b.Nodes[b.Root()].Size)
	countRow("Files", a.NumFiles, b.NumFiles)
	countRow("Folders", a.NumDirs, b.NumDirs)
}

type differ struct {
	m gpscan.Matcher
	o *options
}

type change struct {
	gpscan.Pair
	old, new int64
}

func (c change) delta() int64 { return c.new - c.old }

func biggerChangeFirst(x, y change) bool {
	if dx, dy := abs(x.delta()), abs(y.delta()); dx != dy {
		return dx > dy
	}
	if x.B != y.B {
		return x.B < y.B
	}
	return x.A < y.A
}

func (d *differ) change(p gpscan.Pair) change {
	c := change{Pair: p}
	if p.A != gpscan.None {
		c.old = d.m.A.Nodes[p.A].Size
	}
	if p.B != gpscan.None {
		c.new = d.m.B.Nodes[p.B].Size
	}
	return c
}

func (d *differ) isDir(p gpscan.Pair) bool {
	if p.B != gpscan.None {
		return d.m.B.Nodes[p.B].Dir
	}
	return d.m.A.Nodes[p.A].Dir
}

func (d *differ) path(p gpscan.Pair) string {
	if p.B != gpscan.None {
		return displayPath(d.m.B, p.B)
	}
	return displayPath(d.m.A, p.A)
}

func (d *differ) name(p gpscan.Pair) string {
	var s string
	if p.B != gpscan.None {
		s = d.m.B.Name(p.B)
	} else {
		s = d.m.A.Name(p.A)
	}
	if d.isDir(p) {
		s += "/"
	}
	return s
}

// formatChange returns the "delta  old → new" columns.
func (d *differ) formatChange(c change) string {
	old, new := formatSize(c.old, d.o.bytes), formatSize(c.new, d.o.bytes)
	if c.A == gpscan.None {
		old = "(new)"
	}
	if c.B == gpscan.None {
		new = "(gone)"
	}
	return fmt.Sprintf("%12s  %11s → %-11s", formatDelta(c.delta(), d.o.bytes), old, new)
}

func (d *differ) printList(root gpscan.Pair) {
	top := newTopK(d.o.top, biggerChangeFirst)
	min := int64(d.o.minSize)
	var walk func(p gpscan.Pair, depth int)
	walk = func(p gpscan.Pair, depth int) {
		c := d.change(p)
		if max(c.old, c.new) < min {
			return // no descendant can have changed by min or more
		}
		dir := d.isDir(p)
		if delta := abs(c.delta()); delta != 0 && delta >= min && d.o.include(dir) {
			top.add(c)
		}
		if !dir || !d.o.withinDepth(depth+1) {
			return
		}
		for _, cp := range d.m.ChildPairs(p.A, p.B) {
			walk(cp, depth+1)
		}
	}
	walk(root, 0)
	for _, c := range top.sorted() {
		fmt.Printf("%s  %s\n", d.formatChange(c), d.path(c.Pair))
	}
}

func (d *differ) printTree(root gpscan.Pair) {
	fmt.Printf("%s  %s\n", d.formatChange(d.change(root)), d.path(root))
	var walk func(p gpscan.Pair, depth int, prefix string)
	walk = func(p gpscan.Pair, depth int, prefix string) {
		var kids []change
		for _, cp := range d.m.ChildPairs(p.A, p.B) {
			if c := d.change(cp); c.delta() != 0 && d.o.include(d.isDir(cp)) {
				kids = append(kids, c)
			}
		}
		slices.SortFunc(kids, func(x, y change) int {
			if biggerChangeFirst(x, y) {
				return -1
			}
			return 1
		})
		shown := 0
		for shown < len(kids) && (d.o.top <= 0 || shown < d.o.top) && abs(kids[shown].delta()) >= int64(d.o.minSize) {
			shown++
		}
		var restDelta int64
		for _, k := range kids[shown:] {
			restDelta += k.delta()
		}
		for n, k := range kids[:shown] {
			last := n == shown-1 && shown == len(kids)
			branch, indent := "├── ", "│   "
			if last {
				branch, indent = "└── ", "    "
			}
			fmt.Printf("%s  %s%s%s\n", d.formatChange(k), prefix, branch, d.name(k.Pair))
			if d.isDir(k.Pair) && d.o.withinDepth(depth+1) {
				walk(k.Pair, depth+1, prefix+indent)
			}
		}
		if rest := len(kids) - shown; rest > 0 {
			fmt.Printf("%12s  %11s   %-11s  %s└── … %s more changed (net)\n",
				formatDelta(restDelta, d.o.bytes), "", "", prefix, formatCount(rest))
		}
	}
	if d.isDir(root) && d.o.withinDepth(1) {
		walk(root, 1, "")
	}
}
