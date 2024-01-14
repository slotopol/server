package game

import (
	"errors"
	"io"
	"math/rand"
)

type Biner interface {
	MarshalBin() ([]byte, error)
	UnmarshalBin(b []byte) (int, error)
}

type Reels interface {
	Cols() int          // returns number of columns
	Reel(col int) []int // returns reel at given column, index from
	Reshuffles() int    // returns total number of reshuffles
	Spin(screen Screen) // fill the screen with random hits on those reels
}

type Screen interface {
	Dim() (int, int)                   // returns screen dimensions
	At(x int, y int) int               // returns symbol at position (x, y), starts from (1, 1)
	SetCol(x int, reel []int, pos int) // setup column on screen with given reel at given position
	Biner
}

type WinItem struct {
	Pay  int  `json:"pay,omitempty" yaml:"pay,omitempty" xml:"pay,omitempty,attr"`    // payment with selected bet
	Mult int  `json:"mult,omitempty" yaml:"mult,omitempty" xml:"mult,omitempty,attr"` // multiplier for payment for free spins and other special cases
	Sym  int  `json:"sym,omitempty" yaml:"sym,omitempty" xml:"sym,omitempty,attr"`    // win symbol
	Num  int  `json:"num,omitempty" yaml:"num,omitempty" xml:"num,omitempty,attr"`    // number of win symbol
	Line int  `json:"line,omitempty" yaml:"line,omitempty" xml:"line,omitempty,attr"` // line mumber (0 for scatters and not lined)
	XY   Line `json:"xy" yaml:"xy" xml:"xy"`                                          // symbols positions on screen
	Free int  `json:"free,omitempty" yaml:"free,omitempty" xml:"free,omitempty,attr"` // number of free spins remains
	BID  int  `json:"bid,omitempty" yaml:"bid,omitempty" xml:"bid,omitempty,attr"`    // bonus identifier
	Jack int  `json:"jack,omitempty" yaml:"jack,omitempty" xml:"jack,omitempty,attr"` // jackpot identifier
	Bon  any  `json:"bon,omitempty" yaml:"bon,omitempty" xml:"bon,omitempty"`         // bonus game data
}

type WinScan struct {
	Wins []WinItem `json:"wins" yaml:"wins" xml:"wins"`
}

func (ws *WinScan) SumPay() int {
	var sum int
	for _, wi := range ws.Wins {
		sum += wi.Pay * wi.Mult
	}
	return sum
}

type SlotGame interface {
	NewScreen() Screen                  // returns new empty screen object for this game
	Spin(screen Screen)                 // fill the screen with random hits on those reels
	GetBet() int                        // returns current bet
	SetBet(int) error                   // set bet to given value
	GetLines() SBL                      // returns selected lines indexes
	SetLines(SBL) error                 // setup selected lines indexes
	Scanner(screen Screen, sw *WinScan) // scan given screen and append result to sw
	Spawn(screen Screen, sw *WinScan)   // setup bonus games to win results
}

// Reels for 5-reels slots.
type Reels5x [5][]int

func (r *Reels5x) Cols() int {
	return 5
}

func (r *Reels5x) Reel(col int) []int {
	return r[col-1]
}

func (r *Reels5x) Reshuffles() int {
	return len(r[0]) * len(r[1]) * len(r[2]) * len(r[3]) * len(r[4])
}

func (r *Reels5x) Spin(screen Screen) {
	for x := 1; x <= 5; x++ {
		var reel = r.Reel(x)
		var hit = rand.Intn(len(reel))
		screen.SetCol(x, reel, hit)
	}
}

// Screen for 5x3 slots.
type Screen5x3 [5][3]int

func (s *Screen5x3) Dim() (int, int) {
	return 5, 3
}

func (s *Screen5x3) At(x int, y int) int {
	return s[x-1][y-1]
}

func (s *Screen5x3) SetCol(x int, reel []int, pos int) {
	for y := 0; y < 3; y++ {
		s[x-1][y] = reel[(pos+y)%len(reel)]
	}
}

func (s *Screen5x3) MarshalBin() ([]byte, error) {
	var b [15]byte
	var i int
	for x := 0; x < 5; x++ {
		for y := 0; y < 3; y++ {
			b[i] = byte(s[x][y])
			i++
		}
	}
	return b[:], nil
}

func (s *Screen5x3) UnmarshalBin(b []byte) (int, error) {
	if len(b) < 15 {
		return 0, io.EOF
	}
	var i int
	for x := 0; x < 5; x++ {
		for y := 0; y < 3; y++ {
			s[x][y] = int(b[i])
			i++
		}
	}
	return 15, nil
}

var (
	ErrBetEmpty   = errors.New("bet is empty")
	ErrNoLineset  = errors.New("lines set is empty")
	ErrLinesetOut = errors.New("lines set is out of range bet lines")
)

type Slot5x3 struct {
	SBL SBL // selected bet lines
	Bet int // bet value
	FS  int // free spin number

	Reels    *Reels5x
	BetLines *Lineset5x
}

func (g *Slot5x3) NewScreen() Screen {
	return &Screen5x3{}
}

func (g *Slot5x3) Spin(screen Screen) {
	g.Reels.Spin(screen)
}

func (g *Slot5x3) GetBet() int {
	return g.Bet
}

func (g *Slot5x3) SetBet(bet int) error {
	if bet < 1 {
		return ErrBetEmpty
	}
	g.Bet = bet
	return nil
}

func (g *Slot5x3) GetLines() SBL {
	return g.SBL
}

func (g *Slot5x3) SetLines(sbl SBL) error {
	var mask SBL = (1<<len(*g.BetLines) - 1) << 1
	if sbl == 0 {
		return ErrNoLineset
	}
	if mask&sbl != 0 {
		return ErrLinesetOut
	}
	g.SBL = sbl
	return nil
}

func (g *Slot5x3) Spawn(screen Screen, sw *WinScan) {
}
