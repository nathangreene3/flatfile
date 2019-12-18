package flatfile

import (
	"bytes"
	"strings"
	"testing"
)

// TestFlatFile simulates reading from a flat file where each line consists of
// two fields: last and first names, each being eight characters long. This test
// should cover every exported flat file function.
func TestFlatFile(t *testing.T) {
	// 1. Create a new flat file
	var (
		input = strings.Join(
			[]string{
				"SkywalkeLuke    ",
				"Vader   Darth   ",
				"Kenobi  Obi-Wan ",
				"Solo    Han     ",
			},
			"\n",
		)

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

		ff  = New(LineFmt{"first": NewFormat(8, 8), "last": NewFormat(0, 8)})
		buf = bytes.NewBufferString(input)
	)

	// 2. Read from a buffer (reader)
	if _, err := ff.ReadFrom(buf); err != nil {
		t.Fatalf("\nunexpected error: '%s'\n", err.Error())
	}

	// 3. Append a new line
	ff.Append(Line{"first": "princess", "last": "Leia"})

	// 4. Correct a typo in the new line
	if err := ff.Set(ff.Len()-1, "first", "Princess"); err != nil {
		t.Fatalf("\nunexpected error: '%s'\n", err.Error())
	}

	// 5. Swap the inserted line
	ff.Swap(ff.Len()-1, ff.Len()-2)

	// 6. Write the flat file to the buffer
	if _, err := ff.WriteTo(buf); err != nil {
		t.Fatalf("\nunexpected error: '%s'\n", err.Error())
	}

	if rec := buf.String(); exp != rec {
		t.Fatalf("\nexpected '%s'\nreceived '%s'\n", exp, rec)
	}

	// 7. Get a valid field from a line
	exp = "Luke"
	rec, err := ff.Get(0, "first")
	if err != nil {
		t.Fatalf("\nunexpected '%s'\n", err.Error())
	}

	if exp != rec {
		t.Fatalf("\nexpected '%s'\nreceived '%s'\n", exp, rec)
	}

	// 8. Attempt to get an invalid field from a line
	exp = ""
	rec, err = ff.Get(0, "middle")
	if err == nil {
		t.Fatalf("\nexpected error\n received '%v'\n", err)
	}

	if exp != rec {
		t.Fatalf("\nexpected '%s'\nreceived '%s'\n", exp, rec)
	}
}
