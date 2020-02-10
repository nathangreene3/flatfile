package flatfile

import (
	"bytes"
	"io"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/nathangreene3/table"
)

// FlatFile consists of a slice of lines and a format applied to each line.
type FlatFile struct {
	lineFmt LineFmt
	lines   Lines
}

// New returns a flat file ready to read from a reader.
func New(lf LineFmt) *FlatFile {
	return &FlatFile{lineFmt: lf.Copy(), lines: make(Lines, 0)}
}

// FromTable ...
func FromTable(t *table.Table) *FlatFile {
	// TODO
	return nil
}

// Append appends several lines to a flat file.
func (ff *FlatFile) Append(lns ...Line) *FlatFile {
	ff.grow(len(lns))
	for _, ln := range lns {
		ff.append(ln)
	}

	return ff
}

// append appends a line to a flat file. It is advised the caller grow the flat file first.
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
	return ff
}

// AppendBts ...
func (ff *FlatFile) AppendBts(bts ...[]byte) *FlatFile {
	for _, b := range bts {
		ff.appendBts(b)
	}

	return ff
}

// appendBts formats a byte slice into a line and appends it. It is advised the caller grow the flat file first.
func (ff *FlatFile) appendBts(b []byte) *FlatFile {
	ln := make(Line)
	for k, f := range ff.lineFmt {
		ln[k] = string(bytes.Trim(b[f.index:f.index+f.length], " "))
	}

	ff.lines = append(ff.lines, ln)
	return ff
}

// AppendStrs ...
func (ff *FlatFile) AppendStrs(ss ...string) *FlatFile {
	for _, s := range ss {
		ff.appendBts([]byte(s))
	}

	return ff
}

// Bytes ...
func (ff *FlatFile) Bytes() ([]byte, error) {
	m := len(ff.lines)
	if m < 1 {
		return make([]byte, 0), nil
	}

	n := 1 // Every line except the last line contains at least '\n'
	for _, lf := range ff.lineFmt {
		n += lf.length
	}

	// Builder is more efficient than Buffer, but doesn't provide a Bytes function.
	buf := bytes.NewBuffer(make([]byte, 0, m*n))
	for i := range ff.lines[:m-1] {
		buf.Write(ff.BytesAt(i))
		buf.WriteByte('\n')
	}

	buf.Write(ff.BytesAt(m - 1))
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

// Fields returns a sorted slice of fields from the ith line in a flat file.
func (ff *FlatFile) Fields(i int) Fields {
	fs := make(Fields, 0, len(ff.lines[i]))
	for k, v := range ff.lineFmt {
		fs = append(fs, NewField(ff.lines[i][k], v.index, v.length))
	}

	sort.Slice(fs, func(i, j int) bool { return fs[i].Compare(fs[j]) < 0 })
	return fs
}

// LineFormat returns a copy of the flat file line format.
func (ff *FlatFile) LineFormat() LineFmt {
	return ff.lineFmt.Copy()
}

// Get the ith line from a flat file.
func (ff *FlatFile) Get(i int) Line {
	return ff.lines[i].Copy()
}

// GetField the field associated with a given field name in the ith line of a flat
// file.
func (ff *FlatFile) GetField(i int, fieldName string) (string, error) {
	return ff.lines[i].Get(fieldName)
}

// Grow increases the flat file's capacity. If the given capacity is not greater
// than the current length, then nothing happens.
func (ff *FlatFile) grow(n int) *FlatFile {
	var (
		m = len(ff.lines)
		c = m + n
	)

	if cap(ff.lines) < c {
		cpy := make(Lines, m, c)
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
		return 0, errLineFmtNotInit
	}

	cts, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, err
	}

	lns := bytes.Split(cts, []byte("\n"))
	ff.grow(len(lns))
	for _, ln := range lns {
		if 0 < len(ln) {
			ff.appendBts(ln)
		}
	}

	return int64(len(cts)), nil
}

// Remove and return a line.
func (ff *FlatFile) Remove(i int) Line {
	var ln Line
	if i < len(ff.lines) {
		ln = ff.lines[i]
	}

	ff.lines = append(ff.lines[:i], ff.lines[i+1:]...)
	return ln
}

// Reset ...
func (ff *FlatFile) Reset() *FlatFile {
	ff.lines = make(Lines, 0)
	return ff
}

// Set the ith line in a flat file.
func (ff *FlatFile) Set(i int, ln Line) error {
	oldLn := ff.lines[i].Copy()
	for k, v := range ln {
		if err := ff.SetField(i, k, v); err != nil {
			ff.lines[i] = oldLn // Restores original data
			return err
		}
	}

	return nil
}

// SetField sets a field. Caution: this overwrites any existing field.
func (ff *FlatFile) SetField(i int, fieldName, fieldContents string) error {
	f, ok := ff.lineFmt[fieldName]
	if !ok {
		return errFieldNotExist
	}

	if f.length < len(fieldContents) {
		return errFieldLengthRestriction
	}

	ff.lines[i].Set(fieldName, fieldContents)
	return nil
}

// String returns flat file lines as strings, concatenated into a single string
// by a carriage return.
func (ff *FlatFile) String() string {
	m := len(ff.lines)
	if m < 1 {
		return ""
	}

	n := 1 // Every line except the last line contains at least '\n'
	for _, lf := range ff.lineFmt {
		n += lf.length
	}

	// Builder is more efficient than Buffer, but doesn't provide a Bytes function.
	var sb strings.Builder
	sb.Grow(m * n)
	for i := range ff.lines[:m-1] {
		sb.Write(ff.BytesAt(i))
		sb.WriteByte('\n')
	}

	sb.Write(ff.BytesAt(m - 1))
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
func (ff *FlatFile) Swap(i, j int) *FlatFile {
	ff.lines[i], ff.lines[j] = ff.lines[j], ff.lines[i]
	return ff
}

// Table ...
func (ff *FlatFile) Table() *table.Table {
	// TODO: return a table.
	// Does table implement Reader/Writer?
	return nil
}

// WriteTo a given writer.
func (ff *FlatFile) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write([]byte(ff.String()))
	return int64(n), err
}
