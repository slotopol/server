package wildhorses

// See: https://www.slotsmate.com/software/novomatic/wild-horses

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsBon slot.Reelx

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [12][5]float64{
	{0, 0, 50, 200, 1000}, //  1 wild
	{},                    //  2 scatter (2, 3, 4 reels only)
	{0, 0, 40, 150, 750},  //  3 white
	{0, 0, 40, 150, 750},  //  4 black
	{0, 0, 20, 80, 400},   //  5 blue amulet
	{0, 0, 20, 80, 400},   //  6 red amulet
	{0, 0, 15, 60, 300},   //  7 ace
	{0, 0, 15, 60, 300},   //  8 king
	{0, 0, 10, 40, 200},   //  9 queen
	{0, 0, 10, 40, 200},   // 10 jack
	{0, 0, 5, 20, 100},    // 11 ten
	{0, 0, 5, 20, 100},    // 12 nine
}

// Bet lines
var BetLines = slot.BetLinesNvm5x4[:]

type Game struct {
	slot.Screen5x4 `yaml:",inline"`
	slot.Slotx     `yaml:",inline"`

	SS slot.Sym `json:"ss" yaml:"ss" xml:"ss"` // selected symbol
	NW slot.Pos `json:"nw" yaml:"nw" xml:"nw"` // number of white horses
	NB slot.Pos `json:"nb" yaml:"nb" xml:"nb"` // number of black horses
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame(sel int) *Game {
	return &Game{
		Slotx: slot.Slotx{
			Sel: sel,
			Bet: 1,
		},
	}
}

func (g *Game) Clone() slot.SlotGame {
	var clone = *g
	return &clone
}

const (
	wild, scat   = 1, 2
	white, black = 3, 4
)

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
	if g.FSR == 0 {
		if count := g.SymNum(scat); count >= 3 {
			const fs = 10
			*wins = append(*wins, slot.WinItem{
				Sym: scat,
				Num: count,
				XY:  g.SymPos(scat),
				FS:  fs,
			})
		}
	} else if g.FSR == 1 {
		var nw, nb = g.SymNum(white), g.SymNum(black)
		nw += g.NW
		nb += g.NB
		if (g.SS == white && nw >= nb) || (g.SS == black && nb >= nw) {
			const fs = 10
			*wins = append(*wins, slot.WinItem{
				FS: fs,
			})
		}
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

func (g *Game) Apply(wins slot.Wins) {
	g.Slotx.Apply(wins)

	if g.FSN > 0 {
		g.NW += g.SymNum(white)
		g.NB += g.SymNum(black)
	}
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}

func (g *Game) SetMode(ss int) error {
	if ss != white && ss != black {
		return slot.ErrNoFeature
	}
	if g.FSR != 10 { // only on free games start it can be
		return slot.ErrDisabled
	}
	g.SS = slot.Sym(ss)
	return nil
}
