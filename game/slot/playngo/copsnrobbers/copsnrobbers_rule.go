package copsnrobbers

// See: https://freeslotshub.com/playngo/cop-the-lot/
// See: https://www.slotsmate.com/software/play-n-go/cops-n-robbers

import (
	"math/rand/v2"

	"github.com/slotopol/server/game/slot"
)

var ReelsBon slot.Reelx

var ReelsMap slot.ReelsMap[slot.Reelx]

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
var BetLines = slot.BetLinesPlt5x3[:]

const (
	Efs = 17  // average free spins for ScatRand set
	Pfs = 0.3 // probability of "got away" at free spins
)

type Game struct {
	slot.Grid5x3 `yaml:",inline"`
	slot.Slotx   `yaml:",inline"`
	// multiplier for free spins, if them ended by "got away"
	M float64 `json:"m,omitempty" yaml:"m,omitempty" xml:"m,omitempty"`
}

// Declare conformity with SlotGeneric interface.
var _ slot.SlotGeneric = (*Game)(nil)

func NewGame(sel int) *Game {
	return &Game{
		Slotx: slot.Slotx{
			Sel: sel,
			Bet: 1,
		},
		M: 0,
	}
}

func (g *Game) Clone() slot.SlotGeneric {
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
			var sx = g.LX(x, line)
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
				Pay: g.Bet * payl,
				MP:  mw * mm,
				Sym: syml,
				Num: numl,
				LI:  li + 1,
				XY:  line.HitxL(numl),
			})
		} else if payw > 0 {
			var mm float64 = 1 // mult mode
			if g.FSR > 0 {
				mm = g.M
			}
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * payw,
				MP:  mm,
				Sym: wild,
				Num: numw,
				LI:  li + 1,
				XY:  line.HitxL(numw),
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.SymNum(scat); count >= 2 {
		var pay, fs = ScatPay[count-1], 0
		if count >= 3 {
			fs = ScatRand[rand.N(len(ScatRand))]
		}
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * float64(g.Sel) * pay,
			MP:  1,
			Sym: scat,
			Num: count,
			XY:  g.SymPos(scat),
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
			if wi.FS > 0 {
				g.FSR = wi.FS
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
