package nordicsong

// See: https://www.slotsmate.com/software/ct-interactive/nordic-song

import (
	"github.com/slotopol/server/game/slot"
)

// Remark: bonus reels are not specified in the game rules,
// but are presented as is.
var ReelsBon slot.Reelx

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [11][5]float64{
	{},                     //  1 wild    (2, 3, 4, 5 reels only)
	{},                     //  2 scatter (1, 3, 5 reels only)
	{0, 10, 50, 200, 1000}, //  3 man
	{0, 0, 50, 150, 500},   //  4 woman
	{0, 0, 20, 100, 400},   //  5 owl
	{0, 0, 20, 100, 400},   //  6 dog
	{0, 0, 10, 50, 200},    //  7 ace
	{0, 0, 10, 50, 200},    //  8 king
	{0, 0, 5, 20, 100},     //  9 queen
	{0, 0, 5, 20, 100},     // 10 jack
	{0, 0, 5, 20, 100},     // 11 ten
}

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
		const pay, fs = 5, 12
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * float64(g.Sel) * pay,
			MP:  1,
			Sym: scat,
			Num: count,
			XY:  g.SymPos(scat),
			FS:  fs,
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	if g.FSR == 0 {
		var reels, _ = ReelsMap.FindClosest(mrtp)
		g.SpinReels(reels)
	} else {
		g.SpinReels(ReelsBon)
	}
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrDisabled
}
