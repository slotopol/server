package ultrahot

// See: https://www.slotsmate.com/software/novomatic/ultra-hot

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed ultrahot_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels3x](reels)

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
var BetLines = slot.BetLinesHot3

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

func (g *Game) FillMult() float64 {
	var sym = g.Scr[0][0]
	if sym < 4 || sym > 7 {
		return 1
	}
	if g.Scr[1][0] != sym || g.Scr[2][0] != sym ||
		g.Scr[0][1] != sym || g.Scr[1][1] != sym || g.Scr[2][1] != sym ||
		g.Scr[0][2] != sym || g.Scr[1][2] != sym || g.Scr[2][2] != sym {
		return 1
	}
	return 2
}

func (g *Game) Scanner(wins *slot.Wins) error {
	var fm float64 // fill mult
	for li, line := range BetLines[:g.Sel] {
		var sym1, sym2, sym3 = g.LY(1, line), g.LY(2, line), g.LY(3, line)
		if sym1 == sym2 && sym1 == sym3 {
			if fm == 0 { // lazy calculation
				fm = g.FillMult()
			}
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * LinePay[sym1-1],
				Mult: fm,
				Sym:  sym1,
				Num:  3,
				Line: li + 1,
				XY:   line, // whole line is used
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
	return g.SetSelNum(sel, len(BetLines))
}
