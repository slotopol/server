package oliversbar

// See: https://casino.ru/olivers-bar-novomatic/

import (
	"github.com/slotopol/server/game/slot"
)

const (
	sn         = 13    // number of symbols
	wild, scat = 1, 13 // wild & scatter symbol IDs
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [sn][5]float64{
	{},                     //  1 wild (2, 3, 4 reels only)
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
var BetLines = slot.BetLinesNvm10[:]

type Game struct {
	slot.Grid5x3 `yaml:",inline"`
	slot.Slotx   `yaml:",inline"`
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
				mm = 4
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
	if count := g.SymNum(scat); count >= 2 {
		var mm float64 = 1 // mult mode
		if g.FSR > 0 {
			mm = 4
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
	return g.SetSelNum(sel, len(BetLines))
}
