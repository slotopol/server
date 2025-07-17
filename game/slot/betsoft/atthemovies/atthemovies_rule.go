package atthemovies

// See: https://www.slotsmate.com/software/betsoft/at-the-movies

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed atthemovies_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [10][5]float64{
	{0, 20, 200, 500, 1000}, //  1 oscar
	{0, 10, 100, 250, 500},  //  2 popcorn
	{0, 5, 50, 100, 200},    //  3 poster
	{0, 2, 25, 50, 100},     //  4 a
	{0, 0, 20, 40, 80},      //  5 dummy
	{0, 0, 15, 30, 60},      //  6 maw
	{0, 0, 10, 20, 40},      //  7 starship
	{0, 0, 5, 10, 20},       //  8 heart
	{},                      //  9 masks
	{},                      // 10 projector
}

// Scatters payment.
var ScatPay = [5]float64{0, 2, 0, 0, 0} // 10 projector

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 8, 12, 20} // 10 projector

// Bet lines
var BetLines = slot.BetLinesNetEnt5x3[:25]

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

const wild, scat = 9, 10

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
				mw = 2
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		if pay := LinePay[syml-1][numl-1]; pay > 0 {
			var mm float64 = 1 // mult mode
			if g.FSR > 0 {
				mm = 2
			}
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * pay,
				Mult: mw * mm,
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
			mm = 2
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
