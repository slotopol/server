package game

type Line interface {
	At(x int) int   // returns symbol at position x, starts from 1
	Set(x, val int) // set value at position x
	Len() int       // returns length of line
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
	var count SBL
	for sbl > 0 {
		count += sbl & 1
		sbl >>= 1
	}
	return int(count)
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

func (l *Line5x) At(x int) int {
	return l[x-1]
}

func (l *Line5x) Set(x, val int) {
	l[x-1] = val
}

func (l *Line5x) Len() int {
	return 5
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

var BetLines5x = map[string]Lineset5x{
	"mgj":   BetLinesMgj,
	"nvm9":  BetLinesNvm9,
	"nvm10": BetLinesNvm10,
	"nvm20": BetLinesNvm20,
}
