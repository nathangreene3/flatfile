package flatfile

// Format defines a field within a line.
type Format struct {
	index, length int
}

// NewFormat returns a new field format. The index specifies the index a field
// begins and the length specifies how many characters long it is in a line.
func NewFormat(index, length int) Format {
	return Format{index: index, length: length}
}

// Compare two field formats.
func (f *Format) Compare(format Format) int {
	switch {
	case f.index < format.index:
		return -1
	case format.index < f.index:
		return 1
	case f.length < format.length:
		return -1
	case format.length < f.length:
		return 1
	default:
		return 0
	}
}
