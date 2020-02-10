package flatfile

// Line represents a single line in a flat file. Each key-valued pair represents
// a substring of a line where the keys are the field names and the values are
// the contents (fields) of a subset of a line in a flat file.
type Line map[string]string // field name --> field contents

// Lines represents several lines.
type Lines []Line

// Contains indicates if a field name is found in a line.
func (ln Line) Contains(fieldName string) bool {
	_, ok := ln[fieldName]
	return ok
}

// Copy a line.
func (ln Line) Copy() Line {
	cpy := make(Line)
	for k, v := range ln {
		cpy[k] = v
	}

	return cpy
}

// Delete a field name from a line. Returns an error if the field name is not
// found.
func (ln Line) Delete(fieldName string) error {
	if _, ok := ln[fieldName]; ok {
		delete(ln, fieldName)
		return nil
	}

	return errFieldNotExist
}

// Get a field associated with a field name. Returns an error if the field name
// is not found.
func (ln Line) Get(fieldName string) (string, error) {
	if field, ok := ln[fieldName]; ok {
		return field, nil
	}

	return "", errFieldNotExist
}

// Insert a field into a line. Returns an error if the field name already
// exists. To overwrite an existing key, use Set.
func (ln Line) Insert(fieldName, contents string) error {
	if _, ok := ln[fieldName]; ok {
		return errFieldExists
	}

	ln[fieldName] = contents
	return nil
}

// Len returns the number of fields.
func (ln Line) Len() int {
	return len(ln)
}

// Set a field to a given field name. Caution: this overwrites any existing
// field associated with the field name. To prevent overwriting, use Insert.
func (ln Line) Set(fieldName, contents string) {
	ln[fieldName] = contents
}
