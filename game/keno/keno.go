package keno

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"

	"github.com/slotopol/server/util"
)

type Paytable [11][11]float64

func (kp *Paytable) Pay(sel, hit int) float64 {
	return kp[sel][hit]
}

func (kp *Paytable) HasSel(sel int) bool {
	for _, pay := range kp[sel] {
		if pay > 0 {
			return true
		}
	}
	return false
}

func (kp *Paytable) Scanner(scrn *Screen, wins *Wins, bet float64) error {
	wins.Sel = 0
	wins.Num = 0
	for i := range 80 {
		if scrn[i]&KSsel > 0 {
			wins.Sel++
			if scrn[i]&KShit > 0 {
				wins.Num++
			}
		}
	}
	wins.Pay = kp[wins.Sel][wins.Num] * bet
	return nil
}

func (kp *Paytable) CalcStat(ctx context.Context) float64 {
	var rtp, lines float64
	for n := 1; n <= 10; n++ {
		var nrtp float64
		for r := 0; r <= n; r++ {
			var pay = kp[n][r]
			nrtp += pay * Prob(n, r)
		}
		if nrtp > 0 {
			fmt.Printf("RTP[%2d] = %.6f%%\n", n, nrtp*100)
			rtp += nrtp
			lines++
		}
	}
	rtp *= 100 / lines
	fmt.Printf("RTP[game] = %.6f%%", rtp)
	return rtp
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

type Bitset = util.Bitset128

var MakeBitNum = util.MakeBitNum128

type Wins struct {
	Sel int     `json:"sel" yaml:"sel" xml:"sel,attr"`
	Num int     `json:"num" yaml:"num" xml:"num,attr"`
	Pay float64 `json:"pay" yaml:"pay" xml:"pay,attr"`
}

// KenoGame is common keno interface. Any keno game should implement this interface.
type KenoGame interface {
	Scanner(*Wins) error  // scan given screen and append result to wins, constant function
	Spin(float64)         // fill the screen with random hits on reels closest to given RTP, constant function
	GetBet() float64      // returns current bet, constant function
	SetBet(float64) error // set bet to given value
	GetSel() Bitset       // returns current selected numbers, constant function
	SetSel(Bitset) error  // set current selected numbers
}

var (
	ErrBetEmpty      = errors.New("bet is empty")
	ErrKenoNotEnough = errors.New("no pays with this selected numbers")
	ErrKenoTooMany   = errors.New("too many numbers selected, not more than 10 expected")
	ErrKenoOutRange  = errors.New("some of given number is out of range 1..80")
)

type Keno80 struct {
	Scr Screen  `json:"scr" yaml:"scr" xml:"scr"` // game screen
	Bet float64 `json:"bet" yaml:"bet" xml:"bet"` // bet value
	Sel Bitset  `json:"sel" yaml:"sel" xml:"sel"` // selected numbers
}

func (g *Keno80) Spin(_ float64) {
	var hits [80]int
	for i := range 80 {
		hits[i] = i + 1
	}
	rand.Shuffle(80, func(i, j int) {
		hits[i], hits[j] = hits[j], hits[i]
	})

	clear(g.Scr[:])
	for n := range g.Sel.Bits() {
		g.Scr[n-1] = KSsel
	}
	for i := range 20 {
		g.Scr[hits[i]-1] |= KShit
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

func (g *Keno80) CheckSel(sel Bitset, kp *Paytable) error {
	if sel.Num() > len(kp)-1 {
		return ErrKenoTooMany
	}
	if !kp.HasSel(sel.Num()) {
		return ErrKenoNotEnough
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
