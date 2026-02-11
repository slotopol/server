package hellshot7s

// See: https://www.slotsmate.com/software/ct-interactive/hells-hot-7s

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [10][5]float64{
	{},                    //  1 wild (2, 3, 4 reels only)
	{},                    //  2 scatter
	{0, 0, 50, 150, 1000}, //  3 bar
	{0, 0, 30, 50, 500},   //  4 grapes
	{0, 0, 30, 50, 200},   //  5 apple
	{0, 0, 10, 20, 150},   //  6 pear
	{0, 0, 5, 20, 50},     //  7 plum
	{0, 0, 5, 20, 50},     //  8 orange
	{0, 0, 5, 20, 50},     //  9 lemon
	{0, 0, 5, 20, 50},     // 10 cherry
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 10, 30, 200} // 2 scatter

// Bet lines
var BetLines = slot.BetLinesCT5x3[:]

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

const wild, scat = 1, 2

func (g *Game) Scanner(wins *slot.Wins) error {
	var mw float64 = 1 // mult wild
	for range g.SymNum(wild) {
		mw *= 2
	}

	// Lined symbols calculation.
	for li, line := range BetLines[:g.Sel] {
		var numl slot.Pos = 5
		var syml = g.LX(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = g.LX(x, line)
			if sx != syml && sx != wild {
				numl = x - 1
				break
			}
		}

		if pay := LinePay[syml-1][numl-1]; pay > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * pay,
				MP:  mw,
				Sym: syml,
				Num: numl,
				LI:  li + 1,
				XY:  line.HitxL(numl),
			})
		}
	}

	// Scatters calculation.
	if count := g.SymNum(scat); count >= 3 {
		var pay = ScatPay[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * float64(g.Sel) * pay,
			MP:  mw,
			Sym: scat,
			Num: count,
			XY:  g.SymPos(scat),
		})
	}

	return nil
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.SpinReels(reels)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrDisabled
}
