package scanio_test

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/tomaskraus/scanio"
)

func ExampleNewAheadScanner_basic() {
	// Do not print trailing comma in a comma-separated list.

	r := strings.NewReader("One two three")
	sc := scanio.NewScanner(r)
	// read whole words
	sc.Split(bufio.ScanWords)

	// chain the next scanner
	asc := scanio.NewAheadScanner(sc)

	for asc.Scan() {
		fmt.Printf("%v:%q", asc.NumRead(), asc.Text())
		if !asc.IsLast() {
			fmt.Print(", ")
		}
	}

	// Output:
	// 1:"One", 2:"two", 3:"three"
}

func ExampleNewAheadScanner_consecutive() {
	// Let's find all consecutive sequences of tokens beginning with an "a".
	// Print ranges of these sequences and print also a number of tokens in each sequence.
	// All items will be semicolon-separated but the last one.

	r := strings.NewReader("One apple two amazing apples three ones.")

	// create a rule for a-words
	beginsWithA := func(b []byte) (bool, error) {
		return bytes.HasPrefix(b, []byte("a")), nil
	}

	sc := scanio.NewScanner(r)
	sc.Split(bufio.ScanWords)

	// chain scanners
	asc := scanio.NewAheadScanner(scanio.NewFilterScanner(sc, beginsWithA))

	beginSeq := 0
	for asc.Scan() {
		if asc.IsConsecutiveBegin() {
			beginSeq = asc.NumRead() - 1
		}
		// there is no "else", as the matching-token-sequence can begin and end at the same token
		if asc.IsConsecutiveEnd() {
			fmt.Printf("[%v:%v],%v", beginSeq,
				asc.NumRead(),
				asc.NumConsecutive(),
			)
			if !asc.IsLast() {
				fmt.Print("; ")
			}
		}
	}

	// Output:
	// [1:2],1; [3:5],2
}

func ExampleNewAheadScanner_error() {
	// Let's filter positive integers
	//
	// MatchRule's error stops the scanning, causes that the last error-free token is treated as the last one.

	r := strings.NewReader("123  -456 5 abc 678 173")
	sc := scanio.NewScanner(r)
	// read whole words
	sc.Split(bufio.ScanWords)

	// MatchRule with error
	isPositiveInt := func(b []byte) (res bool, err error) {
		num, err := strconv.ParseInt(string(b), 0, 0)
		if err != nil {
			return false, err
		}
		return (num > 0), nil
	}

	// chain the next scanner
	asc := scanio.NewAheadScanner(
		scanio.NewFilterScanner(sc, isPositiveInt))

	for asc.Scan() {
		fmt.Printf("%v:%q", asc.NumRead(), asc.Text())
		if !asc.IsLast() {
			fmt.Print(", ")
		}
	}
	if asc.Err() != nil {
		fmt.Printf("\n%v", asc.Err())
	}

	// Output:
	// 1:"123", 3:"5"
	// strconv.ParseInt: parsing "abc": invalid syntax
}

func ExampleNewRuleScanner() {
	f := strings.NewReader("\n# comment 1\n  \n#comment2\nsomething")

	scn := scanio.NewAheadScanner(
		scanio.NewOnlyMatchScanner(
			scanio.NewRuleScanner(
				scanio.NewScanner(f),
				func(b []byte) (bool, error) {
					return bytes.HasPrefix(b, []byte("#")), nil
				})))

	for scn.Scan() {
		if scn.IsLast() {
			fmt.Printf("(%d, %q).", scn.NumRead(), scn.Text())
		} else {
			fmt.Printf("(%d, %q), ", scn.NumRead(), scn.Text())
		}
	}
	// Output:
	// (2, "# comment 1"), (4, "#comment2").
}

func ExampleNewAheadFilterScanner() {
	f := strings.NewReader("\n# comment 1\n  \n#comment2\nsomething")

	scn := scanio.NewAheadScanner(
		scanio.NewFilterScanner(
			scanio.NewScanner(f),
			func(b []byte) (bool, error) {
				return bytes.HasPrefix(b, []byte("#")), nil
			}))

	for scn.Scan() {
		if scn.IsLast() {
			fmt.Printf("(%d, %q).", scn.NumRead(), scn.Text())
		} else {
			fmt.Printf("(%d, %q), ", scn.NumRead(), scn.Text())
		}
	}
	// Output:
	// (2, "# comment 1"), (4, "#comment2").
}

// -----------------------------------------------------------------------

func ExampleAheadScanner_NumConsecutive() {

	const input = "34 235 1234 5678 123456 145 1 2 15678 123"
	// let's match words beginning with "1"
	sc := scanio.NewAheadScanner(
		scanio.NewRuleScanner(
			scanio.NewScanner(strings.NewReader(input)),
			func(input []byte) (bool, error) {
				return (input[0] == []byte("1")[0]), nil
			}))

	// Set the split function for the scanning operation.
	sc.Split(bufio.ScanWords)

	for sc.Scan() {
		fmt.Printf("%v: %q, %v, %v, %d\n", sc.NumRead(), sc.Bytes(), sc.IsConsecutiveBegin(), sc.IsConsecutiveEnd(), sc.NumConsecutive())
	}
	fmt.Printf("%v: %q, %v, %v, %d\n", sc.NumRead(), sc.Bytes(), sc.IsConsecutiveBegin(), sc.IsConsecutiveEnd(), sc.NumConsecutive())

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

func ExampleAheadScanner_NumConsecutive2() {

	const input = "34 235 1234"
	// let's match words beginning with "1"
	sc := scanio.NewAheadScanner(
		scanio.NewScanner(strings.NewReader(input)),
	)

	// Set the split function for the scanning operation.
	sc.Split(bufio.ScanWords)

	for sc.Scan() {
		fmt.Printf("%v: %q, %v, %v, %d\n", sc.NumRead(), sc.Bytes(), sc.IsConsecutiveBegin(), sc.IsConsecutiveEnd(), sc.NumConsecutive())
	}
	fmt.Printf("%v: %q, %v, %v, %d\n", sc.NumRead(), sc.Bytes(), sc.IsConsecutiveBegin(), sc.IsConsecutiveEnd(), sc.NumConsecutive())

	// Output:
	// 1: "34", true, false, 1
	// 2: "235", false, false, 2
	// 3: "1234", false, true, 3
	// 3: "", false, false, 0
}

func ExampleNewAheadFilterScanner_smallBuffer() {

	const input = "1 1234 5678 123456"
	// let's filter words beginning with "1"
	scanner := scanio.NewAheadScanner(
		scanio.NewFilterScanner(
			scanio.NewScanner(strings.NewReader(input)),
			func(input []byte) (bool, error) {
				return (input[0] == []byte("1")[0]), nil
			}))

	// Set the split function for the scanning operation.
	scanner.Split(bufio.ScanWords)
	// set buffer too small (here it has a capacity of 2)
	// ensure the aheadScanner.Buffer() has the same behavior as bufio.Scanner.Buffer()
	scanner.Buffer(make([]byte, 0, 2), 0)
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
