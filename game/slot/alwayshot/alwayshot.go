package alwayshot

// See: https://freeslotshub.com/novomatic/always-hot/

import (
	slot "github.com/slotopol/server/game/slot"
)

// Lined payment.
var LinePay = [9][3]float64{
	{0, 0, 60}, // 1 seven
	{0, 0, 40}, // 2 star
	{0, 0, 20}, // 3 melon
	{0, 0, 16}, // 4 grapes
	{0, 0, 16}, // 5 bell
	{0, 0, 8},  // 6 orange
	{0, 0, 8},  // 7 plum
	{0, 0, 8},  // 8 lemon
	{0, 0, 8},  // 9 cherry
}

// Bet lines
var bl = slot.BetLinesHot3

type Game struct {
	slot.Slot3x3 `yaml:",inline"`
}

func NewGame() *Game {
	return &Game{
		Slot3x3: slot.Slot3x3{
			Sel: slot.MakeBitNum(len(bl), 1),
			Bet: 1,
		},
	}
}

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	for li := g.Sel.Next(0); li != -1; li = g.Sel.Next(li) {
		var line = bl[li-1]
		var sym1, sym2, sym3 = screen.Pos(1, line), screen.Pos(2, line), screen.Pos(3, line)
		if sym1 == sym2 && sym1 == sym3 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * LinePay[sym1-1][2],
				Mult: 1,
				Sym:  sym1,
				Num:  3,
				Line: li,
				XY:   line.CopyL(3),
			})
		}
	}
}

func (g *Game) Spin(screen slot.Screen, mrtp float64) {
	var reels, _ = slot.FindReels(ReelsMap, mrtp)
	screen.Spin(reels)
}

func (g *Game) SetSel(sel slot.Bitset) error {
	return slot.ErrNoFeature
}
