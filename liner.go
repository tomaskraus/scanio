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
	Match() bool      // true if a line matches Liner's requirement
	End() bool        // return true if an end of data occured
	Original() string // original string of current line
	Number() int      // number of a current line
}

// reader Liner
type readerLiner struct {
	sc       *bufio.Scanner
	text     string
	match    bool
	end      bool
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
	result := rli.sc.Scan()
	if result {
		rli.number++
		rli.match = true
	} else if rli.sc.Err() == nil {
		rli.end = true
	}
	return result
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

func (rli *readerLiner) End() bool {
	return rli.end
}

func (rli *readerLiner) Original() string {
	return rli.sc.Text()
}

func (rli *readerLiner) Number() int {
	return rli.number
}

// MatchRule for NewMatchLiner
type MatchRule func(input string) (match bool, text string)

type matchLiner struct {
	Liner
	rule   MatchRule
	matchr bool
	textr  string
}

// NewMatchLiner returns new, rule-based Liner
func NewMatchLiner(lin Liner, rule MatchRule) Liner {
	return Liner(&matchLiner{
		lin,
		rule,
		false,
		"",
	})
}

func (fli *matchLiner) Scan() bool {
	scanResult := fli.Liner.Scan()
	if scanResult {
		fli.matchr, fli.textr = fli.rule(fli.Liner.Text())
	} else {
		fli.matchr, fli.textr = false, ""
	}
	return scanResult
}

func (fli *matchLiner) Match() bool {
	return fli.matchr
}

func (fli *matchLiner) Text() string {
	return fli.textr
}

type filterLiner struct {
	Liner
}

// NewFilterLiner returns only lines that the underlying Liner matches
func NewFilterLiner(li Liner) Liner {
	return Liner(&filterLiner{
		li,
	})
}

func (mli *filterLiner) Scan() bool {
	for mli.Liner.Scan() {
		if mli.Liner.Match() != true {
			continue
		} else {
			return true
		}
	}
	return false
}

type noMatchLiner struct {
	Liner
}

// NewNoMatchLiner returns only lines that the underlying Liner does not match
func NewNoMatchLiner(li Liner) Liner {
	return Liner(&noMatchLiner{
		li,
	})
}

func (nli *noMatchLiner) Scan() bool {
	for nli.Liner.Scan() {
		if nli.Liner.Match() {
			continue
		} else {
			return true
		}
	}
	return false
}

func (nli *noMatchLiner) Match() bool {
	if nli.Liner.End() == false {
		return !nli.Liner.Match()
	}
	return false
}

// info Liner info
type info struct {
	Text     string
	Original string
	Number   int
	Match    bool
	End      bool
}

// UpdateInfo updates an Info, reflecting current state of a Liner.
func updateInfo(info *info, li Liner) {
	info.Text, info.Original, info.Number, info.Match, info.End =
		li.Text(),
		li.Original(),
		li.Number(),
		li.Match(),
		li.End()
}

// LastLiner knows if the line is the last one
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

func (lli *lastLiner) Original() string {
	return lli.info.Original
}
func (lli *lastLiner) Number() int {
	return lli.info.Number
}
func (lli *lastLiner) Match() bool {
	return lli.info.Match
}
func (lli *lastLiner) End() bool {
	return lli.info.End
}

func (lli *lastLiner) Last() bool {
	return lli.nextInfo.End
}
