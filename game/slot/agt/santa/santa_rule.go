package santa

// See: https://agtsoftware.com/games/agt/santa

import (
	"github.com/slotopol/server/game/slot"
)

const (
	sn         = 10   // number of symbols
	wild, scat = 2, 1 // wild & scatter symbol IDs
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [sn][4]float64{
	{},                 //  1 scatter (2 reel only)
	{0, 20, 200, 1000}, //  2 wild
	{0, 10, 100, 500},  //  3 gnomes
	{0, 0, 50, 100},    //  4 snowman
	{0, 0, 40, 80},     //  5 christmas tree
	{0, 0, 30, 60},     //  6 socks
	{0, 0, 30, 60},     //  7 balls
	{0, 0, 20, 40},     //  8 sweets
	{0, 0, 10, 20},     //  9 present
	{0, 0, 10, 20},     // 10 bells
}

// Bet lines
var BetLines = slot.BetLinesAgt4x4[:]

type Game struct {
	slot.Grid4x4 `yaml:",inline"`
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
			var sx = g.LX(x, line)
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
	if count := g.SymNum(scat); count > 0 {
		*wins = append(*wins, slot.WinItem{
			Sym: scat,
			Num: 1,
			XY:  g.SymPos(scat),
			FS:  3,
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.SpinReels(reels)
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
