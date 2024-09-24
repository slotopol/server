package greatblue

// See: https://freeslotshub.com/playtech/great-blue/

import (
	"math/rand/v2"

	slot "github.com/slotopol/server/game/slot"
)

// Lined payment.
var LinePay = [13][5]float64{
	{0, 10, 250, 2500, 10000}, //  1 wild
	{0, 2, 25, 125, 750},      //  2 dolphin
	{0, 2, 25, 125, 750},      //  3 turtle
	{0, 0, 20, 100, 400},      //  4 fish
	{0, 0, 15, 75, 250},       //  5 seahorse
	{0, 0, 15, 75, 250},       //  6 starfish
	{0, 0, 10, 50, 150},       //  7 ace
	{0, 0, 10, 50, 150},       //  8 king
	{0, 0, 5, 25, 100},        //  9 queen
	{0, 0, 5, 25, 100},        // 10 jack
	{0, 0, 5, 25, 100},        // 11 ten
	{0, 2, 5, 25, 100},        // 12 nine
	{0, 0, 0, 0, 0},           // 13 scatter
}

// Scatters payment.
var ScatPay = [5]float64{0, 2, 5, 20, 500} // 13 scatter

type Seashells struct {
	Sel1 string  `json:"sel1" yaml:"sel1" xml:"sel1"`
	Sel2 string  `json:"sel2" yaml:"sel2" xml:"sel2"`
	Mult float64 `json:"mult" yaml:"mult" xml:"mult"`
	Free int     `json:"free" yaml:"free" xml:"free"`
}

func (s *Seashells) SetupShell(shell string) {
	switch shell {
	case "x5":
		s.Mult += 5
	case "x8":
		s.Mult += 8
	case "7":
		s.Free += 7
	case "10":
		s.Free += 10
	case "15":
		s.Free += 15
	}
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
			Sel: slot.MakeBitNum(25, 1),
			Bet: 1,
		},
		FS: 0,
		M:  0,
	}
}

const wild, scat = 1, 13

var bl = slot.BetLinesPlt30

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
		var mw float64 = 1 // mult wild
		for x = 1; x <= 5; x++ {
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
		var mm float64 = 1 // mult mode
		if g.FS > 0 {
			mm, fs = g.M, 15
		} else if count >= 3 {
			fs = 8
		}
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

func (g *Game) Spawn(screen slot.Screen, wins slot.Wins) {
	if g.FS > 0 {
		return
	}
	for i := range wins {
		if wi := &wins[i]; wi.Sym == scat {
			var idx = []string{"x5", "x8", "7", "10", "15"}
			rand.Shuffle(len(idx), func(i, j int) {
				idx[i], idx[j] = idx[j], idx[i]
			})
			var bon = Seashells{
				Sel1: idx[0],
				Sel2: idx[1],
				Mult: 2,
				Free: 8,
			}
			bon.SetupShell(idx[0])
			bon.SetupShell(idx[1])
			wi.Mult = 1
			wi.Free = bon.Free
			wi.Bon = bon
		}
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
		if wi.Sym == scat {
			if g.FS > 0 {
				g.FS += wi.Free
			} else {
				var bon = wi.Bon.(Seashells)
				g.FS = bon.Free
				g.M = bon.Mult
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
