package merrychristmas

// See: https://demo.agtsoftware.com/games/agt/christmas

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed merrychristmas_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels3x](reels)

// Lined payment.
var LinePay = [8][3]float64{
	{0, 0, 500}, // 1 wild
	{},          // 2 scatter
	{0, 0, 500}, // 3 snowman
	{0, 0, 250}, // 4 ice
	{0, 0, 100}, // 5 sled
	{0, 0, 20},  // 6 house
	{0, 0, 10},  // 7 bell
	{0, 0, 5},   // 8 deer
}

// Bet lines
var BetLines = slot.BetLinesHot3

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

const wild, scat, wbon = 1, 2, 3

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]

		var numw, numl slot.Pos = 0, 3
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 3; x++ {
			var sx = screen.Pos(x, line)
			if sx == wild || (g.FSR > 0 && sx == wbon) {
				if syml == 0 {
					numw = x
				}
			} else if syml == 0 && sx != scat {
				syml = sx
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		if numw == 3 {
			var pay = LinePay[wild-1][numw-1]
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * pay,
				Mult: 1,
				Sym:  wild,
				Num:  numw,
				Line: li,
				XY:   line.CopyL(numw),
			})
		} else if numl == 3 {
			var pay = LinePay[syml-1][numl-1]
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

	if count := screen.ScatNum(scat); count >= 3 {
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * 10,
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
