package liner

import (
	"fmt"
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

	li := NewLiner(f)

	expected := []result{
		{false, 0, false, ""},
		{false, 0, false, ""},
		{false, 0, false, ""},
		{false, 0, false, ""},
	}
	for _, v := range expected {
		res, num, match, text := li.Scan(), li.Number(), li.Match(), li.Text()

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

	li := NewLiner(f)

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
		res, num, match, text := li.Scan(), li.Number(), li.Match(), li.Text()

		if res != v.canParse || num != v.number || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
		}
	}
}

func TestRuleLiner(t *testing.T) {
	f, err := os.Open("assets/simpleFile.txt")
	defer f.Close()
	if err != nil {
		t.Error(err)
		return
	}

	// matches a line with a # at the begin, trims a #
	li := NewRuleLiner(NewLiner(f), func(in string) bool {
		return strings.HasPrefix(in, "#")
	})

	expected := []result{
		{true, 1, false, "this is a simple file"},
		{true, 2, false, "next line is empty"},
		{true, 3, false, ""},
		{true, 4, false, "next line has two spaces"},
		{true, 5, false, "  "},
		{true, 6, true, "# bash-like comment"},
		{true, 7, false, "line with two trailing spaces  "},
		{true, 8, false, " line with one leading space"},
		{true, 9, false, " line with one leading and one trailing space "},
		{true, 10, true, "# bash-like comment 2 "},
		{true, 11, false, "last line"},
		{false, 11, false, ""},
		{false, 11, false, ""},
		{false, 11, false, ""},
	}
	for _, v := range expected {
		res, num, match, text := li.Scan(), li.Number(), li.Match(), li.Text()

		if res != v.canParse || num != v.number || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
		}
	}
}

func TestRuleLinerEmpty(t *testing.T) {
	f := strings.NewReader("")

	li := NewRuleLiner(NewLiner(f), func(in string) bool {
		return strings.HasPrefix(in, "#")
	})

	expected := []result{
		{false, 0, false, ""},
		{false, 0, false, ""},
		{false, 0, false, ""},
	}
	for _, v := range expected {
		res, num, match, text := li.Scan(), li.Number(), li.Match(), li.Text()

		if res != v.canParse || num != v.number || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
		}
	}
}
func TestRuleLinerFull(t *testing.T) {
	f, err := os.Open("assets/simpleFile.txt")
	defer f.Close()
	if err != nil {
		t.Error(err)
		return
	}

	li := NewRuleLiner(NewLiner(f), func(in string) bool {
		return strings.HasPrefix(in, "#")
	})

	expected := []result{
		{true, 1, false, "this is a simple file"},
		{true, 2, false, "next line is empty"},
		{true, 3, false, ""},
		{true, 4, false, "next line has two spaces"},
		{true, 5, false, "  "},
		{true, 6, true, "# bash-like comment"},
		{true, 7, false, "line with two trailing spaces  "},
		{true, 8, false, " line with one leading space"},
		{true, 9, false, " line with one leading and one trailing space "},
		{true, 10, true, "# bash-like comment 2 "},
		{true, 11, false, "last line"},
		{false, 11, false, ""},
		{false, 11, false, ""},
		{false, 11, false, ""},
	}
	for _, v := range expected {
		res, num, match, text := li.Scan(), li.Number(), li.Match(), li.Text()

		if res != v.canParse || num != v.number || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
		}
	}
}

func TestOnlyMatchLinerEmpty(t *testing.T) {
	f := strings.NewReader("")

	li := NewOnlyMatchLiner(NewLiner(f))

	expected := []result{
		{false, 0, false, ""},
		{false, 0, false, ""},
		{false, 0, false, ""},
	}
	for _, v := range expected {
		res, num, match, text := li.Scan(), li.Number(), li.Match(), li.Text()

		if res != v.canParse || num != v.number || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
		}
	}
}
func TestOnlyMatchLinerFull(t *testing.T) {
	f, err := os.Open("assets/simpleFile.txt")
	defer f.Close()
	if err != nil {
		t.Error(err)
		return
	}

	li := NewOnlyMatchLiner(NewLiner(f))

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
		res, num, match, text := li.Scan(), li.Number(), li.Match(), li.Text()

		if res != v.canParse || num != v.number || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
		}
	}
}

func TestOnlyMatchLinerRuled(t *testing.T) {
	f, err := os.Open("assets/simpleFile.txt")
	defer f.Close()
	if err != nil {
		t.Error(err)
		return
	}

	li := NewOnlyMatchLiner(NewRuleLiner(NewLiner(f), func(in string) bool {
		return strings.HasPrefix(in, "#")
	}))

	expected := []result{
		{true, 6, true, "# bash-like comment"},
		{true, 10, true, "# bash-like comment 2 "},
		{false, 11, false, ""},
		{false, 11, false, ""},
		{false, 11, false, ""},
	}
	for _, v := range expected {
		res, num, match, text := li.Scan(), li.Number(), li.Match(), li.Text()

		if res != v.canParse || num != v.number || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
		}
	}
}

func TestLastLinerEmpty(t *testing.T) {
	f := strings.NewReader("")

	li := NewLastLiner(NewLiner(f))

	expected := []resultL{
		{false, 0, false, "", true},
		{false, 0, false, "", true},
		{false, 0, false, "", true},
	}
	for _, v := range expected {
		res, num, match, text, last := li.Scan(), li.Number(), li.Match(), li.Text(), li.Last()

		if res != v.canParse || num != v.number || match != v.match || text != v.text || last != v.last {
			t.Errorf("should be %v, is %v", v, resultL{res, num, match, text, last})
		}
	}
}
func TestLastLinerOneLine(t *testing.T) {
	f := strings.NewReader("one line")

	li := NewLastLiner(NewLiner(f))

	expected := []resultL{
		{true, 1, true, "one line", true},
		{false, 1, false, "", true},
		{false, 1, false, "", true},
	}
	for _, v := range expected {
		res, num, match, text, last := li.Scan(), li.Number(), li.Match(), li.Text(), li.Last()

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

	li := NewLastLiner(NewLiner(f))

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
		res, num, match, text, last := li.Scan(), li.Number(), li.Match(), li.Text(), li.Last()

		if res != v.canParse || num != v.number || match != v.match || text != v.text || last != v.last {
			t.Errorf("should be %v, is %v", v, resultL{res, num, match, text, last})
		}
	}
}

func ExampleNewRuleLiner() {
	f := strings.NewReader("\n# comment 1\n  \n#comment2\nsomething")

	li := NewLastLiner(
		NewOnlyMatchLiner(
			NewRuleLiner(
				NewLiner(f),
				func(s string) bool {
					return strings.HasPrefix(s, "#")
				})))

	for li.Scan() {
		if li.Last() {
			fmt.Printf("(%d, %q).", li.Number(), li.Text())
		} else {
			fmt.Printf("(%d, %q), ", li.Number(), li.Text())
		}
	}
	// Output:
	// (2, "# comment 1"), (4, "#comment2").
}
