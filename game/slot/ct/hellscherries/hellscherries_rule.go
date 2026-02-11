package hellscherries

// See: https://www.slotsmate.com/software/ct-interactive/hells-cherries

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [10]float64{
	500, //  1 seven
	250, //  2 bar
	200, //  3 melon
	100, //  4 bell
	50,  //  5 apple
	50,  //  6 pear
	50,  //  7 plum
	10,  //  8 lemon
	10,  //  9 orange
	10,  // 10 cherry
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

const wild = 1

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
	return nil
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.SpinReels(reels)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrDisabled
}
