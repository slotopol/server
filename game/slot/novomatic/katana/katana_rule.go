package katana

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsBon slot.Reelx

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [12][5]float64{
	{0, 0, 100, 500, 1000}, // 1  shogun
	{0, 0, 80, 250, 800},   // 2  geisha
	{0, 0, 40, 200, 600},   // 3  general
	{0, 0, 40, 200, 600},   // 4  archers
	{0, 0, 30, 150, 500},   // 5  template
	{0, 0, 30, 150, 500},   // 6  gazebo
	{0, 0, 10, 50, 200},    // 7  ace
	{0, 0, 10, 30, 150},    // 8  king
	{0, 0, 10, 30, 150},    // 9  queen
	{0, 0, 10, 20, 100},    // 10 jack
	{0, 0, 10, 20, 100},    // 11 ten
	{},                     // 12 katana
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 2, 20, 200} // 12 katana

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 10, 10, 10} // 12 katana

// Bet lines
var BetLines = slot.BetLinesNvm5x3v2[:]

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

const wild, scat = 1, 12

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	var reelwild [5]bool
	if g.FSR > 0 {
		for x := range 5 {
			for _, sy := range g.Grid[x] {
				if sy == wild {
					reelwild[x] = true
					break
				}
			}
		}
	}

	for li, line := range BetLines[:g.Sel] {
		var numw, numl slot.Pos = 0, 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = g.LX(x, line)
			if sx == wild || reelwild[x-1] {
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
	if count := g.SymNum(scat); count >= 3 {
		var pay, fs = ScatPay[count-1], ScatFreespin[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * float64(g.Sel) * pay,
			MP:  1,
			Sym: scat,
			Num: count,
			XY:  g.SymPos(scat),
			FS:  fs,
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
