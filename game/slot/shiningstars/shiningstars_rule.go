package shiningstars

import (
	slot "github.com/slotopol/server/game/slot"
)

// Lined payment.
var LinePay = [11][5]float64{
	{},                 //  1 wild (on 2, 3, 4 reels)
	{},                 //  2 scatter1 (on all reels)
	{},                 //  3 scatter2 (on 1, 3, 5 reels)
	{0, 1, 5, 25, 500}, //  4 seven
	{0, 0, 4, 12, 70},  //  5 grape
	{0, 0, 4, 12, 70},  //  6 watermelon
	{0, 0, 2, 4, 20},   //  7 avocado
	{0, 0, 1, 3, 15},   //  8 pomegranate
	{0, 0, 1, 3, 15},   //  9 carambola
	{0, 0, 1, 3, 15},   // 10 maracuya
	{0, 0, 1, 3, 15},   // 11 orange
}

// Scatters payment.
var ScatPay1 = [5]float64{0, 0, 5, 20, 100} // 2 scatter1
var ScatPay2 = [5]float64{0, 0, 20}         // 3 scatter2

// Bet lines
var BetLines = slot.BetLinesAgt5x3[:10]

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

const wild, scat1, scat2 = 1, 2, 3

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	g.ScanLined(screen, wins)
	g.ScanScatters(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen slot.Screen, wins *slot.Wins) {
	var reelwild [5]bool
	var x, y slot.Pos
	for x = 2; x <= 4; x++ {
		for y = 1; y <= 3; y++ {
			if screen.At(x, y) == wild {
				reelwild[x-1] = true
				break
			}
		}
	}

	for li := g.Sel.Next(0); li != -1; li = g.Sel.Next(li) {
		var line = BetLines[li-1]

		var numl slot.Pos = 5
		var syml = screen.Pos(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = screen.Pos(x, line)
			if reelwild[x-1] {
				continue
			} else if sx != syml {
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
	if count := screen.ScatNum(scat1); count >= 3 {
		var pay = ScatPay1[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel.Num()) * pay,
			Mult: 1,
			Sym:  scat1,
			Num:  count,
			XY:   screen.ScatPos(scat1),
		})
	} else if count := screen.ScatNumOdd(scat2); count >= 3 {
		var pay = ScatPay2[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel.Num()) * pay,
			Mult: 1,
			Sym:  scat2,
			Num:  count,
			XY:   screen.ScatPosOdd(scat2),
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
