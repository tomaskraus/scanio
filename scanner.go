// Package scanio provides se of text-line scanning utilities.
//
package scanio

import (
	"bufio"
	"io"
)

// Scanner interface for reading text lines.
// Scan, Err and Text are methods similar to the bufio.Scanner.
// Provides an additional information - current line number etc.
// Text returns current line content
type Scanner interface {
	Scan() bool
	Err() error
	Text() string
	Bytes() []byte
	Match() bool  // true if a line matches Scanner's MatchRule
	LineNum() int // number of a current line
}

// reader Scanner
type readerScanner struct {
	sc      *bufio.Scanner
	match   bool
	lineNum int
}

// New creates a new Scanner, using a scanner.
func New(r io.Reader) Scanner {
	return NewFromScanner(bufio.NewScanner(r))
}

// NewFromScanner creates a new Scanner, using a scanner.
func NewFromScanner(sc *bufio.Scanner) Scanner {
	return Scanner(&readerScanner{
		sc: sc,
	})
}

func (rli *readerScanner) Scan() bool {
	if rli.sc.Scan() {
		rli.lineNum++
		rli.match = true
		return true
	}
	rli.match = false
	return false
}

func (rli *readerScanner) Err() error {
	return rli.sc.Err()
}

func (rli *readerScanner) Text() string {
	return rli.sc.Text()
}

func (rli *readerScanner) Bytes() []byte {
	return rli.sc.Bytes()
}

func (rli *readerScanner) Match() bool {
	return rli.match
}

func (rli *readerScanner) LineNum() int {
	return rli.lineNum
}

// MatchRule for NewRuled.
type MatchRule func(input string) bool

type ruleScanner struct {
	Scanner
	rule   MatchRule
	matchr bool
}

// NewRuled returns new, rule-based Scanner.
func NewRuled(li Scanner, rule MatchRule) Scanner {
	return Scanner(&ruleScanner{
		Scanner: li,
		rule:    rule,
	})
}

func (rli *ruleScanner) Scan() bool {
	if rli.Scanner.Scan() {
		rli.matchr = rli.rule(rli.Scanner.Text())
		return true
	}
	return false
}

func (rli *ruleScanner) Match() bool {
	return rli.matchr
}

type onlyMatchScanner struct {
	Scanner
}

// NewOnlyMatch returns new Scanner.
func NewOnlyMatch(li Scanner) Scanner {
	return Scanner(&onlyMatchScanner{
		Scanner: li,
	})
}

func (omli *onlyMatchScanner) Scan() bool {
	for omli.Scanner.Scan() {
		if omli.Scanner.Match() {
			return true
		}
		continue
	}
	return false
}

type onlyNotMatchScanner struct {
	Scanner
}

// NewOnlyNotMatch returns new Scanner.
func NewOnlyNotMatch(li Scanner) Scanner {
	return Scanner(&onlyNotMatchScanner{
		Scanner: li,
	})
}

func (omli *onlyNotMatchScanner) Scan() bool {
	for omli.Scanner.Scan() {
		if omli.Scanner.Match() {
			continue
		}
		return true
	}
	return false
}

// info Scanner info.
type info struct {
	Text    string
	Bytes   []byte
	LineNum int
	Match   bool
}

const infoBufferCap = 1024

func newInfo(bufferCap int) *info {
	i := info{}
	i.Bytes = make([]byte, bufferCap)
	return &i
}

// updateInfo updates an Info, reflecting current state of a Scanner.
func (info *info) update(li Scanner) {
	info.Text, info.LineNum, info.Match = li.Text(), li.LineNum(), li.Match()
	// copy slices
	length := copy(info.Bytes, li.Bytes())
	info.Bytes = info.Bytes[:length]
}

// LastScanner knows if the current line is the last one.
type LastScanner interface {
	Scanner
	Last() bool
}

type lastScanner struct {
	Scanner
	info, nextInfo *info
	scan, nextScan bool
	last           bool // true if the current line is the last one
	started        bool
}

// NewLast creates a new LastScanner using a Scanner.
func NewLast(li Scanner) LastScanner {
	return LastScanner(&lastScanner{
		Scanner:  li,
		info:     newInfo(infoBufferCap),
		nextInfo: newInfo(infoBufferCap),
		last:     false,
	})
}

func (lli *lastScanner) Scan() bool {
	if !lli.started {
		lli.scan = lli.Scanner.Scan()
		lli.info.update(lli.Scanner)
		lli.nextScan = lli.Scanner.Scan()
		lli.nextInfo.update(lli.Scanner)
		lli.started = true
		return lli.scan
	}
	lli.info, lli.nextInfo = lli.nextInfo, lli.info
	lli.scan, lli.nextScan = lli.nextScan, lli.Scanner.Scan()
	lli.nextInfo.update(lli.Scanner)
	return lli.scan
}

func (lli *lastScanner) Text() string {
	return lli.info.Text
}

func (lli *lastScanner) Bytes() []byte {
	return lli.info.Bytes
}

func (lli *lastScanner) LineNum() int {
	return lli.info.LineNum
}
func (lli *lastScanner) Match() bool {
	return lli.info.Match
}

func (lli *lastScanner) Last() bool {
	return lli.nextScan == false
}

// NewFilter creates a Scanner that produces only lines matched by a rule provided.
func NewFilter(li Scanner, rule MatchRule) Scanner {
	return NewOnlyMatch(NewRuled(li, rule))
}
