package simsalabim

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
	"github.com/slotopol/server/game/slot/netent/diamonddogs"
)

//go:embed simsalabim_bon.yaml
var rbon []byte

var ReelsBon = slot.ReadObj[*slot.Reels5x](rbon)

//go:embed simsalabim_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [11][5]float64{
	{0, 0, 50, 150, 1000},   //  1 hat
	{0, 0, 25, 100, 500},    //  2 chest
	{0, 0, 15, 75, 300},     //  3 cell
	{0, 0, 10, 50, 200},     //  4 cards
	{0, 0, 5, 50, 150},      //  5 ace
	{0, 0, 5, 25, 100},      //  6 king
	{0, 0, 5, 25, 100},      //  7 queen
	{0, 0, 5, 25, 100},      //  8 jack
	{},                      //  9 bonus
	{0, 5, 200, 2000, 7500}, // 10 wild
	{},                      // 11 scatter
}

// Scatters payment.
var ScatPay = [5]float64{0, 2, 4, 25, 200} // 11 scatter

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 10, 20, 30} // 11 scatter

// Bet lines
var BetLines = slot.BetLinesNetEnt5x3[:25]

const (
	ne12 = 1 // bonus ID
)

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

const bon, wild, scat = 9, 10, 11

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
				} else if syml == bon {
					numl = x - 1
					break
				}
			} else if syml == 0 && sx != scat {
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
		if payl > payw {
			var mm float64 = 1 // mult mode
			if g.FSR > 0 {
				mm = 3
			}
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payl,
				Mult: mm,
				Sym:  syml,
				Num:  numl,
				Line: li + 1,
				XY:   line.CopyL(numl),
			})
		} else if payw > 0 {
			var mm float64 = 1 // mult mode
			if g.FSR > 0 {
				mm = 3
			}
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payw,
				Mult: mm,
				Sym:  wild,
				Num:  numw,
				Line: li + 1,
				XY:   line.CopyL(numw),
			})
		} else if syml == bon && numl >= 3 { // appear on regular games only
			*wins = append(*wins, slot.WinItem{
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li + 1,
				XY:   line.CopyL(numl),
				BID:  ne12,
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.ScatNum(scat); count >= 2 {
		var mm float64 = 1 // mult mode
		if g.FSR > 0 {
			mm = 3
		}
		var pay, fs = ScatPay[count-1], ScatFreespin[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * pay,
			Mult: mm,
			Sym:  scat,
			Num:  count,
			XY:   g.ScatPos(scat),
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

func (g *Game) Spawn(wins slot.Wins, fund, mrtp float64) {
	for i, wi := range wins {
		switch wi.BID {
		case ne12:
			wins[i].Bon, wins[i].Pay = diamonddogs.BonusSpawn(g.Bet)
		}
	}
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
