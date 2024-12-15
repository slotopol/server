package jokers100

// See: https://demo.agtsoftware.com/games/agt/jokers100

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
	"github.com/slotopol/server/game/slot/agt/jokers"
)

//go:embed jokers100_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = jokers.LinePay

// Scatters payment.
var ScatPay = jokers.ScatPay

// Bet lines
var BetLines = slot.BetLinesAgt5x4[:]

type Game struct {
	slot.Slot5x4 `yaml:",inline"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slot5x4: slot.Slot5x4{
			Sel: len(BetLines),
			Bet: 1,
		},
	}
}

const wild, scat = 1, 2

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	var scrn5x4 = screen.(*slot.Screen5x4)
	g.ScanLined(scrn5x4, wins)
	g.ScanScatters(scrn5x4, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen *slot.Screen5x4, wins *slot.Wins) {
	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]

		var numl slot.Pos = 5
		var syml = screen.Pos(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = screen.Pos(x, line)
			if sx == wild {
				continue
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
				Line: li,
				XY:   line.CopyL(numl),
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen *slot.Screen5x4, wins *slot.Wins) {
	if count := screen.ScatNum(scat); count >= 3 {
		var pay = ScatPay[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
		})
	}
}

func (g *Game) Spin(screen slot.Screen, mrtp float64) {
	var reels, _ = slot.FindReels(ReelsMap, mrtp)
	screen.Spin(reels)
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
