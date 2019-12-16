package flatfile

import (
	"bytes"
	"io"
	"io/ioutil"
	"sort"

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

// DetermineFormatLengthsFromIndices TODO
func (ff *FlatFile) DetermineFormatLengthsFromIndices() {}

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

// Read a flat file into a slice of bytes. TODO
func (ff *FlatFile) Read(b []byte) (int64, error) {
	return 0, nil
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
				ln[fieldName] = string(line[lf.index : lf.index+lf.length])
			}
		}

		ff.lines = append(ff.lines, ln)
	}

	return int64(len(contents)), nil
}

// Set a field name. Caution: this overwrites any existing field name. To
// prevent overwriting, use Insert.
func (ff *FlatFile) Set(i int, fieldName, field string) error {
	if ff.lineFmt[fieldName].length < len(field) {
		return errors.E(errors.Invalid, "format length restriction")
	}

	ff.lines[i].Set(fieldName, field)
	return nil
}

// Slice ...TODO
func (ff *FlatFile) Slice(i int) []string {
	type temp struct {
		index int
		field string
	}

	var (
		compareTemps = func(t, u temp) int {
			switch {
			case t.index < u.index:
				return -1
			case u.index < t.index:
				return 1
			case t.field < u.field:
				return -1
			case u.field < t.field:
				return 1
			default:
				return 0
			}
		}

		n     = len(ff.lines[i])
		temps = make([]temp, 0, n)
		s     = make([]string, 0, n)
	)

	for k, v := range ff.lines[i] {
		temps = append(temps, temp{index: ff.lineFmt[k].index, field: v})
	}

	sort.Slice(temps, func(i, j int) bool { return compareTemps(temps[i], temps[j]) < 0 })
	for _, t := range temps {
		s = append(s, t.field)
	}

	return s
}
