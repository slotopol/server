package slot

import (
	"errors"
)

// WinItem describes win on each line or by scatters.
type WinItem struct {
	Pay float64 `json:"pay,omitempty" yaml:"pay,omitempty" xml:"pay,omitempty,attr"` // payment with selected bet
	MP  float64 `json:"mp,omitempty" yaml:"mp,omitempty" xml:"mp,omitempty,attr"`    // multiplier of payment for wilds, free spins and other special cases
	Sym Sym     `json:"sym,omitempty" yaml:"sym,omitempty" xml:"sym,omitempty,attr"` // win symbol
	Num Pos     `json:"num,omitempty" yaml:"num,omitempty" xml:"num,omitempty,attr"` // number of win symbols
	LI  int     `json:"li,omitempty" yaml:"li,omitempty" xml:"li,omitempty,attr"`    // line index (0 for scatters and not lined combinations)
	XY  Hitx    `json:"xy,omitempty" yaml:"xy,omitempty,flow" xml:"xy,omitempty"`    // symbols (X, Y) positions on screen
	FS  int     `json:"fs,omitempty" yaml:"fs,omitempty" xml:"fs,omitempty,attr"`    // number of free spins
	BID int     `json:"bid,omitempty" yaml:"bid,omitempty" xml:"bid,omitempty,attr"` // bonus identifier
	Bon any     `json:"bon,omitempty" yaml:"bon,omitempty" xml:"bon,omitempty"`      // bonus game data
	JID int     `json:"jid,omitempty" yaml:"jid,omitempty" xml:"jid,omitempty,attr"` // jackpot identifier
	JR  float64 `json:"jr,omitempty" yaml:"jr,omitempty" xml:"jr,omitempty,attr"`    // jackpot rate (share of the progressive jackpot for this algorithm)
}

// Progressive jackpot calculated as P * Bet / JackBasis * JackFund
// where P - is the reciprocal value of the occurrence probability.
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
		sum += wi.Pay * wi.MP
	}
	return sum
}

// Total jackpot for spin.
func (wins Wins) Jackpot() float64 {
	var sum float64
	for _, wi := range wins {
		sum += wi.JR
	}
	return sum
}

// SlotGame is common slots interface. Any slot game should implement this interface.
type SlotGame interface {
	Clone() SlotGame              // returns full cloned copy of itself
	Scanner(*Wins) error          // scan given screen and append result to wins, constant function
	Cost() float64                // cost of spin on current bet and lines, constant function
	JackFreq(float64) []float64   // returns occurrence frequency set of progressive jackpots if it has, constant function
	FreeMode() bool               // returns true on spins without pay, constant function
	Spin(float64)                 // fill the screen with random hits on reels closest to given RTP, constant function
	Spawn(Wins, float64, float64) // setup bonus games to wins results, constant function
	Prepare()                     // update game state before new spin, screen is unknown yet
	Apply(Wins)                   // update game state to spin results, screen is calculated
	GetGain() float64             // returns gain for double up games, constant function
	SetGain(float64) error        // set gain to given value on double up games
	GetBet() float64              // returns current bet per line, constant function
	SetBet(float64) error         // set bet per line to given value
	GetSel() int                  // returns number of selected bet lines, constant function
	SetSel(int) error             // setup number of selected bet lines
	SetMode(int) error            // change game mode depending on the user's choice
}

// Remark: "Scanner" method should return error only if generated screen does not appear
// to game rules, and can not be calculated. Screen combination with error will be
// dropped out. If screen appear to rules it should be calculated in any case.
// If scanner algorithm receives wrong data - its case for panic.
// The scanning method should never change the game state.

type SlotGeneric interface {
	Screen
	SlotGame
}

var (
	ErrNoWay     = errors.New("no way to here")
	ErrBadParam  = errors.New("wrong parameter")     // parameter is not acceptable
	ErrNoFeature = errors.New("feature unavailable") // feature is unavailable by game logic
	ErrDisabled  = errors.New("feature is disabled") // feature is currently can not be applicable
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

func (g *Slotx) Cost() float64 {
	return g.Bet * float64(g.Sel)
}

func (g *Slotx) JackFreq(mrtp float64) []float64 {
	return nil
}

func (g *Slotx) FreeMode() bool {
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
		if wi.FS > 0 {
			g.FSR += wi.FS
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
		return ErrBadParam
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
	if sel < 1 || sel > bln {
		return ErrBadParam
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
