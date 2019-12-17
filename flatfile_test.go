package flatfile

import (
	"bytes"
	"strings"
	"testing"
)

func TestFlatFile1(t *testing.T) {
	// Simulate reading from a flat file where each line consists of two fields:
	// last and first names, each being eight characters long.

	var (
		exp = strings.Join(
			[]string{
				"SkywalkeLuke    ",
				"Vader   Darth   ",
				"Kenobi  Obi-Wan ",
				"Leia    Princess",
				"Solo    Han     ",
			},
			"\n",
		)

		ff = New(LineFmt{"first": NewFormat(8, 8), "last": NewFormat(0, 8)})
	)

	if _, err := ff.ReadFrom(bytes.NewBufferString(exp)); err != nil {
		t.Fatalf("\nunexpected error: '%s'\n", err.Error())
	}

	buf := bytes.NewBuffer(make([]byte, 0, len(exp)))
	if _, err := ff.WriteTo(buf); err != nil {
		t.Fatalf("\nunexpected error: '%s'\n", err.Error())
	}

	if rec := buf.String(); exp != rec {
		t.Fatalf("\nexpected '%s'\nreceived '%s'\n", exp, rec)
	}
}
