// Package scanio provides bufio.Scanner wrappers with filtering, chaining and last-token ahead reading.
//
package scanio

import (
	"bufio"
	"io"
)

// Scanner interface reflects the set of methods of the bufio.Scanner.
// Adds two more methods: number of the tokens scanned and match status.
type Scanner interface {
	Buffer(buf []byte, max int)
	Bytes() []byte
	Err() error
	Scan() bool
	Split(split bufio.SplitFunc)
	Text() string

	Match() bool // true if a token matches Scanner's MatchRule
	Num() int    // number of a current token
}

//--------------------------------------------------------------------

// reader Scanner
type readerScanner struct {
	scn   *bufio.Scanner
	match bool
	num   int
}

// NewScanner creates a new Scanner using a Reader.
// This Scanner can be used instead of bufio.Scanner
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

//--------------------------------------------------------------------------------

// MatchRule for NewRuleScanner.
type MatchRule func(input string) bool

type ruleScanner struct {
	Scanner
	rule   MatchRule
	matchR bool
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
		sc.matchR = sc.rule(sc.Scanner.Text())
		return true
	}
	return false
}

func (sc *ruleScanner) Match() bool {
	return sc.matchR
}

//--------------------------------------------------------------------------------

// MatchByteRule for NewByteRuleScanner.
type MatchByteRule func(input []byte) bool

type byteRuleScanner struct {
	Scanner
	rule   MatchByteRule
	matchR bool
}

// NewByteRuleScanner returns new, rule-based Scanner.
func NewByteRuleScanner(scn Scanner, rule MatchByteRule) Scanner {
	return Scanner(&byteRuleScanner{
		Scanner: scn,
		rule:    rule,
	})
}

func (sc *byteRuleScanner) Scan() bool {
	if sc.Scanner.Scan() {
		sc.matchR = sc.rule(sc.Scanner.Bytes())
		return true
	}
	return false
}

func (sc *byteRuleScanner) Match() bool {
	return sc.matchR
}

//--------------------------------------------------------------------------------

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

// NewFilterScanner creates a Scanner that outputs only tokens matched by a rule provided.
func NewFilterScanner(scn Scanner, rule MatchRule) Scanner {
	return NewOnlyMatchScanner(NewRuleScanner(scn, rule))
}

// NewByteFilterScanner creates a Scanner that outputs only tokens matched by a rule provided.
func NewByteFilterScanner(scn Scanner, rule MatchByteRule) Scanner {
	return NewOnlyMatchScanner(NewByteRuleScanner(scn, rule))
}

// ---------------------------------------------------------------------------

// AheadScanner can tell if the current token is the last one.
// Does the one token forward-read to achieve this.
type AheadScanner interface {
	Scanner
	Last() bool
	BeginConsecutive() bool // begin of consecutive match sequence (even if its length is 1)
	EndConsecutive() bool   // end of consecutive match sequence (even if its length is 1)
	NumConsecutive() int    // number of consecutive matches
}

// info stores the Scanner's current state.
type info struct {
	ScanResult bool // holds result of Scanner.Scan()
	Text       string
	Bytes      []byte
	Num        int
	Match      bool
}

func newInfo(bufferLen, bufferCap int) *info {
	i := info{}
	i.Bytes = make([]byte, bufferLen, bufferCap)
	return &i
}

// isConsecEnd returns true if the next Info has a consecutive match
func (inf *info) isConsecMatch(nextInf *info) bool {
	if inf.Match && nextInf.Match && inf.Num+1 == nextInf.Num {
		// next token matches
		return true
	}
	return false
}

// update makes a snapshot of Scanner's current state.
func (inf *info) update(scn Scanner, scResult bool) {
	inf.Text, inf.Num, inf.Match, inf.ScanResult = scn.Text(), scn.Num(), scn.Match(), scResult
	// preserve the underlying scanner's buffer
	srcLen := len(scn.Bytes())
	inf.Bytes = inf.Bytes[:srcLen]
	copy(inf.Bytes, scn.Bytes())
}

const (
	startBufSize = 4096 // Size of initial allocation for buffer.   from golang.org/src/bufio/scan.go
)

type aheadScanner struct {
	Scanner
	info, nextInfo         *info
	consecNum              int
	consecBegin, consecEnd bool
	consecMode             bool
	bufSize, bufCap        int
	started                bool
}

// NewAheadScanner creates a new AheadScanner.
func NewAheadScanner(scn Scanner) AheadScanner {
	return AheadScanner(&aheadScanner{
		Scanner: scn,
		bufSize: startBufSize,
		bufCap:  bufio.MaxScanTokenSize,
	})
}

func (asc *aheadScanner) Scan() bool {
	if !asc.started {
		//initialize buffers
		asc.info = newInfo(asc.bufSize, asc.bufCap)
		asc.nextInfo = newInfo(asc.bufSize, asc.bufCap)

		//scan two tokens (one ahead)
		scanRes := asc.Scanner.Scan()
		asc.info.update(asc.Scanner, scanRes)
		nextScanRes := asc.Scanner.Scan()
		asc.nextInfo.update(asc.Scanner, nextScanRes)

		asc.started = true
	} else {
		asc.info, asc.nextInfo = asc.nextInfo, asc.info
		nextScanRes2 := asc.Scanner.Scan()
		asc.nextInfo.update(asc.Scanner, nextScanRes2)
	}

	if !asc.info.ScanResult {
		asc.consecNum = 0
		asc.consecBegin, asc.consecEnd = false, false
		return false
	}

	consecModeHasNowStarted := false
	if !asc.consecMode {
		asc.consecNum = 0
		asc.consecBegin, asc.consecEnd = false, false
		if asc.info.Match {
			consecModeHasNowStarted = true
			asc.consecBegin = true
			asc.consecMode = true
		}
	}
	if asc.consecMode {
		asc.consecNum++
		if !asc.info.isConsecMatch(asc.nextInfo) {
			asc.consecEnd = true
			asc.consecMode = false
		}
		if !consecModeHasNowStarted {
			asc.consecBegin = false
		}
	}

	return true
}

func (asc *aheadScanner) Text() string {
	return asc.info.Text
}

func (asc *aheadScanner) Buffer(buf []byte, max int) {
	asc.Scanner.Buffer(buf, max)
	// memorize size values for future creation of prev/next buffers
	asc.bufSize, asc.bufCap = len(buf), max
	if asc.bufCap < asc.bufSize {
		asc.bufCap = asc.bufSize
	}
}

func (asc *aheadScanner) Bytes() []byte {
	return asc.info.Bytes
}

func (asc *aheadScanner) Num() int {
	return asc.info.Num
}
func (asc *aheadScanner) Match() bool {
	return asc.info.Match
}

func (asc *aheadScanner) Last() bool {
	return asc.nextInfo.ScanResult == false
}

func (asc *aheadScanner) BeginConsecutive() bool {
	return asc.consecBegin
}

func (asc *aheadScanner) EndConsecutive() bool {
	return asc.consecEnd
}

func (asc *aheadScanner) NumConsecutive() int {
	return asc.consecNum
}
