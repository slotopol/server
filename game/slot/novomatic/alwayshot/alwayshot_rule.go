package alwayshot

// See: https://freeslotshub.com/novomatic/always-hot/

import (
	"github.com/slotopol/server/game/slot"
)

const (
	sn = 9 // number of symbols
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [sn]float64{
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
	slot.Grid3x3 `yaml:",inline"`
	slot.Slotx   `yaml:",inline"`
}

// Declare conformity with SlotGeneric interface.
var _ slot.SlotGeneric = (*Game)(nil)

func NewGame(sel int) *Game {
	var g = &Game{
		Slotx: slot.Slotx{
			Sel: sel,
			Bet: 1,
		},
	}
	g.SpinReels(slot.WySRNG, g.GetReels(slot.InitRTP))
	return g
}

func (g *Game) Clone() slot.SlotGeneric {
	var clone = *g
	return &clone
}

func (g *Game) Scanner(wins *slot.Wins) error {
	for li, line := range BetLines[:g.Sel] {
		var sym1, sym2, sym3 = g.LX(1, line), g.LX(2, line), g.LX(3, line)
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

func (g *Game) GetReels(mrtp float64) slot.Reelx {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	return reels
}

func (g *Game) Spin(rng *slot.Rand, mrtp float64) {
	g.SpinReels(rng, g.GetReels(mrtp))
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
