package liner

import (
	"os"
	"strings"
	"testing"
)

type result struct {
	canParse bool
	number   int
	match    bool
	text     string
	original string
	eof      bool
}

type resultL struct {
	canParse bool
	number   int
	match    bool
	text     string
	original string
	eof      bool
	last     bool
}

func TestReaderLinerEmpty(t *testing.T) {
	f := strings.NewReader("")

	lin := NewLiner(f)

	expected := []result{
		{false, 0, false, "", "", true},
		{false, 0, false, "", "", true},
		{false, 0, false, "", "", true},
		{false, 0, false, "", "", true},
	}
	for _, v := range expected {
		res, num, match, text, orig, eof := lin.Scan(), lin.Number(), lin.Match(), lin.Text(), lin.Original(), lin.Eof()

		if res != v.canParse || num != v.number || match != v.match || text != v.text || orig != v.original || eof != v.eof {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text, orig, eof})
		}
	}
}

func TestLinerFile(t *testing.T) {
	f, err := os.Open("assets/simpleFile.txt")
	defer f.Close()
	if err != nil {
		t.Error(err)
		return
	}

	lin := NewLiner(f)

	expected := []result{
		{true, 1, true, "this is a simple file", "this is a simple file", false},
		{true, 2, true, "next line is empty", "next line is empty", false},
		{true, 3, true, "", "", false},
		{true, 4, true, "next line has two spaces", "next line has two spaces", false},
		{true, 5, true, "  ", "  ", false},
		{true, 6, true, "# bash-like comment", "# bash-like comment", false},
		{true, 7, true, "line with two trailing spaces  ", "line with two trailing spaces  ", false},
		{true, 8, true, " line with one leading space", " line with one leading space", false},
		{true, 9, true, " line with one leading and one trailing space ", " line with one leading and one trailing space ", false},
		{true, 10, true, "# bash-like comment 2 ", "# bash-like comment 2 ", false},
		{true, 11, true, "last line", "last line", false},
		{false, 11, false, "", "", true},
		{false, 11, false, "", "", true},
		{false, 11, false, "", "", true},
	}
	for _, v := range expected {
		res, num, match, text, orig, eof := lin.Scan(), lin.Number(), lin.Match(), lin.Text(), lin.Original(), lin.Eof()

		if res != v.canParse || num != v.number || match != v.match || text != v.text || orig != v.original || eof != v.eof {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text, orig, eof})
		}
	}
}

func TestFilterLiner(t *testing.T) {
	f, err := os.Open("assets/simpleFile.txt")
	defer f.Close()
	if err != nil {
		t.Error(err)
		return
	}

	// matches a line with a # at the begin, trims a #
	lin := NewFilterLiner(NewLiner(f), func(in string) (bool, string) {
		return strings.HasPrefix(in, "#"), strings.Replace(in, "#", "", 1)
	})

	expected := []result{
		{true, 1, false, "this is a simple file", "this is a simple file", false},
		{true, 2, false, "next line is empty", "next line is empty", false},
		{true, 3, false, "", "", false},
		{true, 4, false, "next line has two spaces", "next line has two spaces", false},
		{true, 5, false, "  ", "  ", false},
		{true, 6, true, " bash-like comment", "# bash-like comment", false},
		{true, 7, false, "line with two trailing spaces  ", "line with two trailing spaces  ", false},
		{true, 8, false, " line with one leading space", " line with one leading space", false},
		{true, 9, false, " line with one leading and one trailing space ", " line with one leading and one trailing space ", false},
		{true, 10, true, " bash-like comment 2 ", "# bash-like comment 2 ", false},
		{true, 11, false, "last line", "last line", false},
		{false, 11, false, "", "", true},
		{false, 11, false, "", "", true},
		{false, 11, false, "", "", true},
	}
	for _, v := range expected {
		res, num, match, text, orig, eof := lin.Scan(), lin.Number(), lin.Match(), lin.Text(), lin.Original(), lin.Eof()

		if res != v.canParse || num != v.number || match != v.match || text != v.text || orig != v.original || eof != v.eof {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text, orig, eof})
		}
	}
}

func TestMatchLinerEmpty(t *testing.T) {
	f := strings.NewReader("")

	lin := NewMatchLiner(NewLiner(f))

	expected := []result{
		{false, 0, false, "", "", true},
		{false, 0, false, "", "", true},
		{false, 0, false, "", "", true},
	}
	for _, v := range expected {
		res, num, match, text, orig, eof := lin.Scan(), lin.Number(), lin.Match(), lin.Text(), lin.Original(), lin.Eof()

		if res != v.canParse || num != v.number || match != v.match || text != v.text || orig != v.original || eof != v.eof {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text, orig, eof})
		}
	}
}
func TestMatchLinerFull(t *testing.T) {
	f, err := os.Open("assets/simpleFile.txt")
	defer f.Close()
	if err != nil {
		t.Error(err)
		return
	}

	lin := NewMatchLiner(NewLiner(f))

	expected := []result{
		{true, 1, true, "this is a simple file", "this is a simple file", false},
		{true, 2, true, "next line is empty", "next line is empty", false},
		{true, 3, true, "", "", false},
		{true, 4, true, "next line has two spaces", "next line has two spaces", false},
		{true, 5, true, "  ", "  ", false},
		{true, 6, true, "# bash-like comment", "# bash-like comment", false},
		{true, 7, true, "line with two trailing spaces  ", "line with two trailing spaces  ", false},
		{true, 8, true, " line with one leading space", " line with one leading space", false},
		{true, 9, true, " line with one leading and one trailing space ", " line with one leading and one trailing space ", false},
		{true, 10, true, "# bash-like comment 2 ", "# bash-like comment 2 ", false},
		{true, 11, true, "last line", "last line", false},
		{false, 11, false, "", "", true},
		{false, 11, false, "", "", true},
		{false, 11, false, "", "", true},
	}
	for _, v := range expected {
		res, num, match, text, orig, eof := lin.Scan(), lin.Number(), lin.Match(), lin.Text(), lin.Original(), lin.Eof()

		if res != v.canParse || num != v.number || match != v.match || text != v.text || orig != v.original || eof != v.eof {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text, orig, eof})
		}
	}
}
func TestMatchLinerFilter(t *testing.T) {
	f, err := os.Open("assets/simpleFile.txt")
	defer f.Close()
	if err != nil {
		t.Error(err)
		return
	}

	// matches a line with a # at the begin, trims a #
	lin := NewMatchLiner(
		NewFilterLiner(NewLiner(f), func(in string) (bool, string) {
			return strings.HasPrefix(in, "#"), strings.Replace(in, "#", "", 1)
		}))

	expected := []result{
		{true, 6, true, " bash-like comment", "# bash-like comment", false},
		{true, 10, true, " bash-like comment 2 ", "# bash-like comment 2 ", false},
		{false, 11, false, "", "", true},
		{false, 11, false, "", "", true},
		{false, 11, false, "", "", true},
	}
	for _, v := range expected {
		res, num, match, text, orig, eof := lin.Scan(), lin.Number(), lin.Match(), lin.Text(), lin.Original(), lin.Eof()

		if res != v.canParse || num != v.number || match != v.match || text != v.text || orig != v.original || eof != v.eof {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text, orig, eof})
		}
	}
}

func TestNoMatchLinerEmpty(t *testing.T) {
	f := strings.NewReader("")

	lin := NewNoMatchLiner(NewLiner(f))

	expected := []result{
		{false, 0, false, "", "", true},
		{false, 0, false, "", "", true},
		{false, 0, false, "", "", true},
	}
	for _, v := range expected {
		res, num, match, text, orig, eof := lin.Scan(), lin.Number(), lin.Match(), lin.Text(), lin.Original(), lin.Eof()

		if res != v.canParse || num != v.number || match != v.match || text != v.text || orig != v.original || eof != v.eof {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text, orig, eof})
		}
	}
}

func TestNoMatchLinerFull(t *testing.T) {
	f, err := os.Open("assets/simpleFile.txt")
	defer f.Close()
	if err != nil {
		t.Error(err)
		return
	}

	lin := NewNoMatchLiner(NewLiner(f))

	expected := []result{
		{false, 11, false, "", "", true},
		{false, 11, false, "", "", true},
		{false, 11, false, "", "", true},
	}
	for _, v := range expected {
		res, num, match, text, orig, eof := lin.Scan(), lin.Number(), lin.Match(), lin.Text(), lin.Original(), lin.Eof()

		if res != v.canParse || num != v.number || match != v.match || text != v.text || orig != v.original || eof != v.eof {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text, orig, eof})
		}
	}
}
func TestNoMatchLinerFilter(t *testing.T) {
	f, err := os.Open("assets/simpleFile.txt")
	defer f.Close()
	if err != nil {
		t.Error(err)
		return
	}

	// matches a line with a # at the begin, trims a #
	lin := NewNoMatchLiner(
		NewFilterLiner(NewLiner(f), func(in string) (bool, string) {
			return strings.HasPrefix(in, "#"), strings.Replace(in, "#", "", 1)
		}))

	expected := []result{
		{true, 1, true, "this is a simple file", "this is a simple file", false},
		{true, 2, true, "next line is empty", "next line is empty", false},
		{true, 3, true, "", "", false},
		{true, 4, true, "next line has two spaces", "next line has two spaces", false},
		{true, 5, true, "  ", "  ", false},
		{true, 7, true, "line with two trailing spaces  ", "line with two trailing spaces  ", false},
		{true, 8, true, " line with one leading space", " line with one leading space", false},
		{true, 9, true, " line with one leading and one trailing space ", " line with one leading and one trailing space ", false},
		{true, 11, true, "last line", "last line", false},
		{false, 11, false, "", "", true},
		{false, 11, false, "", "", true},
		{false, 11, false, "", "", true},
	}
	for _, v := range expected {
		res, num, match, text, orig, eof := lin.Scan(), lin.Number(), lin.Match(), lin.Text(), lin.Original(), lin.Eof()

		if res != v.canParse || num != v.number || match != v.match || text != v.text || orig != v.original || eof != v.eof {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text, orig, eof})
		}
	}
}

func TestLastLinerEmpty(t *testing.T) {
	f := strings.NewReader("")

	lin := NewLastLiner(NewLiner(f))

	expected := []resultL{
		{false, 0, false, "", "", true, true},
		{false, 0, false, "", "", true, true},
		{false, 0, false, "", "", true, true},
	}
	for _, v := range expected {
		res, num, match, text, orig, eof, last := lin.Scan(), lin.Number(), lin.Match(), lin.Text(), lin.Original(), lin.Eof(), lin.Last()

		if res != v.canParse || num != v.number || match != v.match || text != v.text || orig != v.original || eof != v.eof || last != v.last {
			t.Errorf("should be %v, is %v", v, resultL{res, num, match, text, orig, eof, last})
		}
	}
}
func TestLastLinerOneLine(t *testing.T) {
	f := strings.NewReader("one line")

	lin := NewLastLiner(NewLiner(f))

	expected := []resultL{
		{true, 1, true, "one line", "one line", false, true},
		{false, 1, false, "", "", true, true},
		{false, 1, false, "", "", true, true},
	}
	for _, v := range expected {
		res, num, match, text, orig, eof, last := lin.Scan(), lin.Number(), lin.Match(), lin.Text(), lin.Original(), lin.Eof(), lin.Last()

		if res != v.canParse || num != v.number || match != v.match || text != v.text || orig != v.original || eof != v.eof || last != v.last {
			t.Errorf("should be %v, is %v", v, resultL{res, num, match, text, orig, eof, last})
		}
	}
}

func TestLastLinerFull(t *testing.T) {
	f, err := os.Open("assets/simpleFile.txt")
	defer f.Close()
	if err != nil {
		t.Error(err)
		return
	}

	lin := NewLastLiner(NewMatchLiner(NewLiner(f)))

	expected := []resultL{
		{true, 1, true, "this is a simple file", "this is a simple file", false, false},
		{true, 2, true, "next line is empty", "next line is empty", false, false},
		{true, 3, true, "", "", false, false},
		{true, 4, true, "next line has two spaces", "next line has two spaces", false, false},
		{true, 5, true, "  ", "  ", false, false},
		{true, 6, true, "# bash-like comment", "# bash-like comment", false, false},
		{true, 7, true, "line with two trailing spaces  ", "line with two trailing spaces  ", false, false},
		{true, 8, true, " line with one leading space", " line with one leading space", false, false},
		{true, 9, true, " line with one leading and one trailing space ", " line with one leading and one trailing space ", false, false},
		{true, 10, true, "# bash-like comment 2 ", "# bash-like comment 2 ", false, false},
		{true, 11, true, "last line", "last line", false, true},
		{false, 11, false, "", "", true, true},
		{false, 11, false, "", "", true, true},
		{false, 11, false, "", "", true, true},
	}
	for _, v := range expected {
		res, num, match, text, orig, eof, last := lin.Scan(), lin.Number(), lin.Match(), lin.Text(), lin.Original(), lin.Eof(), lin.Last()

		if res != v.canParse || num != v.number || match != v.match || text != v.text || orig != v.original || eof != v.eof || last != v.last {
			t.Errorf("should be %v, is %v", v, resultL{res, num, match, text, orig, eof, last})
		}
	}
}

func TestLastLinerFilter(t *testing.T) {
	f, err := os.Open("assets/simpleFile.txt")
	defer f.Close()
	if err != nil {
		t.Error(err)
		return
	}

	// matches a line with a # at the begin, trims a #
	lin := NewLastLiner(
		NewMatchLiner(NewFilterLiner(NewLiner(f), func(in string) (bool, string) {
			return strings.HasPrefix(in, "#"), strings.Replace(in, "#", "", 1)
		})))

	expected := []resultL{
		{true, 6, true, " bash-like comment", "# bash-like comment", false, false},
		{true, 10, true, " bash-like comment 2 ", "# bash-like comment 2 ", false, true},
		{false, 11, false, "", "", true, true},
		{false, 11, false, "", "", true, true},
		{false, 11, false, "", "", true, true},
	}
	for _, v := range expected {
		res, num, match, text, orig, eof, last := lin.Scan(), lin.Number(), lin.Match(), lin.Text(), lin.Original(), lin.Eof(), lin.Last()

		if res != v.canParse || num != v.number || match != v.match || text != v.text || orig != v.original || eof != v.eof || last != v.last {
			t.Errorf("should be %v, is %v", v, resultL{res, num, match, text, orig, eof, last})
		}
	}
}
