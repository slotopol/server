package kingofclovers

// See: https://www.slotsmate.com/software/ct-interactive/king-of-clovers

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsBon slot.Reelx

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [10][5]float64{
	{},                   //  1 wild (2, 3, 4 reels only)
	{},                   //  2 scatter
	{0, 0, 25, 50, 1000}, //  3 seven
	{0, 0, 5, 15, 100},   //  4 coin
	{0, 0, 5, 15, 100},   //  5 bell
	{0, 0, 5, 15, 100},   //  6 horseshoe
	{0, 0, 5, 10, 100},   //  7 apple
	{0, 0, 5, 10, 100},   //  8 lemon
	{0, 0, 5, 10, 100},   //  9 plum
	{0, 0, 5, 10, 100},   // 10 cherry
}

// Scatters payment on regular games.
var ScatPayReg = [5]float64{0, 0, 5, 20, 100} // 2 scatter

// Scatters payment on bonus games.
var ScatPayBon = [5]float64{0, 2, 5, 20, 100} // 2 scatter

// Scatter freespins table on regular games.
var ScatFreespinReg = [5]int{0, 0, 14, 14, 14} // 2 scatter

// Scatter freespins table on bonus games.
var ScatFreespinBon = [5]int{0, 5, 14, 14, 14} // 2 scatter

// Bet lines
var BetLines = slot.BetLinesCT5x3[:]

type Game struct {
	slot.Cascade5x3 `yaml:",inline"`
	slot.Slotx      `yaml:",inline"`
}

// Declare conformity with SlotCascade interface.
var _ slot.SlotCascade = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slotx: slot.Slotx{
			Sel: 30,
			Bet: 1,
		},
	}
}

func (g *Game) Clone() slot.SlotGame {
	var clone = *g
	return &clone
}

func (g *Game) FreeMode() bool {
	return g.FSR != 0 || g.Cascade()
}

const wild, scat = 1, 2

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	if g.FSR == 0 {
		g.ScanScattersReg(wins)
	} else {
		g.ScanScattersBon(wins)
	}
	return nil
}

func (g *Game) FillMult() float64 {
	var n int
	if g.Scr[1][0] == wild && g.Scr[1][1] == wild && g.Scr[1][2] == wild {
		n++
	}
	if g.Scr[2][0] == wild && g.Scr[2][1] == wild && g.Scr[2][2] == wild {
		n++
	}
	if g.Scr[3][0] == wild && g.Scr[3][1] == wild && g.Scr[3][2] == wild {
		n++
	}
	if n == 2 {
		return 2
	} else if n == 3 {
		return 10
	}
	return 1
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	var fm float64 // fill mult
	for li, line := range BetLines[:g.Sel] {
		var numl slot.Pos = 5
		var syml = g.LX(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = g.LX(x, line)
			if sx != syml && sx != wild {
				numl = x - 1
				break
			}
		}
		if pay := LinePay[syml-1][numl-1]; pay > 0 {
			if fm == 0 { // lazy calculation
				fm = g.FillMult()
			}
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * pay,
				MP:  fm,
				Sym: syml,
				Num: numl,
				LI:  li + 1,
				XY:  line.HitxL(numl),
			})
		}
	}
}

// Scatters calculation on regular games.
func (g *Game) ScanScattersReg(wins *slot.Wins) {
	if count := g.SymNum(scat); count >= 3 {
		var pay, fs = ScatPayReg[count-1], ScatFreespinReg[count-1]
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

// Scatters calculation on bonus games.
func (g *Game) ScanScattersBon(wins *slot.Wins) {
	if count := g.SymNum(scat); count >= 2 {
		var pay, fs = ScatPayBon[count-1], ScatFreespinBon[count-1]
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
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.SpinReels(reels)
}

func (g *Game) Prepare() {
	g.UntoFall()
}

func (g *Game) Apply(wins slot.Wins) {
	g.Slotx.Apply(wins)
	g.Strike(wins)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
