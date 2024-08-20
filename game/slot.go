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

// Wins is full list of wins by all lines and scatters for some spin.
type Wins []WinItem

// Reset puts lines to pool and set array empty with saved capacity.
func (wins *Wins) Reset() {
	for _, wi := range *wins {
		if wi.XY != nil {
			wi.XY.Free()
		}
	}
	*wins = (*wins)[:0] // set it empty
}

// Total gain for spin.
func (wins Wins) Gain() float64 {
	var sum float64
	for _, wi := range wins {
		sum += wi.Pay * wi.Mult
	}
	return sum
}

type SlotGame interface {
	NewScreen() Screen     // returns new empty screen object for this game, constat function
	Scanner(Screen, *Wins) // scan given screen and append result to wins, constat function
	Spin(Screen)           // fill the screen with random hits on those reels, constat function
	Spawn(Screen, Wins)    // setup bonus games to wins results, constat function
	Apply(Screen, Wins)    // update game state to spin results
	FreeSpins() int        // returns number of free spins remained, constat function
	GetGain() float64      // returns gain for double up games, constat function
	SetGain(float64) error // set gain to given value on double up games
	GetBet() float64       // returns current bet, constat function
	SetBet(float64) error  // set bet to given value
	GetLines() Bitset      // returns selected bet lines indexes, constat function
	SetLines(Bitset) error // setup selected bet lines indexes
	GetRTP() float64       // returns master RTP
	SetRTP(float64) error  // setup master RTP
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
	ErrNoWay      = errors.New("no way to here")
	ErrBetEmpty   = errors.New("bet is empty")
	ErrNoLineset  = errors.New("lines set is empty")
	ErrNoRtp      = errors.New("RTP not given")
	ErrLinesetOut = errors.New("lines set is out of range bet lines")
	ErrNoFeature  = errors.New("feature not available")
	ErrNoReels    = errors.New("no reels for given descriptor")
)

// Slot5x3 is base struct for all slot games with screen 5x3.
type Slot3x3 struct {
	RTP float64 `json:"rtp" yaml:"rtp" xml:"rtp"` // master RTP
	SBL Bitset  `json:"sbl" yaml:"sbl" xml:"sbl"` // selected bet lines
	Bet float64 `json:"bet" yaml:"bet" xml:"bet"` // bet value

	Gain float64 `json:"gain,omitempty" yaml:"gain,omitempty" xml:"gain,omitempty"` // gain for double up games
}

func (g *Slot3x3) NewScreen() Screen {
	return NewScreen3x3()
}

func (g *Slot3x3) Spawn(screen Screen, wins Wins) {
}

func (g *Slot3x3) Apply(screen Screen, wins Wins) {
	g.Gain = wins.Gain()
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
	if bet <= 0 {
		return ErrBetEmpty
	}
	if g.FreeSpins() > 0 {
		return ErrNoFeature
	}
	g.Bet = bet
	return nil
}

func (g *Slot3x3) GetLines() Bitset {
	return g.SBL
}

func (g *Slot3x3) GetRTP() float64 {
	return g.RTP
}

func (g *Slot3x3) SetRTP(rtp float64) error {
	if rtp <= 0 {
		return ErrNoRtp
	}
	g.RTP = rtp
	return nil
}

// Slot5x3 is base struct for all slot games with screen 5x3.
type Slot5x3 struct {
	RTP float64 `json:"rtp" yaml:"rtp" xml:"rtp"` // master RTP
	SBL Bitset  `json:"sbl" yaml:"sbl" xml:"sbl"` // selected bet lines
	Bet float64 `json:"bet" yaml:"bet" xml:"bet"` // bet value

	Gain float64 `json:"gain,omitempty" yaml:"gain,omitempty" xml:"gain,omitempty"` // gain for double up games
}

func (g *Slot5x3) NewScreen() Screen {
	return NewScreen5x3()
}

func (g *Slot5x3) Spawn(screen Screen, wins Wins) {
}

func (g *Slot5x3) Apply(screen Screen, wins Wins) {
	g.Gain = wins.Gain()
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
	if bet <= 0 {
		return ErrBetEmpty
	}
	if g.FreeSpins() > 0 {
		return ErrNoFeature
	}
	g.Bet = bet
	return nil
}

func (g *Slot5x3) GetLines() Bitset {
	return g.SBL
}

func (g *Slot5x3) GetRTP() float64 {
	return g.RTP
}

func (g *Slot5x3) SetRTP(rtp float64) error {
	if rtp <= 0 {
		return ErrNoRtp
	}
	g.RTP = rtp
	return nil
}

// Slot5x4 is base struct for all slot games with screen 5x4.
type Slot5x4 struct {
	RTP float64 `json:"rtp" yaml:"rtp" xml:"rtp"` // master RTP
	SBL Bitset  `json:"sbl" yaml:"sbl" xml:"sbl"` // selected bet lines
	Bet float64 `json:"bet" yaml:"bet" xml:"bet"` // bet value

	Gain float64 `json:"gain,omitempty" yaml:"gain,omitempty" xml:"gain,omitempty"` // gain for double up games
}

func (g *Slot5x4) NewScreen() Screen {
	return NewScreen5x4()
}

func (g *Slot5x4) Spawn(screen Screen, wins Wins) {
}

func (g *Slot5x4) Apply(screen Screen, wins Wins) {
	g.Gain = wins.Gain()
}

func (g *Slot5x4) FreeSpins() int {
	return 0
}

func (g *Slot5x4) GetGain() float64 {
	return g.Gain
}

func (g *Slot5x4) SetGain(gain float64) error {
	g.Gain = gain
	return nil
}

func (g *Slot5x4) GetBet() float64 {
	return g.Bet
}

func (g *Slot5x4) SetBet(bet float64) error {
	if bet <= 0 {
		return ErrBetEmpty
	}
	if g.FreeSpins() > 0 {
		return ErrNoFeature
	}
	g.Bet = bet
	return nil
}

func (g *Slot5x4) GetLines() Bitset {
	return g.SBL
}

func (g *Slot5x4) GetRTP() float64 {
	return g.RTP
}

func (g *Slot5x4) SetRTP(rtp float64) error {
	if rtp <= 0 {
		return ErrNoRtp
	}
	g.RTP = rtp
	return nil
}
