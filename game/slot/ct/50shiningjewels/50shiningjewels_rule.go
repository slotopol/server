package fiftyshiningjewels

// See: https://www.slotsmate.com/software/ct-interactive/50-shining-jewels

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [10][5]float64{
	{},                    //  1 wild (2, 3, 4 reels only)
	{},                    //  2 scatter
	{0, 0, 60, 200, 1000}, //  3 seven
	{0, 0, 40, 100, 300},  //  4 heart
	{0, 0, 20, 80, 200},   //  5 diamond
	{0, 0, 20, 80, 200},   //  6 emerald
	{0, 0, 8, 40, 100},    //  7 gold
	{0, 0, 8, 40, 100},    //  8 blue
	{0, 0, 8, 40, 100},    //  9 yellow
	{0, 0, 8, 40, 100},    // 10 lilac
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 2, 15, 500} // 2 scatter

// Bet lines
var BetLines = slot.BetLinesCT5x4v2[:]

type Game struct {
	slot.Grid5x4 `yaml:",inline"`
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
	return slot.ErrDisabled
}
