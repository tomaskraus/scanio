// Package scanio provides bufio.Scanner wrappers with filtering, chaining and last-token ahead reading.
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
	Buffer(buf []byte, max int)
	Bytes() []byte
	Err() error
	Scan() bool
	Split(split bufio.SplitFunc)
	Text() string

	Match() bool // true if a line matches Scanner's MatchRule
	Num() int    // number of a current line
}

// reader Scanner
type readerScanner struct {
	scn   *bufio.Scanner
	match bool
	num   int
}

// NewScanner creates a new Scanner, using a scanner.
func NewScanner(r io.Reader) Scanner {
	return Scanner(&readerScanner{
		scn: bufio.NewScanner(r),
	})
}

func (sc *readerScanner) Buffer(buf []byte, max int) {
	sc.scn.Buffer(buf, max)
}

func (sc *readerScanner) Scan() bool {
	if sc.scn.Scan() {
		sc.num++
		sc.match = true
		return true
	}
	sc.match = false
	return false
}

func (sc *readerScanner) Err() error {
	return sc.scn.Err()
}

func (sc *readerScanner) Split(split bufio.SplitFunc) {
	sc.scn.Split(split)
}

func (sc *readerScanner) Text() string {
	return sc.scn.Text()
}

func (sc *readerScanner) Bytes() []byte {
	return sc.scn.Bytes()
}

func (sc *readerScanner) Match() bool {
	return sc.match
}

func (sc *readerScanner) Num() int {
	return sc.num
}

// MatchRule for NewRuleScanner.
type MatchRule func(input string) bool

type ruleScanner struct {
	Scanner
	rule   MatchRule
	matchr bool
}

// NewRuleScanner returns new, rule-based Scanner.
func NewRuleScanner(scn Scanner, rule MatchRule) Scanner {
	return Scanner(&ruleScanner{
		Scanner: scn,
		rule:    rule,
	})
}

func (sc *ruleScanner) Scan() bool {
	if sc.Scanner.Scan() {
		sc.matchr = sc.rule(sc.Scanner.Text())
		return true
	}
	return false
}

func (sc *ruleScanner) Match() bool {
	return sc.matchr
}

type onlyMatchScanner struct {
	Scanner
}

// NewOnlyMatchScanner returns new Scanner.
func NewOnlyMatchScanner(scn Scanner) Scanner {
	return Scanner(&onlyMatchScanner{
		Scanner: scn,
	})
}

func (scn *onlyMatchScanner) Scan() bool {
	for scn.Scanner.Scan() {
		if scn.Scanner.Match() {
			return true
		}
		continue
	}
	return false
}

type onlyNotMatchScanner struct {
	Scanner
}

// NewOnlyNotMatchScanner returns new Scanner.
func NewOnlyNotMatchScanner(scn Scanner) Scanner {
	return Scanner(&onlyNotMatchScanner{
		Scanner: scn,
	})
}

func (scn *onlyNotMatchScanner) Scan() bool {
	for scn.Scanner.Scan() {
		if scn.Scanner.Match() {
			continue
		}
		return true
	}
	return false
}

// ---------------------------------------------------------------------------

const (
	startBufSize = 4096 // Size of initial allocation for buffer.   from golang.org/src/bufio/scan.go
)

// info Scanner info.
type info struct {
	Text  string
	Bytes []byte
	Num   int
	Match bool
}

func newInfo(bufferLen, bufferCap int) *info {
	i := info{}
	i.Bytes = make([]byte, bufferLen, bufferCap)
	return &i
}

// updateInfo updates an Info, reflecting current state of a Scanner.
func (info *info) update(scn Scanner) {
	info.Text, info.Num, info.Match = scn.Text(), scn.Num(), scn.Match()
	// copy slices
	length := copy(info.Bytes, scn.Bytes())
	info.Bytes = info.Bytes[:length]
}

// LastScanner knows if the current line is the last one.
type LastScanner interface {
	Scanner
	Last() bool
}

type lastScanner struct {
	Scanner
	info, nextInfo  *info
	scan, nextScan  bool
	last            bool // true if the current line is the last one
	started         bool
	bufSize, bufCap int
}

// NewLastScanner creates a new LastScanner using a Scanner.
func NewLastScanner(scn Scanner) LastScanner {
	return LastScanner(&lastScanner{
		Scanner: scn,
		last:    false,
		bufSize: startBufSize,
		bufCap:  bufio.MaxScanTokenSize,
	})
}

func (lsc *lastScanner) Scan() bool {
	if !lsc.started {
		//initialize buffers
		lsc.info = newInfo(lsc.bufSize, lsc.bufCap)
		lsc.nextInfo = newInfo(lsc.bufSize, lsc.bufCap)

		lsc.scan = lsc.Scanner.Scan()
		lsc.info.update(lsc.Scanner)
		lsc.nextScan = lsc.Scanner.Scan()
		lsc.nextInfo.update(lsc.Scanner)
		lsc.started = true
		return lsc.scan
	}
	lsc.info, lsc.nextInfo = lsc.nextInfo, lsc.info
	lsc.scan, lsc.nextScan = lsc.nextScan, lsc.Scanner.Scan()
	lsc.nextInfo.update(lsc.Scanner)
	return lsc.scan
}

func (lsc *lastScanner) Text() string {
	return lsc.info.Text
}

func (lsc *lastScanner) Buffer(buf []byte, max int) {
	lsc.Scanner.Buffer(buf, max)
	// memorize size values for future creation of prev/next buffers
	lsc.bufSize, lsc.bufCap = len(buf), max
	if lsc.bufCap < lsc.bufSize {
		lsc.bufCap = lsc.bufSize
	}
}

func (lsc *lastScanner) Bytes() []byte {
	return lsc.info.Bytes
}

func (lsc *lastScanner) Num() int {
	return lsc.info.Num
}
func (lsc *lastScanner) Match() bool {
	return lsc.info.Match
}

func (lsc *lastScanner) Last() bool {
	return lsc.nextScan == false
}

// NewFilterScanner creates a Scanner that produces only lines matched by a rule provided.
func NewFilterScanner(scn Scanner, rule MatchRule) Scanner {
	return NewOnlyMatchScanner(NewRuleScanner(scn, rule))
}
