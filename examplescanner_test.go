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
	//set buffer too small (here it has a capacity of 2)
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
