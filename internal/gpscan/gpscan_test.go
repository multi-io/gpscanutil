package gpscan

import (
	"slices"
	"testing"
)

func mustLoad(t *testing.T, path string) *Tree {
	t.Helper()
	tr, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	return tr
}

func TestLoad(t *testing.T) {
	tr := mustLoad(t, "testdata/a.gpscan")
	if got, want := tr.Info.FreeSpace, int64(500000); got != want {
		t.Errorf("FreeSpace = %d, want %d", got, want)
	}
	if tr.Info.SizeMeasure != "physical" || tr.Info.AppVersion != "3.4.2" {
		t.Errorf("Info = %+v", tr.Info)
	}
	if tr.NumFiles != 4 || tr.NumDirs != 4 {
		t.Errorf("NumFiles, NumDirs = %d, %d, want 4, 4", tr.NumFiles, tr.NumDirs)
	}
	sizes := map[string]int64{
		"/":                          1360,
		"/Users":                     1310,
		"/Users/old":                 300,
		"/Users/old/x":               300,
		"/empty":                     0,
		`/Users/Tom & "Jerry" é.txt`: 10,
	}
	for path, want := range sizes {
		i := tr.Lookup(path)
		if i == None {
			t.Errorf("Lookup(%q) = None", path)
			continue
		}
		if got := tr.Nodes[i].Size; got != want {
			t.Errorf("size of %q = %d, want %d", path, got, want)
		}
		if got := tr.Path(i); got != path {
			t.Errorf("Path(Lookup(%q)) = %q", path, got)
		}
	}
	if tr.Lookup("/nope") != None {
		t.Error("Lookup of missing path succeeded")
	}
}

func TestParseErrors(t *testing.T) {
	for _, s := range []string{
		``,
		`<Folder name="">`,
		`<Folder name=""></Folder></Folder>`,
		`<File name="x" size="1"/>`,
		`<Folder name=""><File name="x" size="-1"/></Folder>`,
		`<Folder name=""></Folder><Folder name=""></Folder>`,
	} {
		if _, err := Parse([]byte(s)); err == nil {
			t.Errorf("Parse(%q) succeeded", s)
		}
	}
}

func TestChildPairs(t *testing.T) {
	a, b := mustLoad(t, "testdata/a.gpscan"), mustLoad(t, "testdata/b.gpscan")
	m := &Matcher{A: a, B: b}
	type row struct {
		name     string
		old, new int64
	}
	collect := func(pa, pb int32) []row {
		var rows []row
		for _, p := range m.ChildPairs(pa, pb) {
			r := row{}
			if p.A != None {
				r.name, r.old = a.Name(p.A), a.Nodes[p.A].Size
			}
			if p.B != None {
				r.name, r.new = b.Name(p.B), b.Nodes[p.B].Size
			}
			rows = append(rows, r)
		}
		return rows
	}
	if got, want := collect(a.Lookup("/Users"), b.Lookup("/Users")), []row{
		{`Tom & "Jerry" é.txt`, 10, 10},
		{"big.bin", 1000, 1500},
		{"new", 0, 200},
		{"old", 300, 0},
	}; !slices.Equal(got, want) {
		t.Errorf("ChildPairs(/Users) = %v, want %v", got, want)
	}
	if got, want := collect(a.Lookup("/Users/old"), None), []row{{"x", 300, 0}}; !slices.Equal(got, want) {
		t.Errorf("ChildPairs(/Users/old, None) = %v, want %v", got, want)
	}
}
