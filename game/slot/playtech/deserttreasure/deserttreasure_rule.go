package deserttreasure

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed deserttreasure_bon.yaml
var rbon []byte

var ReelsBon = slot.ReadObj[*slot.Reels5x](rbon)

//go:embed deserttreasure_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [10][5]float64{
	{0, 8, 80, 800, 8000}, //  1 wild
	{},                    //  2 scatter
	{0, 5, 40, 250, 1000}, //  3 shield
	{0, 0, 20, 75, 500},   //  4 swords
	{0, 0, 0, 50, 250},    //  5 lamp
	{0, 2, 10, 30, 150},   //  6 ligature1
	{0, 2, 10, 30, 150},   //  7 ligature2
	{0, 0, 5, 15, 75},     //  8 ligature3
	{0, 0, 5, 15, 75},     //  9 ligature4
	{0, 0, 0, 10, 50},     // 10 ligature5
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 4, 40, 400} // scatter

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 10, 25, 50} // scatter

// Bet lines
var BetLines = []slot.Linex{
	{2, 2, 2, 2, 2}, // 1
	{1, 1, 1, 1, 1}, // 2
	{3, 3, 3, 3, 3}, // 3
	{1, 2, 3, 2, 1}, // 4
	{3, 2, 1, 2, 3}, // 5
	{1, 1, 2, 3, 3}, // 6
	{3, 3, 2, 1, 1}, // 7
	{2, 1, 2, 3, 2}, // 8
	{2, 3, 2, 1, 2}, // 9
}

type Game struct {
	slot.Slot5x3 `yaml:",inline"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slot5x3: slot.Slot5x3{
			Sel: len(BetLines),
			Bet: 1,
		},
	}
}

func (g *Game) Clone() slot.SlotGame {
	var clone = *g
	return &clone
}

const wild, scat = 1, 2

func (g *Game) Scanner(wins *slot.Wins) {
	g.ScanLined(&g.Scrn, wins)
	g.ScanScatters(&g.Scrn, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen slot.Screen, wins *slot.Wins) {
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
			if g.FSR > 0 {
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
			if g.FSR > 0 {
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
		if g.FSR > 0 {
			mm = 3
		}
		var pay, fs = ScatPay[count-1], ScatFreespin[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * pay,
			Mult: mm,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
			Free: fs,
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	if g.FSR == 0 {
		var reels, _ = slot.FindReels(ReelsMap, mrtp)
		g.Scrn.Spin(reels)
	} else {
		g.Scrn.Spin(ReelsBon)
	}
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
