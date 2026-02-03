package ultrahot

// See: https://www.slotsmate.com/software/novomatic/ultra-hot

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [8]float64{
	750, // 1 seven
	200, // 2 star
	60,  // 3 bar
	40,  // 4 plum
	40,  // 5 orange
	40,  // 6 lemon
	40,  // 7 cherry
	5,   // 8 x
}

// Bet lines
var BetLines = slot.BetLinesHot3x3[:]

type Game struct {
	slot.Grid3x3 `yaml:",inline"`
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

func (g *Game) FillMult() float64 {
	var sym = g.Grid[0][0]
	if sym < 4 || sym > 7 {
		return 1
	}
	if g.Grid[1][0] != sym || g.Grid[2][0] != sym ||
		g.Grid[0][1] != sym || g.Grid[1][1] != sym || g.Grid[2][1] != sym ||
		g.Grid[0][2] != sym || g.Grid[1][2] != sym || g.Grid[2][2] != sym {
		return 1
	}
	return 2
}

func (g *Game) Scanner(wins *slot.Wins) error {
	var fm float64 // fill mult
	for li, line := range BetLines[:g.Sel] {
		var sym1, sym2, sym3 = g.LX(1, line), g.LX(2, line), g.LX(3, line)
		if sym1 == sym2 && sym1 == sym3 {
			if fm == 0 { // lazy calculation
				fm = g.FillMult()
			}
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * LinePay[sym1-1],
				MP:  fm,
				Sym: sym1,
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
	return g.SetSelNum(sel, len(BetLines))
}
