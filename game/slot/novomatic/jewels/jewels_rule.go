package jewels

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed jewels_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [7][5]float64{
	{0, 0, 20, 200, 2000}, // 1 crown
	{0, 0, 15, 100, 500},  // 2 gold
	{0, 0, 15, 100, 500},  // 3 money
	{0, 0, 10, 50, 200},   // 4 ruby
	{0, 0, 10, 50, 200},   // 5 sapphire
	{0, 0, 5, 25, 100},    // 6 emerald
	{0, 0, 5, 25, 100},    // 7 amethyst
}

// Bet lines
var BetLines = slot.BetLinesNvm10

type Game struct {
	slot.Screen5x3 `yaml:",inline"`
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

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	for li, line := range BetLines[:g.Sel] {
		var numl slot.Pos = 1
		var syml = g.LY(3, line)
		var xy slot.Linex
		xy.Set(3, line.At(3))
		if g.LY(2, line) == syml {
			xy.Set(2, line.At(2))
			numl++
			if g.LY(1, line) == syml {
				xy.Set(1, line.At(1))
				numl++
			}
		}
		if g.LY(4, line) == syml {
			xy.Set(4, line.At(4))
			numl++
			if g.LY(5, line) == syml {
				xy.Set(5, line.At(5))
				numl++
			}
		}

		if numl >= 3 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * LinePay[syml-1][numl-1],
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li + 1,
				XY:   xy,
			})
		}
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	g.ReelSpin(reels)
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
