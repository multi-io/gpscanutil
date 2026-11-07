package gpscan

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strconv"
	"unicode/utf8"
)

// Load reads a .gpscan file.
func Load(path string) (*Tree, error) {
	data, release, err := mapFile(path)
	if err != nil {
		return nil, err
	}
	defer release()
	t, err := Parse(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return t, nil
}

type attr struct{ key, val []byte }

type parser struct {
	data    []byte
	pos     int
	t       *Tree
	stack   []int32
	attrs   []attr
	scratch []byte
}

// Parse parses the contents of a .gpscan file.
//
// This is a minimal, allocation-free XML scanner rather than encoding/xml,
// which is an order of magnitude too slow for multi-GB scan dumps. It only
// understands what GrandPerspective writes: elements with quoted attributes,
// the XML declaration, and comments. Unknown elements are ignored.
func Parse(data []byte) (*Tree, error) {
	p := &parser{
		data: data,
		t: &Tree{
			Nodes: make([]Node, 0, len(data)/160),
			names: make([]byte, 0, len(data)/12),
		},
	}
	if err := p.run(); err != nil {
		return nil, fmt.Errorf("offset %d: %w", p.pos, err)
	}
	if len(p.stack) != 0 {
		return nil, errors.New("unexpected end of file (unclosed Folder)")
	}
	if len(p.t.Nodes) == 0 {
		return nil, errors.New("no Folder elements found; is this a GrandPerspective scan file?")
	}
	return p.t, nil
}

func (p *parser) run() error {
	d := p.data
	for {
		i := bytes.IndexByte(d[p.pos:], '<')
		if i < 0 {
			return nil
		}
		p.pos += i + 1
		if p.pos >= len(d) {
			return errors.New("truncated tag")
		}
		switch d[p.pos] {
		case '?':
			if err := p.skipPast("?>"); err != nil {
				return err
			}
		case '!':
			end := ">"
			if bytes.HasPrefix(d[p.pos:], []byte("!--")) {
				end = "-->"
			}
			if err := p.skipPast(end); err != nil {
				return err
			}
		case '/':
			j := bytes.IndexByte(d[p.pos:], '>')
			if j < 0 {
				return errors.New("truncated closing tag")
			}
			name := bytes.TrimSpace(d[p.pos+1 : p.pos+j])
			p.pos += j + 1
			if string(name) == "Folder" {
				if err := p.closeFolder(); err != nil {
					return err
				}
			}
		default:
			if err := p.element(); err != nil {
				return err
			}
		}
	}
}

func (p *parser) skipPast(end string) error {
	j := bytes.Index(p.data[p.pos:], []byte(end))
	if j < 0 {
		return fmt.Errorf("missing %q", end)
	}
	p.pos += j + len(end)
	return nil
}

func isSpace(c byte) bool { return c == ' ' || c == '\t' || c == '\n' || c == '\r' }

// element parses a start (or self-closing) tag; p.pos is just after '<'.
func (p *parser) element() error {
	d := p.data
	start := p.pos
	for p.pos < len(d) && !isSpace(d[p.pos]) && d[p.pos] != '>' && d[p.pos] != '/' {
		p.pos++
	}
	tag := d[start:p.pos]
	p.attrs = p.attrs[:0]
	selfClosing := false
	for {
		for p.pos < len(d) && isSpace(d[p.pos]) {
			p.pos++
		}
		if p.pos >= len(d) {
			return errors.New("truncated tag")
		}
		if d[p.pos] == '>' {
			p.pos++
			break
		}
		if d[p.pos] == '/' {
			if p.pos+1 >= len(d) || d[p.pos+1] != '>' {
				return errors.New("malformed tag")
			}
			p.pos += 2
			selfClosing = true
			break
		}
		eq := bytes.IndexByte(d[p.pos:], '=')
		if eq < 0 || p.pos+eq+1 >= len(d) {
			return errors.New("malformed attribute")
		}
		key := bytes.TrimSpace(d[p.pos : p.pos+eq])
		p.pos += eq + 1
		for p.pos < len(d) && isSpace(d[p.pos]) {
			p.pos++
		}
		if p.pos >= len(d) || (d[p.pos] != '"' && d[p.pos] != '\'') {
			return fmt.Errorf("unquoted value for attribute %q", key)
		}
		q := d[p.pos]
		p.pos++
		end := bytes.IndexByte(d[p.pos:], q)
		if end < 0 {
			return errors.New("unterminated attribute value")
		}
		p.attrs = append(p.attrs, attr{key, d[p.pos : p.pos+end]})
		p.pos += end + 1
	}

	switch string(tag) {
	case "File":
		return p.file()
	case "Folder":
		if err := p.openFolder(); err != nil {
			return err
		}
		if selfClosing {
			return p.closeFolder()
		}
	case "ScanInfo":
		p.scanInfo()
	case "GrandPerspectiveScanDump":
		p.t.Info.AppVersion = string(p.attr("appVersion"))
		p.t.Info.FormatVersion = string(p.attr("formatVersion"))
	}
	return nil
}

func (p *parser) attr(key string) []byte {
	for _, a := range p.attrs {
		if string(a.key) == key {
			return a.val
		}
	}
	return nil
}

func (p *parser) scanInfo() {
	info := &p.t.Info
	info.VolumePath = string(p.unescape(p.attr("volumePath")))
	info.VolumeSize, _ = parseUint(p.attr("volumeSize"))
	info.FreeSpace, _ = parseUint(p.attr("freeSpace"))
	info.ScanTime = string(p.attr("scanTime"))
	info.SizeMeasure = string(p.attr("fileSizeMeasure"))
}

func (p *parser) newNode(dir bool) (int32, error) {
	t := p.t
	if len(t.Nodes) >= math.MaxInt32 {
		return None, errors.New("too many nodes")
	}
	parent := None
	if len(p.stack) > 0 {
		parent = p.stack[len(p.stack)-1]
	} else if len(t.Nodes) > 0 || !dir {
		return None, errors.New("element outside of the root Folder")
	}
	off := len(t.names)
	t.names = p.appendUnescaped(t.names, p.attr("name"))
	if len(t.names) > math.MaxUint32 {
		return None, errors.New("names too large")
	}
	idx := int32(len(t.Nodes))
	n := Node{
		nameOff: uint32(off),
		nameLen: uint32(len(t.names) - off),
		Parent:  parent,
		first:   None,
		next:    None,
		Dir:     dir,
	}
	if parent != None {
		n.next = t.Nodes[parent].first
		t.Nodes[parent].first = idx
	}
	t.Nodes = append(t.Nodes, n)
	return idx, nil
}

func (p *parser) openFolder() error {
	idx, err := p.newNode(true)
	if err != nil {
		return err
	}
	p.t.NumDirs++
	p.stack = append(p.stack, idx)
	return nil
}

func (p *parser) closeFolder() error {
	if len(p.stack) == 0 {
		return errors.New("unbalanced </Folder>")
	}
	f := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]
	if parent := p.t.Nodes[f].Parent; parent != None {
		p.t.Nodes[parent].Size += p.t.Nodes[f].Size
	}
	return nil
}

func (p *parser) file() error {
	size, err := parseUint(p.attr("size"))
	if err != nil {
		return fmt.Errorf("File %q: bad size: %w", p.attr("name"), err)
	}
	idx, err := p.newNode(false)
	if err != nil {
		return err
	}
	p.t.NumFiles++
	p.t.Nodes[idx].Size = size
	p.t.Nodes[p.t.Nodes[idx].Parent].Size += size
	return nil
}

func parseUint(b []byte) (int64, error) {
	if len(b) == 0 {
		return 0, errors.New("empty number")
	}
	var n int64
	for _, c := range b {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("invalid number %q", b)
		}
		if n > (math.MaxInt64-9)/10 {
			return 0, fmt.Errorf("number too large %q", b)
		}
		n = n*10 + int64(c-'0')
	}
	return n, nil
}

func (p *parser) unescape(s []byte) []byte {
	p.scratch = p.appendUnescaped(p.scratch[:0], s)
	return p.scratch
}

// appendUnescaped appends s to dst, resolving XML entity and character
// references. Unknown references are kept verbatim.
func (p *parser) appendUnescaped(dst, s []byte) []byte {
	for {
		i := bytes.IndexByte(s, '&')
		if i < 0 {
			return append(dst, s...)
		}
		dst = append(dst, s[:i]...)
		s = s[i:]
		semi := bytes.IndexByte(s, ';')
		if semi < 0 {
			return append(dst, s...)
		}
		ref := s[1:semi]
		switch string(ref) {
		case "amp":
			dst = append(dst, '&')
		case "lt":
			dst = append(dst, '<')
		case "gt":
			dst = append(dst, '>')
		case "quot":
			dst = append(dst, '"')
		case "apos":
			dst = append(dst, '\'')
		default:
			r, ok := charRef(ref)
			if !ok {
				dst = append(dst, s[:semi+1]...)
			} else {
				dst = utf8.AppendRune(dst, r)
			}
		}
		s = s[semi+1:]
	}
}

func charRef(ref []byte) (rune, bool) {
	if len(ref) < 2 || ref[0] != '#' {
		return 0, false
	}
	var v uint64
	var err error
	if ref[1] == 'x' || ref[1] == 'X' {
		v, err = strconv.ParseUint(string(ref[2:]), 16, 32)
	} else {
		v, err = strconv.ParseUint(string(ref[1:]), 10, 32)
	}
	if err != nil || v > utf8.MaxRune {
		return 0, false
	}
	return rune(v), true
}
