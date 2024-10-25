package ultrahot

import (
	slot "github.com/slotopol/server/game/slot"
)

// Lined payment.
var LinePay = [8][3]float64{
	{0, 0, 750}, // 1 seven
	{0, 0, 200}, // 2 star
	{0, 0, 60},  // 3 bar
	{0, 0, 40},  // 4 plum
	{0, 0, 40},  // 5 orange
	{0, 0, 40},  // 6 lemon
	{0, 0, 40},  // 7 cherry
	{0, 0, 5},   // 8 x
}

// Bet lines
var BetLines = slot.BetLinesHot3

type Game struct {
	slot.Slot3x3 `yaml:",inline"`
}

func NewGame() *Game {
	return &Game{
		Slot3x3: slot.Slot3x3{
			Sel: len(BetLines),
			Bet: 1,
		},
	}
}

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]
		var fm float64 = 1 // fill mult
		if sym := screen.FillSym(); sym >= 4 && sym <= 7 {
			fm = 2
		}
		var sym1, sym2, sym3 = screen.Pos(1, line), screen.Pos(2, line), screen.Pos(3, line)
		if sym1 == sym2 && sym1 == sym3 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * LinePay[sym1-1][2],
				Mult: fm,
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

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
