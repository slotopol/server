package iceiceice

// See: https://demo.agtsoftware.com/games/agt/iceiceice

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed iceiceice_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels3x](reels)

// Lined payment.
var LinePay = [7]float64{
	500, // 1 seven
	250, // 2 strawberry
	100, // 3 grapes
	20,  // 4 plum
	10,  // 5 pear
	5,   // 6 cherry
	0,   // 7 star
}

// Bet lines
var BetLines = []slot.Linex{
	{2, 2, 2}, // 1
	{1, 1, 1}, // 2
	{3, 3, 3}, // 3
	{1, 2, 3}, // 4
	{3, 2, 1}, // 5
	{2, 3, 2}, // 6
	{2, 1, 2}, // 7
	{1, 2, 1}, // 8
	{3, 2, 3}, // 9
	{1, 3, 1}, // 10
}

type Game struct {
	slot.Screen3x3 `yaml:",inline"`
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

const wild, scat = 1, 7

func (g *Game) Scanner(wins *slot.Wins) error {
	for li, line := range BetLines[:g.Sel] {
		if g.FSR > 0 {
			var numl slot.Pos = 3
			var syml slot.Sym
			var x slot.Pos
			for x = 1; x <= 3; x++ {
				var sx = g.LY(x, line)
				if sx == wild {
					continue
				} else if syml == 0 {
					syml = sx
				} else if sx != syml {
					numl = x - 1
					break
				}
			}
			if numl == 3 {
				if syml == 0 {
					syml = wild
				}
				*wins = append(*wins, slot.WinItem{
					Pay:  g.Bet * LinePay[syml-1],
					Mult: 1,
					Sym:  syml,
					Num:  3,
					Line: li + 1,
					XY:   line, // whole line is used
				})
			}
		} else { // g.FSR == 0
			var numl slot.Pos = 3
			var syml = g.LY(1, line)
			var x slot.Pos
			for x = 2; x <= 3; x++ {
				var sx = g.LY(x, line)
				if sx != syml {
					numl = x - 1
					break
				}
			}
			if numl == 3 && syml != scat {
				*wins = append(*wins, slot.WinItem{
					Pay:  g.Bet * LinePay[syml-1],
					Mult: 1,
					Sym:  syml,
					Num:  3,
					Line: li + 1,
					XY:   line, // whole line is used
				})
			}
		}
	}

	if count := g.ScatNum(scat); count == 3 {
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * 10,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   g.ScatPos(scat),
			Free: 20,
		})
	}
	return nil
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	g.ReelSpin(reels)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
