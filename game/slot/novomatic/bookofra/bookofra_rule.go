package bookofra

// See: https://www.slotsmate.com/software/novomatic/book-of-ra-deluxe

import (
	"math/rand/v2"

	"github.com/slotopol/server/game/slot"
)

var ReelsBon slot.Reelx

var ReelsMap slot.ReelsMap[slot.Reelx]

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
var BetLines = slot.BetLinesNvm10[:]

type Game struct {
	slot.Screen5x3 `yaml:",inline"`
	slot.Slotx     `yaml:",inline"`
	// Expanding Symbol
	ES slot.Sym `json:"es,omitempty" yaml:"es,omitempty" xml:"es,omitempty"`
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

const wsc = 10

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

func LineES(src slot.Linex, xye slot.Hitx) (dst slot.Hitx) {
	var i slot.Pos
	for i = 0; xye[i][0] > 0; i++ {
		dst[i][0], dst[i][1] = xye[i][0], src[xye[i][0]-1]
	}
	dst[i][0] = 0
	return
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	for li, line := range BetLines[:g.Sel] {
		var numw, numl slot.Pos = 0, 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = g.LX(x, line)
			if sx == wsc {
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

		if numl >= 2 && syml > 0 {
			if payl := LinePay[syml-1][numl-1]; payl > 0 {
				*wins = append(*wins, slot.WinItem{
					Pay: g.Bet * payl,
					MP:  1,
					Sym: syml,
					Num: numl,
					LI:  li + 1,
					XY:  line.HitxL(numl),
				})
			}
		}
	}

	if g.ES == 0 {
		return
	}
	var nume = g.SymNum(g.ES)
	if nume < 2 {
		return
	}
	var paye = LinePay[g.ES-1][nume-1]
	if paye == 0 {
		return
	}
	var xye = g.SymPos(g.ES)
	for li, line := range BetLines[:g.Sel] {
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * paye,
			MP:  1,
			Sym: g.ES,
			Num: nume,
			LI:  li + 1,
			XY:  LineES(line, xye),
		})
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.SymNum(wsc); count >= 3 {
		var pay, fs = ScatPay[count-1], ScatFreespin[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * float64(g.Sel) * pay,
			MP:  1,
			Sym: wsc,
			Num: count,
			XY:  g.SymPos(wsc),
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
