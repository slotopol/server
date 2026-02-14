package iceiceice

// See: https://agtsoftware.com/games/agt/iceiceice

import (
	"github.com/slotopol/server/game/slot"
)

const (
	sn         = 7    // number of symbols
	wild, scat = 1, 7 // wild & scatter symbol IDs
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [sn]float64{
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
	slot.Grid3x3 `yaml:",inline"`
	slot.Slotx   `yaml:",inline"`
}

// Declare conformity with SlotGeneric interface.
var _ slot.SlotGeneric = (*Game)(nil)

func NewGame(sel int) *Game {
	return &Game{
		Slotx: slot.Slotx{
			Sel: sel,
			Bet: 1,
		},
	}
}

func (g *Game) Clone() slot.SlotGeneric {
	var clone = *g
	return &clone
}

func (g *Game) Scanner(wins *slot.Wins) error {
	for li, line := range BetLines[:g.Sel] {
		if g.FSR > 0 {
			var numl slot.Pos = 3
			var syml = g.LX(1, line)
			var x slot.Pos
			for x = 2; x <= 3; x++ {
				var sx = g.LX(x, line)
				if sx == wild {
					continue
				} else if syml == wild {
					syml = sx
				} else if sx != syml {
					numl = x - 1
					break
				}
			}
			if numl == 3 && syml != scat {
				*wins = append(*wins, slot.WinItem{
					Pay: g.Bet * LinePay[syml-1],
					MP:  1,
					Sym: syml,
					Num: 3,
					LI:  li + 1,
					XY:  slot.L2H(line), // whole line is used
				})
			}
		} else { // g.FSR == 0
			var numl slot.Pos = 3
			var syml = g.LX(1, line)
			var x slot.Pos
			for x = 2; x <= 3; x++ {
				var sx = g.LX(x, line)
				if sx != syml {
					numl = x - 1
					break
				}
			}
			if numl == 3 && syml != scat {
				*wins = append(*wins, slot.WinItem{
					Pay: g.Bet * LinePay[syml-1],
					MP:  1,
					Sym: syml,
					Num: 3,
					LI:  li + 1,
					XY:  slot.L2H(line), // whole line is used
				})
			}
		}
	}

	if count := g.SymNum(scat); count == 3 {
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * 10,
			MP:  1,
			Sym: scat,
			Num: count,
			XY:  g.SymPos(scat),
			FS:  20,
		})
	}
	return nil
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.SpinReels(reels)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrDisabled
}
