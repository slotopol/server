package columbus

// See: https://www.slotsmate.com/software/novomatic/columbus-deluxe

import (
	"github.com/slotopol/server/game/slot"
)

const (
	sn         = 10    // number of symbols
	wild, scat = 1, 10 // wild & scatter symbol IDs
)

var ReelsBon slot.Reelx

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [sn][5]float64{
	{0, 10, 100, 1000, 5000}, //  1 columbus
	{0, 5, 50, 200, 1000},    //  2 spain
	{0, 5, 25, 100, 500},     //  3 necklace
	{0, 5, 15, 75, 250},      //  4 sextant
	{0, 0, 10, 40, 150},      //  5 ace
	{0, 0, 10, 40, 150},      //  6 king
	{0, 0, 10, 40, 150},      //  7 queen
	{0, 0, 5, 20, 100},       //  8 jack
	{0, 0, 5, 20, 100},       //  9 ten
	{},                       // 10 ship
}

// Bet lines
var BetLines = slot.BetLinesNvm10[:]

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
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	var wbon slot.Sym
	if g.FSR > 0 {
		wbon = scat
	}

	for li, line := range BetLines[:g.Sel] {
		var numw, numl slot.Pos = 0, 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = g.LX(x, line)
			if sx == wild || sx == wbon {
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
	if count := g.SymNum(scat); count >= 3 {
		*wins = append(*wins, slot.WinItem{
			MP:  1,
			Sym: scat,
			Num: count,
			XY:  g.SymPos(scat),
			FS:  10,
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	if g.FSR == 0 {
		var reels, _ = ReelsMap.FindClosest(mrtp)
		g.SpinReels(reels)
	} else {
		g.SpinReels(ReelsBon)
	}
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
