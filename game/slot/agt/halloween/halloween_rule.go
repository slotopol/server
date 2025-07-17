package halloween

// See: https://demo.agtsoftware.com/games/agt/halloween

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed halloween_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels3x](reels)

// Lined payment.
var LinePay = [8]float64{
	1000, // 1 pumpkin
	500,  // 2 witch
	200,  // 3 castle
	100,  // 4 scarecrow
	30,   // 5 ghost
	20,   // 6 spider
	10,   // 7 skeleton
	5,    // 8 candles
}

// Bet lines
var BetLines = []slot.Linex{
	{2, 2, 2}, // 1
	{1, 1, 1}, // 2
	{3, 3, 3}, // 3
	{1, 2, 3}, // 4
	{3, 2, 1}, // 5
	{1, 2, 1}, // 6
	{2, 3, 2}, // 7
	{2, 1, 2}, // 8
	{3, 2, 3}, // 9
	{2, 2, 1}, // 10
}

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
	return slot.ErrNoFeature
}
