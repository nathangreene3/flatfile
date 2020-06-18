package flatfile

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

// A FlatFile represents a file in which each line contains a fixed number of characters with fields partitioned over the line length. Not all columns need be used. A formatter, determined by the caller, is required to format each line. Usually, flat files have one and only one format for all lines, but here, a formatter allows multiple formats.
//
// * lines: Slice of lines. Each line is formatted using the formats function (Formatter) when read.
//
// * formats: Formatter. This formatter determines how to parse a line from a string.
type FlatFile struct {
	formats Formatter
	lines   []Line
}

const (
	// lf is the line ending for non-Windows machines.
	lf string = "\n"

	// crlf is the line ending for Windows machines.
	crlf string = "\r\n"
)

// New returns a new flat file given a formatter.
func New(f Formatter) *FlatFile {
	return &FlatFile{
		lines:   make([]Line, 0),
		formats: f,
	}
}

// Append several lines to a flat file.
func (ff *FlatFile) Append(lines ...*Line) {
	for i := 0; i < len(lines); i++ {
		ff.lines = append(ff.lines, *lines[i].Copy())
	}
}

// AppendBts appends a line to a flat file.
func (ff *FlatFile) AppendBts(line []byte) error {
	return ff.AppendStr(string(line))
}

// AppendStr appends a line to a flat file. If the line doesn't parse (i.e., the provided formatter returns nil), an error will be returned.
func (ff *FlatFile) AppendStr(line string) error {
	fmts := ff.formats(line)
	if fmts == nil {
		return NewParsingError(line)
	}

	ff.Append(NewLine(line, fmts...))
	return nil
}

// ByteLen returns the entire length of the flat file in bytes. If called on Windows, each line is assumed to end in CRLF and will be larger than on other operating systems.
func (ff *FlatFile) ByteLen() int {
	n := len(ff.lines) // Each line ends in at least LF (\n)
	if runtime.GOOS == "windows" {
		n <<= 1 // Each line ends in CRLF (\r\n)
	}

	for i := 0; i < len(ff.lines); i++ {
		n += ff.lines[i].length
	}

	return n
}

// Bytes returns the flat file as a byte slice. If called on Windows, each line will end in CRLF, including the last line. Otherwise, each line will end in LF.
func (ff *FlatFile) Bytes() []byte {
	lineEnding := lf
	if runtime.GOOS == "windows" {
		lineEnding = crlf
	}

	buf := *bytes.NewBuffer(make([]byte, 0, ff.ByteLen()))
	for i := 0; i < len(ff.lines); i++ {
		buf.WriteString(ff.lines[i].String())
		buf.WriteString(lineEnding)
	}

	return buf.Bytes()
}

// BytesAt returns the ith line as a slice of bytes.
func (ff *FlatFile) BytesAt(i int) []byte {
	return ff.lines[i].Bytes()
}

// Clear removes the contents of a flat file. WARNING: This is irreversible.
func (ff *FlatFile) Clear() {
	ff.lines = make([]Line, 0)
}

// Index returns the index of a key and true if the key is found in the ith line.
func (ff *FlatFile) Index(i int, key string) (int, bool) {
	index, ok := ff.lines[i].keyToIndex[key]
	return index, ok
}

// Field returns the field given a key in the ith line.
func (ff *FlatFile) Field(i int, key string) (Field, error) {
	return ff.lines[i].Field(key)
}

// FieldAt returns the jth field in the ith line.
func (ff *FlatFile) FieldAt(i, j int) Field {
	return ff.lines[i].FieldAt(j)
}

// Formats returns the list of formats given a line. This function parses a line but does not add it to the flat file.
func (ff *FlatFile) Formats(line string) []Format {
	return ff.formats(line)
}

// FormatsAt returns a slice of formats for the ith line.
func (ff *FlatFile) FormatsAt(i int) []Format {
	return ff.lines[i].Formats()
}

// MarshalJSON ...TODO: Make this leaner.
func (ff *FlatFile) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 256)) // TODO: Find an appropriate estimate for capacity
	buf.WriteByte('[')
	for i := 0; i < len(ff.lines); i++ {
		buf.WriteByte('[')
		for j := 0; j < len(ff.lines[i].fields); j++ {
			buf.WriteString(
				"{" +
					"\"key\":\"" + ff.lines[i].fields[j].key + "\"," +
					"\"value\":\"" + ff.lines[i].fields[j].value + "\"," +
					"\"index\":\"" + strconv.Itoa(ff.lines[i].fields[j].index) + "\"," +
					"\"length\":\"" + strconv.Itoa(ff.lines[i].fields[j].length) + "\"" +
					"}",
			)

			if j+1 < len(ff.lines[i].fields) {
				buf.WriteByte(',')
			}
		}

		buf.WriteByte(']')
		if i+1 < len(ff.lines) {
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

// Keys returns the keys of the ith line.
func (ff *FlatFile) Keys(i int) []string {
	keys := make([]string, 0, len(ff.lines[i].fields))
	for j := 0; j < len(ff.lines[i].fields); j++ {
		keys = append(keys, ff.lines[i].fields[j].key)
	}

	return keys
}

// KeyValue returns the jth key-value pair in the ith line.
func (ff *FlatFile) KeyValue(i, j int) (string, string) {
	return ff.lines[i].KeyValueAt(j)
}

// KeyValues returns a map of the keys to their values in the ith line.
func (ff *FlatFile) KeyValues(i int) map[string]string {
	m := make(map[string]string)
	for j := 0; j < len(ff.lines[i].fields); j++ {
		m[ff.lines[i].fields[j].key] = ff.lines[i].fields[j].value
	}

	return m
}

// Len returns the number of lines.
func (ff *FlatFile) Len() int {
	return len(ff.lines)
}

// Line returns the ith line in a flat file.
func (ff *FlatFile) Line(i int) *Line {
	return ff.lines[i].Copy()
}

// Note: io.Read interface will not be supported because the caller MAY not know how big the byte slice should be for the flat file to fill. ByteLen could tell them that, but it would be annoying at best and probably dangerous at the worst.

// ReadFile reads a file into a flat file.  The contents will be appended.
func (ff *FlatFile) ReadFile(filename string) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	_, err = ff.ReadFrom(bytes.NewBuffer(b))
	return err
}

// ReadFrom implements io.ReaderFrom interface. The number of bytes read reflects the number of bytes written to the flat file, not the actual number of bytes in the file. The white space and line endings (LF or CRLF) are removed, but added back in WriteFile.
func (ff *FlatFile) ReadFrom(r io.Reader) (int64, error) {
	var (
		b   = make([]byte, 1<<7)
		buf bytes.Buffer
		n0  = ff.ByteLen()
		n   int
		err error
	)

	for {
		switch n, err = r.Read(b); err {
		case nil:
			buf.Write(b[:n])
		case io.EOF:
			buf.Write(b[:n])
			bts := bytes.Split(buf.Bytes(), []byte{'\n'})
			for i := 0; i < len(bts); i++ {
				bts[i] = bytes.Trim(bts[i], "\r")
				if 0 < len(bts[i]) {
					if _, err := ff.Write(bts[i]); err != nil {
						return int64(ff.ByteLen() - n0), err
					}
				}
			}

			return int64(ff.ByteLen() - n0), nil
		default:
			return int64(ff.ByteLen() - n0), err
		}
	}
}

// Remove the ith line from a flat file.
func (ff *FlatFile) Remove(i int) *Line {
	line := ff.lines[i]
	ff.lines = append(ff.lines[:i], ff.lines[i+1:]...)
	return &line
}

// Set sets the ith line in a flat file.
func (ff *FlatFile) Set(i int, line Line) {
	ff.lines[i] = *line.Copy()
}

// SetStr sets the ith line in a flat file.
func (ff *FlatFile) SetStr(i int, line string) error {
	fmts := ff.formats(line)
	if fmts == nil {
		return NewParsingError(line)
	}

	ff.Set(i, *NewLine(line, fmts...))
	return nil
}

// SetValue sets the value given a key in the ith line.
func (ff *FlatFile) SetValue(i int, key, value string) error {
	return ff.lines[i].Set(key, value)
}

// SetValueAt sets the jth value in the ith line.
func (ff *FlatFile) SetValueAt(i, j int, value string) {
	ff.lines[i].SetAt(j, value)
}

// Sort a flat file's lines given a less-than comparison function.
func (ff *FlatFile) Sort(less func(line0, line1 Line) bool) {
	sort.Slice(ff.lines, func(i, j int) bool { return less(ff.lines[i], ff.lines[j]) })
}

// String returns a string representing a flat file.
func (ff *FlatFile) String() string {
	var sb strings.Builder
	sb.Grow(ff.ByteLen())
	for i := 0; i < len(ff.lines); i++ {
		sb.WriteString(ff.lines[i].String() + "\n")
	}

	return sb.String()
}

// StringAt returns the ith line as a string.
func (ff *FlatFile) StringAt(i int) string {
	return ff.lines[i].String()
}

// Strings returns a slice of strings representing each line in the flat file.
func (ff *FlatFile) Strings() []string {
	ss := make([]string, 0, len(ff.lines))
	for i := 0; i < len(ff.lines); i++ {
		ss = append(ss, ff.StringAt(i))
	}

	return ss
}

// UnmarshalJSON ...TODO
func (ff *FlatFile) UnmarshalJSON(b []byte) error {
	return errors.New("flatfile: FlatFile.UnmarshalJSON not yet implemented")
}

// Value returns the value given a key in the ith line.
func (ff *FlatFile) Value(i int, key string) (string, error) {
	return ff.lines[i].Value(key)
}

// ValueAt returns the jth value in the ith line.
func (ff *FlatFile) ValueAt(i, j int) string {
	return ff.lines[i].ValueAt(j)
}

// Values returns a slice of the values in the ith line.
func (ff *FlatFile) Values(i int) []string {
	values := make([]string, 0, len(ff.lines[i].fields))
	for j := 0; j < len(ff.lines[i].fields); j++ {
		values = append(values, ff.lines[i].fields[j].value)
	}

	return values
}

// Write implements the io.Writer interface. This is equivalent to AppendBts.
func (ff *FlatFile) Write(line []byte) (int, error) {
	if err := ff.AppendBts(line); err != nil {
		return 0, err
	}

	return len(line), nil
}

// WriteFile writes a flat file to file given a file name.
func (ff *FlatFile) WriteFile(fileName string) error {
	return ioutil.WriteFile(fileName, ff.Bytes(), os.ModePerm)
}

// WriteString implements the io.StringWriter interface. It is equivalent to AppendStr.
func (ff *FlatFile) WriteString(line string) (int, error) {
	if err := ff.AppendStr(line); err != nil {
		return 0, err
	}

	return len(line), nil
}
