package ultrasevens

// See: https://www.slotsmate.com/software/novomatic/ultra-sevens

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsMap slot.ReelsMap[*slot.Reels5x]

var JackMap slot.ReelsMap[[3]float64]

// Lined payment.
var LinePay = [7][5]float64{
	{0, 10, 100, 1000, 10000}, // 1 seven
	{0, 0, 40, 200, 500},      // 2 melon
	{0, 0, 40, 200, 500},      // 3 grapes
	{0, 0, 10, 50, 200},       // 4 plum
	{0, 0, 10, 50, 200},       // 5 orange
	{0, 0, 10, 50, 200},       // 6 lemon
	{0, 5, 10, 50, 200},       // 7 cherry
}

// Bet lines
var BetLines = slot.BetLinesNvm5x4[:40]

type Game struct {
	slot.Screen5x4 `yaml:",inline"`
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

const (
	ssj1 = 1
	ssj2 = 2
	ssj3 = 3
)

func (g *Game) Filled() slot.Sym {
	var sym = g.Scr[4][3]
	for x := range 5 {
		for y := range 4 {
			if g.Scr[x][y] != sym {
				return 0
			}
		}
	}
	return sym
}

func (g *Game) Scanner(wins *slot.Wins) (err error) {
	switch sym := g.Filled(); sym {
	case 1:
		*wins = append(*wins, slot.WinItem{
			Sym: sym,
			JID: ssj1,
		})
		return
	case 2, 3:
		*wins = append(*wins, slot.WinItem{
			Sym: sym,
			JID: ssj2,
		})
		return
	case 4, 5, 6, 7:
		*wins = append(*wins, slot.WinItem{
			Sym: sym,
			JID: ssj3,
		})
		return
	}
	for li, line := range BetLines[:g.Sel] {
		var numl slot.Pos = 5
		var syml = g.LY(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = g.LY(x, line)
			if sx != syml {
				numl = x - 1
				break
			}
		}

		if pay := LinePay[syml-1][numl-1]; pay > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * pay,
				MP:  1,
				Sym: syml,
				Num: numl,
				LI:  li + 1,
				XY:  line.HitxL(numl),
			})
		}
	}
	return
}

func (g *Game) JackFreq(mrtp float64) []float64 {
	var bulk, _ = JackMap.FindClosest(mrtp)
	return bulk[:]
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.ReelSpin(reels)
}

func (g *Game) Spawn(wins slot.Wins, fund, mrtp float64) {
	for i, wi := range wins {
		if wi.JID != 0 {
			var bulk, _ = JackMap.FindClosest(mrtp)
			var jf = min(bulk[wi.JID-1]*g.Bet/slot.JackBasis, 1)
			wins[i].JR = jf * fund
		}
	}
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
