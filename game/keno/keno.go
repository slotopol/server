package keno

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"sort"

	"github.com/slotopol/server/game"
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

func (kp *Paytable) Scanner(grid *Grid, wins *Wins, bet float64) error {
	wins.Sel = 0
	wins.Num = 0
	for i := range 80 {
		if grid[i]&KSsel > 0 {
			wins.Sel++
			if grid[i]&KShit > 0 {
				wins.Num++
			}
		}
	}
	wins.Pay = kp[wins.Sel][wins.Num] * bet
	return nil
}

func Print_all(sp *game.ScanPar, ev, D float64) {
	if sp.IsMain() {
		fmt.Printf("RTP = %.8g%%\n", ev*100)
	}
	if sp.IsVI() {
		var sigma = math.Sqrt(D)
		var vi = game.GetZ(sp.Conf) * sigma
		fmt.Printf("sigma = %.6g, VI[%.4g%%] = %.6g (%s)\n", sigma, sp.Conf*100, vi, game.VIname5[game.VIclass5(sigma)])
	}
	if sp.IsCI() && ev < game.RTPconv {
		var sigma = math.Sqrt(D)
		var ci = game.CI(sp.Conf, ev, sigma)
		var BRci = game.BankrollPlayer(sp.Conf, ev, sigma, ci)
		fmt.Printf("CI[%.4g%%] = %d, bankroll[CI] = %.6g\n", sp.Conf*100, int(ci+0.5), BRci)
	}
	if sp.IsSpread() {
		fmt.Println()
		fmt.Printf("RTP spread for spins number with confidence %.4g%%:\n", sp.Conf*100)
		var N = []int{1e3, 1e4, 1e5, 1e6, 1e7}
		var sigma = math.Sqrt(D)
		var vi = game.GetZ(sp.Conf) * sigma
		var ci = game.CI(sp.Conf, ev, sigma)
		if ci < 1e7 {
			N = append(N, int(ci+0.5))
			sort.Ints(N)
		}
		for _, n := range N {
			var Δ = vi / math.Sqrt(float64(n))
			fmt.Printf("%8d: %.2f%% ... %.2f%%\n", n, (ev-Δ)*100, (ev+Δ)*100)
		}
	}
}

func (kp *Paytable) CalcStat(ctx context.Context, sp *game.ScanPar) (float64, float64) {
	var ev5, D5 float64
	for n := 1; n <= 10; n++ {
		var ev, e2 float64
		for k := 0; k <= n; k++ {
			var pay = kp[n][k]
			ev += game.KenoProb(n, k) * pay
			e2 += game.KenoProb(n, k) * pay * pay
		}
		if ev > 0 {
			var D = e2 - ev*ev
			fmt.Printf("\n%d cells selected\n", n)
			Print_all(sp, ev, D)
			if n == 5 {
				ev5, D5 = ev, D
			}
		}
	}
	return ev5, D5
}

// Keno spot type
type KS byte

const (
	KSempty  KS = 0             // empty cell
	KSsel    KS = 0x1           // cell with selection without hit
	KShit    KS = 0x2           // cell with hit without selection
	KSselhit KS = KSsel | KShit // win cell, hit and selection
)

type Grid [80]KS

type Bitset = util.Bitset128

var MakeBitNum = util.MakeBitNum128

type Wins struct {
	Sel int     `json:"sel" yaml:"sel" xml:"sel,attr"`
	Num int     `json:"num" yaml:"num" xml:"num,attr"`
	Pay float64 `json:"pay" yaml:"pay" xml:"pay,attr"`
}

// KenoGame is common keno interface. Any keno game should implement this interface.
type KenoGame interface {
	Scanner(*Wins) error  // scan given grid and append result to wins, constant function
	Spin(float64)         // fill the grid with random hits on reels closest to given RTP, constant function
	GetBet() float64      // returns current bet, constant function
	SetBet(float64) error // set bet to given value
	GetSel() Bitset       // returns current selected numbers, constant function
	SetSel(Bitset) error  // set current selected numbers
}

var (
	ErrBadParam      = errors.New("wrong parameter") // parameter is not acceptable
	ErrKenoNotEnough = errors.New("no pays with this selected numbers")
	ErrKenoTooMany   = errors.New("too many numbers selected, not more than 10 expected")
	ErrKenoOutRange  = errors.New("some of given number is out of range 1..80")
)

type Keno80 struct {
	Grid Grid    `json:"grid" yaml:"grid" xml:"grid"` // game grid
	Bet  float64 `json:"bet" yaml:"bet" xml:"bet"`    // bet value
	Sel  Bitset  `json:"sel" yaml:"sel" xml:"sel"`    // selected numbers
}

func (g *Keno80) Spin(_ float64) {
	var hits [80]int
	for i := range 80 {
		hits[i] = i + 1
	}
	rand.Shuffle(80, func(i, j int) {
		hits[i], hits[j] = hits[j], hits[i]
	})

	clear(g.Grid[:])
	for n := range g.Sel.Bits() {
		g.Grid[n-1] = KSsel
	}
	for i := range 20 {
		g.Grid[hits[i]-1] |= KShit
	}
}

func (g *Keno80) GetBet() float64 {
	return g.Bet
}

func (g *Keno80) SetBet(bet float64) error {
	if bet <= 0 {
		return ErrBadParam
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
