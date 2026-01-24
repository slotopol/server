package doublediamond

// See: https://www.slotsmate.com/software/igt/double-diamond

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [6]float64{
	1000, // 1 diamond
	80,   // 2 seven
	40,   // 3 bar3
	25,   // 4 bar2
	10,   // 5 bar1
	10,   // 6 cherry
}

// Bet lines
var BetLines = []slot.Linex{
	{2, 2, 2}, // 1
}

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

const (
	space   = 0
	diamond = 1
	seven   = 2
	bar3    = 3
	bar2    = 4
	bar1    = 5
	cherry  = 6
)

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
					MP:  2,
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
					MP:  4,
					Sym: sym,
					Num: 3,
					LI:  li + 1,
					XY:  slot.L2H(line),
				})
			}
		} else if m[diamond] == 1 && m[cherry] == 1 { // 1 cherry with diamond
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * 2,
				MP:  2,
				Sym: cherry,
				Num: 1,
				LI:  li + 1,
				XY:  slot.L2H(line),
			})
		} else if m[cherry] == 1 { // 1 cherry
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * 2,
				MP:  1,
				Sym: cherry,
				Num: 1,
				LI:  li + 1,
				XY:  slot.L2H(line),
			})
		} else if m[cherry] == 2 { // 2 cherry
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * 5,
				MP:  1,
				Sym: cherry,
				Num: 2,
				LI:  li + 1,
				XY:  slot.L2H(line),
			})
		} else if m[diamond] == 1 && m[space] == 0 && m[seven] == 0 { // any bar with diamond
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * 5,
				MP:  2,
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
		}
	}
	return nil
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.SpinReels(reels)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrDisabled
}
