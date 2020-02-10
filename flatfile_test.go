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
				"        Yoda    ",
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

	// 3. Remove a line
	lnExp := Line{"first": "Yoda", "last": ""}
	lnRec := ff.Remove(ff.Len() - 1)
	if len(lnExp) != len(lnRec) {
		t.Fatalf("\nexpected '%v'\nreceived '%v'\n", lnExp, lnRec)
	}

	for k, exp := range lnExp {
		if rec := lnRec[k]; exp != rec {
			t.Fatalf("\nexpected '%v'\nreceived '%v'\n", lnExp, lnRec)
		}
	}

	// 4. Append a new line
	// ff.AppendLines(Line{"first": "princess", "last": "Leia"})
	ff.AppendBts([]byte("Leia    princess"))

	// 5. Correct a typo in the new line
	if err := ff.SetField(ff.Len()-1, "first", "Princess"); err != nil {
		t.Fatalf("\nunexpected error: '%s'\n", err.Error())
	}

	// 6. Swap the inserted line
	ff.Swap(ff.Len()-1, ff.Len()-2)

	// 7. Write the flat file to the buffer
	if _, err := ff.WriteTo(buf); err != nil {
		t.Fatalf("\nunexpected error: '%s'\n", err.Error())
	}

	if rec := buf.String(); exp != rec {
		t.Fatalf("\nexpected '%s'\nreceived '%s'\n", exp, rec)
	}

	// 8. Get a valid field from a line
	exp = "Luke"
	rec, err := ff.GetField(0, "first")
	if err != nil {
		t.Fatalf("\nunexpected '%s'\n", err.Error())
	}

	if exp != rec {
		t.Fatalf("\nexpected '%s'\nreceived '%s'\n", exp, rec)
	}

	// 9. Attempt to get an invalid field from a line
	exp = ""
	rec, err = ff.GetField(0, "middle")
	if err == nil {
		t.Fatalf("\nexpected error\n received '%v'\n", err)
	}

	if exp != rec {
		t.Fatalf("\nexpected '%s'\nreceived '%s'\n", exp, rec)
	}
}
