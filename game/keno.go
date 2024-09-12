package game

import (
	"errors"
	"math/rand/v2"
)

type KenoPaytable [9][11]float64

func (kp *KenoPaytable) Pay(sel, win int) float64 {
	return kp[sel-2][win]
}

// Keno ball type
type KB byte

const (
	KBempty  KB = 0             // empty cell
	KBsel    KB = 0x1           // cell with selection without hit
	KBhit    KB = 0x2           // cell with hit without selection
	KBselhit KB = KBsel | KBhit // win cell, hit and selection
)

type KenoScreen [80]KB

type KWins struct {
	Num  int
	Pay  float64
	Hits [20]int
}

type KenoGame interface {
	Scanner(*KenoScreen, *Wins) // scan given screen and set result to wins, constat function
	Spin(*KenoScreen, []int)    // fill the screen with random hits, constat function
	GetBet() float64            // returns current bet, constat function
	SetBet(float64) error       // set bet to given value
	GetSel() []int              // returns current selected numbers, constat function
	SetSel([]int) error         // set current selected numbers
}

var (
	ErrKenoNotEnough = errors.New("not enough numbers selected, minimum 2 expected")
	ErrKenoTooMany   = errors.New("too many numbers selected, not more than 10 expected")
	ErrKenoOutRange  = errors.New("some of given number is out of range 1..80")
	ErrKenoRepeat    = errors.New("some numbers are repeated")
)

type Keno80 struct {
	Bet float64 `json:"bet" yaml:"bet" xml:"bet"` // bet value
	Sel []int   `json:"sel" yaml:"sel" xml:"sel"` // selected numbers
}

func (g *Keno80) Spin(ks *KenoScreen, hits []int) {
	for i := range 80 {
		hits[i] = i + 1
	}
	rand.Shuffle(80, func(i, j int) {
		hits[i], hits[j] = hits[j], hits[i]
	})

	clear(ks[:])
	for _, n := range g.Sel {
		ks[n] = KBsel
	}
	for i := range 20 {
		ks[hits[i]] |= KBhit
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

func (g *Keno80) GetSel() []int {
	return g.Sel
}

func (g *Keno80) SetSel(sel []int) error {
	if len(sel) < 2 {
		return ErrKenoNotEnough
	}
	if len(sel) > 10 {
		return ErrKenoTooMany
	}
	var m = make(map[int]struct{}, len(sel))
	for _, n := range sel {
		if n < 1 || n > 80 {
			return ErrKenoOutRange
		}
		if _, ok := m[n]; ok {
			return ErrKenoRepeat
		}
		m[n] = struct{}{}
	}
	g.Sel = sel
	return nil
}
