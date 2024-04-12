package game

import (
	"errors"
)

type Sym byte // symbol type

type Reels interface {
	Cols() int          // returns number of columns
	Reel(col int) []Sym // returns reel at given column, index from
	Reshuffles() int    // returns total number of reshuffles
}

// WinItem describes win on each line or scatters.
type WinItem struct {
	Pay  float64 `json:"pay,omitempty" yaml:"pay,omitempty" xml:"pay,omitempty,attr"`    // payment with selected bet
	Mult float64 `json:"mult,omitempty" yaml:"mult,omitempty" xml:"mult,omitempty,attr"` // multiplier for payment for free spins and other special cases
	Sym  Sym     `json:"sym,omitempty" yaml:"sym,omitempty" xml:"sym,omitempty,attr"`    // win symbol
	Num  int     `json:"num,omitempty" yaml:"num,omitempty" xml:"num,omitempty,attr"`    // number of win symbol
	Line int     `json:"line,omitempty" yaml:"line,omitempty" xml:"line,omitempty,attr"` // line mumber (0 for scatters and not lined)
	XY   Line    `json:"xy" yaml:"xy" xml:"xy"`                                          // symbols positions on screen
	Free int     `json:"free,omitempty" yaml:"free,omitempty" xml:"free,omitempty,attr"` // number of free spins remains
	BID  int     `json:"bid,omitempty" yaml:"bid,omitempty" xml:"bid,omitempty,attr"`    // bonus identifier
	Jack int     `json:"jack,omitempty" yaml:"jack,omitempty" xml:"jack,omitempty,attr"` // jackpot identifier
	Bon  any     `json:"bon,omitempty" yaml:"bon,omitempty" xml:"bon,omitempty"`         // bonus game data
}

// WinScan is full list of wins by all lines and scatters for some spin.
type WinScan struct {
	Wins []WinItem `json:"wins,omitempty" yaml:"wins,omitempty" xml:"wins,omitempty"`
}

// Reset puts lines to pool and set array empty with saved capacity.
func (ws *WinScan) Reset() {
	for _, wi := range ws.Wins {
		wi.XY.Free()
	}
	ws.Wins = ws.Wins[:0] // set it empty
}

// Total gain for spin.
func (ws *WinScan) Gain() float64 {
	var sum float64
	for _, wi := range ws.Wins {
		sum += wi.Pay * wi.Mult
	}
	return sum
}

type SlotGame interface {
	NewScreen() Screen                  // returns new empty screen object for this game, constat function
	Scanner(screen Screen, sw *WinScan) // scan given screen and append result to sw, constat function
	Spin(screen Screen)                 // fill the screen with random hits on those reels, constat function
	Spawn(screen Screen, sw *WinScan)   // setup bonus games to win results, constat function
	Apply(screen Screen, sw *WinScan)   // update game state to spin results
	FreeSpins() int                     // returns number of free spins remained, constat function
	GetGain() float64                   // returns gain for double up games, constat function
	SetGain(gain float64) error         // set gain to given value on double up games
	GetBet() float64                    // returns current bet, constat function
	SetBet(float64) error               // set bet to given value
	GetLines() SBL                      // returns selected lines indexes, constat function
	SetLines(SBL) error                 // setup selected lines indexes
	GetReels() string                   // returns reels descriptor
	SetReels(rd string) error           // setup reels descriptor
}

// Reels for 3-reels slots.
type Reels3x [3][]Sym

func (r *Reels3x) Cols() int {
	return 3
}

func (r *Reels3x) Reel(col int) []Sym {
	return r[col-1]
}

func (r *Reels3x) Reshuffles() int {
	return len(r[0]) * len(r[1]) * len(r[2])
}

// Reels for 5-reels slots.
type Reels5x [5][]Sym

func (r *Reels5x) Cols() int {
	return 5
}

func (r *Reels5x) Reel(col int) []Sym {
	return r[col-1]
}

func (r *Reels5x) Reshuffles() int {
	return len(r[0]) * len(r[1]) * len(r[2]) * len(r[3]) * len(r[4])
}

var (
	ErrBetEmpty   = errors.New("bet is empty")
	ErrNoLineset  = errors.New("lines set is empty")
	ErrLinesetOut = errors.New("lines set is out of range bet lines")
	ErrNoFeature  = errors.New("feature not available")
	ErrNoReels    = errors.New("no reels for given descriptor")
)

// Slot5x3 is base struct for all slot games with screen 5x3.
type Slot3x3 struct {
	RD  string  `json:"rd" yaml:"rd" xml:"rd"`    // reels descriptor
	SBL SBL     `json:"sbl" yaml:"sbl" xml:"sbl"` // selected bet lines
	Bet float64 `json:"bet" yaml:"bet" xml:"bet"` // bet value

	Gain float64 `json:"gain,omitempty" yaml:"gain,omitempty" xml:"gain,omitempty"` // gain for double up games
}

func (g *Slot3x3) NewScreen() Screen {
	return NewScreen3x3()
}

func (g *Slot3x3) Spawn(screen Screen, sw *WinScan) {
}

func (g *Slot3x3) Apply(screen Screen, sw *WinScan) {
	g.Gain = sw.Gain()
}

func (g *Slot3x3) FreeSpins() int {
	return 0
}

func (g *Slot3x3) GetGain() float64 {
	return g.Gain
}

func (g *Slot3x3) SetGain(gain float64) error {
	g.Gain = gain
	return nil
}

func (g *Slot3x3) GetBet() float64 {
	return g.Bet
}

func (g *Slot3x3) SetBet(bet float64) error {
	if bet < 1 {
		return ErrBetEmpty
	}
	if g.FreeSpins() > 0 {
		return ErrNoFeature
	}
	g.Bet = bet
	return nil
}

func (g *Slot3x3) GetLines() SBL {
	return g.SBL
}

func (g *Slot3x3) SetLines(sbl SBL) error {
	return ErrNoFeature
}

func (g *Slot3x3) GetReels() string {
	return g.RD
}

// Slot5x3 is base struct for all slot games with screen 5x3.
type Slot5x3 struct {
	RD  string  `json:"rd" yaml:"rd" xml:"rd"`    // reels descriptor
	BLI string  `json:"bli" yaml:"bli" xml:"bli"` // bet lines index
	SBL SBL     `json:"sbl" yaml:"sbl" xml:"sbl"` // selected bet lines
	Bet float64 `json:"bet" yaml:"bet" xml:"bet"` // bet value

	Gain float64 `json:"gain,omitempty" yaml:"gain,omitempty" xml:"gain,omitempty"` // gain for double up games
}

func (g *Slot5x3) NewScreen() Screen {
	return NewScreen5x3()
}

func (g *Slot5x3) Spawn(screen Screen, sw *WinScan) {
}

func (g *Slot5x3) Apply(screen Screen, sw *WinScan) {
	g.Gain = sw.Gain()
}

func (g *Slot5x3) FreeSpins() int {
	return 0
}

func (g *Slot5x3) GetGain() float64 {
	return g.Gain
}

func (g *Slot5x3) SetGain(gain float64) error {
	g.Gain = gain
	return nil
}

func (g *Slot5x3) GetBet() float64 {
	return g.Bet
}

func (g *Slot5x3) SetBet(bet float64) error {
	if bet < 1 {
		return ErrBetEmpty
	}
	if g.FreeSpins() > 0 {
		return ErrNoFeature
	}
	g.Bet = bet
	return nil
}

func (g *Slot5x3) GetLines() SBL {
	return g.SBL
}

func (g *Slot5x3) SetLines(sbl SBL) error {
	var bl = BetLines5x[g.BLI]
	var mask SBL = (1<<len(bl) - 1) << 1
	if sbl == 0 {
		return ErrNoLineset
	}
	if sbl&^mask != 0 {
		return ErrLinesetOut
	}
	if g.FreeSpins() > 0 {
		return ErrNoFeature
	}
	g.SBL = sbl
	return nil
}

func (g *Slot5x3) GetReels() string {
	return g.RD
}
