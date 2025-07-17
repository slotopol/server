package cherryhot

// See: https://demo.agtsoftware.com/games/agt/cherryhot

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed cherryhot_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [8][5]float64{
	{0, 0, 100, 1000, 5000}, // 1 strawberry
	{0, 0, 40, 400, 1000},   // 2 blueberry
	{0, 0, 24, 60, 200},     // 3 plum
	{0, 0, 20, 50, 200},     // 4 pear
	{0, 0, 20, 50, 200},     // 5 peach
	{0, 5, 16, 40, 160},     // 6 cherry
	{0, 0, 5, 20, 100},      // 7 apple
	{},                      // 8 scatter
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 2, 12, 60} // 8 scatter

// Bet lines
var BetLines = slot.BetLinesAgt5x3[:5]

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

const scat = 8

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

func (g *Game) FillMult() float64 {
	var sym = g.Scr[0][0]
	if sym < 3 || sym > 6 {
		return 1
	}
	var r *[3]slot.Sym
	var i int
	for i = 0; i < 5; i++ {
		if r = &g.Scr[i]; r[0] != sym || r[1] != sym || r[2] != sym {
			break
		}
	}
	if i < 3 {
		return 1
	}
	return float64(i)
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	var fm float64 // fill mult
	for li, line := range BetLines[:g.Sel] {
		var numl slot.Pos = 5
		var syml = g.LY(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = g.LY(x, line)
			if sx != syml {
				numl = x - 1
				break
			}
		}

		if pay := LinePay[syml-1][numl-1]; pay > 0 {
			if fm == 0 { // lazy calculation
				fm = g.FillMult()
			}
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * pay,
				Mult: fm,
				Sym:  syml,
				Num:  numl,
				Line: li + 1,
				XY:   line.CopyL(numl),
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
	return g.SetSelNum(sel, len(BetLines))
}
