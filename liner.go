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
// Text returns current line content
type Liner interface {
	Scan() bool
	Err() error
	Text() string
	Match() bool  // true if a line matches Liner's MatchRule
	LineNum() int // number of a current line
}

// reader Liner
type readerLiner struct {
	sc      *bufio.Scanner
	text    string
	match   bool
	lineNum int
}

// New creates a new Liner, using a scanner.
func New(r io.Reader) Liner {
	return NewFromScanner(bufio.NewScanner(r))
}

// NewFromScanner creates a new Liner, using a scanner.
func NewFromScanner(sc *bufio.Scanner) Liner {
	return Liner(&readerLiner{
		sc: sc,
	})
}

func (rli *readerLiner) Scan() bool {
	rli.match = false
	if rli.sc.Scan() {
		rli.lineNum++
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

func (rli *readerLiner) LineNum() int {
	return rli.lineNum
}

// MatchRule for NewRuled.
type MatchRule func(input string) bool

type ruleLiner struct {
	Liner
	rule   MatchRule
	matchr bool
}

// NewRuled returns new, rule-based Liner.
func NewRuled(li Liner, rule MatchRule) Liner {
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

type onlyMatchLiner struct {
	Liner
}

// NewOnlyMatch returns new Liner.
func NewOnlyMatch(li Liner) Liner {
	return Liner(&onlyMatchLiner{
		Liner: li,
	})
}

func (omli *onlyMatchLiner) Scan() bool {
	for omli.Liner.Scan() {
		if omli.Liner.Match() {
			return true
		}
		continue
	}
	return false
}

type onlyNotMatchLiner struct {
	Liner
}

// NewOnlyNotMatch returns new Liner.
func NewOnlyNotMatch(li Liner) Liner {
	return Liner(&onlyNotMatchLiner{
		Liner: li,
	})
}

func (omli *onlyNotMatchLiner) Scan() bool {
	for omli.Liner.Scan() {
		if omli.Liner.Match() {
			continue
		}
		return true
	}
	return false
}

// info Liner info.
type info struct {
	Text    string
	LineNum int
	Match   bool
}

// updateInfo updates an Info, reflecting current state of a Liner.
func updateInfo(info *info, li Liner) {
	info.Text, info.LineNum, info.Match = li.Text(), li.LineNum(), li.Match()
}

// LastLiner knows if the current line is the last one.
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

// NewLast creates a new LastLiner using a Liner.
func NewLast(li Liner) LastLiner {
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

func (lli *lastLiner) LineNum() int {
	return lli.info.LineNum
}
func (lli *lastLiner) Match() bool {
	return lli.info.Match
}

func (lli *lastLiner) Last() bool {
	return lli.currentScan == false
}

// NewFilter creates a Liner that produces only lines matched by a rule provided.
func NewFilter(li Liner, rule MatchRule) Liner {
	return NewOnlyMatch(NewRuled(li, rule))
}
