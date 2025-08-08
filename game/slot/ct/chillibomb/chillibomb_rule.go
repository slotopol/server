package chillibomb

// See: https://www.slotsmate.com/software/ct-interactive/chilli-bomb

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed chillibomb_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [12][5]float64{
	{0, 0, 100, 200, 1000}, //  1 wild
	{},                     //  2 chilli
	{},                     //  3 scatter
	{0, 0, 40, 50, 500},    //  4 seven
	{0, 0, 25, 50, 200},    //  5 avocado
	{0, 0, 25, 50, 200},    //  6 peach
	{0, 0, 20, 50, 200},    //  7 apple
	{0, 0, 20, 50, 200},    //  8 watermelon
	{0, 0, 20, 50, 100},    //  9 orange
	{0, 0, 20, 50, 100},    // 10 lemon
	{0, 0, 20, 50, 100},    // 11 plum
	{0, 0, 20, 50, 100},    // 12 cherry
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 10, 20, 500} // 2 scatter

// Bet lines
var BetLines = slot.BetLinesAgt5x3[:20]

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

const wild, chilli, scat = 1, 2, 3

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	var scrnwild slot.Screen5x3 = g.Screen5x3
	for y := range 3 {
		if g.Scr[2][y] == chilli {
			for i := max(0, 1); i <= min(4, 3); i++ {
				for j := max(0, y-1); j <= min(2, y+1); j++ {
					if scrnwild.Scr[i][j] != scat {
						scrnwild.Scr[i][j] = wild
					}
				}
			}
		}
	}

	for li, line := range BetLines[:g.Sel] {
		var numw, numl slot.Pos = 0, 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = scrnwild.LY(x, line)
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
		if numw >= 3 {
			payw = LinePay[wild-1][numw-1]
		}
		if numl >= 3 && syml > 0 {
			payl = LinePay[syml-1][numl-1]
		}
		if payl > payw {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payl,
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li + 1,
				XY:   line.CopyL(numl),
			})
		} else if payw > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payw,
				Mult: 1,
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
	if count := g.ScatNum(scat); count >= 3 {
		var pay = ScatPay[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   g.ScatPos(scat),
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	g.ReelSpin(reels)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
