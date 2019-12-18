# flatfile

A package to read from and write to flat files.

## FlatFile

```go
type FlatFile struct {
    lineFmt LineFmt
    lines   Lines
}
```

A flat file consists of data (Lines) and a format applied to each line.

### New

```go
func New(lf LineFmt) FlatFile
```

Returns a flat file ready to read from a reader.

### FlatFile.Append

```go
func (ff *FlatFile) Append(lns ...Line)
```

Append several lines.

### FlatFile.Fields

```go
func (ff *FlatFile) Fields(i int) Fields
```

Returns a sorted slice of fields from a flat file.

### FlatFile.Format

```go
func (ff *FlatFile) Format() LineFmt
```

Returns a copy of the flat file line format.

### FlatFile.Get

```go
func (ff *FlatFile) Get(i int, fieldName string) (string, error)
```

Returns field associated with a given field name in the ith line of a flat file.

### FlatFile.Grow

```go
func (ff *FlatFile) Grow(cap int)
```

Increases the flat file's capacity. If the given capacity is not greater than the current length, then nothing happens.

### FlatFile.Len

```go
func (ff *FlatFile) Len() int
```

Returns the number of lines in the flat file.

### FlatFile.ReadFrom

```go
func (ff *FlatFile) ReadFrom(r io.Reader) (int64, error)
```

Read from a reader into a flat file.

### FlatFile.Set

```go
func (ff *FlatFile) Set(i int, fieldName, fieldContents string) error
```

Set a field. Caution: this overwrites any existing field.

### FlatFile.String

```go
func (ff *FlatFile) String() string
```

Returns flat file lines as strings, concatenated into a single string by a carriage return.

### FlatFile.StringAt

```go
func (ff *FlatFile) StringAt(i int) string
```

Returns the ith line as a string.

### FlatFile.Swap

```go
func (ff *FlatFile) Swap(i, j int)
```

Swap two lines.

### FlatFile.WriteTo

```go
func (ff *FlatFile) WriteTo(w io.Writer) (int64, error)
```

Write to a given writer.

## Field

```go
type Field struct {
    Format
    contents string
}
```

Extends format by adding contents.

### NewField

```go
func NewField(contents string, index, length int) Field
```

Returns a new field.

### Field.Compare

```go
func (f *Field) Compare(fld Field) int
```

Compare two fields.

## Format

```go
type Format struct {
    index, length int
}
```

Consists of field data used to import from flat files.

### NewFormat

```go
func NewFormat(index, length int) Format
```

Returns a new field format. The index specifies the index a field begins and the length specifies how many characters long it is in a line.

### Format.Compare

```go
func (f *Format) Compare(format Format) int
```

Compare two field formats.

## Line

```go
type Line map[string]string
type Lines []Line
```

Line epresents a single line in a flat file. Each key-valued pair represents a substring of a line where the keys are the field names and the values are the contents (fields) of a subset of a line in a flat file.

### Line.Contains

```go
func (ln *Line) Contains(fieldName string) bool
```

Indicates if a field name is found in a line.

### Line.Copy

```go
func (ln *Line) Copy() Line
```

Copy a line.

### Line.Delete

```go
func (ln *Line) Delete(fieldName string) error
```

Delete a field name from a line. Returns an error if the field name is not found.

### Line.Get

```go
func (ln *Line) Get(fieldName string) (string, error)
```

Get a field associated with a field name. Returns an error if the field name is not found.

### Line.Insert

```go
func (ln *Line) Insert(fieldName, field string) error
```

Insert a field into a line. Returns an error if the field name already exists. To overwrite an existing key, use Set.

### Line.Len

```go
func (ln *Line) Len() int
```

Len returns the number of fields.

### Line.Set

```go
func (ln *Line) Set(fieldName, field string)
```

Set a field to a given field name. Caution: this overwrites any existing field associated with the field name. To prevent overwriting, use Insert.

## LineFmt

```go
type LineFmt map[string]Format
```

LineFmt maps field names to their formats.

### LineFmt.Contains

```go
func (lf *LineFmt) Contains(fieldName string) bool
```

Contains indicates if a field name is found in a line format.

### LineFmt.Copy

```go
func (lf *LineFmt) Copy() LineFmt
```

Copy a line format.

### LineFmt.Delete

```go
func (lf *LineFmt) Delete(fieldName string) error
```

Delete a field name from a line format.

### LineFmt.Get

```go
func (lf *LineFmt) Get(fieldName string) (Format, error)
```

Get a field format associated by a field name.

### LineFmt.Insert

```go
func (lf *LineFmt) Insert(fieldName string, fieldFmt Format) error
```

Insert a field format into a line format. Returns an error if the field name already exists. To overwrite an existing field format associated with the field name, use Set.

### LineFmt.Len

```go
func (lf *LineFmt) Len() int
```

Len returns the number of field names in a line format.

### LineFmt.Set

```go
func (lf *LineFmt) Set(fieldName string, fieldFmt Format)
```

Set a field format to a given field name. Caution: this overwrites any existing field associated with the field name. To prevent overwriting, use Insert.
