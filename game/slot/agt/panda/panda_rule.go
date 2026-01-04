package panda

// See: https://agtsoftware.com/games/agt/panda

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [9]float64{
	500, // 1 wild
	160, // 2 bonsai
	80,  // 3 fish
	40,  // 4 fan
	20,  // 5 lamp
	20,  // 6 pot
	20,  // 7 flower
	10,  // 8 button
	0,   // 9 scatter
}

// Bet lines
var BetLines = slot.BetLinesAgt3x3[:27]

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

const wild, scat = 1, 9

func (g *Game) Scanner(wins *slot.Wins) error {
	for li, line := range BetLines[:g.Sel] {
		var numl slot.Pos = 3
		var syml = g.LX(1, line)
		var x slot.Pos
		for x = 2; x <= 3; x++ {
			var sx = g.LX(x, line)
			if sx == wild {
				continue
			} else if syml == wild {
				syml = sx
			} else if sx != syml {
				numl = x - 1
				break
			}
		}
		if numl == 3 {
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * LinePay[syml-1],
				MP:  1,
				Sym: syml,
				Num: 3,
				LI:  li + 1,
				XY:  slot.L2H(line), // whole line is used
			})
		}
	}

	if count := g.SymNum(scat); count > 0 {
		*wins = append(*wins, slot.WinItem{
			Sym: scat,
			Num: count,
			XY:  g.SymPos(scat),
			FS:  int(count),
		})
	}
	return nil
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.SpinReels(reels)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
