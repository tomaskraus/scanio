package scanio

import (
	"bufio"
	"bytes"
	"fmt"
)

func Example() {
	r := bytes.NewReader([]byte("abcd ef gh"))
	sc := NewScanner(r)
	sc.Split(bufio.ScanWords)

	sc = NewLastScanner(sc)
	buf := make([]byte, 5)
	sc.Buffer(buf, 2)

	scanned := false
	for scanned = sc.Scan(); scanned == true; scanned = sc.Scan() {
		fmt.Printf("%v, %v, %v, %q\n", scanned, sc.Num(), sc.Match(), sc.Bytes())
	}
	fmt.Printf("%v, %v, %v, %q\n", scanned, sc.Num(), sc.Match(), sc.Bytes())

	// Output:
	// true, 1, true, "abcd"
	// true, 2, true, "ef"
	// true, 3, true, "gh"
	// false, 3, false, ""
}
