package alwayshot

// See: https://freeslotshub.com/novomatic/always-hot/

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed alwayshot_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels3x](reels)

// Lined payment.
var LinePay = [9][3]float64{
	{0, 0, 300}, // 1 seven
	{0, 0, 200}, // 2 star
	{0, 0, 100}, // 3 melon
	{0, 0, 80},  // 4 grapes
	{0, 0, 80},  // 5 bell
	{0, 0, 40},  // 6 orange
	{0, 0, 40},  // 7 plum
	{0, 0, 40},  // 8 lemon
	{0, 0, 40},  // 9 cherry
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

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]
		var sym1, sym2, sym3 = screen.Pos(1, line), screen.Pos(2, line), screen.Pos(3, line)
		if sym1 == sym2 && sym1 == sym3 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * LinePay[sym1-1][2],
				Mult: 1,
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
