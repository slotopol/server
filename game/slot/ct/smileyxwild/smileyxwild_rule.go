package smileyxwild

// See: https://www.slotsmate.com/software/ct-interactive/smiley-x-wild

import (
	_ "embed"
	"math/rand/v2"

	"github.com/slotopol/server/game/slot"
)

//go:embed smileyxwild_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [8][5]float64{
	{},                    // 1 wild
	{},                    // 2 scatter
	{0, 0, 35, 100, 1000}, // 3 heart
	{0, 0, 15, 50, 300},   // 4 sun
	{0, 0, 15, 50, 300},   // 5 beer
	{0, 0, 10, 30, 100},   // 6 pizza
	{0, 0, 10, 30, 100},   // 7 bomb
	{0, 0, 10, 30, 100},   // 8 flower
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 5, 20, 100} // 2 scatter

// Bet lines
var BetLines = slot.BetLinesMgj[:20]

type Game struct {
	slot.Screen5x3 `yaml:",inline"`
	slot.Slotx     `yaml:",inline"`
	W2             float64 `json:"w2" yaml:"w2" xml:"w2"` // wild multiplier on 2 reel
	W4             float64 `json:"w4" yaml:"w4" xml:"w4"` // wild multiplier on 4 reel
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
		var mw float64 = 1 // mult wild
		var numl slot.Pos = 5
		var syml = g.LY(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = g.LY(x, line)
			if sx == wild {
				if x == 2 {
					mw = g.W2
				} else {
					mw *= g.W4
				}
				continue
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		if pay := LinePay[syml-1][numl-1]; pay > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * pay,
				Mult: mw,
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
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.ReelSpin(reels)
}

func (g *Game) Prepare() {
	g.W2 = float64(rand.N(4) + 1)
	g.W4 = float64(rand.N(4) + 1)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
