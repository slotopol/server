package piggyriches

import (
	slot "github.com/slotopol/server/game/slot"
)

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
var bl = slot.BetLinesNetEnt15

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
			Sel: slot.MakeBitNum(len(bl), 1),
			Bet: 1,
		},
		FS: 0,
		M:  0,
	}
}

const wild, scat = 1, 2

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	g.ScanLined(screen, wins)
	g.ScanScatters(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen slot.Screen, wins *slot.Wins) {
	for li := g.Sel.Next(0); li != -1; li = g.Sel.Next(li) {
		var line = bl[li-1]

		var mw float64 = 1 // mult wild
		var numl slot.Pos = 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = screen.Pos(x, line)
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
				if g.FS > 0 {
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
func (g *Game) ScanScatters(screen slot.Screen, wins *slot.Wins) {
	var count = screen.ScatNum(scat)
	if g.FS > 0 {
		*wins = append(*wins, slot.WinItem{
			Pay:  0,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
			Free: int(count),
		})
	} else if count >= 2 {
		var pay, fs = ScatPay[count-1], 0
		if count >= 3 {
			fs = 15
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
	var _, reels = slot.FindReels(ReelsMap, mrtp)
	screen.Spin(reels)
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
			if g.M == 0 {
				g.M = 3
			}
		}
	}
	if g.FS == 0 {
		g.M = 0
	}
}

func (g *Game) FreeSpins() int {
	return g.FS
}

func (g *Game) SetSel(sel slot.Bitset) error {
	return g.SetSelNum(sel, len(bl))
}

// Mode can be:
// * 2 with 22 free spins
// * 3 with 15 free spins
// * 5 with 9 free spins
func (g *Game) SetMode(n int) error {
	var fs = g.FS * int(g.M)
	switch n {
	case 2:
		g.FS = fs / 2
	case 3:
		g.FS = fs / 3
	case 5:
		g.FS = fs / 5
	default:
		return slot.ErrDisabled
	}
	return nil
}
