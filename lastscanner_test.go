package scanio

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

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

func ExampleNewRuleScanner() {
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

func ExampleNewFilterScanner() {
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

func ExampleNewByteFilterScanner() {

	const input = "1234 5678 123456"
	// let's filter words beginning with "1"
	scanner := NewByteFilterScanner(
		NewScanner(strings.NewReader(input)),
		func(input []byte) bool {
			return (input[0] == []byte("1")[0])
		})

	// Set the split function for the scanning operation.
	scanner.Split(bufio.ScanWords)
	// Validate the input
	for scanner.Scan() {
		fmt.Printf("%s\n", scanner.Bytes())
	}

	// Output:
	// 1234
	// 123456
}

// -----------------------------------------------------------------------

func ExampleLastScanner_NumConsecutive() {

	const input = "34 235 1234 5678 123456 145 1 2 15678 123"
	// let's match words beginning with "1"
	sc := NewLastScanner(
		NewByteRuleScanner(
			NewScanner(strings.NewReader(input)),
			func(input []byte) bool {
				return (input[0] == []byte("1")[0])
			}))

	// Set the split function for the scanning operation.
	sc.Split(bufio.ScanWords)

	for sc.Scan() {
		fmt.Printf("%v: %q, %v, %v, %d\n", sc.Num(), sc.Bytes(), sc.BeginConsecutive(), sc.EndConsecutive(), sc.NumConsecutive())
	}
	fmt.Printf("%v: %q, %v, %v, %d\n", sc.Num(), sc.Bytes(), sc.BeginConsecutive(), sc.EndConsecutive(), sc.NumConsecutive())

	// Output:
	// 1: "34", false, false, 0
	// 2: "235", false, false, 0
	// 3: "1234", true, true, 1
	// 4: "5678", false, false, 0
	// 5: "123456", true, false, 1
	// 6: "145", false, false, 2
	// 7: "1", false, true, 3
	// 8: "2", false, false, 0
	// 9: "15678", true, false, 1
	// 10: "123", false, true, 2
	// 10: "", false, false, 0
}

func ExampleLastScanner_NumConsecutive2() {

	const input = "34 235 1234"
	// let's match words beginning with "1"
	sc := NewLastScanner(
		NewScanner(strings.NewReader(input)),
	)

	// Set the split function for the scanning operation.
	sc.Split(bufio.ScanWords)

	for sc.Scan() {
		fmt.Printf("%v: %q, %v, %v, %d\n", sc.Num(), sc.Bytes(), sc.BeginConsecutive(), sc.EndConsecutive(), sc.NumConsecutive())
	}
	fmt.Printf("%v: %q, %v, %v, %d\n", sc.Num(), sc.Bytes(), sc.BeginConsecutive(), sc.EndConsecutive(), sc.NumConsecutive())

	// Output:
	// 1: "34", true, false, 1
	// 2: "235", false, false, 2
	// 3: "1234", false, true, 3
	// 3: "", false, false, 0
}

func lastConsecutiveMatch(t *testing.T) {

	const input = "34 235 1234 5678 123456 145 1 2 15678 123"
	lastConsecutives := []bool{false, false, false, false, false, false, false, false, false, false}
	// let's filter words beginning with "1"
	scanner := NewLastScanner(
		NewScanner(strings.NewReader(input)))

	// Set the split function for the scanning operation.
	scanner.Split(bufio.ScanWords)

	for _, v := range lastConsecutives {
		scanner.Scan()
		if scanner.EndConsecutive() != v {
			t.Errorf("at %d: Want %v, is %v\n", scanner.Num(), v, scanner.EndConsecutive())
		}
	}
}

func TestLastConsecutiveMatchRuled(t *testing.T) {

	const input = "34 235 1234 5678 123456 145 1 2 15678 123"
	lastConsecutives := []bool{false, false, true, false, false, false, true, false, false, true}
	// let's filter words beginning with "1"
	scanner := NewLastScanner(
		NewByteRuleScanner(
			NewScanner(strings.NewReader(input)),
			func(input []byte) bool {
				return (input[0] == []byte("1")[0])
			}))

	// Set the split function for the scanning operation.
	scanner.Split(bufio.ScanWords)

	for _, v := range lastConsecutives {
		scanner.Scan()
		if scanner.EndConsecutive() != v {
			t.Errorf("at %d: Want %v, is %v\n", scanner.Num(), v, scanner.EndConsecutive())
		}
	}
}

func TestLastConsecutiveMatchFilter(t *testing.T) {

	const input = "34 235 1234 5678 123456 145 1 2 15678 123"
	lastConsecutives := []bool{true, false, false, true, false, true}
	// let's filter words beginning with "1"
	scanner := NewLastScanner(
		NewByteFilterScanner(
			NewScanner(strings.NewReader(input)),
			func(input []byte) bool {
				return (input[0] == []byte("1")[0])
			}))

	// Set the split function for the scanning operation.
	scanner.Split(bufio.ScanWords)

	for _, v := range lastConsecutives {
		scanner.Scan()
		if scanner.EndConsecutive() != v {
			t.Errorf("at %d: Want %v, is %v\n", scanner.Num(), v, scanner.EndConsecutive())
		}
	}
}
