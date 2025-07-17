package oliversbar

// See: https://casino.ru/olivers-bar-novomatic/

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed oliversbar_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [13][5]float64{
	{},                     //  1 wild
	{0, 5, 100, 500, 5000}, //  2 Oliver
	{0, 0, 25, 200, 1000},  //  3 friends
	{0, 0, 25, 200, 1000},  //  4 couple
	{0, 0, 15, 100, 500},   //  5 sweet-stuffs
	{0, 0, 15, 100, 500},   //  6 cocktails
	{0, 0, 10, 50, 200},    //  7 flower
	{0, 0, 10, 50, 200},    //  8 lime
	{0, 0, 5, 25, 100},     //  9 olives
	{0, 0, 5, 25, 100},     // 10 strawberries
	{0, 0, 5, 25, 100},     // 11 oranges
	{0, 2, 5, 25, 100},     // 12 cherry
	{},                     // 13 scatter
}

// Scatters payment.
var ScatPay = [5]float64{0, 2, 5, 20, 500} // 13 scatter

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 20, 20, 20} // 13 scatter

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

const wild, scat = 1, 13

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	for li, line := range BetLines[:g.Sel] {
		var numl slot.Pos = 5
		var syml = g.LY(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = g.LY(x, line)
			if sx == wild {
				continue
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		if pay := LinePay[syml-1][numl-1]; pay > 0 {
			var mm float64 = 1 // mult mode
			if g.FSR > 0 {
				mm = 4
			}
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * pay,
				Mult: mm,
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
	if count := g.ScatNum(scat); count >= 2 {
		var mm float64 = 1 // mult mode
		if g.FSR > 0 {
			mm = 4
		}
		var pay, fs = ScatPay[count-1], ScatFreespin[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * pay,
			Mult: mm,
			Sym:  scat,
			Num:  count,
			XY:   g.ScatPos(scat),
			Free: fs,
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
