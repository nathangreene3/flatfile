package flatfile

import "errors"

var (
	errFieldExists            = errors.New("field name already exists")
	errFieldNotExist          = errors.New("field name not found")
	errFieldLengthRestriction = errors.New("format length restriction")
	errLineFmtNotInit         = errors.New("line format not initialized")
)
