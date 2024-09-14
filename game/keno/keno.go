package game

import (
	"errors"
	"math/rand/v2"

	"github.com/slotopol/server/util"
)

type KenoPaytable [11][11]float64

func (kp *KenoPaytable) Pay(sel, win int) float64 {
	return kp[sel][win]
}

// Keno spot type
type KS byte

const (
	KSempty  KS = 0             // empty cell
	KSsel    KS = 0x1           // cell with selection without hit
	KShit    KS = 0x2           // cell with hit without selection
	KSselhit KS = KSsel | KShit // win cell, hit and selection
)

type Screen [80]KS

type Wins struct {
	Num  int
	Pay  float64
	Hits [20]int
}

type Bitset = util.Bitset128

// KenoGame is common keno interface. Any keno game should implement this interface.
type KenoGame interface {
	Scanner(*Screen, *Wins) // scan given screen and set result to wins, constat function
	Spin(*Screen, []int)    // fill the screen with random hits, constat function
	GetBet() float64        // returns current bet, constat function
	SetBet(float64) error   // set bet to given value
	GetSel() Bitset         // returns current selected numbers, constat function
	SetSel(Bitset) error    // set current selected numbers
}

var (
	ErrBetEmpty      = errors.New("bet is empty")
	ErrKenoNotEnough = errors.New("not enough numbers selected, minimum 2 expected")
	ErrKenoTooMany   = errors.New("too many numbers selected, not more than 10 expected")
	ErrKenoOutRange  = errors.New("some of given number is out of range 1..80")
)

type Keno80 struct {
	Bet float64 `json:"bet" yaml:"bet" xml:"bet"` // bet value
	Sel Bitset  `json:"sel" yaml:"sel" xml:"sel"` // selected numbers
}

func (g *Keno80) Spin(ks *Screen, hits []int) {
	for i := range 80 {
		hits[i] = i + 1
	}
	rand.Shuffle(80, func(i, j int) {
		hits[i], hits[j] = hits[j], hits[i]
	})

	clear(ks[:])
	for _, n := range g.Sel {
		ks[n] = KSsel
	}
	for i := range 20 {
		ks[hits[i]] |= KShit
	}
}

func (g *Keno80) GetBet() float64 {
	return g.Bet
}

func (g *Keno80) SetBet(bet float64) error {
	if bet <= 0 {
		return ErrBetEmpty
	}
	g.Bet = bet
	return nil
}

func (g *Keno80) GetSel() Bitset {
	return g.Sel
}

func (g *Keno80) SetSel(sel Bitset) error {
	if len(sel) < 2 {
		return ErrKenoNotEnough
	}
	if len(sel) > 10 {
		return ErrKenoTooMany
	}
	for n := range sel.Bits() {
		if n < 1 || n > 80 {
			return ErrKenoOutRange
		}
	}
	g.Sel = sel
	return nil
}

func Combin(n, r int) float64 {
	var mi, mj float64 = 1, 1
	var i, j float64 = float64(n), 1
	for range r {
		mi *= i
		mj *= j
		i--
		j++
	}
	return mi / mj
}

var C_80_20 = Combin(80, 20)

func Prob(n, r int) float64 {
	return Combin(n, r) * Combin(80-n, 20-r) / C_80_20
}
