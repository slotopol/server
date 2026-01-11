package jaguarmoon

// See: https://www.slotsmate.com/software/novomatic/jaguar-moon

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsBon slot.Reelx

var ReelsMap slot.ReelsMap[slot.Reelx]

const (
	sn         = 12   // number of symbols
	wild, scat = 1, 2 // wild & scatter symbol IDs
	linemin    = 3    // minimum line symbols to win
)

// Lined payment.
var LinePay = [sn][5]float64{
	{},                   //  1 wild    (2, 3, 4 reels only)
	{},                   //  2 scatter (1, 2, 3 reels only)
	{0, 0, 40, 200, 800}, //  3 wooman
	{0, 0, 20, 60, 200},  //  4 panther
	{0, 0, 10, 50, 100},  //  5 footprint
	{0, 0, 10, 30, 100},  //  6 rings
	{0, 0, 4, 10, 50},    //  7 ace
	{0, 0, 4, 10, 50},    //  8 king
	{0, 0, 4, 10, 50},    //  9 queen
	{0, 0, 2, 8, 40},     // 10 jack
	{0, 0, 2, 8, 40},     // 11 ten
	{0, 0, 2, 8, 40},     // 12 nine
}

// Scatter freespins table
var ScatFreespin = [6]int{0, 0, 8, 12, 15, 20} // 2 scatter

// Free games multipliers
var FreeMult = [6]float64{0, 0, 2, 3, 4, 5}

type Game struct {
	slot.Screen5x3 `yaml:",inline"`
	slot.Slotx     `yaml:",inline"`
	// multiplier on freespins
	M float64 `json:"m,omitempty" yaml:"m,omitempty" xml:"m,omitempty"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slotx: slot.Slotx{
			Bet: 1,
		},
		M: 0,
	}
}

func (g *Game) Clone() slot.SlotGame {
	var clone = *g
	return &clone
}

func (g *Game) Scanner(wins *slot.Wins) error {
	// Count symbols
	var counts [5 + 1][sn + 1]int
	for x := range 5 {
		var r = g.Scr[x]
		counts[x][r[0]]++
		counts[x][r[1]]++
		counts[x][r[2]]++
	}
	// Ways calculation
	var combs = [sn + 1]int{0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	for x, cx := range counts {
		for sym, c := range combs {
			var n = cx[sym] + cx[wild]
			combs[sym] = c * n
			if x >= linemin && c > 0 && n == 0 {
				var pay = LinePay[sym-1][x-1]
				var mm float64 = 1 // mult mode
				if g.FSR > 0 {
					mm = g.M
				}
				*wins = append(*wins, slot.WinItem{
					Pay: g.Bet * pay,
					MP:  mm * float64(c),
					Sym: slot.Sym(sym),
					Num: slot.Pos(x),
					LI:  243,
					XY:  g.SymPosL2(slot.Pos(x), slot.Sym(sym), wild),
				})
			}
		}
	}
	// Scatters calculation
	if counts[0][scat] > 0 && counts[1][scat] > 0 && counts[2][scat] > 0 {
		var count = counts[0][scat] + counts[1][scat] + counts[2][scat]
		var fs = ScatFreespin[count-1]
		*wins = append(*wins, slot.WinItem{
			Sym: scat,
			Num: slot.Pos(count),
			XY:  g.SymPos(scat),
			FS:  fs,
		})
	}
	return nil
}

func (g *Game) Cost() float64 {
	return g.Bet * 10
}

func (g *Game) Spin(mrtp float64) {
	if g.FSR == 0 {
		var reels, _ = ReelsMap.FindClosest(mrtp)
		g.SpinReels(reels)
	} else {
		g.SpinReels(ReelsBon)
	}
}

func (g *Game) Prepare() {
	if g.FSR == 0 {
		g.M = 0
	}
}

func (g *Game) Apply(wins slot.Wins) {
	g.Slotx.Apply(wins)
	for _, wi := range wins {
		if wi.FS > 0 {
			g.M = FreeMult[wi.Num-1]
		}
	}
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
