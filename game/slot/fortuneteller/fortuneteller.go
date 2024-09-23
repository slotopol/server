package fortuneteller

// See: https://freeslotshub.com/playngo/fortune-teller/
// See: https://www.youtube.com/watch?v=bFQq3cCz9XY

import (
	slot "github.com/slotopol/server/game/slot"
)

// Lined payment.
var LinePay = [12][5]float64{
	{0, 0, 50, 500, 5000}, //  1 wild
	{0, 0, 0, 0, 0},       //  2 cat
	{0, 0, 0, 0, 0},       //  3 bonus
	{0, 0, 25, 250, 1000}, //  4 girl
	{0, 0, 15, 100, 500},  //  5 hand
	{0, 0, 15, 100, 500},  //  6 zodiac
	{0, 0, 15, 75, 250},   //  7 candle
	{0, 0, 10, 50, 150},   //  8 ace
	{0, 0, 10, 50, 150},   //  9 king
	{0, 0, 5, 25, 100},    // 10 queen
	{0, 0, 5, 25, 100},    // 11 jack
	{0, 0, 5, 25, 100},    // 12 ten
}

// Scatters payment.
var ScatPay = [5]float64{0, 2, 3, 30, 300} // 2 cat

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 15, 20, 25} // 2 cat

const (
	cbn = 1 // cards bonus
)

type Game struct {
	slot.Slot5x3 `yaml:",inline"`
	// free spin number
	FS int `json:"fs,omitempty" yaml:"fs,omitempty" xml:"fs,omitempty"`
}

func NewGame() *Game {
	return &Game{
		Slot5x3: slot.Slot5x3{
			Sel: slot.MakeBitNum(20, 1),
			Bet: 1,
		},
		FS: 0,
	}
}

const wild, scat, bon = 1, 2, 3

var bl = []slot.Linex{
	{2, 2, 2, 2, 2}, // 1
	{1, 1, 1, 1, 1}, // 2
	{3, 3, 3, 3, 3}, // 3
	{1, 2, 3, 2, 1}, // 4
	{3, 2, 1, 2, 3}, // 5
	{2, 1, 1, 1, 2}, // 6
	{2, 3, 3, 3, 2}, // 7
	{1, 1, 2, 3, 3}, // 8
	{3, 3, 2, 1, 1}, // 9
	{2, 3, 2, 1, 2}, // 10
	{2, 1, 2, 3, 2}, // 11
	{1, 2, 2, 2, 1}, // 12
	{3, 2, 2, 2, 3}, // 13
	{1, 2, 1, 2, 1}, // 14
	{3, 2, 3, 2, 3}, // 15
	{2, 2, 1, 2, 2}, // 16
	{2, 2, 3, 2, 2}, // 17
	{1, 1, 3, 1, 1}, // 18
	{3, 3, 1, 3, 3}, // 19
	{1, 3, 3, 3, 1}, // 20
}

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	if g.FS == 0 {
		g.ScanLinedReg(screen, wins)
		g.ScanScattersReg(screen, wins)
	} else {
		g.ScanLinedBon(screen, wins)
		g.ScanScattersBon(screen, wins)
	}
}

// Lined symbols calculation on regular games.
func (g *Game) ScanLinedReg(screen slot.Screen, wins *slot.Wins) {
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
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyL(numl),
			})
		} else if payw > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payw,
				Mult: 1,
				Sym:  wild,
				Num:  numw,
				Line: li,
				XY:   line.CopyL(numw),
			})
		}
	}
}

// Lined symbols calculation on free spins.
func (g *Game) ScanLinedBon(screen slot.Screen, wins *slot.Wins) {
	for li := g.Sel.Next(0); li != -1; li = g.Sel.Next(li) {
		var line = bl[li-1]

		var numw, numl slot.Pos = 0, 5
		var syml slot.Sym
		var x slot.Pos
		var cw = true // continues wilds
		for x = 1; x <= 5; x++ {
			var sx = screen.Pos(x, line)
			if sx == wild {
				if cw && syml == 0 {
					numw = x
				}
			} else if sx == scat {
				cw = false
			} else if syml == 0 {
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
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyL(numl),
			})
		} else if payw > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payw,
				Mult: 1,
				Sym:  wild,
				Num:  numw,
				Line: li,
				XY:   line.CopyL(numw),
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScattersReg(screen slot.Screen, wins *slot.Wins) {
	if count := screen.ScatNum(scat); count >= 2 {
		var pay, fs = ScatPay[count-1], ScatFreespin[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel.Num()) * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
			Free: fs,
		})
	}
	if count := screen.ScatNum(bon); count >= 3 {
		*wins = append(*wins, slot.WinItem{
			Mult: 1,
			Sym:  bon,
			Num:  count,
			XY:   screen.ScatPos(bon),
			BID:  cbn,
		})
	}
}

// Scatters calculation.
func (g *Game) ScanScattersBon(screen slot.Screen, wins *slot.Wins) {
	if count := screen.ScatNum(bon); count >= 3 {
		*wins = append(*wins, slot.WinItem{
			Mult: 1,
			Sym:  bon,
			Num:  count,
			XY:   screen.ScatPos(bon),
			BID:  cbn,
		})
	}
}

func (g *Game) Spin(screen slot.Screen, mrtp float64) {
	var _, reels = FindReels(mrtp)
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
