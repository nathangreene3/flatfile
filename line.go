package flatfile

import (
	"bytes"
	"encoding/json"
	"errors"
	"sort"
	"strconv"
	"strings"
)

// A Line is a flat line within a flat file.
//
// * fields:     List of formatted values.
//
// * keyToIndex: Maps the keys to the index in fields (not the index within the line).
//
// * length:     Number of characters the line contains.
type Line struct {
	fields     []Field
	keyToIndex map[string]int
	length     int
}

// NewLine returns a new line given a list of formats.
func NewLine(line string, fmts ...Format) *Line {
	ln := Line{
		fields:     make([]Field, 0, len(fmts)),
		keyToIndex: make(map[string]int),
		length:     len(line), // When printed, the original line length will be preserved, but the contents may be lost if the provided formats don't adequately cover the line
	}

	for i := 0; i < len(fmts); i++ {
		ln.fields = append(ln.fields, Field{Format: fmts[i], value: strings.Trim(line[fmts[i].index:fmts[i].index+fmts[i].length], " ")})
	}

	sort.Slice(ln.fields, func(i, j int) bool { return ln.fields[i].index < ln.fields[j].index })
	for i := 0; i < len(ln.fields); i++ {
		ln.keyToIndex[ln.fields[i].key] = i
	}

	return &ln
}

// Bytes returns a byte slice representing a line.
func (ln *Line) Bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, ln.length))
	for i := 0; i < len(ln.fields); i++ {
		buf.Write(bytes.Repeat([]byte{' '}, ln.fields[i].index-buf.Len()))
		buf.Write(ln.fields[i].Bytes())
	}

	buf.Write(bytes.Repeat([]byte{' '}, ln.length-buf.Len()))
	return buf.Bytes()
}

// Copy a line.
func (ln *Line) Copy() *Line {
	cpy := Line{
		fields:     append(make([]Field, 0, len(ln.fields)), ln.fields...),
		keyToIndex: make(map[string]int),
		length:     ln.length,
	}

	for k, v := range ln.keyToIndex {
		cpy.keyToIndex[k] = v
	}

	return &cpy
}

// Field returns the field given a key.
func (ln *Line) Field(key string) (Field, error) {
	if i, ok := ln.keyToIndex[key]; ok {
		return ln.fields[i], nil
	}

	return Field{}, NewMissingKeyError(key, ln.Formats())
}

// FieldAt returns the ith field in a line.
func (ln *Line) FieldAt(i int) Field {
	return ln.fields[i]
}

// Formats returns a slice of formats in a line.
func (ln *Line) Formats() []Format {
	fmts := make([]Format, 0, len(ln.fields))
	for i := 0; i < len(ln.fields); i++ {
		fmts = append(fmts, ln.fields[i].Format)
	}

	return fmts
}

// MarshalJSON ...
func (ln *Line) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 256))
	buf.WriteByte('[')
	for i := 0; i < len(ln.fields); i++ {
		buf.WriteString(
			"{" +
				"\"key\":\"" + ln.fields[i].key + "\"," +
				"\"value\":\"" + ln.fields[i].value + "\"," +
				"\"index\":\"" + strconv.Itoa(ln.fields[i].index) + "\"," +
				"\"length\":\"" + strconv.Itoa(ln.fields[i].length) + "\"" +
				"}",
		)

		if i+1 < len(ln.fields) {
			buf.WriteByte(',')
		}
	}

	buf.WriteByte(']')

	b := buf.Bytes()
	if !json.Valid(b) {
		return nil, NewMarshalError(b)
	}

	return b, nil
}

// Value a value in a line given a key.
func (ln *Line) Value(key string) (string, error) {
	if i, ok := ln.keyToIndex[key]; ok {
		return ln.fields[i].value, nil
	}

	return "", NewMissingKeyError(key, ln.Formats())
}

// Index returns the index a field begins at in a line given a key.
func (ln *Line) Index(key string) (int, error) {
	if i, ok := ln.keyToIndex[key]; ok {
		return ln.fields[i].index, nil
	}

	return 0, NewMissingKeyError(key, ln.Formats())
}

// IndexAt returns the ith index in a line.
func (ln *Line) IndexAt(i int) int {
	return ln.fields[i].index
}

// Key returns the ith field key.
func (ln *Line) Key(i int) string {
	return ln.fields[i].key
}

// KeyValueAt returns the ith key-value pair.
func (ln *Line) KeyValueAt(i int) (string, string) {
	return ln.fields[i].key, ln.fields[i].value
}

// Length returns the maximum number of characters in a field given a key.
func (ln *Line) Length(key string) (int, error) {
	if i, ok := ln.keyToIndex[key]; ok {
		return ln.fields[i].length, nil
	}

	return 0, NewMissingKeyError(key, ln.Formats())
}

// LengthAt returns the maximum number of characters in the ith field.
func (ln *Line) LengthAt(i int) int {
	return ln.fields[i].length
}

// Set a value in a line given a key.
func (ln *Line) Set(key, value string) error {
	if i, ok := ln.keyToIndex[key]; ok {
		ln.SetAt(i, value)
		return nil
	}

	return NewMissingKeyError(key, ln.Formats())
}

// SetAt sets the ith value in a line.
func (ln *Line) SetAt(i int, value string) {
	if n := ln.fields[i].length; n < len(value) {
		value = value[:ln.fields[i].length]
	}

	ln.fields[i].value = value
}

// String returns a string representing a line.
func (ln *Line) String() string {
	var sb strings.Builder
	sb.Grow(ln.length)

	for i := 0; i < len(ln.fields); i++ {
		sb.WriteString(strings.Repeat(" ", ln.fields[i].index-sb.Len()) + ln.fields[i].String())
	}

	sb.WriteString(strings.Repeat(" ", ln.length-sb.Len()))
	return sb.String()
}

// UnmarshalJSON ...TODO
func (ln *Line) UnmarshalJSON(b []byte) error {
	return errors.New("flatfile: Line.UnmarshalJSON not yet implemented")
}

// ValueAt returns the ith value in a line.
func (ln *Line) ValueAt(i int) string {
	return ln.fields[i].value
}
