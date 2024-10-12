package fruitshop

// See: https://slotsspot.com/online-free-slots/fruit-shop-slot/

import (
	slot "github.com/slotopol/server/game/slot"
)

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
var BetLines = slot.BetLinesNetEnt15

type Game struct {
	slot.Slot5x3 `yaml:",inline"`
	// free spin number
	FS int `json:"fs,omitempty" yaml:"fs,omitempty" xml:"fs,omitempty"`
}

func NewGame() *Game {
	return &Game{
		Slot5x3: slot.Slot5x3{
			Sel: slot.MakeBitNum(len(BetLines), 1),
			Bet: 1,
		},
		FS: 0,
	}
}

const wild = 1

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	g.ScanLined(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen slot.Screen, wins *slot.Wins) {
	for li := g.Sel.Next(0); li != -1; li = g.Sel.Next(li) {
		var line = BetLines[li-1]

		var mw float64 = 1 // mult wild
		var numl slot.Pos = 1
		var syml = screen.Pos(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = screen.Pos(x, line)
			if sx == wild {
				mw = 2
			} else if sx != syml {
				break
			}
			numl++
		}

		if pay := LinePay[syml-1][numl-1]; pay > 0 {
			var mm float64 = 1 // mult mode
			var fs int
			if g.FS > 0 {
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
				Line: li,
				XY:   line.CopyL(numl),
				Free: fs,
			})
		}
	}
}

func (g *Game) Spin(screen slot.Screen, mrtp float64) {
	var reels, _ = slot.FindReels(ReelsMap, mrtp)
	screen.Spin(reels)
}

func (g *Game) Apply(screen slot.Screen, wins slot.Wins) {
	if g.FS > 0 {
		g.Gain += wins.Gain()
	} else {
		g.Gain = wins.Gain()
	}

	if g.FS > 0 {
		g.FS--
	}
	for _, wi := range wins {
		if wi.Free > 0 {
			g.FS += wi.Free
		}
	}
}

func (g *Game) FreeSpins() int {
	return g.FS
}

func (g *Game) SetSel(sel slot.Bitset) error {
	return slot.ErrNoFeature
}
