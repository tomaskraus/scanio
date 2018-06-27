package scanio_test

import (
	"bufio"
	"bytes"
	"fmt"
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
		if !asc.Last() {
			fmt.Print(", ")
		}
	}

	// Output:
	// 0:"One", 1:"two", 2:"three"
}

func ExampleNewAheadScanner_consecutive() {
	// Let's find all consecutive sequences of tokens beginning with an "a".
	// Print at what token index these sequences begins and ends and print also a number of tokens in each sequence.
	// All items will be comma-separated but the last one.

	r := strings.NewReader("One apple two amazing apples three ones.")

	// create a rule for a-words
	beginsWithA := func(b []byte) bool {
		return bytes.HasPrefix(b, []byte("a"))
	}

	sc := scanio.NewScanner(r)
	sc.Split(bufio.ScanWords)

	// chain scanners
	asc := scanio.NewAheadScanner(scanio.NewFilterScanner(sc, beginsWithA))

	beginSeq := 0
	for asc.Scan() {
		if asc.BeginConsecutive() {
			beginSeq = asc.Index()
		}
		// there is no "else", as the matching-token-sequence can begin and end at the same token
		if asc.EndConsecutive() {
			fmt.Printf("%v:%v-%v", beginSeq, asc.Index(), asc.NumConsecutive())
			if !asc.Last() {
				fmt.Print(", ")
			}
		}
	}

	// Output:
	// 1:1-1, 3:4-2
}
