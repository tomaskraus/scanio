package scanio_test

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/tomaskraus/scanio"
)

func ExampleNewFilterScanner_error() {
	// Let's filter positive integers

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
	asc := scanio.NewFilterScanner(sc, isPositiveInt)

	for asc.Scan() {
		fmt.Printf("%v:%q,", asc.NumRead(), asc.Text())
	}
	if asc.Err() != nil {
		fmt.Printf("\n%v", asc.Err())
	}

	// Output:
	// 1:"123",3:"5",
	// strconv.ParseInt: parsing "abc": invalid syntax
}
