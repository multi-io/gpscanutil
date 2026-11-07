package main

import (
	"flag"
	"slices"
	"testing"
)

// TestFlagSpecsInSync checks that the completion specs match the real options.
func TestFlagSpecsInSync(t *testing.T) {
	fs := newFlagSet("test", &options{})
	var real, specs []string
	fs.VisitAll(func(f *flag.Flag) { real = append(real, f.Name) })
	for _, s := range flagSpecs {
		specs = append(specs, s.name)
		f := fs.Lookup(s.name)
		if f == nil {
			continue
		}
		bf, isBool := f.Value.(interface{ IsBoolFlag() bool })
		if isBool = isBool && bf.IsBoolFlag(); isBool != (s.arg == "") {
			t.Errorf("flagSpec %q: arg = %q, but boolean = %v", s.name, s.arg, isBool)
		}
	}
	slices.Sort(real)
	slices.Sort(specs)
	if !slices.Equal(real, specs) {
		t.Errorf("flagSpecs = %v, but options are %v", specs, real)
	}
}
