package flatfile

var (
	// errStrKeyNotFound indicates a field key doesn't exist within a line. Requires key and line formats to be inserted.
	errStrKeyNotFound string = "key '%s' not found in line formatted as %v"

	// errStrFmt indicates a string could not be parsed with a given formatter. Requires line to be inserted.
	errStrFmt string = "formatter could not parse line '%s'"
)
