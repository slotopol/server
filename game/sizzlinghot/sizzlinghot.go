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
	"96": &Reels96,
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
	game.Slot5x3 `yaml:",inline"`
}

func NewGame(ri string) *Game {
	return &Game{
		Slot5x3: game.Slot5x3{
			RI:  ri,
			BLI: "nvm10",
			SBL: game.MakeSBL(1, 2, 3, 4, 5),
			Bet: 1,
		},
	}
}

const scat = 8

func (g *Game) Scanner(screen game.Screen, ws *game.WinScan) {
	g.ScanLined(screen, ws)
	g.ScanScatters(screen, ws)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen game.Screen, ws *game.WinScan) {
	var bl = game.BetLines5x[g.BLI]
	for li := g.SBL.Next(0); li != 0; li = g.SBL.Next(li) {
		var line = bl.Line(li)

		var syml, numl = screen.At(1, line.At(1)), 1
		for x := 2; x <= 5; x++ {
			var sx = screen.At(x, line.At(x))
			if sx == syml {
				numl++
			} else {
				break
			}
		}

		if pay := LinePay[syml-1][numl-1]; pay > 0 {
			var xy = game.NewLine5x()
			for x := 1; x <= numl; x++ {
				xy.Set(x, line.At(x))
			}
			ws.Wins = append(ws.Wins, game.WinItem{
				Pay:  g.Bet * pay,
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   xy,
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen game.Screen, ws *game.WinScan) {
	var count = 0
	for x := 1; x <= 5; x++ {
		if screen.At(x, 1) == scat || screen.At(x, 2) == scat || screen.At(x, 3) == scat {
			count++
		}
	}

	if count >= 3 {
		var pay = ScatPay[count-1]
		var xy = game.NewLine5x()
		for x := 1; x <= 5; x++ {
			for y := 1; y <= 3; y++ {
				if screen.At(x, y) == scat {
					xy.Set(x, y)
					break
				}
			}
		}
		ws.Wins = append(ws.Wins, game.WinItem{
			Pay:  g.Bet * pay, // independent from selected lines
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   xy,
		})
	}
}

func (g *Game) Spin(screen game.Screen) {
	screen.Spin(ReelsMap[g.RI])
}

func (g *Game) SetLines(sbl game.SBL) error {
	return game.ErrNoFeature
}
