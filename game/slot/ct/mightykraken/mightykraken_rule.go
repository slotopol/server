package mightykraken

// See: https://www.slotsmate.com/software/ct-interactive/mighty-kraken
// similar: ct/dancingbananas
// (difference on reels sets)

import (
	"github.com/slotopol/server/game/slot"
)

const (
	sn                 = 11      // number of symbols
	wild, scat1, scat2 = 1, 2, 3 // wild & scatter symbol IDs
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [sn][5]float64{
	{},                     //  1 wild (on 2, 3, 4 reels)
	{},                     //  2 star (on all reels)
	{},                     //  3 dollar (on 1, 3, 5 reels)
	{0, 10, 50, 200, 3000}, //  4 captain
	{0, 0, 40, 100, 500},   //  5 compass
	{0, 0, 40, 100, 500},   //  6 anchor
	{0, 0, 20, 50, 200},    //  7 bell
	{0, 0, 10, 30, 100},    //  8 spades
	{0, 0, 10, 30, 100},    //  9 hearts
	{0, 0, 10, 30, 100},    // 10 diamonds
	{0, 0, 10, 30, 100},    // 11 clubs
}

// Scatters payment.
var ScatPay1 = [5]float64{0, 0, 3, 20, 100} // 2 star
var ScatPay2 = [5]float64{0, 0, 20}         // 3 dollar

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

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	var reelwild [5]bool
	for x := 1; x < 4; x++ { // 2, 3, 4 reels only
		for _, sy := range g.Grid[x] {
			if sy == wild {
				reelwild[x] = true
				break
			}
		}
	}

	for li, line := range BetLines[:g.Sel] {
		var numl slot.Pos = 5
		var syml = g.LX(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = g.LX(x, line)
			if reelwild[x-1] {
				continue
			} else if sx != syml {
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
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.SymNum(scat1); count >= 3 {
		var pay = ScatPay1[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * float64(g.Sel) * pay,
			MP:  1,
			Sym: scat1,
			Num: count,
			XY:  g.SymPos(scat1),
		})
	} else if count := g.SymNum(scat2); count >= 3 {
		var pay = ScatPay2[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * float64(g.Sel) * pay,
			MP:  1,
			Sym: scat2,
			Num: count,
			XY:  g.SymPos(scat2),
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.SpinReels(reels)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrDisabled
}
