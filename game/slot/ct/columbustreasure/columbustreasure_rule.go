package columbustreasure

// See: https://www.livebet.com/casino/slots/ct-interactive/columbus-treasure

import (
	"github.com/slotopol/server/game/slot"
)

const (
	sn                 = 12      // number of symbols
	wild1, wild2, scat = 1, 2, 3 // wild & scatter symbol IDs
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [sn][5]float64{
	{0, 0, 100, 500, 10000}, //  1 wild1
	{},                      //  2 wild2 (2, 4 reels only)
	{},                      //  3 scatter
	{0, 0, 25, 250, 1000},   //  4 cardinal
	{0, 0, 10, 75, 250},     //  5 wizard
	{0, 0, 10, 75, 250},     //  6 sailor
	{0, 0, 7, 50, 100},      //  7 lady
	{0, 0, 7, 50, 100},      //  8 knight
	{0, 0, 5, 15, 25},       //  9 ace
	{0, 0, 5, 15, 25},       // 10 king
	{0, 0, 5, 15, 25},       // 11 queen
	{0, 0, 5, 15, 25},       // 12 jack
}

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 15, 20, 25} // 2 scatter

// Bet lines
var BetLines = slot.BetLinesAgt5x3[:]

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
	for li, line := range BetLines[:g.Sel] {
		var mw float64 = 1 // mult wild
		var numw, numl slot.Pos = 0, 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = g.LX(x, line)
			if sx == wild1 {
				if syml == 0 {
					numw = x
				}
			} else if sx == wild2 {
				if syml == 0 {
					numw = x
				}
				mw = 5
			} else if syml == 0 {
				syml = sx
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		var payw, payl float64
		if numw >= 3 {
			payw = LinePay[wild1-1][numw-1]
		}
		if numl >= 3 && syml > 0 {
			payl = LinePay[syml-1][numl-1]
		}
		if payl*mw > payw {
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * payl,
				MP:  mw,
				Sym: syml,
				Num: numl,
				LI:  li + 1,
				XY:  line.HitxL(numl),
			})
		} else if payw > 0 {
			if numw == 5 {
				mw = 1
			}
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * payw,
				MP:  mw,
				Sym: wild1,
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
		var fs = ScatFreespin[count-1]
		*wins = append(*wins, slot.WinItem{
			Sym: scat,
			Num: count,
			XY:  g.SymPos(scat),
			FS:  fs,
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
