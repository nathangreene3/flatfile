package flatfile

import "bytes"

// Field extends a format.
//
// * Format: Embedded within the field.
//
// * value: Contents of the field, typically cleaned of leading and trailing whitespace.
type Field struct {
	Format
	value string
}

// NewField returns a new field that references a keyed value in a line at an index
// having a maximum allowed length.
func NewField(key, value string, index, length int) Field {
	if length < len(value) {
		value = value[:length]
	}

	return Field{
		Format: Format{
			key:    key,
			index:  index,
			length: length,
		},
		value: value,
	}
}

// Bytes returns a slice of bytes representing a field.
func (fld *Field) Bytes() []byte {
	return append(append(make([]byte, 0, fld.length), []byte(fld.value)...), bytes.Repeat([]byte{' '}, fld.length-len(fld.value))...)
}

// String returns a string representing a field.
func (fld *Field) String() string {
	// TODO: Determine what's more efficient, concatenating strings or converting bytes to string.
	return string(fld.Bytes()) // fld.value + strings.Repeat(" ", fld.length-len(fld.value))
}
