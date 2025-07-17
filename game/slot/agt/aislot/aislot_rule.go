package aislot

// See: https://demo.agtsoftware.com/games/agt/aislot

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed aislot_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [10][5]float64{
	{},                       //  1 scatter
	{0, 10, 100, 1000, 5000}, //  2 man
	{0, 5, 40, 400, 2000},    //  3 mind
	{0, 5, 30, 100, 750},     //  4 internet
	{0, 5, 30, 100, 750},     //  5 eye
	{0, 0, 5, 40, 150},       //  6 ace
	{0, 0, 5, 40, 150},       //  7 king
	{0, 0, 5, 25, 100},       //  8 queen
	{0, 0, 5, 25, 100},       //  9 jack
	{0, 0, 5, 25, 100},       // 10 ten
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 2, 25, 250} // 1 scatter

// Bet lines
var BetLines = slot.BetLinesAgt5x3[:30]

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

const wild, scat = 1, 1

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	for li, line := range BetLines[:g.Sel] {
		var numl slot.Pos = 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = g.LY(x, line)
			if sx == wild {
				continue
			} else if syml == 0 {
				syml = sx
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		if syml > 0 {
			if pay := LinePay[syml-1][numl-1]; pay > 0 {
				*wins = append(*wins, slot.WinItem{
					Pay:  g.Bet * pay,
					Mult: 1,
					Sym:  syml,
					Num:  numl,
					Line: li + 1,
					XY:   line.CopyL(numl),
				})
			}
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
			Free: 12,
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
