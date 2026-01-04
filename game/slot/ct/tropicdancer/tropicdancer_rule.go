package tropicdancer

// See: https://www.slotsmate.com/software/ct-interactive/tropic-dancer

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [10][5]float64{
	{},                    //  1 wild (2, 3, 4, 5 reels only)
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
var BetLines = slot.BetLinesCT5x3[:]

type Game struct {
	slot.Screen5x3 `yaml:",inline"`
	slot.Slotx     `yaml:",inline"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame(sel int) *Game {
	return &Game{
		Slotx: slot.Slotx{
			Sel: sel,
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
		var syml = g.LX(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = g.LX(x, line)
			if sx != syml && sx != wild {
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
				Pay: g.Bet * pay,
				MP:  mm,
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
	if count := g.SymNum(scat); count >= 6 {
		var mm float64 = 1 // mult mode
		if g.FSR > 0 {
			mm = 2
		}
		var pay, fs = ScatPay[count-1], ScatFreespin[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * float64(g.Sel) * pay,
			MP:  mm,
			Sym: scat,
			Num: count,
			XY:  g.SymPos(scat),
			FS:  fs,
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.SpinReels(reels)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
