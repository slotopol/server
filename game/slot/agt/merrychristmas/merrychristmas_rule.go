package merrychristmas

// See: https://agtsoftware.com/games/agt/christmas

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsMap slot.ReelsMap[*slot.Reels3x]

// Lined payment.
var LinePay = [8]float64{
	500, // 1 snowman
	0,   // 2 scatter
	250, // 3 ice
	100, // 4 sled
	20,  // 5 house
	10,  // 6 bell
	5,   // 7 deer
}

// Bet lines
var BetLines = slot.BetLinesHot3

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

const wild, scat = 1, 2

func (g *Game) Scanner(wins *slot.Wins) error {
	for li, line := range BetLines[:g.Sel] {
		var numw, numl slot.Pos = 0, 3
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 3; x++ {
			var sx = g.LY(x, line)
			if g.FSR > 0 && sx == wild {
				if syml == 0 {
					numw = x
				}
			} else if syml == 0 {
				syml = sx
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		if numw == 3 {
			var pay = LinePay[wild-1]
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * pay,
				MP:  1,
				Sym: wild,
				Num: numw,
				LI:  li + 1,
				XY:  line.HitxL(numw),
			})
		} else if numl == 3 {
			var pay = LinePay[syml-1]
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

	if count := g.ScatNum(scat); count >= 3 {
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * float64(g.Sel) * 10,
			MP:  1,
			Sym: scat,
			Num: count,
			XY:  g.ScatPos(scat),
			FS:  20,
		})
	}
	return nil
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.ReelSpin(reels)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
