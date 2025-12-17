package halloweenfruits

// See: https://www.slotsmate.com/software/ct-interactive/ct-gaming-halloween-fruits

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsBon slot.Reelx

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [12][5]float64{
	{},                  //  1 wild (2, 3, 4, 5 reels only)
	{},                  //  2 scatter (2, 3, 4 reels only)
	{0, 0, 20, 50, 300}, //  3 witch
	{0, 0, 15, 30, 100}, //  4 cat
	{0, 0, 15, 30, 100}, //  5 banana
	{0, 0, 15, 30, 100}, //  6 grape
	{0, 0, 10, 15, 50},  //  7 apple
	{0, 0, 10, 15, 50},  //  8 melon
	{0, 0, 10, 15, 30},  //  9 orange
	{0, 0, 10, 15, 30},  // 10 lemon
	{0, 0, 10, 15, 30},  // 11 plum
	{0, 0, 10, 15, 30},  // 12 cherry
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 0, 3, 5} // 2 scatter

// Bet lines
var BetLines = slot.BetLinesCT5x3[:]

type Game struct {
	slot.Screen5x3 `yaml:",inline"`
	slot.Slotx     `yaml:",inline"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame(sel int) *Game {
	return &Game{
		Slotx: slot.Slotx{
			Sel: sel,
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

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
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
	if count := g.SymNum(scat); count >= 4 {
		var pay = ScatPay[min(count-1, 4)]
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * float64(g.Sel) * pay,
			MP:  1,
			Sym: scat,
			Num: count,
			XY:  g.SymPos(scat),
			FS:  15,
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	if g.FSR == 0 {
		var reels, _ = ReelsMap.FindClosest(mrtp)
		g.ReelSpin(reels)
	} else {
		g.ReelSpin(ReelsBon)
	}
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
