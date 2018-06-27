package scanio

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"
)

type result struct {
	canParse bool
	index    int
	match    bool
	text     string
}

type resultL struct {
	canParse bool
	index    int
	match    bool
	text     string
	last     bool
}

func TestScanWords(t *testing.T) {
	f := strings.NewReader("This is example  ")
	scn := NewScanner(f)
	scn.Split(bufio.ScanWords)

	expected := []result{
		{true, 0, true, "This"},
		{true, 1, true, "is"},
		{true, 2, true, "example"},
		{false, 2, false, ""},
		{false, 2, false, ""},
		{false, 2, false, ""},
		{false, 2, false, ""},
	}
	for _, v := range expected {
		res, index, match, text := scn.Scan(), scn.Index(), scn.Match(), scn.Text()

		if res != v.canParse || index != v.index || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, index, match, text})
		}
	}
}

func TestReaderScannerEmpty(t *testing.T) {
	f := strings.NewReader("")

	scn := NewScanner(f)

	expected := []result{
		{false, -1, false, ""},
		{false, -1, false, ""},
		{false, -1, false, ""},
		{false, -1, false, ""},
	}
	for _, v := range expected {
		res, index, match, text := scn.Scan(), scn.Index(), scn.Match(), scn.Text()

		if res != v.canParse || index != v.index || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, index, match, text})
		}
	}
}

func TestScannerFile(t *testing.T) {
	f, err := os.Open("assets/simpleFile.txt")
	defer f.Close()
	if err != nil {
		t.Error(err)
		return
	}

	scn := NewScanner(f)

	expected := []result{
		{true, 0, true, "this is a simple file"},
		{true, 1, true, "next line is empty"},
		{true, 2, true, ""},
		{true, 3, true, "next line has two spaces"},
		{true, 4, true, "  "},
		{true, 5, true, "# bash-like comment"},
		{true, 6, true, "line with two trailing spaces  "},
		{true, 7, true, " line with one leading space"},
		{true, 8, true, " line with one leading and one trailing space "},
		{true, 9, true, "# bash-like comment 2 "},
		{true, 10, true, "last line"},
		{false, 10, false, ""},
		{false, 10, false, ""},
		{false, 10, false, ""},
	}
	for _, v := range expected {
		res, index, match, text := scn.Scan(), scn.Index(), scn.Match(), scn.Text()

		if res != v.canParse || index != v.index || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, index, match, text})
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
	scn := NewRuleScanner(NewScanner(f), func(b []byte) bool {
		return bytes.HasPrefix(b, []byte("#"))
	})

	expected := []result{
		{true, 0, false, "this is a simple file"},
		{true, 1, false, "next line is empty"},
		{true, 2, false, ""},
		{true, 3, false, "next line has two spaces"},
		{true, 4, false, "  "},
		{true, 5, true, "# bash-like comment"},
		{true, 6, false, "line with two trailing spaces  "},
		{true, 7, false, " line with one leading space"},
		{true, 8, false, " line with one leading and one trailing space "},
		{true, 9, true, "# bash-like comment 2 "},
		{true, 10, false, "last line"},
		{false, 10, false, ""},
		{false, 10, false, ""},
		{false, 10, false, ""},
	}
	for _, v := range expected {
		res, index, match, text := scn.Scan(), scn.Index(), scn.Match(), scn.Text()

		if res != v.canParse || index != v.index || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, index, match, text})
		}
	}
}

func TestRuleScannerEmpty(t *testing.T) {
	f := strings.NewReader("")

	scn := NewRuleScanner(NewScanner(f), func(b []byte) bool {
		return bytes.HasPrefix(b, []byte("#"))
	})

	expected := []result{
		{false, -1, false, ""},
		{false, -1, false, ""},
		{false, -1, false, ""},
	}
	for _, v := range expected {
		res, index, match, text := scn.Scan(), scn.Index(), scn.Match(), scn.Text()

		if res != v.canParse || index != v.index || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, index, match, text})
		}
	}
}
func TestRuleScannerFull(t *testing.T) {
	f, err := os.Open("assets/simpleFile.txt")
	defer f.Close()
	if err != nil {
		t.Error(err)
		return
	}

	scn := NewRuleScanner(NewScanner(f), func(b []byte) bool {
		return bytes.HasPrefix(b, []byte("#"))
	})

	expected := []result{
		{true, 0, false, "this is a simple file"},
		{true, 1, false, "next line is empty"},
		{true, 2, false, ""},
		{true, 3, false, "next line has two spaces"},
		{true, 4, false, "  "},
		{true, 5, true, "# bash-like comment"},
		{true, 6, false, "line with two trailing spaces  "},
		{true, 7, false, " line with one leading space"},
		{true, 8, false, " line with one leading and one trailing space "},
		{true, 9, true, "# bash-like comment 2 "},
		{true, 10, false, "last line"},
		{false, 10, false, ""},
		{false, 10, false, ""},
		{false, 10, false, ""},
	}
	for _, v := range expected {
		res, index, match, text := scn.Scan(), scn.Index(), scn.Match(), scn.Text()

		if res != v.canParse || index != v.index || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, index, match, text})
		}
	}
}

func TestOnlyMatchScannerEmpty(t *testing.T) {
	f := strings.NewReader("")

	scn := NewOnlyMatchScanner(NewScanner(f))

	expected := []result{
		{false, -1, false, ""},
		{false, -1, false, ""},
		{false, -1, false, ""},
	}
	for _, v := range expected {
		res, index, match, text := scn.Scan(), scn.Index(), scn.Match(), scn.Text()

		if res != v.canParse || index != v.index || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, index, match, text})
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

	scn := NewOnlyMatchScanner(NewScanner(f))

	expected := []result{
		{true, 0, true, "this is a simple file"},
		{true, 1, true, "next line is empty"},
		{true, 2, true, ""},
		{true, 3, true, "next line has two spaces"},
		{true, 4, true, "  "},
		{true, 5, true, "# bash-like comment"},
		{true, 6, true, "line with two trailing spaces  "},
		{true, 7, true, " line with one leading space"},
		{true, 8, true, " line with one leading and one trailing space "},
		{true, 9, true, "# bash-like comment 2 "},
		{true, 10, true, "last line"},
		{false, 10, false, ""},
		{false, 10, false, ""},
		{false, 10, false, ""},
	}
	for _, v := range expected {
		res, index, match, text := scn.Scan(), scn.Index(), scn.Match(), scn.Text()

		if res != v.canParse || index != v.index || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, index, match, text})
		}
	}
}

func TestOnlyMatchScannerRuled(t *testing.T) {
	f, err := os.Open("assets/simpleFile.txt")
	defer f.Close()
	if err != nil {
		t.Error(err)
		return
	}

	scn := NewOnlyMatchScanner(NewRuleScanner(NewScanner(f), func(b []byte) bool {
		return bytes.HasPrefix(b, []byte("#"))
	}))

	expected := []result{
		{true, 5, true, "# bash-like comment"},
		{true, 9, true, "# bash-like comment 2 "},
		{false, 10, false, ""},
		{false, 10, false, ""},
		{false, 10, false, ""},
	}
	for _, v := range expected {
		res, index, match, text := scn.Scan(), scn.Index(), scn.Match(), scn.Text()

		if res != v.canParse || index != v.index || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, index, match, text})
		}
	}
}
