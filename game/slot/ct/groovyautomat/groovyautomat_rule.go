package groovyautomat

// See: https://www.slotsmate.com/software/ct-interactive/groovy-automat

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed groovyautomat_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [8][5]float64{
	{},                      // 1 scatter
	{0, 0, 100, 1000, 5000}, // 2 dollar
	{0, 0, 50, 200, 1000},   // 3 melon
	{0, 0, 50, 200, 1000},   // 4 apple
	{0, 0, 20, 50, 200},     // 5 orange
	{0, 0, 20, 50, 200},     // 6 lemon
	{0, 0, 20, 40, 200},     // 7 plum
	{0, 5, 20, 40, 200},     // 8 cherry
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 10, 50, 250} // 1 scatter

// Bet lines
var BetLines = slot.BetLinesMgj[:5]

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

const scat = 1

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
			if sx != syml {
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
