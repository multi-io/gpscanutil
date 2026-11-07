package main

import (
	"container/heap"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var units = []string{"B", "kB", "MB", "GB", "TB", "PB", "EB"}

// formatSize formats n with decimal units (as Finder and GrandPerspective do).
func formatSize(n int64, exact bool) string {
	if exact {
		return formatCount(int(n))
	}
	if n < 1000 && n > -1000 {
		return fmt.Sprintf("%d B", n)
	}
	f := float64(n)
	u := 0
	for (f >= 999.95 || f <= -999.95) && u < len(units)-1 {
		f /= 1000
		u++
	}
	return fmt.Sprintf("%.1f %s", f, units[u])
}

func formatDelta(n int64, exact bool) string {
	if n > 0 {
		return "+" + formatSize(n, exact)
	}
	return formatSize(n, exact)
}

// formatCount formats n with thousands separators.
func formatCount(n int) string {
	s := strconv.Itoa(n)
	neg := strings.HasPrefix(s, "-")
	s = strings.TrimPrefix(s, "-")
	var b strings.Builder
	if neg {
		b.WriteByte('-')
	}
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(c)
	}
	return b.String()
}

func formatCountDelta(n int) string {
	if n > 0 {
		return "+" + formatCount(n)
	}
	return formatCount(n)
}

// sizeFlag is a byte size given as e.g. "1.5G", "500M", "4Ki" or "1024".
type sizeFlag int64

func (s *sizeFlag) String() string { return strconv.FormatInt(int64(*s), 10) }

func (s *sizeFlag) Set(v string) error {
	str := strings.TrimSuffix(strings.TrimSpace(v), "B")
	mult := 1.0
	binary := strings.HasSuffix(str, "i")
	str = strings.TrimSuffix(str, "i")
	if str != "" {
		base := 1000.0
		if binary {
			base = 1024
		}
		if i := strings.IndexByte("kmgtpe", lower(str[len(str)-1])); i >= 0 {
			for range i + 1 {
				mult *= base
			}
			str = str[:len(str)-1]
		} else if binary {
			return fmt.Errorf("invalid size %q", v)
		}
	}
	f, err := strconv.ParseFloat(strings.TrimSpace(str), 64)
	if err != nil || f < 0 {
		return fmt.Errorf("invalid size %q", v)
	}
	*s = sizeFlag(f * mult)
	return nil
}

func lower(c byte) byte {
	if 'A' <= c && c <= 'Z' {
		return c + 'a' - 'A'
	}
	return c
}

// topK collects the k best items seen (all of them if k <= 0) without
// keeping the rest in memory.
type topK[T any] struct {
	k      int
	better func(a, b T) bool
	items  []T
}

func newTopK[T any](k int, better func(a, b T) bool) *topK[T] {
	return &topK[T]{k: k, better: better}
}

func (h *topK[T]) Len() int           { return len(h.items) }
func (h *topK[T]) Less(i, j int) bool { return h.better(h.items[j], h.items[i]) } // worst on top
func (h *topK[T]) Swap(i, j int)      { h.items[i], h.items[j] = h.items[j], h.items[i] }
func (h *topK[T]) Push(x any)         { h.items = append(h.items, x.(T)) }
func (h *topK[T]) Pop() any {
	x := h.items[len(h.items)-1]
	h.items = h.items[:len(h.items)-1]
	return x
}

func (h *topK[T]) add(x T) {
	switch {
	case h.k <= 0:
		h.items = append(h.items, x)
	case len(h.items) < h.k:
		heap.Push(h, x)
	case h.better(x, h.items[0]):
		h.items[0] = x
		heap.Fix(h, 0)
	}
}

// sorted returns the collected items, best first.
func (h *topK[T]) sorted() []T {
	slices.SortFunc(h.items, func(a, b T) int {
		switch {
		case h.better(a, b):
			return -1
		case h.better(b, a):
			return 1
		}
		return 0
	})
	return h.items
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
