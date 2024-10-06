package sizzlinghot

import (
	slot "github.com/slotopol/server/game/slot"
)

// Lined payment.
var LinePay = [8][5]float64{
	{0, 0, 100, 1000, 5000}, // seven
	{0, 0, 50, 200, 500},    // melon
	{0, 0, 50, 200, 500},    // grapes
	{0, 0, 20, 50, 200},     // plum
	{0, 0, 20, 50, 200},     // orange
	{0, 0, 20, 50, 200},     // lemon
	{0, 5, 20, 50, 200},     // cherry
	{0, 0, 0, 0, 0},         // star
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 2, 10, 50} // star

const (
	jid = 1 // jackpot ID
)

// Jackpot win combinations.
var Jackpot = [8][5]int{
	{0, 0, 0, 0, 0}, // seven
	{0, 0, 0, 0, 0}, // melon
	{0, 0, 0, 0, 0}, // grapes
	{0, 0, 0, 0, 0}, // plum
	{0, 0, 0, 0, 0}, // orange
	{0, 0, 0, 0, 0}, // lemon
	{0, 0, 0, 0, 0}, // cherry
	{0, 0, 0, 0, 0}, // star
}

// Bet lines
var bl = slot.BetLinesHot5

type Game struct {
	slot.Slot5x3 `yaml:",inline"`
}

func NewGame() *Game {
	return &Game{
		Slot5x3: slot.Slot5x3{
			Sel: slot.MakeBitNum(len(bl), 1),
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
		var line = bl[li-1]

		var numl slot.Pos = 1
		var syml = screen.Pos(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = screen.Pos(x, line)
			if sx != syml {
				break
			}
			numl++
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
	var _, reels = slot.FindReels(ReelsMap, mrtp)
	screen.Spin(reels)
}

func (g *Game) SetSel(sel slot.Bitset) error {
	return slot.ErrNoFeature
}
