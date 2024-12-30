package ultrahot

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

func (g *Game) Clone() slot.SlotGame {
	var clone = *g
	return &clone
}

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	var fm float64 = 1 // fill mult
	if sym := screen.FillSym(); sym >= 4 && sym <= 7 {
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

func (g *Game) Spin(mrtp float64) {
	var reels, _ = slot.FindReels(ReelsMap, mrtp)
	g.Scrn.Spin(reels)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
