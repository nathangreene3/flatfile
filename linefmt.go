package flatfile

import "git.biscorp.local/serverdev/errors"

// LineFmt maps field names to their formats.
type LineFmt map[string]FieldFmt

// Contains indicates if a field name is found in a line format.
func (lf *LineFmt) Contains(fieldName string) bool {
	_, ok := (*lf)[fieldName]
	return ok
}

// Copy a line format.
func (lf *LineFmt) Copy() LineFmt {
	cpy := make(map[string]FieldFmt)
	for k, v := range *lf {
		cpy[k] = v
	}

	return cpy
}

// Delete a field name from a line format.
func (lf *LineFmt) Delete(fieldName string) error {
	if _, ok := (*lf)[fieldName]; ok {
		delete(*lf, fieldName)
		return nil
	}

	return errors.E(errors.NotExist, "field name not found")
}

// Get a field format associated by a field name.
func (lf *LineFmt) Get(fieldName string) (FieldFmt, error) {
	if fieldFmt, ok := (*lf)[fieldName]; ok {
		return fieldFmt, nil
	}

	return FieldFmt{}, errors.E(errors.NotExist, "field name not found")
}

// Insert a field format into a line format. Returns an error if the field name
// already exists. To overwrite an existing field format associated with the
// field name, use Set.
func (lf *LineFmt) Insert(fieldName string, fieldFmt FieldFmt) error {
	if _, ok := (*lf)[fieldName]; ok {
		return errors.E(errors.Exist, "field name already exists")
	}

	(*lf)[fieldName] = fieldFmt
	return nil
}

// Len returns the number of field names in a line format.
func (lf *LineFmt) Len() int {
	return len(*lf)
}

// Set a field format to a given field name. Caution: this overwrites any
// existing field associated with the field name. To prevent overwriting, use
// Insert.
func (lf *LineFmt) Set(fieldName string, fieldFmt FieldFmt) {
	(*lf)[fieldName] = fieldFmt
}
