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
	Text() string
	Match() bool // true if a line matches Liner's MatchRule
	Number() int // number of a current line
}

// reader Liner
type readerLiner struct {
	sc       *bufio.Scanner
	text     string
	match    bool
	original string
	number   int
}

// NewLiner scans from an io.Reader
func NewLiner(r io.Reader) Liner {
	return Liner(&readerLiner{
		sc: bufio.NewScanner(r),
	})
}

func (rli *readerLiner) Scan() bool {
	rli.match = false
	if rli.sc.Scan() {
		rli.number++
		rli.match = true
		return true
	}
	return false
}

func (rli *readerLiner) Err() error {
	return rli.sc.Err()
}

func (rli *readerLiner) Text() string {
	return rli.sc.Text()
}

func (rli *readerLiner) Match() bool {
	return rli.match
}

func (rli *readerLiner) Number() int {
	return rli.number
}

// MatchRule for NewRuleLiner
type MatchRule func(input string) bool

type ruleLiner struct {
	Liner
	rule   MatchRule
	matchr bool
}

// NewRuleLiner returns new, rule-based Liner
func NewRuleLiner(li Liner, rule MatchRule) Liner {
	return Liner(&ruleLiner{
		Liner: li,
		rule:  rule,
	})
}

func (rli *ruleLiner) Scan() bool {
	if rli.Liner.Scan() {
		rli.matchr = rli.rule(rli.Liner.Text())
		return true
	}
	return false
}

func (rli *ruleLiner) Match() bool {
	return rli.matchr
}

// info Liner info.
type info struct {
	Text   string
	Number int
	Match  bool
}

// UpdateInfo updates an Info, reflecting current state of a Liner.
func updateInfo(info *info, li Liner) {
	info.Text, info.Number, info.Match = li.Text(), li.Number(), li.Match()
}

// LastLiner knows if the line is the last one.
type LastLiner interface {
	Liner
	Last() bool
}

type lastLiner struct {
	Liner
	info, nextInfo            *info
	last                      bool // true if the current line is the last one
	started                   bool
	previousScan, currentScan bool
}

// NewLastLiner creates a new LastLiner from liner.
func NewLastLiner(li Liner) LastLiner {
	return LastLiner(&lastLiner{
		Liner:    li,
		info:     &info{},
		nextInfo: &info{},
		last:     false,
	})
}

func (lli *lastLiner) Scan() bool {
	if !lli.started {
		lli.previousScan = lli.Liner.Scan()
		updateInfo(lli.info, lli.Liner)
		lli.currentScan = lli.Liner.Scan()
		updateInfo(lli.nextInfo, lli.Liner)
		lli.started = true
		return lli.previousScan
	}
	lli.info, lli.nextInfo = lli.nextInfo, lli.info
	lli.previousScan, lli.currentScan = lli.currentScan, lli.Liner.Scan()
	updateInfo(lli.nextInfo, lli.Liner)
	return lli.previousScan
}

func (lli *lastLiner) Text() string {
	return lli.info.Text
}

func (lli *lastLiner) Number() int {
	return lli.info.Number
}
func (lli *lastLiner) Match() bool {
	return lli.info.Match
}

func (lli *lastLiner) Last() bool {
	return lli.currentScan == false
}
