package slot

import (
	"errors"
	"math"
)

type (
	Sym byte // symbol type
	Pos int8 // screen or line position
)

type Reels interface {
	Cols() int          // returns number of columns
	Reel(col Pos) []Sym // returns reel at given column, index from
	Reshuffles() uint64 // returns total number of reshuffles
}

// WinItem describes win on each line or by scatters.
type WinItem struct {
	Pay  float64 `json:"pay,omitempty" yaml:"pay,omitempty" xml:"pay,omitempty,attr"`    // payment with selected bet
	Mult float64 `json:"mult,omitempty" yaml:"mult,omitempty" xml:"mult,omitempty,attr"` // multiplier of payment for free spins and other special cases
	Sym  Sym     `json:"sym,omitempty" yaml:"sym,omitempty" xml:"sym,omitempty,attr"`    // win symbol
	Num  Pos     `json:"num,omitempty" yaml:"num,omitempty" xml:"num,omitempty,attr"`    // number of win symbols
	Line int     `json:"line,omitempty" yaml:"line,omitempty" xml:"line,omitempty,attr"` // line mumber (0 for scatters and not lined combinations)
	XY   Linex   `json:"xy" yaml:"xy" xml:"xy"`                                          // symbols positions on screen
	Free int     `json:"free,omitempty" yaml:"free,omitempty" xml:"free,omitempty,attr"` // number of free spins
	BID  int     `json:"bid,omitempty" yaml:"bid,omitempty" xml:"bid,omitempty,attr"`    // bonus identifier
	Jack int     `json:"jack,omitempty" yaml:"jack,omitempty" xml:"jack,omitempty,attr"` // jackpot identifier
	Bon  any     `json:"bon,omitempty" yaml:"bon,omitempty" xml:"bon,omitempty"`         // bonus game data
}

// Wins is full list of wins by all lines and scatters for some spin.
type Wins []WinItem

// Reset puts lines to pool and set array empty with saved capacity.
func (wins *Wins) Reset() {
	*wins = (*wins)[:0] // set it empty without memory reallocation
}

// Total gain for spin.
func (wins Wins) Gain() float64 {
	var sum float64
	for _, wi := range wins {
		sum += wi.Pay * wi.Mult
	}
	return sum
}

// SlotGame is common slots interface. Any slot game should implement this interface.
type SlotGame interface {
	NewScreen() Screen     // returns new empty screen object for this game, constat function
	Scanner(Screen, *Wins) // scan given screen and append result to wins, constat function
	Spin(Screen, float64)  // fill the screen with random hits on reels closest to given RTP, constat function
	Spawn(Screen, Wins)    // setup bonus games to wins results, constat function
	Prepare()              // update game state before new spin
	Apply(Screen, Wins)    // update game state to spin results
	FreeSpins() int        // returns number of free spins remained, constat function
	GetGain() float64      // returns gain for double up games, constat function
	SetGain(float64) error // set gain to given value on double up games
	GetBet() float64       // returns current bet, constat function
	SetBet(float64) error  // set bet to given value
	GetSel() int           // returns number of selected bet lines, constat function
	SetSel(int) error      // setup number of selected bet lines
	SetMode(int) error     // change game mode depending on the user's choice
}

// Reels for 3-reels slots.
type Reels3x [3][]Sym

func (r *Reels3x) Cols() int {
	return 3
}

func (r *Reels3x) Reel(col Pos) []Sym {
	return r[col-1]
}

func (r *Reels3x) Reshuffles() uint64 {
	return uint64(len(r[0])) * uint64(len(r[1])) * uint64(len(r[2]))
}

// Reels for 5-reels slots.
type Reels5x [5][]Sym

func (r *Reels5x) Cols() int {
	return 5
}

func (r *Reels5x) Reel(col Pos) []Sym {
	return r[col-1]
}

func (r *Reels5x) Reshuffles() uint64 {
	return uint64(len(r[0])) * uint64(len(r[1])) * uint64(len(r[2])) * uint64(len(r[3])) * uint64(len(r[4]))
}

func FindReels[T any](reelsmap map[float64]T, mrtp float64) (reels T, rtp float64) {
	for p, r := range reelsmap {
		if math.Abs(mrtp-p) < math.Abs(mrtp-rtp) {
			reels, rtp = r, p
		}
	}
	return
}

var (
	ErrNoWay      = errors.New("no way to here")
	ErrBetEmpty   = errors.New("bet is empty")
	ErrNoLineset  = errors.New("lines set is empty")
	ErrLinesetOut = errors.New("lines set is out of range bet lines")
	ErrNoFeature  = errors.New("feature not available")
	ErrDisabled   = errors.New("feature is disabled")
)

// Slot5x3 is base struct for all slot games with screen 5x3.
type Slot3x3 struct {
	Sel int     `json:"sel" yaml:"sel" xml:"sel"` // selected bet lines
	Bet float64 `json:"bet" yaml:"bet" xml:"bet"` // bet value

	Gain float64 `json:"gain,omitempty" yaml:"gain,omitempty" xml:"gain,omitempty"` // gain for double up games
}

func (g *Slot3x3) NewScreen() Screen {
	return NewScreen3x3()
}

func (g *Slot3x3) Spawn(screen Screen, wins Wins) {
}

func (g *Slot3x3) Prepare() {
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
		return ErrDisabled
	}
	g.Bet = bet
	return nil
}

func (g *Slot3x3) GetSel() int {
	return g.Sel
}

func (g *Slot3x3) SetSelNum(sel int, bln int) error {
	if sel < 1 {
		return ErrNoLineset
	}
	if sel > bln {
		return ErrLinesetOut
	}
	if g.FreeSpins() > 0 {
		return ErrDisabled
	}
	g.Sel = sel
	return nil
}

func (g *Slot3x3) SetMode(int) error {
	return ErrNoFeature
}

// Slot5x3 is base struct for all slot games with screen 5x3.
type Slot5x3 struct {
	Sel int     `json:"sel" yaml:"sel" xml:"sel"` // selected bet lines
	Bet float64 `json:"bet" yaml:"bet" xml:"bet"` // bet value

	Gain float64 `json:"gain,omitempty" yaml:"gain,omitempty" xml:"gain,omitempty"` // gain for double up games
}

func (g *Slot5x3) NewScreen() Screen {
	return NewScreen5x3()
}

func (g *Slot5x3) Spawn(screen Screen, wins Wins) {
}

func (g *Slot5x3) Prepare() {
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
		return ErrDisabled
	}
	g.Bet = bet
	return nil
}

func (g *Slot5x3) GetSel() int {
	return g.Sel
}

func (g *Slot5x3) SetSelNum(sel int, bln int) error {
	if sel < 1 {
		return ErrNoLineset
	}
	if sel > bln {
		return ErrLinesetOut
	}
	if g.FreeSpins() > 0 {
		return ErrDisabled
	}
	g.Sel = sel
	return nil
}

func (g *Slot5x3) SetMode(int) error {
	return ErrNoFeature
}

// Slot5x4 is base struct for all slot games with screen 5x4.
type Slot5x4 struct {
	Sel int     `json:"sel" yaml:"sel" xml:"sel"` // selected bet lines
	Bet float64 `json:"bet" yaml:"bet" xml:"bet"` // bet value

	Gain float64 `json:"gain,omitempty" yaml:"gain,omitempty" xml:"gain,omitempty"` // gain for double up games
}

func (g *Slot5x4) NewScreen() Screen {
	return NewScreen5x4()
}

func (g *Slot5x4) Spawn(screen Screen, wins Wins) {
}

func (g *Slot5x4) Prepare() {
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
		return ErrDisabled
	}
	g.Bet = bet
	return nil
}

func (g *Slot5x4) GetSel() int {
	return g.Sel
}

func (g *Slot5x4) SetSelNum(sel int, bln int) error {
	if sel < 1 {
		return ErrNoLineset
	}
	if sel > bln {
		return ErrLinesetOut
	}
	if g.FreeSpins() > 0 {
		return ErrDisabled
	}
	g.Sel = sel
	return nil
}

func (g *Slot5x4) SetMode(int) error {
	return ErrNoFeature
}
