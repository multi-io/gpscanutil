// Command gpscanutil analyzes and compares GrandPerspective scan files (.gpscan).
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/oklischat/gpscanutil/internal/gpscan"
)

const usage = `gpscanutil analyzes and compares GrandPerspective scan files (.gpscan).

Usage:
  gpscanutil introspect [options] <file>          list the biggest files/folders
  gpscanutil diff       [options] <file1> <file2> list the biggest changes from file1 to file2
  gpscanutil completion zsh|bash                  print a shell completion script

By default a flat list of full paths is printed, sorted by size (introspect)
or by absolute size change (diff). Folders are shown with a trailing "/".
With --tree, an indented tree is printed instead, with each folder's children
sorted the same way.

Options (may also be given after the file arguments):
  --top N        flat: show the N biggest entries (default 100)
                 tree: show the N biggest children per folder (default 10)
                 0 means no limit
  --depth N      only descend N levels below the root (default: unlimited
                 for flat lists, 3 for --tree)
  --min-size S   hide entries smaller than S (introspect) or with a smaller
                 change than S (diff), e.g. 500M, 1.5G, 100k (decimal units;
                 use Ki/Mi/Gi for binary)
  --root PATH    start at PATH instead of the volume root
  --files        flat list: only show files
  --dirs         only show folders
  --tree         print an indented tree instead of a flat list
  --bytes        print sizes as exact byte counts
  -q             don't print progress messages to stderr

Shell completion: add one of these to your ~/.zshrc or ~/.bashrc:
  eval "$(gpscanutil completion zsh)"     # after compinit
  eval "$(gpscanutil completion bash)"
`

type options struct {
	top     int
	depth   int
	minSize sizeFlag
	root    string
	files   bool
	dirs    bool
	tree    bool
	bytes   bool
	quiet   bool
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	var err error
	switch cmd := os.Args[1]; cmd {
	case "introspect":
		err = runIntrospect(os.Args[2:])
	case "diff":
		err = runDiff(os.Args[2:])
	case "completion":
		err = runCompletion(os.Args[2:])
	case "help", "-h", "-help", "--help":
		fmt.Print(usage)
		return
	default:
		err = usageError(fmt.Sprintf("unknown command %q", cmd))
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "gpscanutil:", err)
		var ue usageError
		if errors.As(err, &ue) {
			fmt.Fprintln(os.Stderr, "Run 'gpscanutil help' for usage.")
			os.Exit(2)
		}
		os.Exit(1)
	}
}

type usageError string

func (e usageError) Error() string { return string(e) }

// newFlagSet registers the options. Keep flagSpecs (completion.go) in sync.
func newFlagSet(name string, o *options) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.IntVar(&o.top, "top", -1, "")
	fs.IntVar(&o.depth, "depth", -1, "")
	fs.Var(&o.minSize, "min-size", "")
	fs.StringVar(&o.root, "root", "", "")
	fs.BoolVar(&o.files, "files", false, "")
	fs.BoolVar(&o.dirs, "dirs", false, "")
	fs.BoolVar(&o.tree, "tree", false, "")
	fs.BoolVar(&o.bytes, "bytes", false, "")
	fs.BoolVar(&o.quiet, "q", false, "")
	return fs
}

// parseArgs parses options and returns the positional arguments. Unlike the
// plain flag package, it accepts options after positional arguments too.
func parseArgs(name string, args []string, nargs int) (*options, []string, error) {
	o := &options{}
	fs := newFlagSet(name, o)
	var pos []string
	for {
		if err := fs.Parse(args); err != nil {
			if errors.Is(err, flag.ErrHelp) {
				fmt.Print(usage)
				os.Exit(0)
			}
			return nil, nil, usageError(err.Error())
		}
		if fs.NArg() == 0 {
			break
		}
		pos = append(pos, fs.Arg(0))
		args = fs.Args()[1:]
	}
	if len(pos) != nargs {
		return nil, nil, usageError(fmt.Sprintf("%s expects %d file argument(s), got %d", name, nargs, len(pos)))
	}
	if o.files && o.dirs {
		return nil, nil, usageError("--files and --dirs are mutually exclusive")
	}
	if o.files && o.tree {
		return nil, nil, usageError("--files can't be combined with --tree")
	}
	if o.top < 0 {
		o.top = 100
		if o.tree {
			o.top = 10
		}
	}
	if o.depth < 0 && o.tree {
		o.depth = 3
	}
	return o, pos, nil
}

// include reports whether a node of the given type passes --files/--dirs.
func (o *options) include(dir bool) bool {
	return !(o.files && dir) && !(o.dirs && !dir)
}

// withinDepth reports whether depth (relative to the root) may be shown.
func (o *options) withinDepth(depth int) bool {
	return o.depth < 0 || depth <= o.depth
}

func load(path string, quiet bool) (*gpscan.Tree, error) {
	start := time.Now()
	t, err := gpscan.Load(path)
	if err != nil {
		return nil, err
	}
	if !quiet {
		fmt.Fprintf(os.Stderr, "read %s (%s items) in %.1fs\n",
			path, formatCount(t.NumFiles+t.NumDirs), time.Since(start).Seconds())
	}
	return t, nil
}

// resolveRoot returns the node for --root, or the tree root.
func resolveRoot(t *gpscan.Tree, root string) int32 {
	if root == "" {
		return t.Root()
	}
	return t.Lookup(root)
}
