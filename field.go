package flatfile

// Field extends FieldFmt by adding contents.
type Field struct {
	Format
	contents string
}

// Fields is a slice of fields.
type Fields []Field

// NewField returns a new field.
func NewField(contents string, index, length int) Field {
	return Field{contents: contents, Format: NewFormat(index, length)}
}

// Compare two fields.
func (f *Field) Compare(fld Field) int {
	r := f.Format.Compare(fld.Format)
	switch {
	case r != 0:
		return r
	case f.contents < fld.contents:
		return -1
	case fld.contents < f.contents:
		return 1
	default:
		return 0
	}
}
