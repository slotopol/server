package vikingsfun

// See: https://www.slotsmate.com/software/ct-interactive/vikings-fun

import (
	"github.com/slotopol/server/game/slot"
)

const (
	sn           = 11   // number of symbols
	wild, scat   = 1, 2 // wild & scatter symbol IDs
	wild2, wild3 = 3, 4
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [sn][5]float64{
	{0, 10, 140, 800, 10000}, //  1 wild
	{},                       //  2 scatter
	{0, 5, 15, 80, 625},      //  3 red
	{0, 5, 15, 80, 625},      //  4 blonde
	{0, 0, 15, 100, 625},     //  5 beer
	{0, 0, 15, 100, 625},     //  6 ham
	{0, 0, 10, 60, 200},      //  7 ace
	{0, 0, 10, 60, 200},      //  8 king
	{0, 0, 5, 30, 100},       //  9 queen
	{0, 0, 5, 30, 100},       // 10 jack
	{0, 0, 5, 30, 100},       // 11 ten
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 7, 20, 100} // 2 scatter

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
	if g.FSR == 0 {
		g.ScanLinedReg(wins)
	} else {
		g.ScanLinedBon(wins)
	}
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation for regular games.
func (g *Game) ScanLinedReg(wins *slot.Wins) {
	for li, line := range BetLines[:g.Sel] {
		var numw, numl slot.Pos = 0, 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
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

// Lined symbols calculation for bonus games.
func (g *Game) ScanLinedBon(wins *slot.Wins) {
	for li, line := range BetLines[:g.Sel] {
		var numw, numl slot.Pos = 0, 5
		var symw, syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = g.LX(x, line)
			if sx == wild || sx == wild2 || sx == wild3 {
				if syml == 0 {
					if symw == 0 {
						symw = sx
					} else if sx != symw {
						break // wilds can't substitute each others
					}
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
			payw = LinePay[symw-1][numw-1]
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
				Sym: symw,
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
		var pay = ScatPay[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * float64(g.Sel) * pay,
			MP:  1,
			Sym: scat,
			Num: count,
			XY:  g.SymPos(scat),
			FS:  15,
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.SpinReels(reels)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrDisabled
}
