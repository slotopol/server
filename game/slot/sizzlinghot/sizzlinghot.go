package sizzlinghot

import (
	"math"

	slot "github.com/slotopol/server/game/slot"
	"github.com/slotopol/server/util"
)

// Original reels.
// reels lengths [25, 25, 25, 25, 25], total reshuffles 9765625
// RTP = 89.653(lined) + 6.0024(scatter) = 95.655629%
var Reels96 = slot.Reels5x{
	{1, 4, 4, 4, 3, 1, 3, 3, 8, 6, 6, 6, 7, 7, 7, 6, 6, 2, 2, 5, 2, 5, 5, 5, 4},
	{1, 6, 6, 6, 2, 2, 1, 2, 7, 7, 7, 7, 8, 4, 4, 4, 4, 5, 5, 5, 3, 5, 3, 3, 6},
	{1, 6, 7, 7, 7, 8, 5, 5, 5, 1, 5, 2, 2, 4, 2, 4, 4, 4, 3, 3, 7, 3, 6, 6, 6},
	{1, 5, 5, 5, 5, 1, 5, 4, 4, 4, 8, 3, 3, 6, 6, 6, 7, 6, 7, 7, 7, 4, 4, 2, 2},
	{1, 4, 4, 6, 6, 6, 2, 2, 5, 8, 5, 5, 5, 8, 5, 4, 4, 4, 6, 1, 7, 7, 7, 3, 3},
}

// Map with available reels.
var ReelsMap = map[float64]*slot.Reels5x{
	95.655629: &Reels96,
}

func FindReels(mrtp float64) (rtp float64, reels slot.Reels) {
	for p, r := range ReelsMap {
		if math.Abs(mrtp-p) < math.Abs(mrtp-rtp) {
			rtp, reels = p, r
		}
	}
	return
}

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

type Game struct {
	slot.Slot5x3 `yaml:",inline"`
}

func NewGame() *Game {
	return &Game{
		Slot5x3: slot.Slot5x3{
			Sel: util.MakeBitNum(5, 1),
			Bet: 1,
		},
	}
}

const scat = 8

var bl = slot.BetLinesNvm10

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	g.ScanLined(screen, wins)
	g.ScanScatters(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen slot.Screen, wins *slot.Wins) {
	for li := range g.Sel.Bits() {
		var line = bl.Line(li)

		var syml, numl = screen.Pos(1, line), 1
		for x := 2; x <= 5; x++ {
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
	var _, reels = FindReels(mrtp)
	screen.Spin(reels)
}

func (g *Game) SetSel(sel slot.Bitset) error {
	return slot.ErrNoFeature
}
