package halloweenhot

// See: https://www.slotsmate.com/software/ct-interactive/halloween-hot

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsMap slot.ReelsMap[*slot.Reels5x]

// Lined payment.
var LinePay = [9][5]float64{
	{},                    // 1 wild
	{},                    // 2 scatter
	{0, 0, 20, 100, 1000}, // 3 seven
	{0, 0, 10, 100, 150},  // 4 dead
	{0, 0, 5, 20, 80},     // 5 cat
	{0, 0, 5, 20, 80},     // 6 vampire
	{0, 0, 5, 20, 80},     // 7 pot
	{0, 0, 5, 20, 80},     // 8 hat
	{0, 0, 5, 20, 80},     // 9 scull
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 20, 50, 500} // 2 scatter

// Bet lines
var BetLines = slot.BetLinesCT5x3[:]

type Game struct {
	slot.Screen5x3 `yaml:",inline"`
	slot.Slotx     `yaml:",inline"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slotx: slot.Slotx{
			Sel: len(BetLines),
			Bet: 1,
		},
	}
}

func (g *Game) Clone() slot.SlotGame {
	var clone = *g
	return &clone
}

const wild, scat = 1, 2

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

func (g *Game) FillMult() float64 {
	var sym = g.Scr[2][1] // center symbol
	var r *[3]slot.Sym
	if r = &g.Scr[2]; r[0] != sym || r[2] != sym {
		return 1
	}
	var n = 1
	if r = &g.Scr[1]; r[0] == sym && r[1] == sym && r[2] == sym {
		n++
		if r = &g.Scr[0]; r[0] == sym && r[1] == sym && r[2] == sym {
			n++
		}
	}
	if r = &g.Scr[3]; r[0] == sym && r[1] == sym && r[2] == sym {
		n++
		if r = &g.Scr[4]; r[0] == sym && r[1] == sym && r[2] == sym {
			n++
		}
	}
	if n < 3 {
		return 1
	}
	return float64(n)
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	var fm float64 // fill mult
	for li, line := range BetLines[:g.Sel] {
		var numl slot.Pos = 5
		var syml = g.LY(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = g.LY(x, line)
			if sx != syml && sx != wild {
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
	g.ReelSpin(reels)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
