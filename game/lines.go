package game

import (
	"math/bits"
	"sync"
)

type Line interface {
	At(x int) int   // returns symbol at position x, starts from 1
	Set(x, val int) // set value at position x
	Len() int       // returns length of line
	Free()          // put object to pool

	CopyN(num int) Line
}

type Lineset interface {
	Cols() int     // returns number of columns
	Line(int) Line // returns line with given number, starts from 1
	Num() int      // returns number lines in set
}

// SBL is selected bet lines bitset.
type SBL uint

// MakeSBL creates lines set from slice of line indexes.
func MakeSBL(lines ...int) SBL {
	var sbl SBL
	for _, n := range lines {
		sbl |= 1 << n
	}
	return sbl
}

// Num returns number of selected lines in set.
func (sbl SBL) Num() int {
	return bits.OnesCount(uint(sbl))
}

// Next helps iterate lines numbers as followed:
//
//	for n := sbl.Next(0); n != 0; n = sbl.Next(n) {}
func (sbl SBL) Next(n int) int {
	sbl >>= n + 1
	for sbl > 0 {
		n++
		if sbl&1 > 0 {
			return n
		}
		sbl >>= 1
	}
	return 0
}

// Is checks that line with given number is set.
func (sbl SBL) Is(n int) bool {
	return sbl&1<<n > 0
}

// Set line with given number.
func (sbl *SBL) Set(n int) {
	*sbl |= 1 << n
}

// Toggle line with given number and return whether it set.
func (sbl *SBL) Toggle(n int) bool {
	var bit SBL = 1 << n
	*sbl ^= bit
	return *sbl&bit > 0
}

type Line5x [5]int

var pooll5x = sync.Pool{
	New: func() any {
		return &Line5x{}
	},
}

func NewLine5x() *Line5x {
	return pooll5x.Get().(*Line5x)
}

func (l *Line5x) Free() {
	pooll5x.Put(l)
}

func (l *Line5x) At(x int) int {
	return l[x-1]
}

func (l *Line5x) Set(x, val int) {
	l[x-1] = val
}

func (l *Line5x) Len() int {
	return 5
}

func (l *Line5x) CopyN(num int) Line {
	var dst = NewLine5x()
	copy(dst[:num], l[:num])
	for i := num; i < 5; i++ {
		dst[i] = 0
	}
	return dst
}

type Lineset5x []Line5x

func (ls Lineset5x) Cols() int {
	return 5
}

func (ls Lineset5x) Line(n int) Line {
	return &ls[n-1]
}

func (ls Lineset5x) Num() int {
	return len(ls)
}

// (1 ,1) symbol is on left top corner

// Megajack 21 bet lines
var BetLinesMgj = Lineset5x{
	{2, 2, 2, 2, 2}, // 1
	{1, 1, 1, 1, 1}, // 2
	{3, 3, 3, 3, 3}, // 3
	{1, 2, 3, 2, 1}, // 4
	{3, 2, 1, 2, 3}, // 5
	{1, 1, 2, 3, 3}, // 6
	{3, 3, 2, 1, 1}, // 7
	{2, 1, 2, 3, 2}, // 8
	{2, 3, 2, 1, 2}, // 9
	{1, 2, 1, 2, 1}, // 10
	{3, 2, 3, 2, 3}, // 11
	{1, 2, 2, 2, 2}, // 12
	{3, 2, 2, 2, 2}, // 13
	{2, 2, 1, 1, 1}, // 14
	{2, 2, 3, 3, 3}, // 15
	{2, 1, 1, 1, 1}, // 16
	{2, 3, 3, 3, 3}, // 17
	{1, 1, 1, 1, 2}, // 18
	{3, 3, 3, 3, 2}, // 19
	{3, 3, 2, 3, 3}, // 20
	{1, 1, 2, 1, 1}, // 21
}

// Novomatic 9 bet lines (old versions of games)
var BetLinesNvm9 = Lineset5x{
	{2, 2, 2, 2, 2}, // 1
	{1, 1, 1, 1, 1}, // 2
	{3, 3, 3, 3, 3}, // 3
	{1, 2, 3, 2, 1}, // 4
	{3, 2, 1, 2, 3}, // 5
	{2, 3, 3, 3, 2}, // 6
	{2, 1, 1, 1, 2}, // 7
	{3, 3, 2, 1, 1}, // 8
	{1, 1, 2, 3, 3}, // 9
}

// Novomatic 10 bet lines (deluxe versions of games)
var BetLinesNvm10 = Lineset5x{
	{2, 2, 2, 2, 2}, // 1
	{1, 1, 1, 1, 1}, // 2
	{3, 3, 3, 3, 3}, // 3
	{1, 2, 3, 2, 1}, // 4
	{3, 2, 1, 2, 3}, // 5
	{2, 3, 3, 3, 2}, // 6
	{2, 1, 1, 1, 2}, // 7
	{3, 3, 2, 1, 1}, // 8
	{1, 1, 2, 3, 3}, // 9
	{3, 2, 2, 2, 1}, // 10
}

// Novomatic 20 bet lines (new games)
var BetLinesNvm20 = Lineset5x{
	{2, 2, 2, 2, 2}, // 1
	{1, 1, 1, 1, 1}, // 2
	{3, 3, 3, 3, 3}, // 3
	{1, 2, 3, 2, 1}, // 4
	{3, 2, 1, 2, 3}, // 5
	{1, 1, 2, 3, 3}, // 6
	{3, 3, 2, 1, 1}, // 7
	{2, 1, 1, 1, 2}, // 8
	{2, 3, 3, 3, 2}, // 9
	{1, 2, 2, 2, 1}, // 10
	{3, 2, 2, 2, 3}, // 11
	{2, 1, 2, 3, 2}, // 12
	{2, 3, 2, 1, 2}, // 13
	{1, 3, 1, 3, 1}, // 14
	{3, 1, 3, 1, 3}, // 15
	{1, 2, 1, 2, 1}, // 16
	{3, 2, 3, 2, 3}, // 17
	{1, 3, 3, 3, 3}, // 18
	{3, 1, 1, 1, 1}, // 19
	{2, 2, 1, 2, 2}, // 20
}

// BetSoft 25 bet lines
var BetLinesBetSoft25 = Lineset5x{
	{2, 2, 2, 2, 2}, // 1
	{1, 1, 1, 1, 1}, // 2
	{3, 3, 3, 3, 3}, // 3
	{1, 2, 3, 2, 1}, // 4
	{3, 2, 1, 2, 3}, // 5
	{1, 1, 2, 1, 1}, // 6
	{3, 3, 2, 3, 3}, // 7
	{2, 3, 3, 3, 2}, // 8
	{2, 1, 1, 1, 2}, // 9
	{2, 1, 2, 1, 2}, // 10
	{2, 3, 2, 3, 2}, // 11
	{1, 2, 1, 2, 1}, // 12
	{3, 2, 3, 2, 3}, // 13
	{2, 2, 1, 2, 2}, // 14
	{2, 2, 3, 2, 2}, // 15
	{1, 2, 2, 2, 1}, // 16
	{3, 2, 2, 2, 3}, // 17
	{1, 2, 3, 3, 3}, // 18
	{3, 2, 1, 1, 1}, // 19
	{1, 3, 1, 3, 1}, // 20
	{3, 1, 3, 1, 3}, // 21
	{1, 3, 3, 3, 1}, // 22
	{3, 1, 1, 1, 3}, // 23
	{1, 1, 3, 1, 1}, // 24
	{3, 3, 1, 3, 3}, // 25
}

// BetSoft 30 bet lines
var BetLinesBetSoft30 = Lineset5x{
	{2, 2, 2, 2, 2}, // 1
	{1, 1, 1, 1, 1}, // 2
	{3, 3, 3, 3, 3}, // 3
	{1, 2, 3, 2, 1}, // 4
	{3, 2, 1, 2, 3}, // 5
	{1, 1, 2, 1, 1}, // 6
	{3, 3, 2, 3, 3}, // 7
	{2, 3, 3, 3, 2}, // 8
	{2, 1, 1, 1, 2}, // 9
	{2, 1, 2, 1, 2}, // 10
	{2, 3, 2, 3, 2}, // 11
	{1, 2, 1, 2, 1}, // 12
	{3, 2, 3, 2, 3}, // 13
	{2, 2, 1, 2, 2}, // 14
	{2, 2, 3, 2, 2}, // 15
	{1, 2, 2, 2, 1}, // 16
	{3, 2, 2, 2, 3}, // 17
	{1, 2, 3, 3, 3}, // 18
	{3, 2, 1, 1, 1}, // 19
	{1, 3, 1, 3, 1}, // 20
	{3, 1, 3, 1, 3}, // 21
	{1, 3, 3, 3, 1}, // 22
	{3, 1, 1, 1, 3}, // 23
	{1, 1, 3, 1, 1}, // 24
	{3, 3, 1, 3, 3}, // 25
	{1, 3, 2, 1, 3}, // 26
	{3, 1, 2, 3, 1}, // 27
	{2, 1, 3, 2, 3}, // 28
	{1, 3, 2, 3, 2}, // 29
	{3, 2, 1, 1, 2}, // 30
}

// NetEnt 10 bet lines
var BetLinesNetEnt10 = Lineset5x{
	{2, 2, 2, 2, 2}, // 1
	{1, 1, 1, 1, 1}, // 2
	{3, 3, 3, 3, 3}, // 3
	{1, 2, 3, 2, 1}, // 4
	{3, 2, 1, 2, 3}, // 5
	{1, 1, 2, 1, 1}, // 6
	{3, 3, 2, 3, 3}, // 7
	{2, 3, 3, 3, 2}, // 8
	{2, 1, 1, 1, 2}, // 9
	{2, 1, 2, 1, 2}, // 10
}

var BetLines5x = map[string]Lineset5x{
	"mgj":   BetLinesMgj,
	"nvm9":  BetLinesNvm9,
	"nvm10": BetLinesNvm10,
	"nvm20": BetLinesNvm20,
	"bs25":  BetLinesBetSoft25,
	"bs30":  BetLinesBetSoft30,
	"ne10":  BetLinesNetEnt10,
}
