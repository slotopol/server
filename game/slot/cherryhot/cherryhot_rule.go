package cherryhot

// See: https://demo.agtsoftware.com/games/agt/cherryhot

import (
	slot "github.com/slotopol/server/game/slot"
)

// Lined payment.
var LinePay = [8][5]float64{
	{0, 0, 20, 200, 1000}, // 1 strawberry
	{0, 0, 8, 80, 200},    // 2 blueberry
	{0, 0, 4.8, 12, 40},   // 3 plum
	{0, 0, 4, 10, 40},     // 4 pear
	{0, 0, 4, 10, 40},     // 5 peach
	{0, 1, 3.2, 8, 32},    // 6 cherry
	{0, 0, 1, 4, 20},      // 7 apple
	{},                    // 8 scatter
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 2, 12, 60} // 8 scatter

// Bet lines
var BetLines = slot.BetLinesAgt5x3[:5]

type Game struct {
	slot.Slot5x3 `yaml:",inline"`
}

func NewGame() *Game {
	return &Game{
		Slot5x3: slot.Slot5x3{
			Sel: len(BetLines),
			Bet: 1,
		},
	}
}

const scat = 8

var Special = map[slot.Sym]bool{
	1: false, // 1 strawberry
	2: false, // 2 blueberry
	3: true,  // 3 plum
	4: true,  // 4 pear
	5: true,  // 5 peach
	6: true,  // 6 cherry
	7: false, // 7 apple
	8: false, // 8 scatter
}

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	g.ScanLined(screen, wins)
	g.ScanScatters(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen slot.Screen, wins *slot.Wins) {
	var ms float64 = 1 // mult screen
	if symm := screen.At(1, 1); Special[symm] {
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			if screen.At(x, 1) != symm || screen.At(x, 2) != symm || screen.At(x, 3) != symm {
				break
			}
		}
		if x > 3 {
			ms = float64(x - 1)
		}
	}
	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]

		var numl slot.Pos = 5
		var syml = screen.Pos(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = screen.Pos(x, line)
			if sx != syml {
				numl = x - 1
				break
			}
		}

		if pay := LinePay[syml-1][numl-1]; pay > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * pay,
				Mult: ms,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyL(numl),
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen slot.Screen, wins *slot.Wins) {
	if count := screen.ScatNum(scat); count >= 3 {
		var pay = ScatPay[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
		})
	}
}

func (g *Game) Spin(screen slot.Screen, mrtp float64) {
	var reels, _ = slot.FindReels(ReelsMap, mrtp)
	screen.Spin(reels)
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
