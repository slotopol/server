package chicago

import (
	"math/rand/v2"

	slot "github.com/slotopol/server/game/slot"
)

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
	{0, 0, 0, 0, 0},           // 13 cadillac
}

// Scatters payment.
var ScatPay = [5]float64{0, 2, 5, 20, 100} // 13 cadillac

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 12, 12, 12} // 13 cadillac

const (
	jid = 1 // jackpot ID
)

// Jackpot win combinations.
var Jackpot = [13][5]int{
	{0, 0, 0, 0, 0}, //  1 chicago
	{0, 0, 0, 0, 0}, //  2 capone
	{0, 0, 0, 0, 0}, //  3 ness
	{0, 0, 0, 0, 0}, //  4 woman
	{0, 0, 0, 0, 0}, //  5 policeman
	{0, 0, 0, 0, 0}, //  6 newsboy
	{0, 0, 0, 0, 0}, //  7 ace
	{0, 0, 0, 0, 0}, //  8 king
	{0, 0, 0, 0, 0}, //  9 queen
	{0, 0, 0, 0, 0}, // 10 jack
	{0, 0, 0, 0, 0}, // 11 ten
	{0, 0, 0, 0, 0}, // 12 nine
	{0, 0, 0, 0, 0}, // 13 cadillac
}

type Game struct {
	slot.Slot5x3 `yaml:",inline"`
	// free spin number
	FS int `json:"fs,omitempty" yaml:"fs,omitempty" xml:"fs,omitempty"`
	// multiplier on freespins
	M float64 `json:"m,omitempty" yaml:"m,omitempty" xml:"m,omitempty"`
}

func NewGame() *Game {
	return &Game{
		Slot5x3: slot.Slot5x3{
			Sel: slot.MakeBitNum(20, 1),
			Bet: 1,
		},
		FS: 0,
		M:  0,
	}
}

const wild, scat = 1, 13

var bl = slot.BetLinesNvm20

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	g.ScanLined(screen, wins)
	g.ScanScatters(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen slot.Screen, wins *slot.Wins) {
	var mm float64 = 1 // mult mode
	if g.FS > 0 {
		mm = g.M
		*wins = append(*wins, slot.WinItem{
			Mult: mm,
		})
	}

	for li := g.Sel.Next(0); li != -1; li = g.Sel.Next(li) {
		var line = bl[li-1]

		var numw, numl slot.Pos = 0, 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = screen.Pos(x, line)
			if sx == wild {
				if syml == 0 {
					numw = x
				}
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
		if payl > payw {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payl,
				Mult: mm,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyL(numl),
			})
		} else if payw > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payw,
				Mult: mm,
				Sym:  wild,
				Num:  numw,
				Line: li,
				XY:   line.CopyL(numw),
				Jack: Jackpot[wild-1][numw-1],
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen slot.Screen, wins *slot.Wins) {
	if count := screen.ScatNum(scat); count >= 2 {
		var mm float64 = 1 // mult mode
		if g.FS > 0 {
			mm = g.M
		}
		var pay, fs = ScatPay[count-1], ScatFreespin[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel.Num()) * pay,
			Mult: mm,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
			Free: fs,
		})
	}
}

func (g *Game) Spin(screen slot.Screen, mrtp float64) {
	var _, reels = slot.FindReels(ReelsMap, mrtp)
	screen.Spin(reels)
}

var MultChoose = []float64{1, 1, 1, 2, 2, 2, 3, 3, 5, 10} // E = 3.0

func (g *Game) Prepare() {
	if g.FS > 0 {
		g.M = MultChoose[rand.N(len(MultChoose))]
	} else {
		g.M = 0
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
	}
	for _, wi := range wins {
		if wi.Free > 0 {
			g.FS += wi.Free
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
		return slot.ErrDisabled
	}
	g.Sel = sel
	return nil
}
