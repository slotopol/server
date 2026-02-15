package triplediamond

// See: https://www.slotsmate.com/software/igt/triple-diamond

import (
	"github.com/slotopol/server/game/slot"
)

const (
	sn      = 5 // number of symbols
	space   = 0
	diamond = 1
	seven   = 2
	bar3    = 3
	bar2    = 4
	bar1    = 5
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [sn]float64{
	1199, // 1 diamond
	100,  // 2 seven
	40,   // 3 bar3
	20,   // 4 bar2
	10,   // 5 bar1
}

// Bet lines
var BetLines = []slot.Linex{
	{2, 2, 2}, // 1
	{1, 1, 1}, // 2
	{3, 3, 3}, // 3
	{1, 2, 3}, // 4
	{3, 2, 1}, // 5
	{2, 1, 2}, // 6
	{2, 3, 2}, // 7
	{3, 2, 3}, // 8
	{1, 2, 1}, // 9
}

type Game struct {
	slot.Grid3x3 `yaml:",inline"`
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
	for li, line := range BetLines[:g.Sel] {
		var m = map[slot.Sym]int{}
		m[g.LX(1, line)]++
		m[g.LX(2, line)]++
		m[g.LX(3, line)]++
		if len(m) == 1 && m[space] == 0 { // 3 symbols
			for sym := range m {
				*wins = append(*wins, slot.WinItem{
					Pay: g.Bet * LinePay[sym-1],
					MP:  1,
					Sym: sym,
					Num: 3,
					LI:  li + 1,
					XY:  slot.L2H(line),
				})
			}
		} else if len(m) == 2 && m[diamond] == 1 && m[space] == 0 { // 2 symbols and diamond
			for sym := range m {
				if sym == diamond {
					continue
				}
				*wins = append(*wins, slot.WinItem{
					Pay: g.Bet * LinePay[sym-1],
					MP:  3,
					Sym: sym,
					Num: 3,
					LI:  li + 1,
					XY:  slot.L2H(line),
				})
			}
		} else if len(m) == 2 && m[diamond] == 2 && m[space] == 0 { // 1 symbol and 2 diamonds
			for sym := range m {
				if sym == diamond {
					continue
				}
				*wins = append(*wins, slot.WinItem{
					Pay: g.Bet * LinePay[sym-1],
					MP:  9,
					Sym: sym,
					Num: 3,
					LI:  li + 1,
					XY:  slot.L2H(line),
				})
			}
		} else if m[diamond] == 1 && m[space] == 0 && m[seven] == 0 { // any bar with diamond
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * 5,
				MP:  3,
				Sym: 0,
				Num: 3,
				LI:  li + 1,
				XY:  slot.L2H(line),
			})
		} else if m[diamond] == 0 && m[space] == 0 && m[seven] == 0 { // any bar without diamond
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * 5,
				MP:  1,
				Sym: 0,
				Num: 3,
				LI:  li + 1,
				XY:  slot.L2H(line),
			})
		} else if m[diamond] == 1 { // 1 diamond
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * 2,
				MP:  1,
				Sym: diamond,
				Num: 1,
				LI:  li + 1,
				XY:  slot.L2H(line),
			})
		} else if m[diamond] == 2 { // 2 diamonds
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * 10,
				MP:  1,
				Sym: diamond,
				Num: 2,
				LI:  li + 1,
				XY:  slot.L2H(line),
			})
		}
	}
	return nil
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.SpinReels(reels)
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
