package scanio

import "bufio"

// AheadScanner can tell if the current token is the last one.
// Does the one token forward-read to achieve this.
type AheadScanner interface {
	Scanner
	IsLast() bool
	IsConsecutiveBegin() bool // begin of consecutive positive-match sequence (even if its length is 1)
	IsConsecutiveEnd() bool   // end of consecutive positive-match sequence (even if its length is 1)
	NumConsecutive() int      // number of consecutive positive matches
}

// info stores the Scanner's current state.
type info struct {
	ScanRes bool // holds result of Scanner.Scan()
	Text    string
	Bytes   []byte
	Index   int
	IsMatch bool
}

func newInfo(bufLen, bufCap int) *info {
	i := info{}
	i.Bytes = make([]byte, bufLen, bufCap)
	i.Index = -1
	return &i
}

// isConsecEnd returns true if the next Info has a consecutive match
func (i *info) isConsecMatch(nextI *info) bool {
	if i.IsMatch && nextI.IsMatch && i.Index+1 == nextI.Index {
		// next token matches
		return true
	}
	return false
}

// update makes a snapshot of Scanner's current state.
func (i *info) update(sc Scanner, scResult bool) {
	i.Text, i.Index, i.IsMatch, i.ScanRes = sc.Text(), sc.Index(), sc.IsMatch(), scResult
	// preserve the underlying scanner's buffer
	srcLen := len(sc.Bytes())
	i.Bytes = i.Bytes[:srcLen]
	copy(i.Bytes, sc.Bytes())
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
func NewAheadScanner(sc Scanner) AheadScanner {
	return AheadScanner(&aheadScanner{
		Scanner: sc,
		bufSize: startBufSize,
		bufCap:  bufio.MaxScanTokenSize,
	})
}

func (sc *aheadScanner) Scan() bool {
	if !sc.started {
		//initialize buffers
		sc.info = newInfo(sc.bufSize, sc.bufCap)
		sc.nextInfo = newInfo(sc.bufSize, sc.bufCap)

		//scan two tokens (one ahead)
		scanRes := sc.Scanner.Scan()
		sc.info.update(sc.Scanner, scanRes)
		nextScanRes := sc.Scanner.Scan()
		sc.nextInfo.update(sc.Scanner, nextScanRes)

		sc.started = true
	} else {
		sc.info, sc.nextInfo = sc.nextInfo, sc.info
		nextScanRes2 := sc.Scanner.Scan()
		sc.nextInfo.update(sc.Scanner, nextScanRes2)
	}

	if !sc.info.ScanRes {
		sc.consecNum = 0
		sc.consecBegin, sc.consecEnd = false, false
		return false
	}

	if !sc.consecMode {
		sc.consecNum = 0
		sc.consecBegin, sc.consecEnd = false, false
		if sc.info.IsMatch {
			sc.consecBegin = true
			sc.consecMode = true
		}
	}
	if sc.consecMode {
		sc.consecNum++
		if !sc.info.isConsecMatch(sc.nextInfo) {
			sc.consecEnd = true
			sc.consecMode = false
		}
		// preserve consecBegin flag if consecEnd has occurred at the same token
		if sc.consecNum > 1 {
			sc.consecBegin = false
		}
	}

	return true
}

func (sc *aheadScanner) Text() string {
	return sc.info.Text
}

func (sc *aheadScanner) Buffer(buf []byte, max int) {
	sc.Scanner.Buffer(buf, max)
	// memorize size values for future creation of prev/next buffers
	sc.bufSize, sc.bufCap = len(buf), max
	if sc.bufCap < sc.bufSize {
		sc.bufCap = sc.bufSize
	}
}

func (sc *aheadScanner) Bytes() []byte {
	return sc.info.Bytes
}

func (sc *aheadScanner) Index() int {
	return sc.info.Index
}
func (sc *aheadScanner) IsMatch() bool {
	return sc.info.IsMatch
}

func (sc *aheadScanner) IsLast() bool {
	return sc.nextInfo.ScanRes == false
}

func (sc *aheadScanner) IsConsecutiveBegin() bool {
	return sc.consecBegin
}

func (sc *aheadScanner) IsConsecutiveEnd() bool {
	return sc.consecEnd
}

func (sc *aheadScanner) NumConsecutive() int {
	return sc.consecNum
}
