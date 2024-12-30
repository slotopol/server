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

const (
	jid = 1 // jackpot ID
)

// Jackpot win combinations.
var Jackpot = [8][5]int{
	{0, 0, 0, 0, 0}, //  1 crown
	{0, 0, 0, 0, 0}, //  2 gold
	{0, 0, 0, 0, 0}, //  3 money
	{0, 0, 0, 0, 0}, //  4 ruby
	{0, 0, 0, 0, 0}, //  5 sapphire
	{0, 0, 0, 0, 0}, //  6 emerald
	{0, 0, 0, 0, 0}, //  7 amethyst
	{0, 0, 0, 0, 0}, //  8 euro
}

// Bet lines
var BetLines = slot.BetLinesNvm10

type Game struct {
	slot.Slot5x3 `yaml:",inline"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slot5x3: slot.Slot5x3{
			Sel: 5,
			Bet: 1,
		},
	}
}

func (g *Game) Clone() slot.SlotGame {
	var clone = *g
	return &clone
}

const wild = 8

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	g.ScanLined(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen slot.Screen, wins *slot.Wins) {
	var scrnwild slot.Screen5x3 = *screen.(*slot.Screen5x3)
	var x, y slot.Pos
	for x = 1; x <= 5; x++ {
		for y = 1; y <= 3; y++ {
			if screen.At(x, y) == wild {
				for i := max(0, x-2); i <= min(4, x); i++ {
					for j := max(0, y-2); j <= min(2, y); j++ {
						scrnwild[i][j] = wild
					}
				}
			}
		}
	}

	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]

		var num slot.Pos = 1
		var sym3 = scrnwild.Pos(3, line)
		var xy slot.Linex
		xy.Set(3, line.At(3))
		if sym2 := scrnwild.Pos(2, line); sym2 == sym3 || sym2 == wild || sym3 == wild {
			if sym3 == wild {
				sym3 = sym2
			}
			xy.Set(2, line.At(2))
			num++
			if sym1 := scrnwild.Pos(1, line); sym1 == sym3 || sym1 == wild || sym3 == wild {
				if sym3 == wild {
					sym3 = sym1
				}
				xy.Set(1, line.At(1))
				num++
			}
		}
		if sym4 := scrnwild.Pos(4, line); sym4 == sym3 || sym4 == wild || sym3 == wild {
			if sym3 == wild {
				sym3 = sym4
			}
			xy.Set(4, line.At(4))
			num++
			if sym5 := scrnwild.Pos(5, line); sym5 == sym3 || sym5 == wild || sym3 == wild {
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
				Line: li,
				XY:   xy,
			})
		}
	}
}

func (g *Game) Spin(mrtp float64) {
	g.Scrn.Spin(Reels)
	var _, wc = slot.FindReels(ChanceMap, mrtp) // wild chance
	if rand.Float64() < wc {
		var x, y = rand.N[slot.Pos](5) + 1, rand.N[slot.Pos](3) + 1
		g.Scrn.Set(x, y, wild)
	}
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
