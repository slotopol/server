package extraspin

// See: https://demo.agtsoftware.com/games/agt/extraspin

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed extraspin_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [9][5]float64{
	{0, 10, 50, 250, 1000}, // 1 wild
	{},                     // 2 scatter
	{0, 0, 40, 120, 500},   // 3 strawberry
	{0, 0, 20, 40, 200},    // 4 papaya
	{0, 0, 20, 40, 200},    // 5 grapes
	{0, 0, 10, 20, 100},    // 6 orange
	{0, 0, 10, 20, 100},    // 7 plum
	{0, 0, 5, 10, 50},      // 8 cherry
	{0, 0, 5, 10, 50},      // 9 pear
}

// Bet lines
var BetLines = slot.BetLinesAgt5x3[:10]

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

const wild, scat = 1, 2

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	for li, line := range BetLines[:g.Sel] {
		var numw, numl slot.Pos = 0, 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = g.LY(x, line)
			if sx == wild {
				if syml == 0 {
					numw = x
				}
			} else if syml == 0 {
				syml = sx
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		var payw, payl float64
		if numw >= 2 {
			payw = LinePay[wild-1][numw-1]
		}
		if numl >= 3 && syml > 0 {
			payl = LinePay[syml-1][numl-1]
		}
		if payl > payw {
			var mm float64 = 1 // mult mode
			if g.FSR > 0 {
				mm = 3
			}
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payl,
				Mult: mm,
				Sym:  syml,
				Num:  numl,
				Line: li + 1,
				XY:   line.CopyL(numl),
			})
		} else if payw > 0 {
			var mm float64 = 1 // mult mode
			if g.FSR > 0 {
				mm = 3
			}
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payw,
				Mult: mm,
				Sym:  wild,
				Num:  numw,
				Line: li + 1,
				XY:   line.CopyL(numw),
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.ScatNum(scat); count >= 1 {
		*wins = append(*wins, slot.WinItem{
			Sym:  scat,
			Num:  count,
			XY:   g.ScatPos(scat),
			Free: int(count),
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	g.ReelSpin(reels)
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
