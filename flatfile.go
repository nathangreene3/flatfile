package flatfile

import (
	"bytes"
	"io"
	"io/ioutil"
	"sort"
	"strings"

	"git.biscorp.local/serverdev/errors"
)

// FlatFile ...
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
	for k, v := range ln {
		ln[k] = strings.Trim(v[:min(len(v), ff.lineFmt[k].length)], " ")
	}

	ff.lines = append(ff.lines, ln)
}

// Fields returns a sorted slice of fields from a flat file.
func (ff *FlatFile) Fields(i int) Fields {
	fs := make(Fields, 0, len(ff.lines[i]))
	for k, v := range ff.lines[i] {
		fs = append(fs, NewField(v, ff.lineFmt[k].index, ff.lineFmt[k].length))
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

// Insert a field name into a flat file. If a field name already exists, an
// error will be returned. To overwrite an existing field name, use Set.
func (ff *FlatFile) Insert(i int, fieldName, field string) error {
	return ff.lines[i].Insert(fieldName, field)
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

	var (
		lines = bytes.Split(contents, []byte("\n"))
		cpy   = make(Lines, len(ff.lines), len(ff.lines)+len(lines))
	)

	copy(cpy, ff.lines)
	ff.lines = cpy

	for _, line := range lines {
		ln := make(Line)
		if 0 < len(line) {
			for fieldName, lf := range ff.lineFmt {
				ln[fieldName] = strings.Trim(string(line[lf.index:lf.index+lf.length]), " ")
			}
		}

		ff.lines = append(ff.lines, ln)
	}

	return int64(len(contents)), nil
}

// Set a field name. Caution: this overwrites any existing field name. To
// prevent overwriting, use Insert. To insert a new line, use Append.
func (ff *FlatFile) Set(i int, fieldName, field string) error {
	if ff.lineFmt[fieldName].length < len(field) {
		return errors.E(errors.Invalid, "format length restriction")
	}

	ff.lines[i].Set(fieldName, field)
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

// WriteTo a given writer.
func (ff *FlatFile) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write([]byte(ff.String()))
	return int64(n), err
}
