package tikiwonders

// See: https://www.youtube.com/watch?v=ERbHFjYyu_E

import (
	slot "github.com/slotopol/server/game/slot"
)

// Lined payment.
var LinePay = [11][5]float64{
	{0, 5, 200, 2000, 10000}, //  1 wild
	{0, 0, 0, 0, 0},          //  2 scatter
	{0, 0, 50, 125, 1000},    //  3 dolphin
	{0, 0, 20, 75, 250},      //  4 guana
	{0, 0, 20, 75, 250},      //  5 turtle
	{0, 0, 15, 50, 200},      //  6 parrot
	{0, 0, 5, 25, 125},       //  7 ace
	{0, 0, 5, 25, 100},       //  8 king
	{0, 0, 3, 20, 100},       //  9 queen
	{0, 0, 2, 15, 75},        // 10 jack
	{0, 0, 2, 15, 75},        // 11 ten
}

// Scatters payment.
var ScatPay = [5]float64{0, 2, 4, 50, 400} // 2 scatter

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 10, 20, 30} // 2 scatter

type Game struct {
	slot.Slot5x3 `yaml:",inline"`
	// free spin number
	FS int `json:"fs,omitempty" yaml:"fs,omitempty" xml:"fs,omitempty"`
}

func NewGame() *Game {
	return &Game{
		Slot5x3: slot.Slot5x3{
			Sel: slot.MakeBitNum(30, 1),
			Bet: 1,
		},
		FS: 0,
	}
}

const wild, scat = 1, 2

var bl = slot.BetLinesBetSoft30

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	g.ScanLined(screen, wins)
	g.ScanScatters(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen slot.Screen, wins *slot.Wins) {
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
			var mm float64 = 1 // mult mode
			if g.FS > 0 {
				mm = 3
			}
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payl,
				Mult: mm,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyL(numl),
			})
		} else if payw > 0 {
			var mm float64 = 1 // mult mode
			if g.FS > 0 {
				mm = 3
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
		var mm float64 = 1 // mult mode
		if g.FS > 0 {
			mm = 3
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
