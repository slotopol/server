package colibriwild

// See: https://www.slotsmate.com/software/ct-interactive/colibri-wild

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsMap slot.ReelsMap[*slot.Reels5x]

// Lined payment.
var LinePay = [10][5]float64{
	{},                  //  1 wild
	{},                  //  2 scatter
	{0, 0, 20, 70, 750}, //  3 lychee
	{0, 0, 20, 30, 100}, //  4 carambola
	{0, 0, 20, 30, 100}, //  5 kiwi
	{0, 0, 15, 30, 100}, //  6 figs
	{0, 0, 15, 20, 100}, //  7 orange
	{0, 0, 15, 20, 100}, //  8 tangerine
	{0, 0, 10, 20, 50},  //  9 pineapple
	{0, 0, 10, 20, 50},  // 10 coconut
}

// Bet lines
var BetLines = slot.BetLinesNetEnt5x4[:40]

type Game struct {
	slot.Screen5x4 `yaml:",inline"`
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

const wild, scat = 1, 2

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	var scrnwild = g.Screen5x4
	for x := 0; x < 5; x += 3 { // 1, 3 reels only
		for y := range 4 {
			if g.Scr[x][y] == wild {
				for i := x; i < x+3; i++ { // 1, 2, 3 or 3, 4, 5 reels
					for j := y; j < 4; j++ { // down only
						if scrnwild.Scr[i][j] != scat {
							scrnwild.Scr[i][j] = wild
						}
					}
				}
				break
			}
		}
	}

	for li, line := range BetLines[:g.Sel] {
		var numl slot.Pos = 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = scrnwild.LY(x, line)
			if sx == wild {
				continue
			} else if syml == 0 {
				syml = sx
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
				Line: li + 1,
				XY:   line.CopyL(numl),
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.ScatNum(scat); count >= 3 {
		const pay = 20
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   g.ScatPos(scat),
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.ReelSpin(reels)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
