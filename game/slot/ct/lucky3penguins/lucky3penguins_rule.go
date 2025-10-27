package lucky3penguins

// See: https://www.slotsmate.com/software/ct-interactive/lucky-3-penguins
// similar: ct/jollybelugawhales, ct/penguinparty

import (
	"github.com/slotopol/server/game/slot"
	"github.com/slotopol/server/game/slot/ct/penguinparty"
)

// Copy data from ct/penguinparty.
var (
	ReelsMap = &penguinparty.ReelsMap
	LinePay  = penguinparty.LinePay
)

// Bet lines
var BetLines = slot.BetLinesMgj[:20]

type Game struct {
	slot.Screen5x3 `yaml:",inline"`
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

const wild = 1

func (g *Game) Scanner(wins *slot.Wins) error {
	// Lined symbols calculation.

	var reelwild [5]bool
	for x := 1; x < 4; x++ { // 2, 3, 4 reels only
		for y := range 3 {
			if g.Scr[x][y] == wild {
				reelwild[x] = true
				break
			}
		}
	}

	for li, line := range BetLines[:g.Sel] {
		var numl slot.Pos = 5
		var syml = g.LY(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = g.LY(x, line)
			if reelwild[x-1] {
				continue
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

	// Scatters calculation.

	if reelwild[1] && reelwild[2] {
		*wins = append(*wins, slot.WinItem{
			Sym: wild,
			Num: 2,
			FS:  15,
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
