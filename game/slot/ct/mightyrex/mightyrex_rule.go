package mightyrex

// See: https://www.slotsmate.com/software/ct-interactive/mighty-rex

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsBon *slot.Reels5x

var ReelsMap slot.ReelsMap[*slot.Reels5x]

// Lined payment.
var LinePay = [12][5]float64{
	{0, 0, 150, 600, 15000}, //  1 wild
	{},                      //  2 scatter (on 3, 4, 5 reels)
	{0, 0, 25, 350, 1250},   //  3 einiosaurus
	{0, 0, 25, 350, 1250},   //  4 kentrosaurus
	{0, 0, 15, 75, 500},     //  5 troodon
	{0, 0, 15, 75, 500},     //  6 spinosaurus
	{0, 0, 10, 50, 300},     //  7 cretoxyrhina
	{0, 0, 10, 50, 300},     //  8 ammonite
	{0, 0, 5, 30, 200},      //  9 ace
	{0, 0, 5, 30, 200},      // 10 king
	{0, 0, 5, 30, 150},      // 11 queen
	{0, 0, 5, 30, 150},      // 12 jack
}

// Bet lines
var BetLines = slot.BetLinesMgj[:25]

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

const wild, scat = 1, 2

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	for li, line := range BetLines[:g.Sel] {
		var numw, numl slot.Pos = 0, 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = g.LY(x, line)
			if sx == wild {
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

		var payw, payl float64
		if numw >= 3 {
			payw = LinePay[wild-1][numw-1]
		}
		if numl >= 3 && syml > 0 {
			payl = LinePay[syml-1][numl-1]
		}
		if payl > payw {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payl,
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li + 1,
				XY:   line.CopyL(numl),
			})
		} else if payw > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payw,
				Mult: 1,
				Sym:  wild,
				Num:  numw,
				Line: li + 1,
				XY:   line.CopyL(numw),
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.ScatNum(scat); count >= 3 {
		const pay = 5
		var fs int
		if g.FSR == 0 {
			fs = 15
		} else if g.FSR+g.FSN < 100 {
			fs = 100
		}
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   g.ScatPos(scat),
			Free: fs,
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	if g.FSR == 0 {
		var reels, _ = ReelsMap.FindClosest(mrtp)
		g.ReelSpin(reels)
	} else {
		g.ReelSpin(ReelsBon)
	}
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
