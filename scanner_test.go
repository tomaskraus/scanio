package scanio_test

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/tomaskraus/scanio"
)

type result struct {
	canParse bool
	num      int
	isMatch  bool
	text     string
}

type resultL struct {
	canParse bool
	num      int
	isMatch  bool
	text     string
	isLast   bool
}

func TestScanWords(t *testing.T) {
	f := strings.NewReader("This is example  ")
	scn := scanio.NewScanner(f)
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
		res, num, isMatch, text := scn.Scan(), scn.NumRead(), scn.IsMatch(), scn.Text()

		if res != v.canParse || num != v.num || isMatch != v.isMatch || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, isMatch, text})
		}
	}
}

func TestReaderScannerEmpty(t *testing.T) {
	f := strings.NewReader("")

	scn := scanio.NewScanner(f)

	expected := []result{
		{false, 0, false, ""},
		{false, 0, false, ""},
		{false, 0, false, ""},
		{false, 0, false, ""},
	}
	for _, v := range expected {
		res, num, isMatch, text := scn.Scan(), scn.NumRead(), scn.IsMatch(), scn.Text()

		if res != v.canParse || num != v.num || isMatch != v.isMatch || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, isMatch, text})
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

	scn := scanio.NewScanner(f)

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
		res, num, isMatch, text := scn.Scan(), scn.NumRead(), scn.IsMatch(), scn.Text()

		if res != v.canParse || num != v.num || isMatch != v.isMatch || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, isMatch, text})
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
	scn := scanio.NewRuleScanner(scanio.NewScanner(f), func(b []byte) (bool, error) {
		return bytes.HasPrefix(b, []byte("#")), nil
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
		res, num, isMatch, text := scn.Scan(), scn.NumRead(), scn.IsMatch(), scn.Text()

		if res != v.canParse || num != v.num || isMatch != v.isMatch || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, isMatch, text})
		}
	}
}

func TestRuleScannerEmpty(t *testing.T) {
	f := strings.NewReader("")

	scn := scanio.NewRuleScanner(scanio.NewScanner(f), func(b []byte) (bool, error) {
		return bytes.HasPrefix(b, []byte("#")), nil
	})

	expected := []result{
		{false, 0, false, ""},
		{false, 0, false, ""},
		{false, 0, false, ""},
	}
	for _, v := range expected {
		res, num, isMatch, text := scn.Scan(), scn.NumRead(), scn.IsMatch(), scn.Text()

		if res != v.canParse || num != v.num || isMatch != v.isMatch || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, isMatch, text})
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

	scn := scanio.NewRuleScanner(scanio.NewScanner(f), func(b []byte) (bool, error) {
		return bytes.HasPrefix(b, []byte("#")), nil
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
		res, num, isMatch, text := scn.Scan(), scn.NumRead(), scn.IsMatch(), scn.Text()

		if res != v.canParse || num != v.num || isMatch != v.isMatch || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, isMatch, text})
		}
	}
}

func TestOnlyMatchScannerEmpty(t *testing.T) {
	f := strings.NewReader("")

	scn := scanio.NewOnlyMatchScanner(scanio.NewScanner(f))

	expected := []result{
		{false, 0, false, ""},
		{false, 0, false, ""},
		{false, 0, false, ""},
	}
	for _, v := range expected {
		res, num, isMatch, text := scn.Scan(), scn.NumRead(), scn.IsMatch(), scn.Text()

		if res != v.canParse || num != v.num || isMatch != v.isMatch || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, isMatch, text})
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

	scn := scanio.NewOnlyMatchScanner(scanio.NewScanner(f))

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
		res, num, isMatch, text := scn.Scan(), scn.NumRead(), scn.IsMatch(), scn.Text()

		if res != v.canParse || num != v.num || isMatch != v.isMatch || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, isMatch, text})
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

	scn := scanio.NewOnlyMatchScanner(scanio.NewRuleScanner(scanio.NewScanner(f), func(b []byte) (bool, error) {
		return bytes.HasPrefix(b, []byte("#")), nil
	}))

	expected := []result{
		{true, 6, true, "# bash-like comment"},
		{true, 10, true, "# bash-like comment 2 "},
		{false, 11, false, ""},
		{false, 11, false, ""},
		{false, 11, false, ""},
	}
	for _, v := range expected {
		res, num, isMatch, text := scn.Scan(), scn.NumRead(), scn.IsMatch(), scn.Text()

		if res != v.canParse || num != v.num || isMatch != v.isMatch || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, isMatch, text})
		}
	}
}

func TestOnlyNotMatchScannerEmpty(t *testing.T) {
	f := strings.NewReader("")

	scn := scanio.NewOnlyNotMatchScanner(scanio.NewScanner(f))

	expected := []result{
		{false, 0, false, ""},
		{false, 0, false, ""},
		{false, 0, false, ""},
	}
	for _, v := range expected {
		res, num, isMatch, text := scn.Scan(), scn.NumRead(), scn.IsMatch(), scn.Text()

		if res != v.canParse || num != v.num || isMatch != v.isMatch || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, isMatch, text})
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

	scn := scanio.NewOnlyNotMatchScanner(scanio.NewScanner(f))

	expected := []result{
		{false, 11, false, ""},
		{false, 11, false, ""},
		{false, 11, false, ""},
	}
	for _, v := range expected {
		res, num, isMatch, text := scn.Scan(), scn.NumRead(), scn.IsMatch(), scn.Text()

		if res != v.canParse || num != v.num || isMatch != v.isMatch || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, isMatch, text})
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

	scn := scanio.NewOnlyNotMatchScanner(scanio.NewRuleScanner(scanio.NewScanner(f), func(b []byte) (bool, error) {
		return bytes.HasPrefix(b, []byte("#")), nil
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
		res, num, isMatch, text := scn.Scan(), scn.NumRead(), scn.IsMatch(), scn.Text()

		if res != v.canParse || num != v.num || isMatch != v.isMatch || text != v.text {
			t.Errorf("should be %v, is %v", v, result{res, num, isMatch, text})
		}
	}
}

func ExampleNewFilterScanner() {

	const input = "1234 5678 123456"
	// let's filter words beginning with "1"
	scanner := scanio.NewFilterScanner(
		scanio.NewScanner(strings.NewReader(input)),
		func(input []byte) (bool, error) {
			return (input[0] == []byte("1")[0]), nil
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

func ExampleNewFilterScanner_smallBuffer() {

	const input = "1 1234 5678 123456"
	// let's filter words beginning with "1"
	scanner := scanio.NewFilterScanner(
		scanio.NewScanner(strings.NewReader(input)),
		func(input []byte) (bool, error) {
			return (input[0] == []byte("1")[0]), nil
		})

	// Set the split function for the scanning operation.
	scanner.Split(bufio.ScanWords)
	//set buffer too small
	scanner.Buffer(make([]byte, 2), 2)
	// Validate the input
	for scanner.Scan() {
		fmt.Printf("%s\n", scanner.Bytes())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	// Output:
	// 1
	// bufio.Scanner: token too long
}
