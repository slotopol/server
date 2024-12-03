package groovysixties

// See: https://www.youtube.com/watch?v=qINQD6wRhpY

import (
	"github.com/slotopol/server/game/slot"
)

// Lined payment.
var LinePay = [12][5]float64{
	{0, 0, 50, 200, 1000}, //  1 wild
	{},                    //  2 scatter
	{0, 0, 25, 100, 400},  //  3 car
	{0, 0, 25, 100, 400},  //  4 tv
	{0, 0, 20, 75, 250},   //  5 recorder
	{0, 0, 20, 75, 250},   //  6 projector
	{0, 0, 5, 50, 150},    //  7 boots
	{0, 0, 5, 50, 150},    //  8 column
	{0, 0, 5, 20, 100},    //  9 ace
	{0, 0, 5, 20, 100},    // 10 king
	{0, 0, 5, 20, 100},    // 11 queen
	{0, 0, 5, 20, 100},    // 12 jack
}

// Bet lines
var BetLines = slot.BetLinesNetEnt5x4[:40]

type Game struct {
	slot.Slot5x4 `yaml:",inline"`
	// free spin number
	FS int `json:"fs,omitempty" yaml:"fs,omitempty" xml:"fs,omitempty"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slot5x4: slot.Slot5x4{
			Sel: len(BetLines),
			Bet: 1,
		},
		FS: 0,
	}
}

const wild, scat = 1, 2

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	var scrn5x4 = screen.(*slot.Screen5x4)
	g.ScanLined(scrn5x4, wins)
	g.ScanScatters(scrn5x4, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen *slot.Screen5x4, wins *slot.Wins) {
	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]

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
				mm = 2
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
				mm = 2
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
func (g *Game) ScanScatters(screen *slot.Screen5x4, wins *slot.Wins) {
	if count := screen.ScatNum(scat); count >= 3 {
		var mm float64 = 1 // mult mode
		if g.FS > 0 {
			mm = 2
		}
		var pay float64 = 2
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * pay,
			Mult: mm,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
			Free: 5,
		})
	}
}

func (g *Game) Spin(screen slot.Screen, mrtp float64) {
	var reels, _ = slot.FindReels(ReelsMap, mrtp)
	screen.Spin(reels)
}

func (g *Game) Apply(screen slot.Screen, wins slot.Wins) {
	if g.FS != 0 {
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

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}