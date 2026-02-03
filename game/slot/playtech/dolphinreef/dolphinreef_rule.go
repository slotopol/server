package dolphinreef

// See: https://www.slotsmate.com/software/playtech/dolphin-reef

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [12][5]float64{
	{},                     //  1 wild (2, 4 reels only)
	{},                     //  2 scatter
	{0, 5, 100, 500, 5000}, //  3 parrot
	{0, 3, 30, 150, 1000},  //  4 turtle
	{0, 2, 20, 100, 500},   //  5 seastar
	{0, 2, 20, 100, 500},   //  6 seahorse
	{0, 0, 15, 50, 250},    //  7 ace
	{0, 0, 15, 50, 250},    //  8 king
	{0, 0, 10, 25, 150},    //  9 queen
	{0, 0, 10, 25, 150},    // 10 jack
	{0, 0, 5, 20, 100},     // 11 ten
	{0, 0, 5, 20, 100},     // 12 nine
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 3, 10, 100} // 2 scatter

// Bet lines
var BetLines = slot.BetLinesPlt5x3[:]

type Game struct {
	slot.Grid5x3 `yaml:",inline"`
	slot.Slotx     `yaml:",inline"`
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

const wild, scat = 1, 2

var reelwild = [5]bool{false, true, false, true, false}

func (g *Game) Scanner(wins *slot.Wins) error {
	if g.FSR == 0 {
		g.ScanLinedReg(wins)
		g.ScanScattersReg(wins)
	} else {
		g.ScanLinedBon(wins)
		g.ScanScattersBon(wins)
	}
	return nil
}

// Lined symbols calculation on regular games.
func (g *Game) ScanLinedReg(wins *slot.Wins) {
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
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * pay,
				MP:  1,
				Sym: syml,
				Num: numl,
				LI:  li + 1,
				XY:  line.HitxL(numl),
			})
		}
	}
}

// Lined symbols calculation on bonus games.
func (g *Game) ScanLinedBon(wins *slot.Wins) {
	for li, line := range BetLines[:g.Sel] {
		var numl slot.Pos = 5
		var syml = g.LX(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = g.LX(x, line)
			if sx != syml && !reelwild[x-1] {
				numl = x - 1
				break
			}
		}

		if pay := LinePay[syml-1][numl-1]; pay > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * pay,
				MP:  1,
				Sym: syml,
				Num: numl,
				LI:  li + 1,
				XY:  line.HitxL(numl),
			})
		}
	}
}

func (g *Game) ScatNumFG() (n slot.Pos) {
	for x := range 5 {
		if x == 1 || x == 3 {
			n++
			continue
		}
		var r = g.Grid[x]
		for y := range 3 {
			if r[y] == scat {
				n++
			}
		}
	}
	return
}

func (g *Game) ScatPosFG() (c slot.Hitx) {
	var x, y, i slot.Pos
	for x = range 5 {
		if x == 1 || x == 3 {
			c[i][0], c[i][1] = x+1, 0
			i++
			continue
		}
		var r = g.Grid[x]
		for y = range 3 {
			if r[y] == scat {
				c[i][0], c[i][1] = x+1, y+1
				i++
			}
		}
	}
	return
}

// Scatters calculation on regular games.
func (g *Game) ScanScattersReg(wins *slot.Wins) {
	var ns, nw = g.SymNum2(scat, wild)
	if ns+nw >= 3 {
		var pay = ScatPay[ns+nw-1]
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * float64(g.Sel) * pay,
			MP:  1,
			Sym: scat,
			Num: ns + nw,
			XY:  g.SymPos2(scat, wild),
		})
	}
	if nw >= 2 {
		*wins = append(*wins, slot.WinItem{
			MP:  1,
			Sym: wild,
			Num: nw,
			XY:  g.SymPos(wild),
			FS:  5,
		})
	}
}

// Scatters calculation on bonus games.
func (g *Game) ScanScattersBon(wins *slot.Wins) {
	if count := g.ScatNumFG(); count >= 3 {
		var pay = ScatPay[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * float64(g.Sel) * pay,
			MP:  1,
			Sym: scat,
			Num: count,
			XY:  g.ScatPosFG(),
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
