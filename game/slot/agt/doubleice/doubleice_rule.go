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
	slot.Slot3x3 `yaml:",inline"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slot3x3: slot.Slot3x3{
			Sel: len(BetLines),
			Bet: 1,
		},
	}
}

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	var fm float64 = 1 // fill mult
	if sym := screen.FillSym(); sym >= 6 {
		fm = 2
	}
	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]
		var sym1, sym2, sym3 = screen.Pos(1, line), screen.Pos(2, line), screen.Pos(3, line)
		if sym1 == sym2 && sym1 == sym3 {
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

func (g *Game) Spin(screen slot.Screen, mrtp float64) {
	var reels, _ = slot.FindReels(ReelsMap, mrtp)
	screen.Spin(reels)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
