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
}

type resultL struct {
	canParse bool
	number   int
	match    bool
	text     string
	last     bool
}

func TestReaderLinerEmpty(t *testing.T) {
	f := strings.NewReader("")

	lin := NewLiner(f)

	expected := []result{
		{false, 0, false, ""},
		{false, 0, false, ""},
		{false, 0, false, ""},
		{false, 0, false, ""},
	}
	for _, v := range expected {
		res, num, match, text := lin.Scan(), lin.Number(), lin.Match(), lin.Text()

		if res != v.canParse || num != v.number || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
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
		{true, 1, true, "this is a simple file"},
		{true, 2, true, "next line is empty"},
		{true, 3, true, ""},
		{true, 4, true, "next line has two spaces"},
		{true, 5, true, "  "},
		{true, 6, true, "# bash-like comment"},
		{true, 7, true, "line with two trailing spaces  "},
		{true, 8, true, " line with one leading space"},
		{true, 9, true, " line with one leading and one trailing space "},
		{true, 10, true, "# bash-like comment 2 "},
		{true, 11, true, "last line"},
		{false, 11, false, ""},
		{false, 11, false, ""},
		{false, 11, false, ""},
	}
	for _, v := range expected {
		res, num, match, text := lin.Scan(), lin.Number(), lin.Match(), lin.Text()

		if res != v.canParse || num != v.number || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
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
	lin := NewMatchLiner(NewLiner(f), func(in string, inf Info) (bool, string, bool) {
		return strings.HasPrefix(in, "#"), strings.Replace(in, "#", "", 1), false
	})

	expected := []result{
		{true, 1, false, "this is a simple file"},
		{true, 2, false, "next line is empty"},
		{true, 3, false, ""},
		{true, 4, false, "next line has two spaces"},
		{true, 5, false, "  "},
		{true, 6, true, " bash-like comment"},
		{true, 7, false, "line with two trailing spaces  "},
		{true, 8, false, " line with one leading space"},
		{true, 9, false, " line with one leading and one trailing space "},
		{true, 10, true, " bash-like comment 2 "},
		{true, 11, false, "last line"},
		{false, 11, false, ""},
		{false, 11, false, ""},
		{false, 11, false, ""},
	}
	for _, v := range expected {
		res, num, match, text := lin.Scan(), lin.Number(), lin.Match(), lin.Text()

		if res != v.canParse || num != v.number || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
		}
	}
}

func TestMatchLinerEmpty(t *testing.T) {
	f := strings.NewReader("")

	lin := NewFilterLiner(NewLiner(f))

	expected := []result{
		{false, 0, false, ""},
		{false, 0, false, ""},
		{false, 0, false, ""},
	}
	for _, v := range expected {
		res, num, match, text := lin.Scan(), lin.Number(), lin.Match(), lin.Text()

		if res != v.canParse || num != v.number || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
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

	lin := NewFilterLiner(NewLiner(f))

	expected := []result{
		{true, 1, true, "this is a simple file"},
		{true, 2, true, "next line is empty"},
		{true, 3, true, ""},
		{true, 4, true, "next line has two spaces"},
		{true, 5, true, "  "},
		{true, 6, true, "# bash-like comment"},
		{true, 7, true, "line with two trailing spaces  "},
		{true, 8, true, " line with one leading space"},
		{true, 9, true, " line with one leading and one trailing space "},
		{true, 10, true, "# bash-like comment 2 "},
		{true, 11, true, "last line"},
		{false, 11, false, ""},
		{false, 11, false, ""},
		{false, 11, false, ""},
	}
	for _, v := range expected {
		res, num, match, text := lin.Scan(), lin.Number(), lin.Match(), lin.Text()

		if res != v.canParse || num != v.number || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
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
	lin := NewFilterLiner(
		NewMatchLiner(NewLiner(f), func(in string, inf Info) (bool, string, bool) {
			return strings.HasPrefix(in, "#"), strings.Replace(in, "#", "", 1), false
		}))

	expected := []result{
		{true, 6, true, " bash-like comment"},
		{true, 10, true, " bash-like comment 2 "},
		{false, 11, false, ""},
		{false, 11, false, ""},
		{false, 11, false, ""},
	}
	for _, v := range expected {
		res, num, match, text := lin.Scan(), lin.Number(), lin.Match(), lin.Text()

		if res != v.canParse || num != v.number || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
		}
	}
}

func TestLastLinerEmpty(t *testing.T) {
	f := strings.NewReader("")

	lin := NewLastLiner(NewLiner(f))

	expected := []resultL{
		{false, 0, false, "", true},
		{false, 0, false, "", true},
		{false, 0, false, "", true},
	}
	for _, v := range expected {
		res, num, match, text, last := lin.Scan(), lin.Number(), lin.Match(), lin.Text(), lin.Last()

		if res != v.canParse || num != v.number || match != v.match || text != v.text || last != v.last {
			t.Errorf("should be %v, is %v", v, resultL{res, num, match, text, last})
		}
	}
}
func TestLastLinerOneLine(t *testing.T) {
	f := strings.NewReader("one line")

	lin := NewLastLiner(NewLiner(f))

	expected := []resultL{
		{true, 1, true, "one line", true},
		{false, 1, false, "", true},
		{false, 1, false, "", true},
	}
	for _, v := range expected {
		res, num, match, text, last := lin.Scan(), lin.Number(), lin.Match(), lin.Text(), lin.Last()

		if res != v.canParse || num != v.number || match != v.match || text != v.text || last != v.last {
			t.Errorf("should be %v, is %v", v, resultL{res, num, match, text, last})
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

	lin := NewLastLiner(NewFilterLiner(NewLiner(f)))

	expected := []resultL{
		{true, 1, true, "this is a simple file", false},
		{true, 2, true, "next line is empty", false},
		{true, 3, true, "", false},
		{true, 4, true, "next line has two spaces", false},
		{true, 5, true, "  ", false},
		{true, 6, true, "# bash-like comment", false},
		{true, 7, true, "line with two trailing spaces  ", false},
		{true, 8, true, " line with one leading space", false},
		{true, 9, true, " line with one leading and one trailing space ", false},
		{true, 10, true, "# bash-like comment 2 ", false},
		{true, 11, true, "last line", true},
		{false, 11, false, "", true},
		{false, 11, false, "", true},
		{false, 11, false, "", true},
	}
	for _, v := range expected {
		res, num, match, text, last := lin.Scan(), lin.Number(), lin.Match(), lin.Text(), lin.Last()

		if res != v.canParse || num != v.number || match != v.match || text != v.text || last != v.last {
			t.Errorf("should be %v, is %v", v, resultL{res, num, match, text, last})
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
		NewFilterLiner(NewMatchLiner(NewLiner(f), func(in string, inf Info) (bool, string, bool) {
			return strings.HasPrefix(in, "#"), strings.Replace(in, "#", "", 1), false
		})))

	expected := []resultL{
		{true, 6, true, " bash-like comment", false},
		{true, 10, true, " bash-like comment 2 ", true},
		{false, 11, false, "", true},
		{false, 11, false, "", true},
		{false, 11, false, "", true},
	}
	for _, v := range expected {
		res, num, match, text, last := lin.Scan(), lin.Number(), lin.Match(), lin.Text(), lin.Last()
		if res != v.canParse || num != v.number || match != v.match || text != v.text || last != v.last {
			t.Errorf("should be %v, is %v", v, resultL{res, num, match, text, last})
		}
	}
}
