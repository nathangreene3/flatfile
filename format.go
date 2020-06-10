package flatfile

// Formatter returns the formats that will be used in parsing a given line. If a line doesn't parse, it should return nil.
type Formatter func(line string) []Format

// Format contains information related to a value in a line.
//
// * key:    A label for looking up the value. This value should be unique within the line.
//
// * index:  Indicates where the value begins in a line. This value is recommended, but not required to be unique within the line to prevent overlapping fields.
//
// * length: The maximum number of characters the value can be when written to a line. The value may be shorter than the format length. When written to a line, the remaining space will be filled in with white space (' ').
type Format struct {
	key           string
	index, length int
}

// NewFormat returns a new format.
func NewFormat(key string, index, length int) Format {
	return Format{
		key:    key,
		index:  index,
		length: length,
	}
}

// Index returns the index a value begins at within a line.
func (fmt *Format) Index() int {
	return fmt.index
}

// Key returns the key describing a value within a line.
func (fmt *Format) Key() string {
	return fmt.key
}

// Length returns the maximum number of characters the value can be within a line.
func (fmt *Format) Length() int {
	return fmt.length
}
