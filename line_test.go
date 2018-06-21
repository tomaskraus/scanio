package line

import (
	"os"
	"strings"
	"testing"
)

type result struct {
	stop     bool
	number   int
	match    bool
	text     string
	original string
}

func TestReaderScannerEmpty(t *testing.T) {
	f := strings.NewReader("")

	sc := NewReaderScanner(f)

	expected := []result{
		{false, 0, false, "", ""},
		{false, 0, false, "", ""},
		{false, 0, false, "", ""},
		{false, 0, false, "", ""},
	}
	for _, v := range expected {
		res, num, match, text, orig := sc.Scan(), sc.Number(), sc.Match(), sc.Text(), sc.Original()

		if res != v.stop || num != v.number || match != v.match || text != v.text || orig != v.original {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text, orig})
		}
	}
}

func TestReaderScannerFile(t *testing.T) {
	f, err := os.Open("assets/simpleFile.txt")
	defer f.Close()
	if err != nil {
		t.Error(err)
		return
	}

	sc := NewReaderScanner(f)

	expected := []result{
		{true, 1, true, "this is a simple file", "this is a simple file"},
		{true, 2, true, "next line is empty", "next line is empty"},
		{true, 3, true, "", ""},
		{true, 4, true, "next line has two spaces", "next line has two spaces"},
		{true, 5, true, "  ", "  "},
		{true, 6, true, "# bash-like comment", "# bash-like comment"},
		{true, 7, true, "line with two trailing spaces  ", "line with two trailing spaces  "},
		{true, 8, true, " line with one leading space", " line with one leading space"},
		{true, 9, true, " line with one leading and one trailing space ", " line with one leading and one trailing space "},
		{true, 10, true, "# bash-like comment 2 ", "# bash-like comment 2 "},
		{true, 11, true, "last line", "last line"},
		{false, 11, false, "", ""},
		{false, 11, false, "", ""},
		{false, 11, false, "", ""},
	}
	for _, v := range expected {
		res, num, match, text, orig := sc.Scan(), sc.Number(), sc.Match(), sc.Text(), sc.Original()

		if res != v.stop || num != v.number || match != v.match || text != v.text || orig != v.original {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text, orig})
		}
	}
}

func TestRuleScanner(t *testing.T) {
	f, err := os.Open("assets/simpleFile.txt")
	defer f.Close()
	if err != nil {
		t.Error(err)
		return
	}

	// matches a line with a # at the begin, trims a #
	sc := NewFilterScanner(NewReaderScanner(f), func(in string) (bool, string) {
		return strings.HasPrefix(in, "#"), strings.Replace(in, "#", "", 1)
	})

	expected := []result{
		{true, 1, false, "this is a simple file", "this is a simple file"},
		{true, 2, false, "next line is empty", "next line is empty"},
		{true, 3, false, "", ""},
		{true, 4, false, "next line has two spaces", "next line has two spaces"},
		{true, 5, false, "  ", "  "},
		{true, 6, true, " bash-like comment", "# bash-like comment"},
		{true, 7, false, "line with two trailing spaces  ", "line with two trailing spaces  "},
		{true, 8, false, " line with one leading space", " line with one leading space"},
		{true, 9, false, " line with one leading and one trailing space ", " line with one leading and one trailing space "},
		{true, 10, true, " bash-like comment 2 ", "# bash-like comment 2 "},
		{true, 11, false, "last line", "last line"},
		{false, 11, false, "", ""},
		{false, 11, false, "", ""},
		{false, 11, false, "", ""},
	}
	for _, v := range expected {
		res, num, match, text, orig := sc.Scan(), sc.Number(), sc.Match(), sc.Text(), sc.Original()

		if res != v.stop || num != v.number || match != v.match || text != v.text || orig != v.original {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text, orig})
		}
	}
}

func TestOnlyMatchScannerEmpty(t *testing.T) {
	f := strings.NewReader("")

	sc := NewOnlyMatchScanner(NewReaderScanner(f))

	expected := []result{
		{false, 0, false, "", ""},
		{false, 0, false, "", ""},
		{false, 0, false, "", ""},
	}
	for _, v := range expected {
		res, num, match, text, orig := sc.Scan(), sc.Number(), sc.Match(), sc.Text(), sc.Original()

		if res != v.stop || num != v.number || match != v.match || text != v.text || orig != v.original {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text, orig})
		}
	}
}
func TestOnlyMatchScannerFull(t *testing.T) {
	f, err := os.Open("assets/simpleFile.txt")
	defer f.Close()
	if err != nil {
		t.Error(err)
		return
	}

	sc := NewOnlyMatchScanner(NewReaderScanner(f))

	expected := []result{
		{true, 1, true, "this is a simple file", "this is a simple file"},
		{true, 2, true, "next line is empty", "next line is empty"},
		{true, 3, true, "", ""},
		{true, 4, true, "next line has two spaces", "next line has two spaces"},
		{true, 5, true, "  ", "  "},
		{true, 6, true, "# bash-like comment", "# bash-like comment"},
		{true, 7, true, "line with two trailing spaces  ", "line with two trailing spaces  "},
		{true, 8, true, " line with one leading space", " line with one leading space"},
		{true, 9, true, " line with one leading and one trailing space ", " line with one leading and one trailing space "},
		{true, 10, true, "# bash-like comment 2 ", "# bash-like comment 2 "},
		{true, 11, true, "last line", "last line"},
		{false, 11, false, "", ""},
		{false, 11, false, "", ""},
		{false, 11, false, "", ""},
	}
	for _, v := range expected {
		res, num, match, text, orig := sc.Scan(), sc.Number(), sc.Match(), sc.Text(), sc.Original()

		if res != v.stop || num != v.number || match != v.match || text != v.text || orig != v.original {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text, orig})
		}
	}
}
func TestOnlyMatchScannerFilter(t *testing.T) {
	f, err := os.Open("assets/simpleFile.txt")
	defer f.Close()
	if err != nil {
		t.Error(err)
		return
	}

	// matches a line with a # at the begin, trims a #
	sc := NewOnlyMatchScanner(
		NewFilterScanner(NewReaderScanner(f), func(in string) (bool, string) {
			return strings.HasPrefix(in, "#"), strings.Replace(in, "#", "", 1)
		}))

	expected := []result{
		{true, 6, true, " bash-like comment", "# bash-like comment"},
		{true, 10, true, " bash-like comment 2 ", "# bash-like comment 2 "},
		{false, 11, false, "", ""},
		{false, 11, false, "", ""},
		{false, 11, false, "", ""},
	}
	for _, v := range expected {
		res, num, match, text, orig := sc.Scan(), sc.Number(), sc.Match(), sc.Text(), sc.Original()

		if res != v.stop || num != v.number || match != v.match || text != v.text || orig != v.original {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text, orig})
		}
	}
}
