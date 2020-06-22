package flatfile

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"

	"github.com/tidwall/gjson"
)

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
func NewField(key, value string, index, length int, jsonType JSONType) Field {
	if length < len(value) {
		value = value[:length]
	}

	return Field{
		Format: Format{
			key:      key,
			index:    index,
			length:   length,
			jsonType: jsonType,
		},
		value: value,
	}
}

// Bytes returns a slice of bytes representing a field.
func (fld *Field) Bytes() []byte {
	return append(append(make([]byte, 0, fld.length), []byte(fld.value)...), bytes.Repeat([]byte{' '}, fld.length-len(fld.value))...)
}

// MarshalJSON ...
func (fld *Field) MarshalJSON() ([]byte, error) {
	var b []byte
	switch fld.jsonType {
	case String:
		b = []byte(
			"{" +
				"\"key\":\"" + fld.key + "\"," +
				"\"value\":\"" + fld.value + "\"," +
				"\"index\":" + strconv.Itoa(fld.index) + "," +
				"\"length\":" + strconv.Itoa(fld.length) + "," +
				"\"jsonType\":" + strconv.Itoa(int(fld.jsonType)) +
				"}",
		)
	case Number, Boolean:
		b = []byte(
			"{" +
				"\"key\":\"" + fld.key + "\"," +
				"\"value\":" + fld.value + "," +
				"\"index\":" + strconv.Itoa(fld.index) + "," +
				"\"length\":" + strconv.Itoa(fld.length) + "," +
				"\"jsonType\":" + strconv.Itoa(int(fld.jsonType)) +
				"}",
		)
	default:
		return nil, errors.New("Undefined json type") // TODO: generate an appropriate custom error
	}

	if !json.Valid(b) {
		return nil, NewMarshalError(b)
	}

	return b, nil
}

// String returns a string representing a field.
func (fld *Field) String() string {
	// TODO: Determine what's more efficient, concatenating strings or converting bytes to string.
	return string(fld.Bytes()) // fld.value + strings.Repeat(" ", fld.length-len(fld.value))
}

// UnmarshalJSON decodes json-encoded bytes.
func (fld *Field) UnmarshalJSON(b []byte) error {
	index := gjson.GetBytes(b, "index").Num
	if float64(int(index)) != index {
		return &json.UnmarshalTypeError{
			Value:  "integer",
			Type:   reflect.TypeOf(index),
			Offset: 0, // TODO
			Struct: "Field",
			Field:  "index",
		}
	}

	length := gjson.GetBytes(b, "length").Num
	if float64(int(length)) != length {
		return &json.UnmarshalTypeError{
			Value:  "integer",
			Type:   reflect.TypeOf(length),
			Offset: 0, // TODO
			Struct: "Field",
			Field:  "length",
		}
	}

	jsonType := gjson.GetBytes(b, "jsonType").Num
	if float64(int(jsonType)) != jsonType {
		return &json.UnmarshalTypeError{
			Value:  "integer",
			Type:   reflect.TypeOf(jsonType),
			Offset: 0, // TODO
			Struct: "Field",
			Field:  "jsonType",
		}
	}

	fld.key = gjson.GetBytes(b, "key").Str
	fld.value = gjson.GetBytes(b, "value").Str
	fld.index = int(index)
	fld.length = int(length)
	fld.jsonType = JSONType(jsonType)
	return nil
}
