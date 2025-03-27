package piggyriches

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed piggyriches_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [12][5]float64{
	{},                    //  1 wild
	{},                    //  2 scatter
	{0, 5, 25, 300, 2000}, //  3 money bag
	{0, 0, 25, 150, 1000}, //  4 banknotes
	{0, 0, 20, 125, 750},  //  5 keys
	{0, 0, 20, 75, 400},   //  6 wallet
	{0, 0, 15, 75, 200},   //  7 piggy bank
	{0, 0, 15, 50, 125},   //  8 ace
	{0, 0, 10, 25, 100},   //  9 king
	{0, 0, 5, 20, 75},     // 10 queen
	{0, 0, 5, 15, 60},     // 11 jack
	{0, 0, 5, 10, 50},     // 12 ten
}

// Scatters payment.
var ScatPay = [5]float64{0, 2, 4, 15, 100} // 2 scatter

// Bet lines
var BetLines = slot.BetLinesNetEnt5x3[:15]

type Game struct {
	slot.Slotx[slot.Screen5x3] `yaml:",inline"`
	// multiplier on freespins
	M float64 `json:"m,omitempty" yaml:"m,omitempty" xml:"m,omitempty"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slotx: slot.Slotx[slot.Screen5x3]{
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

const wild, scat = 1, 2

func (g *Game) Scanner(wins *slot.Wins) {
	g.ScanLined(wins)
	g.ScanScatters(wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]

		var mw float64 = 1 // mult wild
		var numl slot.Pos = 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = g.Scr.LY(x, line)
			if sx == wild {
				mw = 3
			} else if syml == 0 && sx != scat {
				syml = sx
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		if numl >= 2 && syml > 0 {
			if pay := LinePay[syml-1][numl-1]; pay > 0 {
				var mm float64 = 1 // mult mode
				if g.FSR > 0 {
					mm = g.M
				}
				*wins = append(*wins, slot.WinItem{
					Pay:  g.Bet * pay,
					Mult: mw * mm,
					Sym:  syml,
					Num:  numl,
					Line: li,
					XY:   line.CopyL(numl),
				})
			}
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	var count = g.Scr.ScatNum(scat)
	if g.FSR > 0 {
		*wins = append(*wins, slot.WinItem{
			Pay:  0,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   g.Scr.ScatPos(scat),
			Free: int(count),
		})
	} else if count >= 2 {
		var pay, fs = ScatPay[count-1], 0
		if count >= 3 {
			fs = 15
		}
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
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	g.Scr.Spin(reels)
}

func (g *Game) Prepare() {
	if g.FSR == 0 {
		g.M = 0
	}
}

func (g *Game) Apply(wins slot.Wins) {
	g.Slotx.Apply(wins)

	if g.FSR > 0 && g.M == 0 {
		g.M = 3
	}
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}

// Mode can be:
// * 2 with 22 free spins
// * 3 with 15 free spins
// * 5 with 9 free spins
func (g *Game) SetMode(n int) error {
	var fs = g.FSR * int(g.M)
	switch n {
	case 2:
		g.FSR = fs / 2
	case 3:
		g.FSR = fs / 3
	case 5:
		g.FSR = fs / 5
	default:
		return slot.ErrDisabled
	}
	return nil
}
