package royaldynasty

// See: https://www.slotsmate.com/software/novomatic/royal-dynasty

import (
	_ "embed"
	"math/rand/v2"

	"github.com/slotopol/server/game/slot"
)

//go:embed royaldynasty_bon.yaml
var rbon []byte

var ReelsBon = slot.ReadObj[*slot.Reels5x](rbon)

//go:embed royaldynasty_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [14][5]float64{
	{0, 10, 200, 1200, 5000}, //  1 oldking
	{0, 0, 25, 100, 500},     //  2 lord
	{0, 0, 20, 80, 300},      //  3 gryphon
	{0, 0, 20, 80, 300},      //  4 leon
	{0, 0, 15, 60, 200},      //  5 money
	{0, 0, 15, 60, 200},      //  6 shield
	{0, 0, 10, 50, 150},      //  7 ace
	{0, 0, 10, 50, 150},      //  8 king
	{0, 0, 5, 25, 100},       //  9 queen
	{0, 0, 5, 25, 100},       // 10 jack
	{0, 0, 5, 25, 100},       // 11 ten
	{0, 0, 5, 25, 100},       // 12 nine
	{},                       // 13 princess
	{},                       // 14 prince
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 5, 25, 500} // scatter

var Freegames = [...]int{25, 30, 35, 40, 45}

// Bet lines
var BetLines = slot.BetLinesNvm20v2

type Game struct {
	slot.Screen5x3 `yaml:",inline"`
	slot.Slotx     `yaml:",inline"`
	// Scatter on freegames that triggers freegames
	TS slot.Sym `json:"ts" yaml:"ts" xml:"ts"`
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

const wild, scat1, scat2 = 1, 13, 14

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
			payw = LinePay[wild-1][numw-1]
		}
		if numl >= 3 && syml > 0 {
			payl = LinePay[syml-1][numl-1]
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
	var ps slot.Sym // pays scatter
	if g.TS == scat1 {
		ps = scat2
	} else if g.TS == scat2 {
		ps = scat1
	}
	for li, line := range BetLines[:g.Sel] {
		var mw float64 = 1 // mult wild
		var numw, numl slot.Pos = 0, 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = g.LY(x, line)
			if sx == wild || sx == ps {
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
			payw = LinePay[wild-1][numw-1]
		}
		if numl >= 3 && syml > 0 {
			payl = LinePay[syml-1][numl-1]
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
		var fs = Freegames[rand.N(len(Freegames))]
		*wins = append(*wins, slot.WinItem{
			Sym:  scat1,
			Num:  count,
			XY:   g.ScatPos(scat1),
			Free: fs,
		})
	} else if count := g.ScatNum(scat2); count >= 3 {
		var fs = Freegames[rand.N(len(Freegames))]
		*wins = append(*wins, slot.WinItem{
			Sym:  scat2,
			Num:  count,
			XY:   g.ScatPos(scat2),
			Free: fs,
		})
	}
}

// Scatters calculation on bonus games.
func (g *Game) ScanScattersBon(wins *slot.Wins) {
	if count := g.ScatNum(g.TS); count >= 3 {
		var pay = ScatPay[count-1]
		var fs = Freegames[rand.N(len(Freegames))]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * pay,
			Mult: 1,
			Sym:  g.TS,
			Num:  count,
			XY:   g.ScatPos(g.TS),
			Free: fs,
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

func (g *Game) Prepare() {
	if g.FSR == 0 {
		g.TS = 0 // reset trigger scatter
	}
}

func (g *Game) Apply(wins slot.Wins) {
	g.Slotx.Apply(wins)

	if g.FSR > 0 && g.FSN == 0 {
		for _, wi := range wins {
			if wi.Free > 0 {
				g.TS = wi.Sym // set trigger scatter
			}
		}
	}
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
