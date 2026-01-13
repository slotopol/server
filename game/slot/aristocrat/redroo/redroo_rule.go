package redroo

// See: https://freeslotshub.com/aristocrat/big-red/

import (
	"math/rand/v2"

	"github.com/slotopol/server/game/slot"
)

const (
	sn         = 13   // number of symbols
	wild, scat = 1, 2 // wild & scatter symbol IDs
	linemin    = 2    // minimum line symbols to win
	prob2x     = 0.5  // probability of 2x multiplier for wild at free games
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [sn][5]float64{
	{},                     //  1 wild (2, 3, 4 reels only)
	{},                     //  2 scatter
	{0, 50, 150, 200, 250}, //  3 redroo
	{0, 20, 80, 150, 200},  //  4 girl
	{0, 20, 80, 150, 200},  //  5 boy
	{0, 10, 40, 100, 150},  //  6 dog
	{0, 10, 40, 100, 150},  //  7 parrot
	{0, 0, 10, 50, 140},    //  8 ace
	{0, 0, 10, 50, 140},    //  9 king
	{0, 0, 5, 40, 120},     // 10 queen
	{0, 0, 5, 40, 120},     // 11 jack
	{0, 0, 5, 20, 100},     // 12 ten
	{0, 2, 5, 20, 100},     // 13 nine
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 80, 400, 800} // scatter

// Scatter freespins table on regular games
var ScatFreespinReg = [5]int{0, 0, 8, 15, 20} // scatter

// Scatter freespins table on bonus games
var ScatFreespinBon = [5]int{0, 5, 8, 15, 20} // scatter

type Game struct {
	slot.Screen5x4 `yaml:",inline"`
	slot.Slotx     `yaml:",inline"`
	// wild multipliers
	MW [3]float64 `json:"mw" yaml:"mw" xml:"mw"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slotx: slot.Slotx{
			Bet: 1,
		},
		MW: [3]float64{1, 1, 1},
	}
}

func (g *Game) Clone() slot.SlotGame {
	var clone = *g
	return &clone
}

func (g *Game) Scanner(wins *slot.Wins) error {
	// Count symbols
	var counts [5 + 1][sn + 1]int // symbol counts per reel
	var active [sn + 1]bool       // symbols present on 1st reel
	for x := range 5 {
		var cx = &counts[x]
		for _, sym := range g.Scr[x] {
			cx[sym]++
		}
	}
	for _, sym := range g.Scr[0] {
		active[sym] = true
	}
	// Ways calculation
	var sym slot.Sym
	for sym = 3; sym <= sn; sym++ { // ignore wild & scatter
		if !active[sym] {
			continue
		}
		var c float64 = 1 // current ways
		for x, cx := range counts {
			var n = float64(cx[sym])
			if x >= 1 && x <= 3 {
				n += float64(cx[wild]) * g.MW[x-1]
			}
			if n == 0 {
				if x >= linemin {
					if pay := LinePay[sym-1][x-1]; pay > 0 {
						*wins = append(*wins, slot.WinItem{
							Pay: g.Bet * pay,
							MP:  c,
							Sym: sym,
							Num: slot.Pos(x),
							LI:  1024,
							XY:  g.SymPosL2(slot.Pos(x), sym, wild),
						})
					}
				}
				break
			}
			c *= n
		}
	}
	// Scatters calculation
	if count := counts[0][scat] + counts[1][scat] + counts[2][scat] + counts[3][scat] + counts[4][scat]; count >= 2 {
		var pay = ScatPay[count-1]
		var fs int
		if g.FSR > 0 {
			fs = ScatFreespinBon[count-1]
		} else {
			fs = ScatFreespinReg[count-1]
		}
		if pay > 0 || fs > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * pay,
				MP:  1,
				Sym: scat,
				Num: slot.Pos(count),
				XY:  g.SymPos(scat),
				FS:  fs,
			})
		}
	}
	return nil
}

func (g *Game) Cost() float64 {
	return g.Bet * 50
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.SpinReels(reels)
}

func (g *Game) Prepare() {
	if g.FSR > 0 {
		for x := range g.MW {
			if rand.Float64() < prob2x {
				g.MW[x] = 2
			} else {
				g.MW[x] = 3
			}
		}
	} else {
		g.MW = [3]float64{1, 1, 1}
	}
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
