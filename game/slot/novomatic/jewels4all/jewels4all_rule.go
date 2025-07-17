package jewels4all

import (
	_ "embed"
	"math/rand/v2"

	"github.com/slotopol/server/game/slot"
)

//go:embed jewels4all_reel.yaml
var reels []byte

var Reels = slot.ReadObj[*slot.Reels5x](reels)

//go:embed jewels4all_chance.yaml
var chance []byte

var ChanceMap = slot.ReadMap[float64](chance)

// Lined payment.
var LinePay = [8][5]float64{
	{0, 0, 20, 100, 1000}, // 1 crown
	{0, 0, 10, 60, 500},   // 2 gold
	{0, 0, 10, 60, 500},   // 3 money
	{0, 0, 5, 40, 200},    // 4 ruby
	{0, 0, 5, 40, 200},    // 5 sapphire
	{0, 0, 5, 20, 100},    // 6 emerald
	{0, 0, 5, 20, 100},    // 7 amethyst
	{},                    // 8 euro
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

const wild = 8

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	var scrnwild slot.Screen5x3 = g.Screen5x3
	var x, y slot.Pos
	for x = 1; x <= 5; x++ {
		for y = 1; y <= 3; y++ {
			if g.At(x, y) == wild {
				for i := max(0, x-2); i <= min(4, x); i++ {
					for j := max(0, y-2); j <= min(2, y); j++ {
						scrnwild.Scr[i][j] = wild
					}
				}
			}
		}
	}

	for li, line := range BetLines[:g.Sel] {
		var num slot.Pos = 1
		var sym3 = scrnwild.LY(3, line)
		var xy slot.Linex
		xy.Set(3, line.At(3))
		if sym2 := scrnwild.LY(2, line); sym2 == sym3 || sym2 == wild || sym3 == wild {
			if sym3 == wild {
				sym3 = sym2
			}
			xy.Set(2, line.At(2))
			num++
			if sym1 := scrnwild.LY(1, line); sym1 == sym3 || sym1 == wild || sym3 == wild {
				if sym3 == wild {
					sym3 = sym1
				}
				xy.Set(1, line.At(1))
				num++
			}
		}
		if sym4 := scrnwild.LY(4, line); sym4 == sym3 || sym4 == wild || sym3 == wild {
			if sym3 == wild {
				sym3 = sym4
			}
			xy.Set(4, line.At(4))
			num++
			if sym5 := scrnwild.LY(5, line); sym5 == sym3 || sym5 == wild || sym3 == wild {
				if sym3 == wild {
					sym3 = sym5
				}
				xy.Set(5, line.At(5))
				num++
			}
		}

		if num >= 3 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * LinePay[sym3-1][num-1],
				Mult: 1,
				Sym:  sym3,
				Num:  num,
				Line: li + 1,
				XY:   xy,
			})
		}
	}
}

func (g *Game) Spin(mrtp float64) {
	g.ReelSpin(Reels)
	var _, wc = slot.FindClosest(ChanceMap, mrtp) // wild chance
	if rand.Float64() < wc {
		var x, y = rand.N[slot.Pos](5) + 1, rand.N[slot.Pos](3) + 1
		g.SetSym(x, y, wild)
	}
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
