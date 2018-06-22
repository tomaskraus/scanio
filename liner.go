// Package liner provides se of text-line scanning utilities.
//
package liner

import (
	"bufio"
	"io"
)

// Liner interface for reading text lines.
// Scan, Err and Text are methods similar to the bufio.Scanner.
// Provides an additional information - current line number etc.
// Text returns an original current line content
type Liner interface {
	Scan() bool
	Err() error
	Match() bool
	Text() string
	Original() string // original string of current line
	Number() int      // number of a current line
}

// reader Liner
type readerLiner struct {
	sc       *bufio.Scanner
	original string
	text     string
	number   int
	match    bool
}

// NewLiner scans from an io.Reader
func NewLiner(r io.Reader) Liner {
	return Liner(&readerLiner{
		sc: bufio.NewScanner(r),
	})
}

func (rs *readerLiner) Scan() bool {
	rs.match = false
	result := rs.sc.Scan()
	if result {
		rs.number++
		rs.match = true
	}
	return result
}

func (rs *readerLiner) Err() error {
	return rs.sc.Err()
}
func (rs *readerLiner) Match() bool {
	return rs.match
}

func (rs *readerLiner) Text() string {
	return rs.sc.Text()
}

func (rs *readerLiner) Original() string {
	return rs.sc.Text()
}

func (rs *readerLiner) Number() int {
	return rs.number
}

// MatchRule for NewFilterLiner
type MatchRule func(input string) (match bool, text string)

type filterLiner struct {
	Liner
	rule   MatchRule
	matchf bool
	textf  string
}

// NewFilterLiner returns new, rule-based Liner
func NewFilterLiner(sc Liner, rule MatchRule) Liner {
	return Liner(&filterLiner{
		sc,
		rule,
		false,
		"",
	})
}

func (fsc *filterLiner) Scan() bool {
	scanResult := fsc.Liner.Scan()
	if scanResult {
		fsc.matchf, fsc.textf = fsc.rule(fsc.Liner.Text())
	} else {
		fsc.matchf, fsc.textf = false, ""
	}
	return scanResult
}

func (fsc *filterLiner) Match() bool {
	return fsc.matchf
}

func (fsc *filterLiner) Text() string {
	return fsc.textf
}

type matchLiner struct {
	Liner
}

// NewOnlyMatchLiner returns only lines that underlying Liner matches
func NewMatchLiner(sc Liner) Liner {
	return Liner(&matchLiner{
		sc,
	})
}

func (fsc *matchLiner) Scan() bool {
	for fsc.Liner.Scan() {
		if fsc.Liner.Match() != true {
			continue
		} else {
			return true
		}
	}
	return false
}
