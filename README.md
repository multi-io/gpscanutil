# gputil

Command line tool for analyzing and comparing [GrandPerspective](https://grandperspectiv.sourceforge.net/)
scan files (`.gpscan`, saved via *File → Save Scan Data* in GrandPerspective).

```
go build -o gputil ./cmd/gputil      # or: go install ./cmd/gputil
```

## Usage

```
gputil introspect [options] <file>            # biggest files/folders
gputil diff       [options] <file1> <file2>   # biggest changes from file1 to file2
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
gputil introspect scan.gpscan --files --top 20
gputil introspect scan.gpscan --tree --root ~/Library --depth 2
gputil diff old.gpscan new.gpscan --min-size 100M
gputil diff old.gpscan new.gpscan --tree --top 5
```

## Notes

Scan files of ~1.3 GB (~10M items) take about 5 s to read, and need about 1 GB
of memory per scan. The file is memory-mapped and parsed with a minimal
purpose-built XML scanner (`internal/gpscan`), since `encoding/xml` is far
too slow for files of this size.
