package scanio

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

type result struct {
	canParse bool
	lineNum  int
	match    bool
	text     string
}

type resultL struct {
	canParse bool
	lineNum  int
	match    bool
	text     string
	last     bool
}

// func TestNewFromScanner(t *testing.T) {
// 	f := strings.NewReader("This is example  ")
// 	sc := bufio.NewScanner(f)
// 	sc.Split(bufio.ScanWords)
// 	li := NewFromScanner(sc)

// 	expected := []result{
// 		{true, 1, true, "This"},
// 		{true, 2, true, "is"},
// 		{true, 3, true, "example"},
// 		{false, 3, false, ""},
// 		{false, 3, false, ""},
// 		{false, 3, false, ""},
// 		{false, 3, false, ""},
// 	}
// 	for _, v := range expected {
// 		res, num, match, text := li.Scan(), li.LineNum(), li.Match(), li.Text()

// 		if res != v.canParse || num != v.lineNum || match != v.match || text != v.text {
// 			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
// 		}
// 	}
// }

func TestReaderScannerEmpty(t *testing.T) {
	f := strings.NewReader("")

	li := NewScanner(f)

	expected := []result{
		{false, 0, false, ""},
		{false, 0, false, ""},
		{false, 0, false, ""},
		{false, 0, false, ""},
	}
	for _, v := range expected {
		res, num, match, text := li.Scan(), li.LineNum(), li.Match(), li.Text()

		if res != v.canParse || num != v.lineNum || match != v.match || text != v.text {
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

	li := NewScanner(f)

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
		res, num, match, text := li.Scan(), li.LineNum(), li.Match(), li.Text()

		if res != v.canParse || num != v.lineNum || match != v.match || text != v.text {
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
	li := NewRuled(NewScanner(f), func(in string) bool {
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
		res, num, match, text := li.Scan(), li.LineNum(), li.Match(), li.Text()

		if res != v.canParse || num != v.lineNum || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
		}
	}
}

func TestRuleScannerEmpty(t *testing.T) {
	f := strings.NewReader("")

	li := NewRuled(NewScanner(f), func(in string) bool {
		return strings.HasPrefix(in, "#")
	})

	expected := []result{
		{false, 0, false, ""},
		{false, 0, false, ""},
		{false, 0, false, ""},
	}
	for _, v := range expected {
		res, num, match, text := li.Scan(), li.LineNum(), li.Match(), li.Text()

		if res != v.canParse || num != v.lineNum || match != v.match || text != v.text {
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

	li := NewRuled(NewScanner(f), func(in string) bool {
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
		res, num, match, text := li.Scan(), li.LineNum(), li.Match(), li.Text()

		if res != v.canParse || num != v.lineNum || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
		}
	}
}

func TestOnlyMatchScannerEmpty(t *testing.T) {
	f := strings.NewReader("")

	li := NewOnlyMatch(NewScanner(f))

	expected := []result{
		{false, 0, false, ""},
		{false, 0, false, ""},
		{false, 0, false, ""},
	}
	for _, v := range expected {
		res, num, match, text := li.Scan(), li.LineNum(), li.Match(), li.Text()

		if res != v.canParse || num != v.lineNum || match != v.match || text != v.text {
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

	li := NewOnlyMatch(NewScanner(f))

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
		res, num, match, text := li.Scan(), li.LineNum(), li.Match(), li.Text()

		if res != v.canParse || num != v.lineNum || match != v.match || text != v.text {
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

	li := NewOnlyMatch(NewRuled(NewScanner(f), func(in string) bool {
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
		res, num, match, text := li.Scan(), li.LineNum(), li.Match(), li.Text()

		if res != v.canParse || num != v.lineNum || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
		}
	}
}

func TestLastScannerEmpty(t *testing.T) {
	f := strings.NewReader("")

	li := NewLast(NewScanner(f))

	expected := []resultL{
		{false, 0, false, "", true},
		{false, 0, false, "", true},
		{false, 0, false, "", true},
	}
	for _, v := range expected {
		res, num, match, text, last := li.Scan(), li.LineNum(), li.Match(), li.Text(), li.Last()

		if res != v.canParse || num != v.lineNum || match != v.match || text != v.text || last != v.last {
			t.Errorf("should be %v, is %v", v, resultL{res, num, match, text, last})
		}
	}
}
func TestLastScannerOneLine(t *testing.T) {
	f := strings.NewReader("one line")

	li := NewLast(NewScanner(f))

	expected := []resultL{
		{true, 1, true, "one line", true},
		{false, 1, false, "", true},
		{false, 1, false, "", true},
	}
	for _, v := range expected {
		res, num, match, text, last := li.Scan(), li.LineNum(), li.Match(), li.Text(), li.Last()

		if res != v.canParse || num != v.lineNum || match != v.match || text != v.text || last != v.last {
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

	li := NewLast(NewScanner(f))

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
		res, num, match, text, last := li.Scan(), li.LineNum(), li.Match(), li.Text(), li.Last()

		if res != v.canParse || num != v.lineNum || match != v.match || text != v.text || last != v.last {
			t.Errorf("should be %v, is %v", v, resultL{res, num, match, text, last})
		}
	}
}

func ExampleNewRuled() {
	f := strings.NewReader("\n# comment 1\n  \n#comment2\nsomething")

	li := NewLast(
		NewOnlyMatch(
			NewRuled(
				NewScanner(f),
				func(s string) bool {
					return strings.HasPrefix(s, "#")
				})))

	for li.Scan() {
		if li.Last() {
			fmt.Printf("(%d, %q).", li.LineNum(), li.Text())
		} else {
			fmt.Printf("(%d, %q), ", li.LineNum(), li.Text())
		}
	}
	// Output:
	// (2, "# comment 1"), (4, "#comment2").
}

func ExampleNewFilter() {
	f := strings.NewReader("\n# comment 1\n  \n#comment2\nsomething")

	li := NewLast(
		NewFilter(
			NewScanner(f),
			func(s string) bool {
				return strings.HasPrefix(s, "#")
			}))

	for li.Scan() {
		if li.Last() {
			fmt.Printf("(%d, %q).", li.LineNum(), li.Text())
		} else {
			fmt.Printf("(%d, %q), ", li.LineNum(), li.Text())
		}
	}
	// Output:
	// (2, "# comment 1"), (4, "#comment2").
}

func TestOnlyNotMatchScannerEmpty(t *testing.T) {
	f := strings.NewReader("")

	li := NewOnlyNotMatch(NewScanner(f))

	expected := []result{
		{false, 0, false, ""},
		{false, 0, false, ""},
		{false, 0, false, ""},
	}
	for _, v := range expected {
		res, num, match, text := li.Scan(), li.LineNum(), li.Match(), li.Text()

		if res != v.canParse || num != v.lineNum || match != v.match || text != v.text {
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

	li := NewOnlyNotMatch(NewScanner(f))

	expected := []result{
		{false, 11, false, ""},
		{false, 11, false, ""},
		{false, 11, false, ""},
	}
	for _, v := range expected {
		res, num, match, text := li.Scan(), li.LineNum(), li.Match(), li.Text()

		if res != v.canParse || num != v.lineNum || match != v.match || text != v.text {
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

	li := NewOnlyNotMatch(NewRuled(NewScanner(f), func(in string) bool {
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
		res, num, match, text := li.Scan(), li.LineNum(), li.Match(), li.Text()

		if res != v.canParse || num != v.lineNum || match != v.match || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, match, text})
		}
	}
}
