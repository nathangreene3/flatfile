package flatfile

// Format consists of field data used to import from flat files.
type Format struct {
	index, length int
}

// NewFormat returns a new field format. The index specifies the index a field
// begins and the length specifies how many characters long it is in a line.
func NewFormat(index, length int) Format {
	return Format{index: index, length: length}
}

// Compare two field formats.
func (f *Format) Compare(fieldFmt Format) int {
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
