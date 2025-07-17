package copsnrobbers

// See: https://freeslotshub.com/playngo/cop-the-lot/
// See: https://www.slotsmate.com/software/play-n-go/cops-n-robbers

import (
	_ "embed"
	"math/rand/v2"

	"github.com/slotopol/server/game/slot"
)

//go:embed copsnrobbers_bon.yaml
var rbon []byte

var ReelsBon = slot.ReadObj[*slot.Reels5x](rbon)

//go:embed copsnrobbers_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [11][5]float64{
	{0, 3, 30, 300, 3000}, //  1 wild
	{},                    //  2 scatter
	{0, 2, 25, 150, 750},  //  3 money bag
	{0, 2, 20, 100, 500},  //  4 diamonds
	{0, 2, 15, 75, 500},   //  5 robbery
	{0, 0, 15, 75, 250},   //  6 picture
	{0, 0, 10, 75, 250},   //  7 watch
	{0, 0, 5, 50, 150},    //  8 cop
	{0, 0, 5, 50, 125},    //  9 jail
	{0, 0, 5, 25, 100},    // 10 thief
	{0, 0, 5, 25, 100},    // 11 handcuffs
}

// Scatters payment.
var ScatPay = [5]float64{0, 2, 3, 25, 250} // 2 scatter

var ScatRand = []int{10, 15, 15, 20, 25}

// Bet lines
var BetLines = []slot.Linex{
	{2, 2, 2, 2, 2}, // 1
	{1, 1, 1, 1, 1}, // 2
	{3, 3, 3, 3, 3}, // 3
	{1, 2, 3, 2, 1}, // 4
	{3, 2, 1, 2, 3}, // 5
	{2, 1, 1, 1, 2}, // 6
	{2, 3, 3, 3, 2}, // 7
	{1, 1, 2, 3, 3}, // 8
	{3, 3, 2, 1, 1}, // 9
}

const (
	Efs = 17  // average free spins for ScatRand set
	Pfs = 0.3 // probability of "got away" at free spins
)

type Game struct {
	slot.Screen5x3 `yaml:",inline"`
	slot.Slotx     `yaml:",inline"`
	// multiplier for free spins, if them ended by "got away"
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

const wild, scat = 1, 2

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	if g.FSR == 0 {
		g.ScanScatters(wins)
	}
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
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
		if numl >= 2 && syml > 0 {
			payl = LinePay[syml-1][numl-1]
		}
		if payl*mw > payw {
			var mm float64 = 1 // mult mode
			if g.FSR > 0 {
				mm = g.M
			}
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payl,
				Mult: mw * mm,
				Sym:  syml,
				Num:  numl,
				Line: li + 1,
				XY:   line.CopyL(numl),
			})
		} else if payw > 0 {
			var mm float64 = 1 // mult mode
			if g.FSR > 0 {
				mm = g.M
			}
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
		var pay, fs = ScatPay[count-1], 0
		if count >= 3 {
			fs = ScatRand[rand.N(len(ScatRand))]
		}
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
		g.M = 0 // no multiplier on regular games
	}
}

func (g *Game) Apply(wins slot.Wins) {
	if g.FSR != 0 {
		g.Gain += wins.Gain()
		g.FSN++
	} else {
		g.Gain = wins.Gain()
		g.FSN = 0
	}

	if g.FSR > 0 {
		g.FSR--
	} else { // free spins can not be nested
		for _, wi := range wins {
			if wi.Free > 0 {
				g.FSR = wi.Free
				if rand.Float64() <= Pfs {
					g.M = 2
				} else {
					g.M = 1
				}
			}
		}
	}
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
