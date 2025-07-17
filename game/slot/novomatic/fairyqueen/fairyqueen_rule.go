package fairyqueen

// See: https://www.slotsmate.com/software/novomatic/fairy-queen

import (
	_ "embed"
	"math/rand/v2"

	"github.com/slotopol/server/game/slot"
)

//go:embed fairyqueen_es.yaml
var esreel []byte

var ExpSymReel = slot.ReadObj[[]slot.Sym](esreel)

//go:embed fairyqueen_bon.yaml
var rbon []byte

var ReelsBon = slot.ReadObj[*slot.Reels5x](rbon)

//go:embed fairyqueen_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [13][5]float64{
	{0, 10, 250, 2500, 9000}, //  1 princess
	{0, 2, 25, 125, 750},     //  2 troll1
	{0, 2, 25, 125, 750},     //  3 troll2
	{0, 0, 20, 100, 400},     //  4 dragon
	{0, 0, 15, 75, 250},      //  5 mushroom
	{0, 0, 15, 75, 250},      //  6 herb
	{0, 0, 10, 50, 125},      //  7 ace
	{0, 0, 10, 50, 125},      //  8 king
	{0, 0, 5, 25, 100},       //  9 queen
	{0, 0, 5, 25, 100},       // 10 jack
	{0, 0, 5, 25, 100},       // 11 ten
	{0, 2, 5, 25, 100},       // 12 nine
	{},                       // 13 light
}

// Number of reels filled by expanding symbols on bonus games.
var ReelNumBon = [13]slot.Pos{
	2, //  1 princess
	2, //  2 troll1
	2, //  3 troll2
	3, //  4 dragon
	3, //  5 mushroom
	3, //  6 herb
	3, //  7 ace
	3, //  8 king
	3, //  9 queen
	3, // 10 jack
	3, // 11 ten
	2, // 12 nine
	2, // 13 light
}

// Scatters payment.
var ScatPay = [5]float64{0, 2, 5, 75, 750} // 13 light

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 10, 10, 10} // 13 light

// Bet lines
var BetLines = slot.BetLinesNvm10

type Game struct {
	slot.Screen5x3 `yaml:",inline"`
	slot.Slotx     `yaml:",inline"`
	// Expanding Symbol
	ES slot.Sym `json:"es,omitempty" yaml:"es,omitempty" xml:"es,omitempty"`
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

const wild, scat = 1, 13

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
				Pay:  g.Bet * payl,
				Mult: 1,
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

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.ScatNum(scat); count >= 2 {
		var pay, fs = ScatPay[count-1], ScatFreespin[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   g.ScatPos(scat),
			Free: fs,
		})
	}
}

func (g *Game) SpinBon() {
	var num = ReelNumBon[g.ES]
	var x slot.Pos
	for x = 1; x <= num; x++ {
		var r = &g.Scr[x-1]
		if g.ES != scat {
			r[0], r[1], r[2] = g.ES, g.ES, g.ES
		} else {
			r[0], r[1], r[2] = 0, scat, 0
		}
	}
	for x = num + 1; x <= 5; x++ {
		var reel = ReelsBon.Reel(x)
		var hit = rand.N(len(reel))
		g.SetCol(x, reel, hit)
	}
}

func (g *Game) Spin(mrtp float64) {
	if g.FSR == 0 {
		var reels, _ = slot.FindClosest(ReelsMap, mrtp)
		g.ReelSpin(reels)
	} else {
		g.SpinBon()
	}
}

func (g *Game) Prepare() {
	if g.FSR > 0 { // setup expanding symbol
		g.ES = ExpSymReel[rand.N(len(ExpSymReel))]
	} else {
		g.ES = 0
	}
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
