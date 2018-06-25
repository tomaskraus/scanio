package scanio

import (
	"bufio"
	"bytes"
	"fmt"
)

type resultB struct {
	canParse bool
	lineNum  int
	match    bool
	bytes    []byte
}

func Example() {
	r := bytes.NewReader([]byte("abcd ef gh"))
	sc := NewScanner(r)
	buf := make([]byte, 5)
	sc.Buffer(buf, 2)
	sc.Split(bufio.ScanWords)

	li := NewLast(sc)

	scn := false
	for scn = li.Scan(); scn == true; scn = li.Scan() {
		fmt.Printf("%v, %v, %v, %q\n", scn, li.Num(), li.Match(), li.Bytes())
	}
	fmt.Printf("%v, %v, %v, %q\n", scn, li.Num(), li.Match(), li.Bytes())

	// Output:
	// true, 1, true, "abcd"
	// true, 2, true, "ef"
	// true, 3, true, "gh"
	// false, 3, false, ""
}
