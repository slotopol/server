package winstorm

// See: https://www.slotsmate.com/software/ct-interactive/win-storm

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed winstorm_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [10][5]float64{
	{},                    //  1 wild
	{},                    //  2 scatter
	{0, 0, 35, 100, 1000}, //  3 seven
	{0, 0, 5, 25, 100},    //  4 coin
	{0, 0, 5, 25, 100},    //  5 horseshoe
	{0, 0, 5, 25, 100},    //  6 bell
	{0, 0, 5, 10, 100},    //  7 ace
	{0, 0, 5, 10, 100},    //  8 king
	{0, 0, 5, 10, 100},    //  9 queen
	{0, 0, 5, 10, 100},    // 10 jack
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 20, 50, 500} // 2 scatter

// Bet lines
var BetLines = slot.BetLinesMgj[:30]

type Game struct {
	slot.Cascade5x3 `yaml:",inline"`
	slot.Slotx      `yaml:",inline"`
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

func (g *Game) Free() bool {
	return g.FSR != 0 || g.Cascade()
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

func (g *Game) Prepare() {
	g.NewFall()
}

func (g *Game) Apply(wins slot.Wins) {
	g.Slotx.Apply(wins)
	g.Strike(wins)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
