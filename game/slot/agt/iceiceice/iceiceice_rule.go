package iceiceice

// See: https://demo.agtsoftware.com/games/agt/iceiceice

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed iceiceice_reel.yaml
var reels []byte

var ReelsMap = slot.ReadReelsMap[*slot.Reels3x](reels)

// Lined payment.
var LinePay = [7][3]float64{
	{0, 0, 500}, // 1 seven
	{0, 0, 250}, // 2 strawberry
	{0, 0, 100}, // 3 grapes
	{0, 0, 20},  // 4 plum
	{0, 0, 10},  // 5 pear
	{0, 0, 5},   // 6 cherry
	{},          // 7 star
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
	slot.Slot3x3 `yaml:",inline"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slot3x3: slot.Slot3x3{
			Sel: len(BetLines),
			Bet: 1,
		},
	}
}

const wild, scat = 1, 7

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]

		if g.FSR > 0 {
			var numl slot.Pos = 3
			var syml slot.Sym
			var x slot.Pos
			for x = 1; x <= 3; x++ {
				var sx = screen.Pos(x, line)
				if sx == wild {
					continue
				} else if syml == 0 && sx != scat {
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
					Pay:  g.Bet * LinePay[syml-1][2],
					Mult: 1,
					Sym:  syml,
					Num:  3,
					Line: li,
					XY:   line, // whole line is used
				})
			}
		} else { // g.FSR == 0
			var numl slot.Pos = 3
			var syml = screen.Pos(1, line)
			var x slot.Pos
			for x = 2; x <= 3; x++ {
				var sx = screen.Pos(x, line)
				if sx != syml {
					numl = x - 1
					break
				}
			}
			if numl == 3 && syml != scat {
				*wins = append(*wins, slot.WinItem{
					Pay:  g.Bet * LinePay[syml-1][2],
					Mult: 1,
					Sym:  syml,
					Num:  3,
					Line: li,
					XY:   line, // whole line is used
				})
			}
		}
	}

	if count := screen.ScatNum(scat); count == 3 {
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * 10,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
			Free: 20,
		})
	}
}

func (g *Game) Spin(screen slot.Screen, mrtp float64) {
	var reels, _ = slot.FindReels(ReelsMap, mrtp)
	screen.Spin(reels)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
