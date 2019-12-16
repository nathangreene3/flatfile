package flatfile

import (
	"bytes"
	"testing"
)

func TestFlatFile1(t *testing.T) {
	var (
		contents = "field1field 2field  3\nfield1field 2field  3\nfield1field 2field  3\n"
		buf      = bytes.NewBuffer(make([]byte, 0, len(contents)))
		format   = LineFmt{
			"1": NewFieldFmt(0, 6),
			"2": NewFieldFmt(6, 7),
			"3": NewFieldFmt(13, 8),
		}

		ff = New(format)
	)

	buf.WriteString(contents)
	ff.ReadFrom(buf)

	t.Fatalf("\nreceived %v\n", ff.lines)
}
