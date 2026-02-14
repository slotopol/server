package oceanlegends

// See: https://www.slotsmate.com/software/ct-interactive/ocean-legends

import (
	"github.com/slotopol/server/game/slot"
)

const (
	sn         = 10   // number of symbols
	wild, scat = 1, 2 // wild & scatter symbol IDs
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [sn][5]float64{
	{},                    //  1 wild (2, 3, 4, 5 reels only)
	{},                    //  2 scatter
	{0, 0, 50, 200, 1000}, //  3 aryan
	{0, 0, 15, 75, 200},   //  4 nymph
	{0, 0, 5, 50, 150},    //  5 pot
	{0, 0, 5, 50, 150},    //  6 vase
	{0, 0, 5, 15, 100},    //  7 ace
	{0, 0, 5, 15, 100},    //  8 king
	{0, 0, 5, 15, 100},    //  9 queen
	{0, 0, 5, 15, 100},    // 10 jack
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 5, 10, 50} // 2 scatter

// Bet lines
var BetLines = slot.BetLinesCT5x3[:]

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
	if count := g.SymNum(scat); count >= 3 {
		var mm float64 = 1 // mult mode
		if g.FSR > 0 {
			mm = 2
		}
		var pay = ScatPay[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * float64(g.Sel) * pay,
			MP:  mm,
			Sym: scat,
			Num: count,
			XY:  g.SymPos(scat),
			FS:  15,
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.SpinReels(reels)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrDisabled
}
