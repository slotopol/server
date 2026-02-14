package lucky3penguins

// See: https://www.slotsmate.com/software/ct-interactive/lucky-3-penguins
// similar: ct/jollybelugawhales, ct/penguinparty

import (
	"github.com/slotopol/server/game/slot"
	"github.com/slotopol/server/game/slot/ct/penguinparty"
)

const (
	sn   = 10 // number of symbols
	wild = 1  // wild & scatter symbol IDs
)

// Copy data from ct/penguinparty.
var (
	ReelsMap = &penguinparty.ReelsMap
	LinePay  = penguinparty.LinePay
)

// Bet lines
var BetLines = slot.BetLinesCT5x3[:]

type Game struct {
	slot.Grid5x3 `yaml:",inline"`
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
	// Lined symbols calculation.

	var reelwild [5]bool
	for x := 1; x < 4; x++ { // 2, 3, 4 reels only
		for _, sy := range g.Grid[x] {
			if sy == wild {
				reelwild[x] = true
				break
			}
		}
	}

	for li, line := range BetLines[:g.Sel] {
		var numl slot.Pos = 5
		var syml = g.LX(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = g.LX(x, line)
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
	g.SpinReels(reels)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrDisabled
}
