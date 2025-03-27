package bookofra

// See: https://www.slotsmate.com/software/novomatic/book-of-ra-deluxe

import (
	_ "embed"
	"math/rand/v2"

	"github.com/slotopol/server/game/slot"
)

//go:embed bookofra_bon.yaml
var rbon []byte

var ReelsBon = slot.ReadObj[*slot.Reels5x](rbon)

//go:embed bookofra_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [10][5]float64{
	{0, 10, 100, 1000, 5000}, //  1 explorer
	{0, 5, 40, 400, 2000},    //  2 mummy
	{0, 5, 30, 100, 750},     //  3 statue
	{0, 5, 30, 100, 750},     //  4 scarab
	{0, 0, 5, 40, 150},       //  5 ace
	{0, 0, 5, 40, 150},       //  6 king
	{0, 0, 5, 25, 100},       //  7 queen
	{0, 0, 5, 25, 100},       //  8 jack
	{0, 0, 5, 25, 100},       //  9 ten
	{},                       // 10 tomb
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 2, 20, 200} // 10 tomb

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 10, 10, 10} // 10 tomb

// Bet lines
var BetLines = slot.BetLinesNvm10

type Game struct {
	slot.Slotx[slot.Screen5x3] `yaml:",inline"`
	// Expanding Symbol
	ES slot.Sym `json:"es,omitempty" yaml:"es,omitempty" xml:"es,omitempty"`
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

const book = 10

func (g *Game) Scanner(wins *slot.Wins) {
	g.ScanLined(wins)
	g.ScanScatters(wins)
}

func LineES(src, xye slot.Linex) (dst slot.Linex) {
	for x := range 5 {
		if xye[x] > 0 {
			dst[x] = src[x]
		}
	}
	return
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]

		var numw, numl slot.Pos = 0, 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = g.Scr.LY(x, line)
			if sx == book {
				if syml == 0 {
					numw = x
				}
			} else if syml == 0 {
				if sx == g.ES && numw == 0 {
					break
				}
				syml = sx
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		if syml > 0 {
			if payl := LinePay[syml-1][numl-1]; payl > 0 {
				*wins = append(*wins, slot.WinItem{
					Pay:  g.Bet * payl,
					Mult: 1,
					Sym:  syml,
					Num:  numl,
					Line: li,
					XY:   line.CopyL(numl),
				})
			}
		}
	}

	if g.ES == 0 {
		return
	}
	var nume = g.Scr.ScatNum(g.ES)
	if nume < 2 {
		return
	}
	var paye = LinePay[g.ES-1][nume-1]
	if paye == 0 {
		return
	}
	var xye = g.Scr.ScatPos(g.ES)
	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * paye,
			Mult: 1,
			Sym:  g.ES,
			Num:  nume,
			Line: li,
			XY:   LineES(line, xye),
		})
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.Scr.ScatNum(book); count >= 3 {
		var pay, fs = ScatPay[count-1], ScatFreespin[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * pay,
			Mult: 1,
			Sym:  book,
			Num:  count,
			XY:   g.Scr.ScatPos(book),
			Free: fs,
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	if g.FSR == 0 {
		var reels, _ = slot.FindClosest(ReelsMap, mrtp)
		g.Scr.Spin(reels)
	} else {
		g.Scr.Spin(ReelsBon)
	}
}

func (g *Game) Prepare() {
	if g.FSR == 0 {
		g.ES = 0 // clear expanding symbol if no freegames
	}
}

func (g *Game) Apply(wins slot.Wins) {
	g.Slotx.Apply(wins)

	if g.FSR > 0 && g.ES == 0 { // setup expanding symbol
		g.ES = rand.N[slot.Sym](9) + 1
	}
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
