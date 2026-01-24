package merrychristmas

// See: https://agtsoftware.com/games/agt/christmas

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [7]float64{
	500, // 1 snowman
	0,   // 2 scatter
	250, // 3 ice
	100, // 4 sled
	20,  // 5 house
	10,  // 6 bell
	5,   // 7 deer
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

const wild, scat = 1, 2

func (g *Game) Scanner(wins *slot.Wins) error {
	for li, line := range BetLines[:g.Sel] {
		var numl slot.Pos = 3
		var syml = g.LX(1, line)
		var x slot.Pos
		if g.FSR > 0 {
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
		} else {
			for x = 2; x <= 3; x++ {
				var sx = g.LX(x, line)
				if sx != syml {
					numl = x - 1
					break
				}
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

	if count := g.SymNum(scat); count >= 3 {
		const pay, fs = 10, 20
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * float64(g.Sel) * pay,
			MP:  1,
			Sym: scat,
			Num: count,
			XY:  g.SymPos(scat),
			FS:  fs,
		})
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
