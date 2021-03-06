// Package scanio provides bufio.Scanner chainable wrappers with read-ahead and filtering.
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

	IsMatch() bool // true if current token matches Scanner's MatchRule
	NumRead() int  // number of all tokens read
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

func (sc *readerScanner) Buffer(b []byte, max int) {
	sc.scn.Buffer(b, max)
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

func (sc *readerScanner) IsMatch() bool {
	return sc.match
}

func (sc *readerScanner) NumRead() int {
	return sc.num
}

//--------------------------------------------------------------------------------

// MatchRule for NewRuleScanner.
type MatchRule func(token []byte) (matched bool, err error)

type ruleScanner struct {
	Scanner
	rule    MatchRule
	matched bool
	err     error
}

// NewRuleScanner returns new, rule-based Scanner.
func NewRuleScanner(sc Scanner, rule MatchRule) Scanner {
	return Scanner(&ruleScanner{
		Scanner: sc,
		rule:    rule,
	})
}

func (sc *ruleScanner) Scan() bool {
	if sc.Scanner.Scan() {
		sc.matched, sc.err = sc.rule(sc.Scanner.Bytes())
		if sc.err != nil {
			return false
		}
		return true
	}
	return false
}

func (sc *ruleScanner) IsMatch() bool {
	return sc.matched
}

func (sc *ruleScanner) Err() error {
	if sc.err != nil {
		return sc.err
	}
	return sc.Scanner.Err()
}

//--------------------------------------------------------------------------------

type onlyMatchScanner struct {
	Scanner
}

// NewOnlyMatchScanner returns new Scanner.
func NewOnlyMatchScanner(sc Scanner) Scanner {
	return Scanner(&onlyMatchScanner{
		Scanner: sc,
	})
}

func (sc *onlyMatchScanner) Scan() bool {
	for sc.Scanner.Scan() {
		if sc.Scanner.IsMatch() {
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
func NewOnlyNotMatchScanner(sc Scanner) Scanner {
	return Scanner(&onlyNotMatchScanner{
		Scanner: sc,
	})
}

func (sc *onlyNotMatchScanner) Scan() bool {
	for sc.Scanner.Scan() {
		if sc.Scanner.IsMatch() {
			continue
		}
		return true
	}
	return false
}

// ---------------------------------------------------------------------------

// NewFilterScanner creates a Scanner that outputs only tokens matched by a rule provided.
func NewFilterScanner(sc Scanner, rule MatchRule) Scanner {
	return NewOnlyMatchScanner(NewRuleScanner(sc, rule))
}
