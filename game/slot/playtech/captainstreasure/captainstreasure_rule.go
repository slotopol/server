package captainstreasure

// See: https://freeslotshub.com/playtech/captains-treasure/
// See: https://www.slotsmate.com/software/playtech/captain-treasure

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [11][5]float64{
	{},                      //  1 wild (2, 3, 4 reels only)
	{},                      //  2 scatter
	{2, 10, 100, 500, 5000}, //  3 sabers
	{0, 5, 50, 250, 2500},   //  4 map
	{0, 3, 20, 100, 1000},   //  5 anchor
	{0, 0, 10, 30, 500},     //  6 ace
	{0, 0, 5, 25, 300},      //  7 king
	{0, 0, 5, 20, 200},      //  8 queen
	{0, 0, 5, 20, 200},      //  9 jack
	{0, 0, 5, 15, 100},      // 10 ten
	{0, 0, 5, 15, 100},      // 11 nine
}

// Scatters payment.
var ScatPay = [5]float64{0, 1, 5, 10, 100} //  2 suitcase

// Bet lines
var BetLines = []slot.Linex{
	{2, 2, 2, 2, 2}, // 1
	{1, 1, 1, 1, 1}, // 2
	{3, 3, 3, 3, 3}, // 3
	{1, 2, 3, 2, 1}, // 4
	{3, 2, 1, 2, 3}, // 5
	{1, 1, 2, 1, 1}, // 6
	{3, 3, 2, 3, 3}, // 7
	{2, 1, 1, 1, 2}, // 8
	{2, 3, 3, 3, 2}, // 9
}

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
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	for li, line := range BetLines[:g.Sel] {
		var mw float64 = 1 // mult wild
		var numl slot.Pos = 5
		var syml = g.LX(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = g.LX(x, line)
			if sx == wild {
				mw = 2
			} else if sx != syml {
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

		if numl < 5 {
			var mw float64 = 1 // mult wild
			var numr slot.Pos = 4
			var symr = g.LX(5, line)
			var x slot.Pos
			for x = 4; x > 1; x-- {
				var sx = g.LX(x, line)
				if sx == wild {
					mw = 2
				} else if sx != symr {
					numr = 5 - x
					break
				}
			}

			if pay := LinePay[symr-1][numr-1]; pay > 0 {
				*wins = append(*wins, slot.WinItem{
					Pay: g.Bet * pay,
					MP:  mw,
					Sym: symr,
					Num: numr,
					LI:  li + 1,
					XY:  line.HitxR5(numr),
				})
			}
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.SymNum(scat); count >= 2 {
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
