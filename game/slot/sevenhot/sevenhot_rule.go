package sevenhot

// See: https://demo.agtsoftware.com/games/agt/sevenhot20

import (
	slot "github.com/slotopol/server/game/slot"
)

// Lined payment.
var LinePay = [8][5]float64{
	{0, 0, 10, 100, 500}, // 1 seven
	{0, 0, 5, 20, 50},    // 2 blueberry
	{0, 0, 5, 20, 40},    // 3 strawberry
	{0, 0, 2.5, 6, 22},   // 4 plum
	{0, 0, 2, 5, 20},     // 5 pear
	{0, 0, 2, 4, 18},     // 6 peach
	{0, 0.5, 1.5, 4, 18}, // 7 cherry
	{},                   // 8 bell
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 2, 10, 60} // 8 bell

// Bet lines
var BetLines = slot.BetLinesAgt5x3[:20]

type Game struct {
	slot.Slot5x3 `yaml:",inline"`
}

func NewGame() *Game {
	return &Game{
		Slot5x3: slot.Slot5x3{
			Sel: slot.MakeBitNum(len(BetLines), 1),
			Bet: 1,
		},
	}
}

const scat = 8

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	g.ScanLined(screen, wins)
	g.ScanScatters(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen slot.Screen, wins *slot.Wins) {
	for li := g.Sel.Next(0); li != -1; li = g.Sel.Next(li) {
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
				Mult: 1,
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
			Pay:  g.Bet * float64(g.Sel.Num()) * pay,
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

func (g *Game) SetSel(sel slot.Bitset) error {
	return g.SetSelNum(sel, len(BetLines))
}
