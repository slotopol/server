package doubleice

// See: https://demo.agtsoftware.com/games/agt/doubleice

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed doubleice_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels3x](reels)

// Lined payment.
var LinePay = [9]float64{
	600, // 1 seven
	400, // 2 strawberry
	200, // 3 bell
	160, // 4 star
	160, // 5 lemon
	20,  // 6 blueberry
	20,  // 7 plum
	20,  // 8 orange
	20,  // 9 cherry
}

// Bet lines
var BetLines = slot.BetLinesAgt3x3[:]

type Game struct {
	slot.Slotx[slot.Screen3x3] `yaml:",inline"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slotx: slot.Slotx[slot.Screen3x3]{
			Sel: len(BetLines),
			Bet: 1,
		},
	}
}

func (g *Game) Clone() slot.SlotGame {
	var clone = *g
	return &clone
}

func FillMult(screen *slot.Screen3x3) float64 {
	var sym = screen[0][0]
	if sym < 6 {
		return 1
	}
	if screen[1][0] != sym || screen[2][0] != sym ||
		screen[0][1] != sym || screen[1][1] != sym || screen[2][1] != sym ||
		screen[0][2] != sym || screen[1][2] != sym || screen[2][2] != sym {
		return 1
	}
	return 2
}

func (g *Game) Scanner(wins *slot.Wins) {
	var fm float64 // fill mult
	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]
		var sym1, sym2, sym3 = g.Scr.Pos(1, line), g.Scr.Pos(2, line), g.Scr.Pos(3, line)
		if sym1 == sym2 && sym1 == sym3 {
			if fm == 0 { // lazy calculation
				fm = FillMult(&g.Scr)
			}
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * LinePay[sym1-1],
				Mult: fm,
				Sym:  sym1,
				Num:  3,
				Line: li,
				XY:   line, // whole line is used
			})
		}
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	g.Scr.Spin(reels)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
