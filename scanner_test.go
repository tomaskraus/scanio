package scanio

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

type result struct {
	canParse bool
	num      int
	match    bool
	text     string
}

type resultL struct {
	canParse bool
	num      int
	match    bool
	text     string
	last     bool
}

func TestScanWords(t *testing.T) {
	f := strings.NewReader("This is example  ")
	scn := NewScanner(f)
	scn.Split(bufio.ScanWords)

	expected := []result{
		{true, 1, true, "This"},
		{true, 2, true, "is"},
		{true, 3, true, "example"},
		{false, 3, false, ""},
		{false, 3, false, ""},
		{false, 3, false, ""},
		{false, 3, false, ""},
	}
	for _, v := range expected {
		res, num, match, text := scn.Scan(), scn.Num(), scn.Match(), scn.Text()

		if res != v.canParse || num != v.num || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
		}
	}
}

func TestReaderScannerEmpty(t *testing.T) {
	f := strings.NewReader("")

	scn := NewScanner(f)

	expected := []result{
		{false, 0, false, ""},
		{false, 0, false, ""},
		{false, 0, false, ""},
		{false, 0, false, ""},
	}
	for _, v := range expected {
		res, num, match, text := scn.Scan(), scn.Num(), scn.Match(), scn.Text()

		if res != v.canParse || num != v.num || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
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
		res, num, match, text := scn.Scan(), scn.Num(), scn.Match(), scn.Text()

		if res != v.canParse || num != v.num || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
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
	scn := NewRuleScanner(NewScanner(f), func(in string) bool {
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
		res, num, match, text := scn.Scan(), scn.Num(), scn.Match(), scn.Text()

		if res != v.canParse || num != v.num || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
		}
	}
}

func TestRuleScannerEmpty(t *testing.T) {
	f := strings.NewReader("")

	scn := NewRuleScanner(NewScanner(f), func(in string) bool {
		return strings.HasPrefix(in, "#")
	})

	expected := []result{
		{false, 0, false, ""},
		{false, 0, false, ""},
		{false, 0, false, ""},
	}
	for _, v := range expected {
		res, num, match, text := scn.Scan(), scn.Num(), scn.Match(), scn.Text()

		if res != v.canParse || num != v.num || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
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

	scn := NewRuleScanner(NewScanner(f), func(in string) bool {
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
		res, num, match, text := scn.Scan(), scn.Num(), scn.Match(), scn.Text()

		if res != v.canParse || num != v.num || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
		}
	}
}

func TestOnlyMatchScannerEmpty(t *testing.T) {
	f := strings.NewReader("")

	scn := NewOnlyMatchScanner(NewScanner(f))

	expected := []result{
		{false, 0, false, ""},
		{false, 0, false, ""},
		{false, 0, false, ""},
	}
	for _, v := range expected {
		res, num, match, text := scn.Scan(), scn.Num(), scn.Match(), scn.Text()

		if res != v.canParse || num != v.num || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
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
		res, num, match, text := scn.Scan(), scn.Num(), scn.Match(), scn.Text()

		if res != v.canParse || num != v.num || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
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

	scn := NewOnlyMatchScanner(NewRuleScanner(NewScanner(f), func(in string) bool {
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
		res, num, match, text := scn.Scan(), scn.Num(), scn.Match(), scn.Text()

		if res != v.canParse || num != v.num || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
		}
	}
}

func TestLastScannerEmpty(t *testing.T) {
	f := strings.NewReader("")

	scn := NewLastScanner(NewScanner(f))

	expected := []resultL{
		{false, 0, false, "", true},
		{false, 0, false, "", true},
		{false, 0, false, "", true},
	}
	for _, v := range expected {
		res, num, match, text, last := scn.Scan(), scn.Num(), scn.Match(), scn.Text(), scn.Last()

		if res != v.canParse || num != v.num || match != v.match || text != v.text || last != v.last {
			t.Errorf("should be %v, is %v", v, resultL{res, num, match, text, last})
		}
	}
}
func TestLastScannerOneLine(t *testing.T) {
	f := strings.NewReader("one line")

	scn := NewLastScanner(NewScanner(f))

	expected := []resultL{
		{true, 1, true, "one line", true},
		{false, 1, false, "", true},
		{false, 1, false, "", true},
	}
	for _, v := range expected {
		res, num, match, text, last := scn.Scan(), scn.Num(), scn.Match(), scn.Text(), scn.Last()

		if res != v.canParse || num != v.num || match != v.match || text != v.text || last != v.last {
			t.Errorf("should be %v, is %v", v, resultL{res, num, match, text, last})
		}
	}
}

func TestLastScannerFull(t *testing.T) {
	f, err := os.Open("assets/simpleFile.txt")
	defer f.Close()
	if err != nil {
		t.Error(err)
		return
	}

	scn := NewLastScanner(NewScanner(f))

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
		res, num, match, text, last := scn.Scan(), scn.Num(), scn.Match(), scn.Text(), scn.Last()

		if res != v.canParse || num != v.num || match != v.match || text != v.text || last != v.last {
			t.Errorf("should be %v, is %v", v, resultL{res, num, match, text, last})
		}
	}
}

func ExampleNewRuled() {
	f := strings.NewReader("\n# comment 1\n  \n#comment2\nsomething")

	scn := NewLastScanner(
		NewOnlyMatchScanner(
			NewRuleScanner(
				NewScanner(f),
				func(s string) bool {
					return strings.HasPrefix(s, "#")
				})))

	for scn.Scan() {
		if scn.Last() {
			fmt.Printf("(%d, %q).", scn.Num(), scn.Text())
		} else {
			fmt.Printf("(%d, %q), ", scn.Num(), scn.Text())
		}
	}
	// Output:
	// (2, "# comment 1"), (4, "#comment2").
}

func ExampleNewFilter() {
	f := strings.NewReader("\n# comment 1\n  \n#comment2\nsomething")

	scn := NewLastScanner(
		NewFilterScanner(
			NewScanner(f),
			func(s string) bool {
				return strings.HasPrefix(s, "#")
			}))

	for scn.Scan() {
		if scn.Last() {
			fmt.Printf("(%d, %q).", scn.Num(), scn.Text())
		} else {
			fmt.Printf("(%d, %q), ", scn.Num(), scn.Text())
		}
	}
	// Output:
	// (2, "# comment 1"), (4, "#comment2").
}

func TestOnlyNotMatchScannerEmpty(t *testing.T) {
	f := strings.NewReader("")

	scn := NewOnlyNotMatchScanner(NewScanner(f))

	expected := []result{
		{false, 0, false, ""},
		{false, 0, false, ""},
		{false, 0, false, ""},
	}
	for _, v := range expected {
		res, num, match, text := scn.Scan(), scn.Num(), scn.Match(), scn.Text()

		if res != v.canParse || num != v.num || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
		}
	}
}
func TestOnlyNotMatchScannerFull(t *testing.T) {
	f, err := os.Open("assets/simpleFile.txt")
	defer f.Close()
	if err != nil {
		t.Error(err)
		return
	}

	scn := NewOnlyNotMatchScanner(NewScanner(f))

	expected := []result{
		{false, 11, false, ""},
		{false, 11, false, ""},
		{false, 11, false, ""},
	}
	for _, v := range expected {
		res, num, match, text := scn.Scan(), scn.Num(), scn.Match(), scn.Text()

		if res != v.canParse || num != v.num || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
		}
	}
}

func TestOnlyNotMatchScannerRuled(t *testing.T) {
	f, err := os.Open("assets/simpleFile.txt")
	defer f.Close()
	if err != nil {
		t.Error(err)
		return
	}

	scn := NewOnlyNotMatchScanner(NewRuleScanner(NewScanner(f), func(in string) bool {
		return strings.HasPrefix(in, "#")
	}))

	expected := []result{
		{true, 1, false, "this is a simple file"},
		{true, 2, false, "next line is empty"},
		{true, 3, false, ""},
		{true, 4, false, "next line has two spaces"},
		{true, 5, false, "  "},
		{true, 7, false, "line with two trailing spaces  "},
		{true, 8, false, " line with one leading space"},
		{true, 9, false, " line with one leading and one trailing space "},
		{true, 11, false, "last line"},
		{false, 11, false, ""},
		{false, 11, false, ""},
		{false, 11, false, ""},
	}
	for _, v := range expected {
		res, num, match, text := scn.Scan(), scn.Num(), scn.Match(), scn.Text()

		if res != v.canParse || num != v.num || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
		}
	}
}
