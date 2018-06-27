package scanio_test

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/tomaskraus/scanio"
)

func Example() {
	r := strings.NewReader("One two three")
	sc := scanio.NewScanner(r)
	sc.Split(bufio.ScanWords)

	// chain the next scanner
	asc := scanio.NewAheadScanner(sc)

	for asc.Scan() {
		if asc.Last() {
			fmt.Printf("%v:%q!", asc.Num(), asc.Text())
		} else {
			fmt.Printf("%v:%q, ", asc.Num(), asc.Text())
		}
	}

	// Output:
	// 1:"One", 2:"two", 3:"three"!
}
