package flatfile

import (
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/tidwall/gjson"
)

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
	jsonType      JSONType
}

// JSONType dictates how a field is formatted.
type JSONType int

const (
	// String is the default json type.
	String JSONType = iota

	// Number (integer or float) data type.
	Number

	// Boolean data type.
	Boolean
)

// NewFormat returns a new format.
func NewFormat(key string, index, length int, jsonType JSONType) Format {
	return Format{
		key:      key,
		index:    index,
		length:   length,
		jsonType: jsonType,
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

// MarshalJSON ...
func (fmt *Format) MarshalJSON() ([]byte, error) {
	b := []byte(
		"{" +
			"\"key\":\"" + fmt.key + "\"," +
			"\"index\":" + strconv.Itoa(fmt.index) + "," +
			"\"length\":" + strconv.Itoa(fmt.length) + "," +
			"\"jsonType\":" + strconv.Itoa(int(fmt.jsonType)) +
			"}",
	)

	if !json.Valid(b) {
		return nil, NewMarshalError(b)
	}

	return b, nil
}

// UnmarshalJSON ...
func (fmt *Format) UnmarshalJSON(b []byte) error {
	index := gjson.GetBytes(b, "index").Num
	if float64(int(index)) != index {
		return &json.UnmarshalTypeError{
			Value:  "integer",
			Type:   reflect.TypeOf(index),
			Offset: 0, // TODO
			Struct: "Format",
			Field:  "index",
		}
	}

	length := gjson.GetBytes(b, "length").Num
	if float64(int(length)) != length {
		return &json.UnmarshalTypeError{
			Value:  "integer",
			Type:   reflect.TypeOf(length),
			Offset: 0, // TODO
			Struct: "Format",
			Field:  "length",
		}
	}

	fmt.key = gjson.GetBytes(b, "key").Str
	fmt.index = int(index)
	fmt.length = int(length)
	return nil
}
