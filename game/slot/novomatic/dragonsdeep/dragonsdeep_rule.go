package dragonsdeep

// See: https://www.slotsmate.com/software/novomatic/dragons-deep

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed dragonsdeep_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [13][5]float64{
	{0, 10, 50, 300, 1000}, //  1 dragon
	{0, 5, 25, 250, 500},   //  2 viking1
	{0, 5, 25, 250, 500},   //  3 viking2
	{0, 5, 25, 250, 500},   //  4 viking3
	{0, 0, 10, 50, 200},    //  5 armor
	{0, 0, 10, 50, 200},    //  6 chest
	{0, 0, 5, 25, 100},     //  7 ace
	{0, 0, 5, 25, 100},     //  8 king
	{0, 0, 5, 25, 100},     //  9 queen
	{0, 0, 5, 25, 100},     // 10 jack
	{0, 0, 5, 25, 100},     // 11 ten
	{0, 0, 5, 25, 100},     // 12 nine
	{},                     // 13 egg
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 5, 20, 100} // 13 egg

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 5, 10, 25} // 13 egg

// Bet lines
var BetLines = slot.BetLinesNvm25

type Game struct {
	slot.Slotx[slot.Screen5x3] `yaml:",inline"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slotx: slot.Slotx[slot.Screen5x3]{
			Sel: len(BetLines),
			Bet: 1,
		},
	}
}

func (g *Game) Clone() slot.SlotGame {
	var clone = *g
	return &clone
}

const wild, scat = 1, 13
const viking1, viking2, viking3 = 2, 3, 4

func (g *Game) Scanner(wins *slot.Wins) {
	if g.FSR == 0 {
		g.ScanLinedReg(wins)
	} else {
		g.ScanLinedBon(wins)
	}
	g.ScanScatters(wins)
}

// Lined symbols calculation for regular games.
func (g *Game) ScanLinedReg(wins *slot.Wins) {
	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]

		var numw, numl slot.Pos = 0, 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = g.Scr.Pos(x, line)
			if sx == wild {
				if syml == 0 {
					numw = x
				}
			} else if syml == 0 && sx != scat {
				syml = sx
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		var payw, payl float64
		if numw > 0 {
			payw = LinePay[wild-1][numw-1]
		}
		if numl > 0 && syml > 0 {
			payl = LinePay[syml-1][numl-1]
		}
		if payl > payw {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payl,
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyL(numl),
			})
		} else if payw > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payw,
				Mult: 1,
				Sym:  wild,
				Num:  numw,
				Line: li,
				XY:   line.CopyL(numw),
			})
		}
	}
}

// Lined symbols calculation for bonus games.
func (g *Game) ScanLinedBon(wins *slot.Wins) {
	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]

		var numw, numl slot.Pos = 0, 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = g.Scr.Pos(x, line)
			if sx == wild || sx == viking1 || sx == viking2 || sx == viking3 {
				if syml == 0 {
					numw = x
				}
			} else if syml == 0 && sx != scat {
				syml = sx
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		var payw, payl float64
		if numw > 0 {
			payw = LinePay[wild-1][numw-1]
		}
		if numl > 0 && syml > 0 {
			payl = LinePay[syml-1][numl-1]
		}
		if payl > payw {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payl,
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyL(numl),
			})
		} else if payw > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payw,
				Mult: 1,
				Sym:  wild,
				Num:  numw,
				Line: li,
				XY:   line.CopyL(numw),
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.Scr.ScatNum(scat); count >= 3 {
		var pay, fs = ScatPay[count-1], ScatFreespin[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   g.Scr.ScatPos(scat),
			Free: fs,
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = slot.FindReels(ReelsMap, mrtp)
	g.Scr.Spin(reels)
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
