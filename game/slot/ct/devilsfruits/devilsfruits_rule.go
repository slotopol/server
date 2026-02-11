package devilsfruits

// See: https://www.slotsmate.com/software/ct-interactive/devils-fruits

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [10]float64{
	1500, //  1 wild
	100,  //  2 seven
	35,   //  3 pike
	25,   //  4 bell
	25,   //  5 orange
	25,   //  6 plum
	25,   //  7 bar3
	20,   //  8 bar2
	15,   //  9 bar1
	10,   // 10 cherry
}

// Bet lines
var BetLines = slot.BetLinesHot3x3[:]

type Game struct {
	slot.Grid3x3 `yaml:",inline"`
	slot.Slotx   `yaml:",inline"`
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
	wild  = 1
	bar1  = 9
	bar2  = 8
	bar3  = 7
	space = 0
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
		} else if len(m) == 2 && m[wild] == 1 && m[space] == 0 { // 2 symbols and wild
			for sym := range m {
				if sym == wild {
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
		} else if len(m) == 2 && m[wild] == 2 && m[space] == 0 { // 1 symbol and 2 wilds
			for sym := range m {
				if sym == wild {
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
		} else if m[wild] == 1 && (m[bar1]+m[bar2]+m[bar3] == 2) { // 2 bars and wild
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * 5,
				MP:  2,
				Sym: bar1,
				Num: 3,
				LI:  li + 1,
				XY:  slot.L2H(line),
			})
		} else if m[bar1]+m[bar2]+m[bar3] == 3 {
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * 5,
				MP:  1,
				Sym: bar1,
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
