package flatfile

import (
	"bytes"
	"io"
	"io/ioutil"
	"sort"
	"strings"

	"git.biscorp.local/serverdev/errors"
)

// FlatFile consists of a slice of lines and a format applied to each line.
type FlatFile struct {
	lineFmt LineFmt
	lines   Lines
}

// New returns a flat file ready to read from a reader.
func New(lf LineFmt) FlatFile {
	return FlatFile{lineFmt: lf.Copy(), lines: make(Lines, 0)}
}

// Append several lines.
func (ff *FlatFile) Append(lns ...Line) {
	for _, ln := range lns {
		ff.append(ln)
	}
}

// append a line.
func (ff *FlatFile) append(ln Line) {
	ln = ln.Copy()
	for k, v := range ff.lineFmt {
		if s, ok := ln[k]; ok {
			ln[k] = strings.Trim(s[:min(len(s), v.length)], " ")
		} else {
			ln[k] = ""
		}
	}

	ff.lines = append(ff.lines, ln)
}

// appendBytes ...
func (ff *FlatFile) appendBytes(b []byte) {
	ln := make(Line)
	for k, v := range ff.lineFmt {
		ln[k] = string(bytes.Trim(b[v.index:v.index+v.length], " "))
	}

	ff.lines = append(ff.lines, ln)
}

// Fields returns a sorted slice of fields from a flat file.
func (ff *FlatFile) Fields(i int) Fields {
	fs := make(Fields, 0, len(ff.lines[i]))
	for k, v := range ff.lineFmt {
		fs = append(fs, NewField(ff.lines[i][k], v.index, v.length))
	}

	sort.Slice(fs, func(i, j int) bool { return fs[i].Compare(fs[j]) < 0 })
	return fs
}

// Format returns a copy of the flat file line format.
func (ff *FlatFile) Format() LineFmt {
	return ff.lineFmt.Copy()
}

// Get the field associated with a given field name in the ith line of a flat
// file.
func (ff *FlatFile) Get(i int, fieldName string) (string, error) {
	return ff.lines[i].Get(fieldName)
}

// Grow increases the flat file's capacity. If the given capacity is not greater
// than the current length, then nothing happens.
func (ff *FlatFile) Grow(cap int) {
	if len(ff.lines) < cap {
		cpy := make(Lines, len(ff.lines), cap)
		copy(cpy, ff.lines)
		ff.lines = cpy
	}
}

// Len returns the number of lines in the flat file.
func (ff *FlatFile) Len() int {
	return len(ff.lines)
}

// ReadFrom a reader into a flat file.
func (ff *FlatFile) ReadFrom(r io.Reader) (int64, error) {
	if len(ff.lineFmt) == 0 {
		return 0, errors.E(errors.Invalid, "line format not initialized")
	}

	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, err
	}

	lines := bytes.Split(contents, []byte("\n"))
	ff.Grow(len(ff.lines) + len(lines))
	for _, line := range lines {
		if 0 < len(line) {
			ff.appendBytes(line)
		}
	}

	return int64(len(contents)), nil
}

// Set a field. Caution: this overwrites any existing field.
func (ff *FlatFile) Set(i int, fieldName, fieldContents string) error {
	fmt, ok := ff.lineFmt[fieldName]
	if !ok {
		return errors.E(errors.Invalid, "invalid field name")
	}

	if fmt.length < len(fieldContents) {
		return errors.E(errors.Invalid, "format length restriction")
	}

	ff.lines[i].Set(fieldName, fieldContents)
	return nil
}

// String returns flat file lines as strings, concatenated into a single string
// by a carriage return.
func (ff *FlatFile) String() string {
	ss := make([]string, 0, len(ff.lines))
	for i := range ff.lines {
		ss = append(ss, ff.StringAt(i))
	}

	return strings.Join(ss, "\n")
}

// StringAt returns the ith line as a string.
func (ff *FlatFile) StringAt(i int) string {
	s := make([]string, 0, len(ff.lines[i]))
	for _, f := range ff.Fields(i) {
		s = append(s, f.contents+strings.Repeat(" ", f.length-len(f.contents)))
	}

	return strings.Join(s, "")
}

// Swap two lines.
func (ff *FlatFile) Swap(i, j int) {
	ff.lines[i], ff.lines[j] = ff.lines[j], ff.lines[i]
}

// WriteTo a given writer.
func (ff *FlatFile) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write([]byte(ff.String()))
	return int64(n), err
}
