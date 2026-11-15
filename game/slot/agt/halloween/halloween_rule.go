package halloween

// See: https://agtsoftware.com/games/agt/halloween

import (
	"github.com/slotopol/server/game/slot"
)

const (
	sn = 8 // number of symbols
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [sn]float64{
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
	return slot.ErrDisabled
}
