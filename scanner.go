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
