package scanio

import (
	"bufio"
	"fmt"
	"strings"
)

func Example() {
	r := strings.NewReader("One two three")
	sc := NewScanner(r)
	sc.Split(bufio.ScanWords)

	// chain the next scanner
	lsc := NewLastScanner(sc)

	for lsc.Scan() {
		if lsc.Last() {
			fmt.Printf("%v:%q!", lsc.Num(), lsc.Text())
		} else {
			fmt.Printf("%v:%q, ", lsc.Num(), lsc.Text())
		}
	}

	// Output:
	// 1:"One", 2:"two", 3:"three"!
}
