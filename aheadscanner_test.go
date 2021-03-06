package scanio_test

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/tomaskraus/scanio"
)

func TestAheadScannerEmpty(t *testing.T) {
	f := strings.NewReader("")

	scn := scanio.NewAheadScanner(scanio.NewScanner(f))

	expected := []resultL{
		{false, 0, false, "", true},
		{false, 0, false, "", true},
		{false, 0, false, "", true},
	}
	for _, v := range expected {
		res, num, isMatch, text, isLast := scn.Scan(), scn.NumRead(), scn.IsMatch(), scn.Text(), scn.IsLast()

		if res != v.canParse || num != v.num || isMatch != v.isMatch || text != v.text || isLast != v.isLast {
			t.Errorf("should be %v, is %v", v, resultL{res, num, isMatch, text, isLast})
		}
	}
}
func TestAheadScannerOneLine(t *testing.T) {
	f := strings.NewReader("one line")

	scn := scanio.NewAheadScanner(scanio.NewScanner(f))

	expected := []resultL{
		{true, 1, true, "one line", true},
		{false, 1, false, "", true},
		{false, 1, false, "", true},
	}
	for _, v := range expected {
		res, num, isMatch, text, isLast := scn.Scan(), scn.NumRead(), scn.IsMatch(), scn.Text(), scn.IsLast()

		if res != v.canParse || num != v.num || isMatch != v.isMatch || text != v.text || isLast != v.isLast {
			t.Errorf("should be %v, is %v", v, resultL{res, num, isMatch, text, isLast})
		}
	}
}

func TestAheadScannerFull(t *testing.T) {
	f, err := os.Open("assets/simpleFile.txt")
	defer f.Close()
	if err != nil {
		t.Error(err)
		return
	}

	scn := scanio.NewAheadScanner(scanio.NewScanner(f))

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
		res, num, isMatch, text, isLast := scn.Scan(), scn.NumRead(), scn.IsMatch(), scn.Text(), scn.IsLast()

		if res != v.canParse || num != v.num || isMatch != v.isMatch || text != v.text || isLast != v.isLast {
			t.Errorf("should be %v, is %v", v, resultL{res, num, isMatch, text, isLast})
		}
	}
}

func TestlastConsecutiveMatch(t *testing.T) {

	const input = "34 235 1234 5678 123456 145 1 2 15678 123"
	lastConsecutives := []bool{false, false, false, false, false, false, false, false, false, false}
	// let's filter words beginning with "1"
	scanner := scanio.NewAheadScanner(
		scanio.NewScanner(strings.NewReader(input)))

	// Set the split function for the scanning operation.
	scanner.Split(bufio.ScanWords)

	for _, v := range lastConsecutives {
		scanner.Scan()
		if scanner.IsConsecutiveEnd() != v {
			t.Errorf("at %d: Want %v, is %v\n", scanner.NumRead(), v, scanner.IsConsecutiveEnd())
		}
	}
}

func TestLastConsecutiveMatchRuled(t *testing.T) {

	const input = "34 235 1234 5678 123456 145 1 2 15678 123"
	lastConsecutives := []bool{false, false, true, false, false, false, true, false, false, true}
	// let's filter words beginning with "1"
	scanner := scanio.NewAheadScanner(
		scanio.NewRuleScanner(
			scanio.NewScanner(strings.NewReader(input)),
			func(input []byte) (bool, error) {
				return (input[0] == []byte("1")[0]), nil
			}))

	// Set the split function for the scanning operation.
	scanner.Split(bufio.ScanWords)

	for _, v := range lastConsecutives {
		scanner.Scan()
		if scanner.IsConsecutiveEnd() != v {
			t.Errorf("at %d: Want %v, is %v\n", scanner.NumRead(), v, scanner.IsConsecutiveEnd())
		}
	}
}

func TestLastConsecutiveMatchFilter(t *testing.T) {

	const input = "34 235 1234 5678 123456 145 1 2 15678 123"
	lastConsecutives := []bool{true, false, false, true, false, true}
	// let's filter words beginning with "1"
	scanner := scanio.NewAheadScanner(
		scanio.NewFilterScanner(
			scanio.NewScanner(strings.NewReader(input)),
			func(input []byte) (bool, error) {
				return (input[0] == []byte("1")[0]), nil
			}))

	// Set the split function for the scanning operation.
	scanner.Split(bufio.ScanWords)

	for _, v := range lastConsecutives {
		scanner.Scan()
		if scanner.IsConsecutiveEnd() != v {
			t.Errorf("at %d: Want %v, is %v\n", scanner.NumRead(), v, scanner.IsConsecutiveEnd())
		}
	}
}
