package cherrycrown

// See: https://www.slotsmate.com/software/ct-interactive/cherry-crown

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [10][5]float64{
	{},                  //  1 wild
	{},                  //  2 scatter
	{0, 0, 20, 60, 500}, //  3 seven
	{0, 0, 20, 30, 100}, //  4 melon
	{0, 0, 20, 30, 100}, //  5 apple
	{0, 0, 15, 30, 100}, //  6 pear
	{0, 0, 15, 20, 100}, //  7 orange
	{0, 0, 15, 20, 100}, //  8 lemon
	{0, 0, 10, 20, 50},  //  9 plum
	{0, 0, 10, 20, 50},  // 10 cherry
}

// Bet lines
var BetLines = slot.BetLinesAgt5x3[:]

type Game struct {
	slot.Screen5x3 `yaml:",inline"`
	slot.Slotx     `yaml:",inline"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame(sel int) *Game {
	return &Game{
		Slotx: slot.Slotx{
			Sel: sel,
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
	var scrnwild = g.Screen5x3
	for x := 0; x < 5; x += 3 { // 1, 3 reels only
		for y := range 3 {
			if g.Scr[x][y] == wild {
				for i := x; i < x+3; i++ { // 1, 2, 3 or 3, 4, 5 reels
					for j := y; j < 3; j++ { // down only
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
			var sx = scrnwild.LX(x, line)
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
				Pay: g.Bet * pay,
				MP:  1,
				Sym: syml,
				Num: numl,
				LI:  li + 1,
				XY:  line.HitxL(numl),
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.SymNum(scat); count >= 3 {
		const pay = 20
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * float64(g.Sel) * pay,
			MP:  1,
			Sym: scat,
			Num: count,
			XY:  g.SymPos(scat),
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.SpinReels(reels)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
