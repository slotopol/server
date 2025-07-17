package alwayshot

// See: https://freeslotshub.com/novomatic/always-hot/

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed alwayshot_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels3x](reels)

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
var BetLines = slot.BetLinesHot3

type Game struct {
	slot.Screen3x3 `yaml:",inline"`
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

func (g *Game) Scanner(wins *slot.Wins) error {
	for li, line := range BetLines[:g.Sel] {
		var sym1, sym2, sym3 = g.LY(1, line), g.LY(2, line), g.LY(3, line)
		if sym1 == sym2 && sym1 == sym3 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * LinePay[sym1-1],
				Mult: 1,
				Sym:  sym1,
				Num:  3,
				Line: li + 1,
				XY:   line, // whole line is used
			})
		}
	}
	return nil
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	g.ReelSpin(reels)
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
