# gpscanutil

Command line tool for analyzing and comparing [GrandPerspective](https://grandperspectiv.sourceforge.net/)
scan files (`.gpscan`, saved via *File → Save Scan Data* in GrandPerspective).

```
go build -o gpscanutil ./cmd/gpscanutil      # or: go install ./cmd/gpscanutil
```

## Usage

```
gpscanutil introspect [options] <file>            # biggest files/folders
gpscanutil diff       [options] <file1> <file2>   # biggest changes from file1 to file2
```

By default a flat list of full paths is printed, sorted by size (`introspect`)
or by absolute size change (`diff`). Folders have a trailing `/`. In `diff`,
`(new)` / `(gone)` mark items that exist in only one of the scans. A moved or
renamed item shows up as one removal plus one addition.

Options (may be given before or after the file arguments):

| Option         | Meaning |
|----------------|---------|
| `--top N`      | flat: show the N biggest entries (default 100); tree: the N biggest children per folder (default 10); 0 = no limit |
| `--depth N`    | only descend N levels below the root (default: unlimited, 3 with `--tree`) |
| `--min-size S` | hide entries smaller than S / with a smaller change than S, e.g. `500M`, `1.5G` (decimal; `Gi` etc. for binary) |
| `--root PATH`  | start at PATH instead of the volume root |
| `--files`      | flat list: only files |
| `--dirs`       | only folders |
| `--tree`       | indented tree instead of a flat list |
| `--bytes`      | exact byte counts |
| `-q`           | no progress messages on stderr |

Examples:

```
gpscanutil introspect scan.gpscan --files --top 20
gpscanutil introspect scan.gpscan --tree --root ~/Library --depth 2
gpscanutil diff old.gpscan new.gpscan --min-size 100M
gpscanutil diff old.gpscan new.gpscan --tree --top 5
```

## Shell completion

`gpscanutil completion zsh|bash` prints a completion script. Add one of these
to your shell's rc file:

```
eval "$(gpscanutil completion zsh)"     # ~/.zshrc, after compinit
eval "$(gpscanutil completion bash)"    # ~/.bashrc (works with macOS's bash 3.2)
```

For zsh you can instead install it on your `$fpath`:
`gpscanutil completion zsh > "${fpath[1]}/_gpscanutil"`.

The scripts are generated from `flagSpecs` in `cmd/gpscanutil/completion.go`,
and a test ensures that list matches the actual options.

## Notes

Scan files of ~1.3 GB (~10M items) take about 5 s to read, and need about 1 GB
of memory per scan. The file is memory-mapped and parsed with a minimal
purpose-built XML scanner (`internal/gpscan`), since `encoding/xml` is far
too slow for files of this size.
