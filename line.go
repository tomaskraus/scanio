// Package line provides se of text-line scanning utilities.
//
package line

import (
	"bufio"
	"io"
)

// Scanner interface for reading text lines.
// Scan, Err and Text are methods similar to the bufio.Scanner.
// Provides an additional information - current line number etc.
// Text returns an original current line content
type Scanner interface {
	Scan() bool
	Err() error
	Match() bool
	Text() string
	Original() string // original string of current line
	Number() int      // number of a current line
}

// reader scanner
type readerScanner struct {
	sc       *bufio.Scanner
	original string
	text     string
	number   int
	match    bool
}

// NewReaderScanner scans from an io.Reader
func NewReaderScanner(r io.Reader) Scanner {
	return Scanner(&readerScanner{
		sc: bufio.NewScanner(r),
	})
}

func (rs *readerScanner) Scan() bool {
	rs.match = false
	result := rs.sc.Scan()
	if result {
		rs.number++
		rs.match = true
	}
	return result
}

func (rs *readerScanner) Err() error {
	return rs.sc.Err()
}
func (rs *readerScanner) Match() bool {
	return rs.match
}

func (rs *readerScanner) Text() string {
	return rs.sc.Text()
}

func (rs *readerScanner) Original() string {
	return rs.sc.Text()
}

func (rs *readerScanner) Number() int {
	return rs.number
}

// MatchRule for NewFilterScanner
type MatchRule func(input string) (match bool, text string)

type filterScanner struct {
	Scanner
	rule   MatchRule
	matchf bool
	textf  string
}

// NewFilterScanner returns new, rule-based scanner
func NewFilterScanner(sc Scanner, rule MatchRule) Scanner {
	return Scanner(&filterScanner{
		sc,
		rule,
		false,
		"",
	})
}

func (fsc *filterScanner) Scan() bool {
	scanResult := fsc.Scanner.Scan()
	if scanResult {
		fsc.matchf, fsc.textf = fsc.rule(fsc.Scanner.Text())
	} else {
		fsc.matchf, fsc.textf = false, ""
	}
	return scanResult
}

func (fsc *filterScanner) Match() bool {
	return fsc.matchf
}

func (fsc *filterScanner) Text() string {
	return fsc.textf
}

type onlyMatchScanner struct {
	Scanner
}

// NewOnlyMatchScanner returns only lines that underlying scanner matches
func NewOnlyMatchScanner(sc Scanner) Scanner {
	return Scanner(&onlyMatchScanner{
		sc,
	})
}

func (fsc *onlyMatchScanner) Scan() bool {
	for fsc.Scanner.Scan() {
		if fsc.Scanner.Match() != true {
			continue
		} else {
			return true
		}
	}
	return false
}
