package richgirl

// See: https://www.slotsmate.com/software/igt/rich-girl

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed richgirl_bon.yaml
var rbon []byte

var ReelsBon = slot.ReadObj[*slot.Reels5x](rbon)

//go:embed richgirl_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment on regular games.
var LinePayReg = [12][5]float64{
	{0, 5, 50, 500, 10000}, //  1 wild
	{0, 0, 10, 100, 500},   //  2 girl
	{0, 0, 10, 100, 500},   //  3 father
	{0, 0, 10, 50, 200},    //  4 doggy
	{0, 0, 10, 50, 200},    //  5 kitty
	{0, 0, 5, 25, 100},     //  6 watermelon
	{0, 0, 5, 25, 100},     //  7 peach
	{0, 0, 5, 25, 100},     //  8 plum
	{0, 2, 5, 25, 100},     //  9 lemon
	{0, 2, 5, 25, 100},     // 10 cherry
	{},                     // 11 scatter
	{},                     // 12 diamond
}

// Lined payment on bonus games.
var LinePayBon = [5][5]float64{
	{0, 0, 50, 250, 1000}, // 1 diamond
	{0, 2, 5, 20, 50},     // 2 emerald
	{0, 2, 5, 20, 50},     // 3 sapphire
	{0, 2, 5, 20, 50},     // 4 ruby
	{0, 2, 5, 20, 50},     // 5 heliodor
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 2, 10, 25} // 11 scatter

// Bet lines
var BetLines = slot.BetLinesIgt5x3[:9]

type Game struct {
	slot.Screen5x3 `yaml:",inline"`
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

const wild, scat1, scat2 = 1, 11, 12

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

// Lined symbols calculation regular games.
func (g *Game) ScanLinedReg(wins *slot.Wins) {
	for li, line := range BetLines[:g.Sel] {
		var mw float64 = 1 // mult wild
		var numw, numl slot.Pos = 0, 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = g.LY(x, line)
			if sx == wild {
				if syml == 0 {
					numw = x
				}
				mw = 2
			} else if syml == 0 {
				syml = sx
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		var payw, payl float64
		if numw >= 2 {
			payw = LinePayReg[wild-1][numw-1]
		}
		if numl >= 2 && syml > 0 {
			payl = LinePayReg[syml-1][numl-1]
		}
		if payl*mw > payw {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payl,
				Mult: mw,
				Sym:  syml,
				Num:  numl,
				Line: li + 1,
				XY:   line.CopyL(numl),
			})
		} else if payw > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payw,
				Mult: 1,
				Sym:  wild,
				Num:  numw,
				Line: li + 1,
				XY:   line.CopyL(numw),
			})
		}
	}
}

// Lined symbols calculation on bonus games.
func (g *Game) ScanLinedBon(wins *slot.Wins) {
	for li, line := range BetLines[:g.Sel] {
		var mw float64 = 1 // mult wild
		var numw, numl slot.Pos = 0, 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = g.LY(x, line)
			if sx == wild {
				if syml == 0 {
					numw = x
				}
				mw = 2
			} else if syml == 0 {
				syml = sx
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		var payw, payl float64
		if numw >= 2 {
			payw = LinePayBon[wild-1][numw-1]
		}
		if numl >= 2 && syml > 0 {
			payl = LinePayBon[syml-1][numl-1]
		}
		if payl*mw > payw {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payl,
				Mult: mw,
				Sym:  syml,
				Num:  numl,
				Line: li + 1,
				XY:   line.CopyL(numl),
			})
		} else if payw > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payw,
				Mult: 1,
				Sym:  wild,
				Num:  numw,
				Line: li + 1,
				XY:   line.CopyL(numw),
			})
		}
	}
}

// Scatters calculation on regular games.
func (g *Game) ScanScattersReg(wins *slot.Wins) {
	if count := g.ScatNum(scat1); count >= 3 {
		var pay = ScatPay[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * pay,
			Mult: 1,
			Sym:  scat1,
			Num:  count,
			XY:   g.ScatPos(scat1),
		})
	} else if count := g.ScatNum(scat2); count >= 3 {
		*wins = append(*wins, slot.WinItem{
			Sym:  scat2,
			Num:  count,
			XY:   g.ScatPos(scat2),
			Free: 3,
		})
	}
}

// Scatters calculation on bonus games.
func (g *Game) ScanScattersBon(wins *slot.Wins) {
	if count := g.ScatNum(wild); count >= 1 {
		*wins = append(*wins, slot.WinItem{
			Sym:  wild,
			Num:  count,
			XY:   g.ScatPos(wild),
			Free: 1,
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	if g.FSR == 0 {
		var reels, _ = slot.FindClosest(ReelsMap, mrtp)
		g.ReelSpin(reels)
	} else {
		g.ReelSpin(ReelsBon)
	}
}

func (g *Game) Apply(wins slot.Wins) {
	g.Slotx.Apply(wins)
	if g.FSN >= 100 {
		g.FSR = 0
	}
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
