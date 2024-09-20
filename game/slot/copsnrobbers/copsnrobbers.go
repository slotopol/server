package copsnrobbers

// See: https://freeslotshub.com/playngo/cop-the-lot/

import (
	"math/rand/v2"

	slot "github.com/slotopol/server/game/slot"
)

// Lined payment.
var LinePay = [11][5]float64{
	{0, 3, 30, 300, 3000}, //  1 wild
	{0, 0, 0, 0, 0},       //  2 scatter
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

const (
	Efs = 17  // average free spins for ScatRand set
	Pfs = 0.3 // probability of "got away" at free spins
)

type Game struct {
	slot.Slot5x3 `yaml:",inline"`
	// free spin number
	FS int `json:"fs,omitempty" yaml:"fs,omitempty" xml:"fs,omitempty"`
	// multiplier for free spins, if them ended by "got away"
	M float64 `json:"m,omitempty" yaml:"m,omitempty" xml:"m,omitempty"`
}

func NewGame() *Game {
	return &Game{
		Slot5x3: slot.Slot5x3{
			Sel: slot.MakeBitNum(9, 1),
			Bet: 1,
		},
		FS: 0,
		M:  0,
	}
}

const wild, scat = 1, 2

var bl = slot.Lineset5x{
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

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	g.ScanLined(screen, wins)
	if g.FS == 0 {
		g.ScanScatters(screen, wins)
	}
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen slot.Screen, wins *slot.Wins) {
	for li := range g.Sel.Bits() {
		var line = bl.Line(li)

		var numw, numl = 0, 5
		var syml slot.Sym
		var mw float64 = 1 // mult wild
		for x := 1; x <= 5; x++ {
			var sx = screen.Pos(x, line)
			if sx == wild {
				if syml == 0 {
					numw = x
				}
				mw = 2
			} else if syml == 0 && sx != scat {
				syml = sx
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		var payw, payl float64
		if numw > 0 {
			payw = LinePay[wild-1][numw-1]
		}
		if numl > 0 && syml > 0 {
			payl = LinePay[syml-1][numl-1]
		}
		if payl*mw > payw {
			var mm float64 = 1 // mult mode
			if g.FS > 0 {
				mm = g.M
			}
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payl,
				Mult: mw * mm,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyL(numl),
			})
		} else if payw > 0 {
			var mm float64 = 1 // mult mode
			if g.FS > 0 {
				mm = g.M
			}
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payw,
				Mult: mm,
				Sym:  wild,
				Num:  numw,
				Line: li,
				XY:   line.CopyL(numw),
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen slot.Screen, wins *slot.Wins) {
	if count := screen.ScatNum(scat); count >= 2 {
		var pay, fs = ScatPay[count-1], 0
		if count >= 3 {
			fs = ScatRand[rand.N(len(ScatRand))]
		}
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel.Num()) * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
			Free: fs,
		})
	}
}

func (g *Game) Spin(screen slot.Screen, mrtp float64) {
	if g.FS == 0 {
		var _, reels = FindReels(mrtp)
		screen.Spin(reels)
	} else {
		screen.Spin(&ReelsBon)
	}
}

func (g *Game) Apply(screen slot.Screen, wins slot.Wins) {
	if g.FS > 0 {
		g.Gain += wins.Gain()
	} else {
		g.Gain = wins.Gain()
	}

	if g.FS > 0 {
		g.FS--
		if g.FS == 0 {
			g.M = 0 // no multiplier on regular games
		}
	} else { // free spins can not be nested
		for _, wi := range wins {
			if wi.Free > 0 {
				g.FS = wi.Free
				if rand.Float64() <= Pfs {
					g.M = 2
				} else {
					g.M = 1
				}
			}
		}
	}
}

func (g *Game) FreeSpins() int {
	return g.FS
}

func (g *Game) SetSel(sel slot.Bitset) error {
	if sel.IsZero() {
		return slot.ErrNoLineset
	}
	if bs := sel; !bs.AndNot(slot.MakeBitNum(len(bl), 1)).IsZero() {
		return slot.ErrLinesetOut
	}
	if g.FreeSpins() > 0 {
		return slot.ErrNoFeature
	}
	g.Sel = sel
	return nil
}
