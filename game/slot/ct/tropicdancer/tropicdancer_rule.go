package tropicdancer

// See: https://www.slotsmate.com/software/ct-interactive/tropic-dancer

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed tropicdancer_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [10][5]float64{
	{},                    //  1 wild
	{},                    //  2 scatter
	{0, 0, 20, 200, 1000}, //  3 singer
	{0, 0, 15, 75, 150},   //  4 dancer man
	{0, 0, 5, 50, 150},    //  5 dancer girl 1
	{0, 0, 5, 50, 150},    //  6 dancer girl 2
	{0, 0, 5, 15, 100},    //  7 ace
	{0, 0, 5, 15, 100},    //  8 king
	{0, 0, 5, 15, 100},    //  9 queen
	{0, 0, 5, 15, 100},    // 10 jack
}

// Scatters payment.
var ScatPay = [15]float64{0, 0, 0, 0, 0, 5, 10, 20, 22, 40, 50, 52, 56, 80, 100} // 2 scatter

// Scatter freespins table
var ScatFreespin = [15]int{0, 0, 0, 0, 0, 15, 17, 19, 21, 23, 25, 27, 29, 31, 33} // 2 scatter

// Bet lines
var BetLines = slot.BetLinesMgj[:25]

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
			var mm float64 = 1 // mult mode
			if g.FSR > 0 {
				mm = 2
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
	if count := g.SymNum(scat); count >= 6 {
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
			Free: fs,
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
