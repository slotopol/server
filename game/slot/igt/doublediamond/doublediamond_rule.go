package doublediamond

// See: https://www.slotsmate.com/software/igt/double-diamond

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed doublediamond_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels3x](reels)

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
		var m = map[slot.Sym]int{}
		m[g.LY(1, line)]++
		m[g.LY(2, line)]++
		m[g.LY(3, line)]++
		if len(m) == 1 && m[0] == 0 { // 3 symbols
			for sym := range m {
				*wins = append(*wins, slot.WinItem{
					Pay:  g.Bet * LinePay[sym-1],
					Mult: 1,
					Sym:  sym,
					Num:  3,
					Line: li + 1,
					XY:   line,
				})
			}
		} else if len(m) == 2 && m[1] == 1 && m[0] == 0 { // 2 symbols and diamond
			for sym := range m {
				if sym == 1 {
					continue
				}
				*wins = append(*wins, slot.WinItem{
					Pay:  g.Bet * LinePay[sym-1],
					Mult: 2,
					Sym:  sym,
					Num:  3,
					Line: li + 1,
					XY:   line,
				})
			}
		} else if len(m) == 2 && m[1] == 2 && m[0] == 0 { // 1 symbol and 2 diamonds
			for sym := range m {
				if sym == 1 {
					continue
				}
				*wins = append(*wins, slot.WinItem{
					Pay:  g.Bet * LinePay[sym-1],
					Mult: 4,
					Sym:  sym,
					Num:  3,
					Line: li + 1,
					XY:   line,
				})
			}
		} else if m[1] == 1 && m[6] == 1 { // 1 cherry with diamond
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * 2,
				Mult: 2,
				Sym:  6,
				Num:  1,
				Line: li + 1,
				XY:   line,
			})
		} else if m[6] == 1 { // 1 cherry
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * 2,
				Mult: 1,
				Sym:  6,
				Num:  1,
				Line: li + 1,
				XY:   line,
			})
		} else if m[6] == 2 { // 2 cherry
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * 5,
				Mult: 1,
				Sym:  6,
				Num:  2,
				Line: li + 1,
				XY:   line,
			})
		} else if m[1] == 1 && m[0] == 0 && m[2] == 0 { // any bar with diamond
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * 5,
				Mult: 2,
				Sym:  0,
				Num:  3,
				Line: li + 1,
				XY:   line,
			})
		} else if m[1] == 0 && m[0] == 0 && m[2] == 0 { // any bar without diamond
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * 5,
				Mult: 1,
				Sym:  0,
				Num:  3,
				Line: li + 1,
				XY:   line,
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
