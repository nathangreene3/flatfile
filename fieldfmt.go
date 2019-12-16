package flatfile

// FieldFmt consists of field data used to import from flat files.
type FieldFmt struct {
	index, length int
}

// NewFieldFmt returns a new field format. The index specifies the index a field
// begins and the length specifies how many characters long it is in a line.
func NewFieldFmt(index, length int) FieldFmt {
	return FieldFmt{index: index, length: length}
}

// Compare two field formats.
func (f *FieldFmt) Compare(fieldFmt FieldFmt) int {
	switch {
	case f.index < fieldFmt.index:
		return -1
	case fieldFmt.index < f.index:
		return 1
	case f.length < fieldFmt.length:
		return -1
	case fieldFmt.length < f.length:
		return 1
	default:
		return 0
	}
}
