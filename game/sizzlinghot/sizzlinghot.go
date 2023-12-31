package sizzlinghot

import (
	"github.com/slotopol/server/game"
)

// Original reels.
// reels lengths [25, 25, 25, 25, 25], total reshuffles 9765625
// RTP = 89.653(lined) + 6.0024(scatter) = 95.655629%
var Reels96 = game.Reels5x{
	{1, 4, 4, 4, 3, 1, 3, 3, 8, 6, 6, 6, 7, 7, 7, 6, 6, 2, 2, 5, 2, 5, 5, 5, 4},
	{1, 6, 6, 6, 2, 2, 1, 2, 7, 7, 7, 7, 8, 4, 4, 4, 4, 5, 5, 5, 3, 5, 3, 3, 6},
	{1, 6, 7, 7, 7, 8, 5, 5, 5, 1, 5, 2, 2, 4, 2, 4, 4, 4, 3, 3, 7, 3, 6, 6, 6},
	{1, 5, 5, 5, 5, 1, 5, 4, 4, 4, 8, 3, 3, 6, 6, 6, 7, 6, 7, 7, 7, 4, 4, 2, 2},
	{1, 4, 4, 6, 6, 6, 2, 2, 5, 8, 5, 5, 5, 8, 5, 4, 4, 4, 6, 1, 7, 7, 7, 3, 3},
}

// Map with available reels.
var ReelsMap = map[string]*game.Reels5x{
	"95.66": &Reels96,
}

// Lined payment.
var LinePay = [8][5]int{
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
var ScatPay = [5]int{0, 0, 2, 10, 50} // star

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

type Game struct {
	game.Slot5x3
	LinePay *[8][5]int
	ScatPay *[5]int
}

func NewGame(reels *game.Reels5x) *Game {
	return &Game{
		Slot5x3: game.Slot5x3{
			SBL:      []int{1, 2, 3, 4, 5},
			Bet:      1,
			FS:       0,
			Reels:    reels,
			BetLines: &game.BetLinesNvm10,
		},
		LinePay: &LinePay,
		ScatPay: &ScatPay,
	}
}

const scat = 8

func (g *Game) Scanner(screen game.Screen, ws *game.WinScan) {
	g.ScanLined(screen, ws)
	g.ScanScatters(screen, ws)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen game.Screen, ws *game.WinScan) {
	for _, li := range g.SBL {
		var line = g.BetLines.Line(li)

		var sc = screen.At(1, line.At(1))
		var xy, count = game.Line5x{line.At(1), 0, 0, 0, 0}, 1
		for x := 2; x <= 5; x++ {
			var sx = screen.At(x, line.At(x))
			if sx == sc {
				count++
			} else {
				break
			}
			xy.Set(x, line.At(x))
		}

		if pay := LinePay[sc-1][count-1]; pay > 0 {
			ws.Wins = append(ws.Wins, game.WinItem{
				Pay:  g.Bet * pay,
				Mult: 1,
				Sym:  sc,
				Num:  count,
				Line: li,
				XY:   &xy,
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen game.Screen, ws *game.WinScan) {
	var xy game.Line5x
	var count = 0
	for x := 1; x <= 5; x++ {
		for y := 1; y <= 3; y++ {
			if screen.At(x, y) == scat {
				xy.Set(x, y)
				count++
				break
			}
		}
	}

	if count > 0 {
		if pay := g.ScatPay[count-1]; pay > 0 {
			ws.Wins = append(ws.Wins, game.WinItem{
				Pay:  g.Bet * pay, // independent from selected lines
				Mult: 1,
				Sym:  scat,
				Num:  count,
				XY:   &xy,
			})
		}
	}
}
