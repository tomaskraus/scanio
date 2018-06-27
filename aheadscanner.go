package scanio

import "bufio"

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
