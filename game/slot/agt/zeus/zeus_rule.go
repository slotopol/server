package zeus

// See: https://agtsoftware.org/games/agt/zeus

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsMap slot.ReelsMap[*slot.Reels4x]

// Lined payment.
var LinePay = [10][5]float64{
	{},                 //  1 scatter
	{0, 20, 200, 1000}, //  2 wild
	{0, 10, 100, 500},  //  3 nymph
	{0, 0, 50, 100},    //  4 armor
	{0, 0, 40, 80},     //  5 boots
	{0, 0, 40, 80},     //  6 crown
	{0, 0, 30, 60},     //  7 staff
	{0, 0, 20, 40},     //  8 pantheon
	{0, 0, 20, 40},     //  9 shield
	{0, 0, 10, 20},     // 10 glove
}

// Bet lines
var BetLines = slot.BetLinesAgt4x4[:10]

type Game struct {
	slot.Screen4x4 `yaml:",inline"`
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

const scat, wild = 1, 2

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	for li, line := range BetLines[:g.Sel] {
		var numw, numl slot.Pos = 0, 4
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 4; x++ {
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
		if numw >= 2 {
			payw = LinePay[wild-1][numw-1]
		}
		if numl >= 2 && syml > 0 {
			payl = LinePay[syml-1][numl-1]
		}
		if payl > payw {
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * payl,
				MP:  1,
				Sym: syml,
				Num: numl,
				LI:  li + 1,
				XY:  line.HitxL(numl),
			})
		} else if payw > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * payw,
				MP:  1,
				Sym: wild,
				Num: numw,
				LI:  li + 1,
				XY:  line.HitxL(numw),
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.SymNum(scat); count == 4 {
		const pay, fs = 40, 10
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * float64(g.Sel) * pay,
			MP:  1,
			Sym: scat,
			Num: 1,
			XY:  g.SymPos(scat),
			FS:  fs,
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.ReelSpin(reels)
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
