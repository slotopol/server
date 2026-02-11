package groovypowers

// See: https://www.slotsmate.com/software/ct-interactive/groovy-powers

import (
	"math/rand/v2"

	"github.com/slotopol/server/game/slot"
)

// Bonus mode probability.
// Determines how often spins are played in the bonus mode.
const Pbm = 1.0 / 15.0

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [10][5]float64{
	{},                    //  1 wild    (2, 3, 4 reels only)
	{},                    //  2 scatter
	{0, 0, 50, 300, 1000}, //  3 glasses
	{0, 0, 30, 50, 200},   //  4 blonde
	{0, 0, 30, 50, 200},   //  5 curly
	{0, 0, 10, 50, 200},   //  6 bald
	{0, 0, 10, 30, 100},   //  7 ace
	{0, 0, 10, 30, 100},   //  8 king
	{0, 0, 10, 30, 100},   //  9 queen
	{0, 0, 10, 30, 100},   // 10 jack
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 20, 50, 500} // 2 scatter

// Bet lines
var BetLines = slot.BetLinesCT5x3[:]

type Game struct {
	slot.Grid5x3 `yaml:",inline"`
	slot.Slotx   `yaml:",inline"`
	BM           bool `json:"bm" yaml:"bm" xml:"bm"` // bonus mode
}

// Declare conformity with SlotGeneric interface.
var _ slot.SlotGeneric = (*Game)(nil)

func NewGame(sel int) *Game {
	return &Game{
		Slotx: slot.Slotx{
			Sel: sel,
			Bet: 1,
		},
	}
}

func (g *Game) Clone() slot.SlotGeneric {
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
	if g.BM {
		for x := 1; x < 4; x++ { // 2, 3, 4 reels only
			for _, sy := range g.Grid[x] {
				if sy == wild {
					reelwild[x] = true
					break
				}
			}
		}
	}

	for li, line := range BetLines[:g.Sel] {
		var numl slot.Pos = 5
		var syml = g.LX(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = g.LX(x, line)
			if reelwild[x-1] || sx == wild {
				continue
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		if pay := LinePay[syml-1][numl-1]; pay > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * pay,
				MP:  1,
				Sym: syml,
				Num: numl,
				LI:  li + 1,
				XY:  line.HitxL(numl),
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.SymNum(scat); count >= 3 {
		var pay = ScatPay[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * float64(g.Sel) * pay,
			MP:  1,
			Sym: scat,
			Num: count,
			XY:  g.SymPos(scat),
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.SpinReels(reels)
}

func (g *Game) Prepare() {
	g.BM = rand.Float64() < Pbm
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrDisabled
}
