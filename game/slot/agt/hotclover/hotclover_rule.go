package hotclover

// See: https://demo.agtsoftware.com/games/agt/hotclover100

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed hotclover_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [10][5]float64{
	{},                     //  1 wild
	{},                     //  2 scatter
	{0, 10, 50, 200, 3000}, //  3 seven
	{0, 0, 40, 100, 600},   //  4 strawberry
	{0, 0, 40, 100, 400},   //  5 grapes
	{0, 0, 20, 50, 200},    //  6 bar
	{0, 0, 12, 30, 100},    //  7 plum
	{0, 0, 12, 30, 100},    //  8 orange
	{0, 0, 8, 28, 80},      //  9 lemon
	{0, 0, 8, 28, 80},      // 10 cherry
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 3, 20, 100} // 2 scatter

// Bet lines
var BetLines = slot.BetLinesAgt5x4[:100]

type Game struct {
	slot.Screen5x4 `yaml:",inline"`
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
	var reelwild [5]bool
	for x := 1; x < 4; x++ { // 2, 3, 4 reel only
		for y := 0; y < 3; y++ {
			if g.Scr[x][y] == wild {
				reelwild[x] = true
				break
			}
		}
	}

	for li, line := range BetLines[:g.Sel] {
		var numl slot.Pos = 5
		var syml = g.LY(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = g.LY(x, line)
			if reelwild[x-1] {
				continue
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

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
