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
		fmt.Printf("%v:%q", asc.Index(), asc.Text())
		if !asc.IsLast() {
			fmt.Print(", ")
		}
	}

	// Output:
	// 0:"One", 1:"two", 2:"three"
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
			beginSeq = asc.Index()
		}
		// there is no "else", as the matching-token-sequence can begin and end at the same token
		if asc.IsConsecutiveEnd() {
			fmt.Printf("[%v:%v],%v", beginSeq,
				asc.Index()+1, // slice-like range
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
		fmt.Printf("%v:%q", asc.Index(), asc.Text())
		if !asc.IsLast() {
			fmt.Print(", ")
		}
	}
	if asc.Err() != nil {
		fmt.Printf("\n%v", asc.Err())
	}

	// Output:
	// 0:"123", 2:"5"
	// strconv.ParseInt: parsing "abc": invalid syntax
}
