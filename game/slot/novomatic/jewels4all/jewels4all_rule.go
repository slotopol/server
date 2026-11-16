package jewels4all

import (
	"github.com/slotopol/server/game/slot"
)

const (
	sn   = 8 // number of symbols
	wild = 8 // wild symbol IDs
)

var Reels slot.Reelx

var ChanceMap slot.ReelsMap[float64]

// Lined payment.
var LinePay = [sn][5]float64{
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
var BetLines = slot.BetLinesNvm10[:]

type Game struct {
	slot.Grid5x3 `yaml:",inline"`
	slot.Slotx   `yaml:",inline"`
}

// Declare conformity with SlotGeneric interface.
var _ slot.SlotGeneric = (*Game)(nil)

func NewGame(sel int) *Game {
	var g = &Game{
		Slotx: slot.Slotx{
			Sel: sel,
			Bet: 1,
		},
	}
	g.SpinReels(slot.WySRNG, g.GetReels(slot.InitRTP))
	return g
}

func (g *Game) Clone() slot.SlotGeneric {
	var clone = *g
	return &clone
}

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	var gridwild slot.Grid5x3 = g.Grid5x3
	for x := range 5 {
		for y := range 3 {
			if g.Grid[x][y] == wild {
				for i := max(0, x-1); i <= min(4, x+1); i++ {
					for j := max(0, y-1); j <= min(2, y+1); j++ {
						gridwild.Grid[i][j] = wild
					}
				}
			}
		}
	}

	for li, line := range BetLines[:g.Sel] {
		var num slot.Pos = 1
		var sym3 = gridwild.LX(3, line)
		var xy slot.Linex
		xy.Set(3, line.At(3))
		if sym2 := gridwild.LX(2, line); sym2 == sym3 || sym2 == wild || sym3 == wild {
			if sym3 == wild {
				sym3 = sym2
			}
			xy.Set(2, line.At(2))
			num++
			if sym1 := gridwild.LX(1, line); sym1 == sym3 || sym1 == wild || sym3 == wild {
				if sym3 == wild {
					sym3 = sym1
				}
				xy.Set(1, line.At(1))
				num++
			}
		}
		if sym4 := gridwild.LX(4, line); sym4 == sym3 || sym4 == wild || sym3 == wild {
			if sym3 == wild {
				sym3 = sym4
			}
			xy.Set(4, line.At(4))
			num++
			if sym5 := gridwild.LX(5, line); sym5 == sym3 || sym5 == wild || sym3 == wild {
				if sym3 == wild {
					sym3 = sym5
				}
				xy.Set(5, line.At(5))
				num++
			}
		}

		if num >= 3 {
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * LinePay[sym3-1][num-1],
				MP:  1,
				Sym: sym3,
				Num: num,
				LI:  li + 1,
				XY:  slot.L2H(xy),
			})
		}
	}
}

func (g *Game) GetReels(mrtp float64) slot.Reelx {
	return Reels
}

func (g *Game) Spin(rng *slot.Rand, mrtp float64) {
	g.SpinReels(rng, Reels)
	var _, wc = ChanceMap.FindClosest(mrtp) // wild chance
	if rng.Float64() < wc {
		var x, y = rng.N[slot.Pos](5) + 1, rng.N[slot.Pos](3) + 1
		g.SetSym(x, y, wild)
	}
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
