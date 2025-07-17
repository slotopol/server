package fruitshop

// See: https://games.netent.com/video-slots/fruit-shop/
// See: https://www.slotsmate.com/software/netent/fruit-shop

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed fruitshop_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [11][5]float64{
	{},                    //  1 wild
	{0, 5, 25, 300, 2000}, //  2 cherry
	{0, 0, 25, 150, 1000}, //  3 plum
	{0, 0, 20, 125, 750},  //  4 lemon
	{0, 0, 20, 100, 500},  //  5 orange
	{0, 0, 15, 75, 200},   //  6 melon
	{0, 0, 15, 50, 150},   //  7 ace
	{0, 0, 10, 25, 100},   //  8 king
	{0, 0, 5, 20, 75},     //  9 queen
	{0, 0, 5, 15, 60},     // 10 jack
	{0, 0, 5, 10, 50},     // 11 ten
}

// Line freespins table on regular games
var LineFreespinReg = [11][5]int{
	{},              //  1 wild
	{0, 1, 1, 2, 5}, //  2 cherry
	{0, 0, 1, 2, 5}, //  3 plum
	{0, 0, 1, 2, 5}, //  4 lemon
	{0, 0, 1, 2, 5}, //  5 orange
	{0, 0, 1, 2, 5}, //  6 melon
	{0, 0, 0, 0, 0}, //  7 ace
	{0, 0, 0, 0, 0}, //  8 king
	{0, 0, 0, 0, 0}, //  9 queen
	{0, 0, 0, 0, 0}, // 10 jack
	{0, 0, 0, 0, 0}, // 11 ten
}

// Line freespins table on bonus games
var LineFreespinBon = [11][5]int{
	{},              //  1 wild
	{0, 1, 1, 2, 5}, //  2 cherry
	{0, 0, 1, 2, 5}, //  3 plum
	{0, 0, 1, 2, 5}, //  4 lemon
	{0, 0, 1, 2, 5}, //  5 orange
	{0, 0, 1, 2, 5}, //  6 melon
	{0, 0, 1, 2, 5}, //  7 ace
	{0, 0, 1, 2, 5}, //  8 king
	{0, 0, 1, 2, 5}, //  9 queen
	{0, 0, 1, 2, 5}, // 10 jack
	{0, 0, 1, 2, 5}, // 11 ten
}

// Bet lines
var BetLines = slot.BetLinesNetEnt5x3[:15]

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

const wild = 1

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
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
			var fs int
			if g.FSR > 0 {
				mm = 2
				fs = LineFreespinBon[syml-1][numl-1]
			} else {
				fs = LineFreespinReg[syml-1][numl-1]
			}
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * pay,
				Mult: mw * mm,
				Sym:  syml,
				Num:  numl,
				Line: li + 1,
				XY:   line.CopyL(numl),
				Free: fs,
			})
		}
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	g.ReelSpin(reels)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
