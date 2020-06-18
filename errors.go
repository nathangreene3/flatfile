package flatfile

import "fmt"

// MarshalError reports that a flat file failed to produce valid json bytes when marshaled.
type MarshalError struct {
	b []byte
}

// NewMarshalError returns a new marshal error reference.
func NewMarshalError(b []byte) *MarshalError {
	cpy := make([]byte, len(b))
	copy(cpy, b)
	return &MarshalError{b: cpy}
}

// Error implements the Error interface.
func (e *MarshalError) Error() string {
	return fmt.Sprintf("flatfile: MarshalJSON interface implementation produced invalid json %s", string(e.b))
}

// MissingKeyError reports that a key was not found in a list of formats.
type MissingKeyError struct {
	key     string
	formats []Format
}

// NewMissingKeyError returns a new missing-key error reference.
func NewMissingKeyError(key string, formats []Format) *MissingKeyError {
	cpy := make([]Format, len(formats))
	copy(cpy, formats)
	return &MissingKeyError{key: key, formats: cpy}
}

// Error implements the Error interface.
func (e *MissingKeyError) Error() string {
	return fmt.Sprintf("flatfile: key %q not found in line formatted as %v", e.key, e.formats)
}

// ParsingError ...TODO
type ParsingError struct {
	line string
}

// NewParsingError reports that a line could not be parsed with the given formatter.
func NewParsingError(line string) *ParsingError {
	return &ParsingError{line: line}
}

func (e *ParsingError) Error() string {
	return fmt.Sprintf("flatfile: formatter could not parse line '%s'", e.line)
}
