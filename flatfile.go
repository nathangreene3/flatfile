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
	length  int
}

// New returns a flat file ready to read from a reader.
func New(lf LineFmt, lns ...Line) *FlatFile {
	ff := &FlatFile{lineFmt: lf.Copy(), lines: make(Lines, 0)}
	return ff.Append(lns...)
}

// Append several lines.
func (ff *FlatFile) Append(lns ...Line) *FlatFile {
	for _, ln := range lns {
		ff.append(ln)
	}

	return ff
}

// append a line.
func (ff *FlatFile) append(ln Line) *FlatFile {
	ln = ln.Copy()
	for k, v := range ff.lineFmt {
		if s, ok := ln[k]; ok {
			ln[k] = strings.Trim(s[:min(len(s), v.length)], " ")
		} else {
			ln[k] = ""
		}
	}

	ff.lines = append(ff.lines, ln)
	ff.length++
	return ff
}

// appendBytes formats a byte slice into a line and appends it.
func (ff *FlatFile) appendBytes(b []byte) *FlatFile {
	ln := make(Line)
	for k, v := range ff.lineFmt {
		ln[k] = string(bytes.Trim(b[v.index:v.index+v.length], " "))
	}

	ff.lines = append(ff.lines, ln)
	ff.length++
	return ff
}

// Fields returns a sorted slice of fields from the ith line in a flat file.
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
	// TODO: Rename to LineFmt?
	return ff.lineFmt.Copy()
}

// Get the field associated with a given field name in the ith line of a flat
// file.
func (ff *FlatFile) Get(i int, fieldName string) (string, error) {
	return ff.lines[i].Get(fieldName)
}

// Grow increases the flat file's capacity. If the given capacity is not greater
// than the current length, then nothing happens.
func (ff *FlatFile) Grow(cap int) *FlatFile {
	if len(ff.lines) < cap {
		cpy := make(Lines, len(ff.lines), cap)
		copy(cpy, ff.lines)
		ff.lines = cpy
	}

	return ff
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

// Remove and return a line.
func (ff *FlatFile) Remove(i int) Line {
	var ln Line
	if i < len(ff.lines) {
		ln = ff.lines[i]
	}

	ff.lines = append(ff.lines[:i], ff.lines[i+1:]...)
	ff.length--
	return ln
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

// Bytes ...
func (ff *FlatFile) Bytes() ([]byte, error) {
	if ff.length < 1 {
		return make([]byte, 0), nil
	}

	lineLen := 1 // Every line except the last line contains at least '\n'
	for _, lf := range ff.lineFmt {
		lineLen += lf.length
	}

	// Builder is more efficient than Buffer, but doesn't provide a Bytes function.
	buf := bytes.NewBuffer(make([]byte, 0, lineLen*ff.length))
	for i := range ff.lines[:ff.length-1] {
		buf.Write(ff.BytesAt(i))
		buf.WriteByte('\n')
	}

	buf.Write(ff.BytesAt(ff.length - 1))
	return buf.Bytes(), nil
}

// BytesAt ...
func (ff *FlatFile) BytesAt(i int) []byte {
	var n int
	for _, lf := range ff.lineFmt {
		n += lf.length
	}

	buf := bytes.NewBuffer(make([]byte, 0, n))
	for _, f := range ff.Fields(i) {
		buf.WriteString(f.contents + strings.Repeat(" ", f.length-len(f.contents)))
	}

	return buf.Bytes()
}

// String returns flat file lines as strings, concatenated into a single string
// by a carriage return.
func (ff *FlatFile) String() string {
	if ff.length < 1 {
		return ""
	}

	lineLen := 1 // Every line except the last line contains at least '\n'
	for _, lf := range ff.lineFmt {
		lineLen += lf.length
	}

	// Builder is more efficient than Buffer, but doesn't provide a Bytes function.
	var sb strings.Builder
	sb.Grow(ff.length * lineLen)
	for i := range ff.lines[:ff.length-1] {
		sb.Write(ff.BytesAt(i))
		sb.WriteByte('\n')
	}

	sb.Write(ff.BytesAt(ff.length - 1))
	return sb.String()
}

// StringAt returns the ith line as a string.
func (ff *FlatFile) StringAt(i int) string {
	var lineLen int
	for _, f := range ff.Fields(i) {
		lineLen += f.length
	}

	var sb strings.Builder
	sb.Grow(lineLen)
	for _, f := range ff.Fields(i) {
		sb.WriteString(f.contents + strings.Repeat(" ", f.length-len(f.contents)))
	}

	return sb.String()
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
