package flatfile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestFlatFile(t *testing.T) {
	// equalFiles determines if files contain the same bytes ignoring LF and CRLF endings.
	equalFiles := func(fileName0, fileName1 string) (bool, error) {
		b0, err := ioutil.ReadFile(fileName0)
		if err != nil {
			return false, err
		}

		b1, err := ioutil.ReadFile(fileName1)
		if err != nil {
			return false, err
		}

		bts0, bts1 := bytes.Split(b0, []byte{'\n'}), bytes.Split(b1, []byte{'\n'})
		if len(bts0) != len(bts1) {
			return false, nil
		}

		for i := 0; i < len(bts0); i++ {
			bts0[i], bts1[i] = bytes.Trim(bts0[i], "\r"), bytes.Trim(bts1[i], "\r")
			if !bytes.Equal(bts0[i], bts1[i]) {
				return false, nil
			}
		}

		return true, nil
	}

	// fmts is the formatter given to the flat filer.
	fmts := func(line string) []Format {
		switch len(line) {
		case 8: // Single name
			return []Format{
				NewFormat("name", 0, 8),
			}
		case 16: // Two names
			return []Format{
				NewFormat("first", 0, 8),
				NewFormat("last", 8, 8),
			}
		case 24: // Three names
			return []Format{
				NewFormat("title", 0, 8),
				NewFormat("first", 8, 8),
				NewFormat("last", 16, 8),
			}
		default:
			return nil
		}
	}

	ff := New(fmts)

	{ // Append, set, set value, and removal
		// One name
		// Improper format
		name := "Yoda"
		if err := ff.AppendStr(name); err == nil {
			t.Fatalf("\nexpected error: %q\n", NewParsingError(name).Error())
		}

		// Proper format
		if err := ff.AppendStr("Yoda    "); err != nil {
			t.Fatalf("\nunexpected error: %q\n", err.Error())
		}

		// Two names
		// Improper format
		name = "Luke    Skywalker"
		if err := ff.AppendStr(name); err == nil {
			t.Fatalf("\nexpected error: %q\n", NewParsingError(name).Error())
		}

		// Proper format
		if err := ff.AppendStr("Luke    Skywalke"); err != nil {
			t.Fatalf("\nunexpected error: %q\n", err.Error())
		}

		// Three names
		// Improper format
		name = "PrincessLeiaOrgana"
		if err := ff.AppendStr(name); err == nil {
			t.Fatalf("\nexpected error: %q\n", NewParsingError(name).Error())
		}

		// Proper format
		if err := ff.AppendStr("PrincessLeia    Organa  "); err != nil {
			t.Fatalf("\nunexpected error: %q\n", err.Error())
		}

		// Names are correct lengths, but S bleeds into the firstname
		if err := ff.AppendStr("Han    Solo     "); err != nil {
			t.Fatalf("\nunexpected error: %q\n", err.Error())
		}

		// Set wrong value
		key := "name"
		if err := ff.SetValue(ff.Len()-1, key, "Han"); err == nil {
			t.Fatalf("\nexpected error: %q\n", NewMissingKeyError(key, ff.FormatsAt(ff.Len()-1)).Error())
		}

		if err := ff.SetValue(ff.Len()-1, "first", "Han"); err != nil {
			t.Fatalf("\nunexpected error: %q\n", err.Error())
		}

		if err := ff.SetValue(ff.Len()-1, "last", "Solo"); err != nil {
			t.Fatalf("\nunexpected error: %q\n", err.Error())
		}

		// Remove last line entirely and append the corrected one to the end
		ff.Remove(ff.Len() - 1)
		if err := ff.AppendStr("Han     Solo    "); err != nil {
			t.Fatalf("\nunexpected error: %q\n", err.Error())
		}

		// Sort longest lines to shortest
		ff.Sort(func(ln0, ln1 Line) bool { return ln1.length < ln0.length })
		if err := ff.WriteFile("starwars_sorted.txt"); err != nil {
			t.Fatalf("\nunexpected error: %q\n", err.Error())
		}

		ff.Clear()
	}

	{ // Read from/write to file
		// Read from reader (a file)
		file, err := os.OpenFile("starwars_1.txt", os.O_RDONLY, os.ModePerm)
		if err != nil {
			t.Fatalf("\nunexpected error: %q\n", err.Error())
		}

		if _, err := ff.ReadFrom(file); err != nil {
			t.Fatalf("\nunexpected error: %q\n", err.Error())
		}

		if err := ff.WriteFile("starwars_2.txt"); err != nil {
			t.Fatalf("\nunexpected error: %q\n", err.Error())
		}

		r, err := equalFiles("starwars_1.txt", "starwars_2.txt")
		if err != nil {
			t.Fatalf("\nunexpected error: %q\n", err.Error())
		}

		if !r {
			t.Fatalf("\nexpected %q and %q to be equal\n", "starwars_1.txt", "starwars_2.txt")
		}

		ff.Clear()

		// Just read file like you're supposed to
		if err := ff.ReadFile("starwars_1.txt"); err != nil {
			t.Fatalf("\nunexpected error: %q\n", err.Error())
		}

		if err := ff.WriteFile("starwars_2.txt"); err != nil {
			t.Fatalf("\nunexpected error: %q\n", err.Error())
		}

		r, err = equalFiles("starwars_1.txt", "starwars_2.txt")
		if err != nil {
			t.Fatalf("\nunexpected error: %q\n", err.Error())
		}

		if !r {
			t.Fatalf("\nexpected %q and %q to be equal\n", "starwars_1.txt", "starwars_2.txt")
		}

		ff.Clear()
	}
}

func TestJSON(t *testing.T) {
	fmts := func(line string) []Format {
		switch len(line) {
		case 8: // Single name
			return []Format{
				NewFormat("name", 0, 8),
			}
		case 16: // Two names
			return []Format{
				NewFormat("first", 0, 8),
				NewFormat("last", 8, 8),
			}
		case 24: // Three names
			return []Format{
				NewFormat("title", 0, 8),
				NewFormat("first", 8, 8),
				NewFormat("last", 16, 8),
			}
		default:
			return nil
		}
	}

	ff := New(fmts)
	if err := ff.ReadFile("starwars_1.txt"); err != nil {
		t.Fatalf("\nunexpected error: %q\n", err.Error())
	}

	b, err := json.Marshal(ff)
	if err != nil {
		t.Fatalf("\nunexpected error: %q\n", err.Error())
	}

	if !json.Valid(b) {
		t.Fatalf("\nexpected valid json, received %q\n", string(b))
	}

	ff.Clear()
	ff.AppendStr("        ")
	ff.AppendStr("12345678")

	b, err = json.Marshal(ff)
	if err != nil {
		t.Fatalf("\nunexpected error: %q\n", err.Error())
	}

	if !json.Valid(b) {
		t.Fatalf("\nexpected valid json, received %q\n", string(b))
	}
}

func BenchmarkJSON(b *testing.B) {
	fmts := func(line string) []Format {
		switch len(line) {
		case 8: // Single name
			return []Format{
				NewFormat("name", 0, 8),
			}
		case 16: // Two names
			return []Format{
				NewFormat("first", 0, 8),
				NewFormat("last", 8, 8),
			}
		case 24: // Three names
			return []Format{
				NewFormat("title", 0, 8),
				NewFormat("first", 8, 8),
				NewFormat("last", 16, 8),
			}
		default:
			return nil
		}
	}

	var (
		ff = New(fmts)
		ln = "12345678"
	)

	benchmarkMarshalFlatFile(b, ff, fmt.Sprintf("%d lines of %q", ff.Len(), ln))
	for i := 0; i < 10; i++ {
		if err := ff.AppendStr(ln); err != nil {
			b.Fatalf("\nunexpected error: %q", err.Error())
		}

		benchmarkMarshalFlatFile(b, ff, fmt.Sprintf("%d lines of %q", ff.Len(), ln))
	}
}

func benchmarkMarshalFlatFile(b *testing.B, ff *FlatFile, name string) bool {
	f := func(b0 *testing.B) {
		for i := 0; i < b0.N; i++ {
			_, _ = json.Marshal(ff)
		}
	}

	return b.Run(name, f)
}
