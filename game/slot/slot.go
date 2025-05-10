package slot

import (
	"errors"
	"fmt"
)

type (
	Sym byte // symbol type
	Pos int8 // screen or line position
)

type Reels interface {
	Cols() int          // returns number of columns
	Reel(col Pos) []Sym // returns reel at given column, index from
	Reshuffles() uint64 // returns total number of reshuffles
	fmt.Stringer
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
	Bon  any     `json:"bon,omitempty" yaml:"bon,omitempty" xml:"bon,omitempty"`         // bonus game data
	JID  int     `json:"jid,omitempty" yaml:"jid,omitempty" xml:"jid,omitempty,attr"`    // jackpot identifier
	Jack float64 `json:"jack,omitempty" yaml:"jack,omitempty" xml:"jack,omitempty,attr"` // jackpot win
}

// Progressive jackpot calculated as P * Bet / JackBasis * JackFund
// where P - is the reciprocal of the probability of occurrence.
// Maximum P=25000000 with maximum Bet=10.
const JackBasis = 250_000_000

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

// Total jackpot for spin.
func (wins Wins) Jackpot() float64 {
	var sum float64
	for _, wi := range wins {
		sum += wi.Jack
	}
	return sum
}

// SlotGame is common slots interface. Any slot game should implement this interface.
type SlotGame interface {
	Clone() SlotGame              // returns full cloned copy of itself
	Scanner(*Wins) error          // scan given screen and append result to wins, constant function
	Cost() (float64, bool)        // cost of spin on current bet and lines, and has it jackpot rate, constant function
	Free() bool                   // returns true on spins without pay, constant function
	Spin(float64)                 // fill the screen with random hits on reels closest to given RTP, constant function
	Spawn(Wins, float64, float64) // setup bonus games to wins results, constant function
	Prepare()                     // update game state before new spin
	Apply(Wins)                   // update game state to spin results
	GetGain() float64             // returns gain for double up games, constant function
	SetGain(float64) error        // set gain to given value on double up games
	GetBet() float64              // returns current bet, constant function
	SetBet(float64) error         // set bet to given value
	GetSel() int                  // returns number of selected bet lines, constant function
	SetSel(int) error             // setup number of selected bet lines
	SetMode(int) error            // change game mode depending on the user's choice
}

type ClassicSlot interface {
	Screen
	SlotGame
}

// Reels for 3-reels slots.
type Reels3x [3][]Sym

// Declare conformity with Reels interface.
var _ Reels = (*Reels3x)(nil)

func (r *Reels3x) Cols() int {
	return 3
}

func (r *Reels3x) Reel(col Pos) []Sym {
	return r[col-1]
}

func (r *Reels3x) Reshuffles() uint64 {
	return uint64(len(r[0])) * uint64(len(r[1])) * uint64(len(r[2]))
}

func (r *Reels3x) String() string {
	return fmt.Sprintf("[%d, %d, %d]", len(r[0]), len(r[1]), len(r[2]))
}

// Reels for 4-reels slots.
type Reels4x [4][]Sym

// Declare conformity with Reels interface.
var _ Reels = (*Reels4x)(nil)

func (r *Reels4x) Cols() int {
	return 4
}

func (r *Reels4x) Reel(col Pos) []Sym {
	return r[col-1]
}

func (r *Reels4x) Reshuffles() uint64 {
	return uint64(len(r[0])) * uint64(len(r[1])) * uint64(len(r[2])) * uint64(len(r[3]))
}

func (r *Reels4x) String() string {
	return fmt.Sprintf("[%d, %d, %d, %d]", len(r[0]), len(r[1]), len(r[2]), len(r[3]))
}

// Reels for 5-reels slots.
type Reels5x [5][]Sym

// Declare conformity with Reels interface.
var _ Reels = (*Reels5x)(nil)

func (r *Reels5x) Cols() int {
	return 5
}

func (r *Reels5x) Reel(col Pos) []Sym {
	return r[col-1]
}

func (r *Reels5x) Reshuffles() uint64 {
	return uint64(len(r[0])) * uint64(len(r[1])) * uint64(len(r[2])) * uint64(len(r[3])) * uint64(len(r[4]))
}

func (r *Reels5x) String() string {
	return fmt.Sprintf("[%d, %d, %d, %d, %d]", len(r[0]), len(r[1]), len(r[2]), len(r[3]), len(r[4]))
}

// Reels for 6-reels slots.
type Reels6x [6][]Sym

// Declare conformity with Reels interface.
var _ Reels = (*Reels6x)(nil)

func (r *Reels6x) Cols() int {
	return 6
}

func (r *Reels6x) Reel(col Pos) []Sym {
	return r[col-1]
}

func (r *Reels6x) Reshuffles() uint64 {
	return uint64(len(r[0])) * uint64(len(r[1])) * uint64(len(r[2])) * uint64(len(r[3])) * uint64(len(r[4])) * uint64(len(r[5]))
}

func (r *Reels6x) String() string {
	return fmt.Sprintf("[%d, %d, %d, %d, %d, %d]", len(r[0]), len(r[1]), len(r[2]), len(r[3]), len(r[4]), len(r[5]))
}

var (
	ErrNoWay      = errors.New("no way to here")
	ErrBetEmpty   = errors.New("bet is empty")
	ErrNoLineset  = errors.New("lines set is empty")
	ErrLinesetOut = errors.New("lines set is out of range bet lines")
	ErrNoFeature  = errors.New("feature not available")
	ErrDisabled   = errors.New("feature is disabled")
)

// Slotx is base struct for all slot games with subsequent screen.
type Slotx struct {
	Sel int     `json:"sel" yaml:"sel" xml:"sel"` // selected bet lines
	Bet float64 `json:"bet" yaml:"bet" xml:"bet"` // bet value

	// gain for double up games
	Gain float64 `json:"gain,omitempty" yaml:"gain,omitempty" xml:"gain,omitempty"`
	// free spin number
	FSN int `json:"fsn,omitempty" yaml:"fsn,omitempty" xml:"fsn,omitempty"`
	// free spin remains
	FSR int `json:"fsr,omitempty" yaml:"fsr,omitempty" xml:"fsr,omitempty"`
}

func (g *Slotx) Cost() (float64, bool) {
	return g.Bet * float64(g.Sel), false
}

func (g *Slotx) Free() bool {
	return g.FSR != 0
}

func (g *Slotx) Spawn(wins Wins, fund, mrtp float64) {
}

func (g *Slotx) Prepare() {
}

func (g *Slotx) Apply(wins Wins) {
	if g.FSR != 0 {
		g.Gain += wins.Gain()
		g.FSN++
	} else {
		g.Gain = wins.Gain()
		g.FSN = 0
	}

	if g.FSR > 0 {
		g.FSR--
	}
	for _, wi := range wins {
		if wi.Free > 0 {
			g.FSR += wi.Free
		}
	}
}

func (g *Slotx) GetGain() float64 {
	return g.Gain
}

func (g *Slotx) SetGain(gain float64) error {
	g.Gain = gain
	return nil
}

func (g *Slotx) GetBet() float64 {
	return g.Bet
}

func (g *Slotx) SetBet(bet float64) error {
	if bet <= 0 {
		return ErrBetEmpty
	}
	if bet == g.Bet {
		return nil
	}
	if g.FSR != 0 {
		return ErrDisabled
	}
	g.Bet = bet
	return nil
}

func (g *Slotx) GetSel() int {
	return g.Sel
}

func (g *Slotx) SetSelNum(sel int, bln int) error {
	if sel < 1 {
		return ErrNoLineset
	}
	if sel > bln {
		return ErrLinesetOut
	}
	if sel == g.Sel {
		return nil
	}
	if g.FSR != 0 {
		return ErrDisabled
	}
	g.Sel = sel
	return nil
}

func (g *Slotx) SetMode(int) error {
	return ErrNoFeature
}
