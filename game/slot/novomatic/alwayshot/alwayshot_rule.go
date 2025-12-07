package alwayshot

// See: https://freeslotshub.com/novomatic/always-hot/

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsMap slot.ReelsMap[*slot.Reels3x]

// Lined payment.
var LinePay = [9]float64{
	300, // 1 seven
	200, // 2 star
	100, // 3 melon
	80,  // 4 grapes
	80,  // 5 bell
	40,  // 6 orange
	40,  // 7 plum
	40,  // 8 lemon
	40,  // 9 cherry
}

// Bet lines
var BetLines = slot.BetLinesHot3x3[:]

type Game struct {
	slot.Screen3x3 `yaml:",inline"`
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

func (g *Game) Scanner(wins *slot.Wins) error {
	for li, line := range BetLines[:g.Sel] {
		var sym1, sym2, sym3 = g.LY(1, line), g.LY(2, line), g.LY(3, line)
		if sym1 == sym2 && sym1 == sym3 {
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * LinePay[sym1-1],
				MP:  1,
				Sym: sym1,
				Num: 3,
				LI:  li + 1,
				XY:  slot.L2H(line), // whole line is used
			})
		}
	}
	return nil
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.ReelSpin(reels)
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
