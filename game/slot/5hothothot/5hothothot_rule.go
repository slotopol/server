package hothothot

// See: https://demo.agtsoftware.com/games/agt/hothothot5

import (
	"github.com/slotopol/server/game/slot"
)

// Lined payment.
var LinePay = [7][3]float64{
	{0, 0, 100}, // 1 seven
	{0, 0, 50},  // 2 strawberry
	{0, 0, 20},  // 3 grapes
	{0, 0, 4},   // 4 plum
	{0, 0, 2},   // 5 pear
	{0, 0, 1},   // 6 cherry
	{},          // 7 star
}

// Bet lines
var BetLines = slot.BetLinesHot3

type Game struct {
	slot.Slot3x3 `yaml:",inline"`
	// free spin number
	FS int `json:"fs,omitempty" yaml:"fs,omitempty" xml:"fs,omitempty"`
}

func NewGame() *Game {
	return &Game{
		Slot3x3: slot.Slot3x3{
			Sel: len(BetLines),
			Bet: 1,
		},
		FS: 0,
	}
}

const wild, scat = 1, 7

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]

		if g.FS > 0 {
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
					XY:   line.CopyL(3),
				})
			}
		} else { // g.FS == 0
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
					XY:   line.CopyL(3),
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
