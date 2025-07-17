package chicago

// See: https://www.slotsmate.com/software/novomatic/novomatic-chicago

import (
	_ "embed"
	"math/rand/v2"

	"github.com/slotopol/server/game/slot"
)

//go:embed chicago_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [13][5]float64{
	{0, 10, 250, 2500, 10000}, //  1 chicago
	{0, 0, 50, 500, 2000},     //  2 capone
	{0, 0, 50, 500, 2000},     //  3 ness
	{0, 0, 30, 200, 1000},     //  4 woman
	{0, 0, 20, 100, 500},      //  5 policeman
	{0, 0, 20, 100, 500},      //  6 newsboy
	{0, 0, 10, 50, 250},       //  7 ace
	{0, 0, 10, 50, 250},       //  8 king
	{0, 0, 5, 20, 100},        //  9 queen
	{0, 0, 5, 20, 100},        // 10 jack
	{0, 0, 5, 20, 100},        // 11 ten
	{0, 0, 5, 20, 100},        // 12 nine
	{},                        // 13 cadillac
}

// Scatters payment.
var ScatPay = [5]float64{0, 2, 5, 20, 100} // 13 cadillac

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 12, 12, 12} // 13 cadillac

// Bet lines
var BetLines = slot.BetLinesNvm20v1

type Game struct {
	slot.Screen5x3 `yaml:",inline"`
	slot.Slotx     `yaml:",inline"`
	// multiplier on freespins
	M float64 `json:"m,omitempty" yaml:"m,omitempty" xml:"m,omitempty"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slotx: slot.Slotx{
			Sel: len(BetLines),
			Bet: 1,
		},
		M: 0,
	}
}

func (g *Game) Clone() slot.SlotGame {
	var clone = *g
	return &clone
}

const wild, scat = 1, 13

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	var mm float64 = 1 // mult mode
	if g.FSR > 0 {
		mm = g.M
		*wins = append(*wins, slot.WinItem{
			Mult: mm,
		})
	}

	for li, line := range BetLines[:g.Sel] {
		var numw, numl slot.Pos = 0, 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = g.LY(x, line)
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
		if numl >= 3 && syml > 0 {
			payl = LinePay[syml-1][numl-1]
		}
		if payl > payw {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payl,
				Mult: mm,
				Sym:  syml,
				Num:  numl,
				Line: li + 1,
				XY:   line.CopyL(numl),
			})
		} else if payw > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payw,
				Mult: mm,
				Sym:  wild,
				Num:  numw,
				Line: li + 1,
				XY:   line.CopyL(numw),
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.ScatNum(scat); count >= 2 {
		var mm float64 = 1 // mult mode
		if g.FSR > 0 {
			mm = g.M
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
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	g.ReelSpin(reels)
}

var MultChoose = []float64{1, 2, 3, 5, 10} // E = 4.2
const Emc = 4.2

func (g *Game) Prepare() {
	if g.FSR > 0 {
		g.M = MultChoose[rand.N(len(MultChoose))]
	} else {
		g.M = 0
	}
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
