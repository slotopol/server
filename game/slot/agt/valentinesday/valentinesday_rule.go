package valentinesday

// See: https://agtsoftware.com/games/agt/valentine

import (
	"github.com/slotopol/server/game/slot"
)

const (
	sn   = 8 // number of symbols
	scat = 1 // scatter symbol IDs
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [sn][5]float64{
	{},                      // 1 scatter
	{0, 0, 100, 1000, 5000}, // 2 angel
	{0, 0, 40, 400, 1000},   // 3 nymph
	{0, 0, 24, 60, 200},     // 4 soul
	{0, 0, 20, 50, 200},     // 5 toy
	{0, 0, 20, 50, 200},     // 6 balloon
	{0, 5, 16, 40, 160},     // 7 hearts
	{0, 0, 5, 20, 100},      // 8 medallion
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 2, 12, 60} // 1 scatter

// Bet lines
var BetLines = slot.BetLinesAgt5x3[:]

type Game struct {
	slot.Grid5x3 `yaml:",inline"`
	slot.Slotx   `yaml:",inline"`
}

// Declare conformity with SlotGeneric interface.
var _ slot.SlotGeneric = (*Game)(nil)

func NewGame(sel int) *Game {
	return &Game{
		Slotx: slot.Slotx{
			Sel: sel,
			Bet: 1,
		},
	}
}

func (g *Game) Clone() slot.SlotGeneric {
	var clone = *g
	return &clone
}

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

func (g *Game) FillMult() float64 {
	var sym = g.Grid[2][1] // center symbol
	if sym < 4 || sym > 7 {
		return 1
	}
	var r *[3]slot.Sym
	if r = &g.Grid[2]; r[0] != sym || r[2] != sym {
		return 1
	}
	var n = 1
	if r = &g.Grid[1]; r[0] == sym && r[1] == sym && r[2] == sym {
		n++
		if r = &g.Grid[0]; r[0] == sym && r[1] == sym && r[2] == sym {
			n++
		}
	}
	if r = &g.Grid[3]; r[0] == sym && r[1] == sym && r[2] == sym {
		n++
		if r = &g.Grid[4]; r[0] == sym && r[1] == sym && r[2] == sym {
			n++
		}
	}
	if n < 3 {
		return 1
	}
	return float64(n - 1)
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	var fm float64 // fill mult
	for li, line := range BetLines[:g.Sel] {
		var numl slot.Pos = 5
		var syml = g.LX(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = g.LX(x, line)
			if sx != syml {
				numl = x - 1
				break
			}
		}

		if pay := LinePay[syml-1][numl-1]; pay > 0 {
			if fm == 0 { // lazy calculation
				fm = g.FillMult()
			}
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * pay,
				MP:  fm,
				Sym: syml,
				Num: numl,
				LI:  li + 1,
				XY:  line.HitxL(numl),
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.SymNum(scat); count >= 3 {
		var pay = ScatPay[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * float64(g.Sel) * pay,
			MP:  1,
			Sym: scat,
			Num: count,
			XY:  g.SymPos(scat),
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.SpinReels(reels)
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
