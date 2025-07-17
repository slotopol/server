package flowers

// See: https://games.netent.com/video-slots/flowers/
// See: https://www.slotsmate.com/software/netent/flowers

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed flowers_bon.yaml
var rbon []byte

var ReelsBon = slot.ReadObj[*slot.Reels5x](rbon)

//go:embed flowers_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [17][10]float64{
	{0, 0, 250, 1000, 5000}, //  1 wild
	{},                      //  2 scatter
	{},                      //  3 scatter2
	{0, 0, 20, 40, 160, 250, 400, 600, 1000, 2000}, //  4 red
	{}, //  5 red2
	{0, 0, 15, 35, 140, 225, 350, 550, 900, 1800}, //  6 yellow
	{}, //  7 yellow2
	{0, 0, 15, 30, 120, 200, 300, 500, 800, 1600}, //  8 green
	{}, //  9 green2
	{0, 0, 10, 25, 100, 175, 250, 450, 700, 1400}, // 10 pink
	{}, // 11 pink2
	{0, 0, 10, 20, 80, 150, 200, 400, 600, 1200}, // 12 blue
	{},                 // 13 blue2
	{0, 0, 5, 20, 200}, // 14 ace
	{0, 0, 5, 20, 150}, // 15 king
	{0, 0, 5, 15, 125}, // 16 queen
	{0, 0, 5, 15, 100}, // 17 jack
}

var DoubleSym = [17]slot.Sym{
	0,  //  1 wild
	0,  //  2 scatter
	2,  //  3 scatter2
	0,  //  4 red
	4,  //  5 red2
	0,  //  6 yellow
	6,  //  7 yellow2
	0,  //  8 green
	8,  //  9 green2
	0,  // 10 pink
	10, // 11 pink2
	0,  // 12 blue
	12, // 13 blue2
	0,  // 14 ace
	0,  // 15 king
	0,  // 16 queen
	0,  // 17 jack
}

// Scatters payment.
var ScatPay = [10]float64{0, 0, 2, 2, 2, 2, 4, 10} // 2 scatter

// Scatter freespins table
var ScatFreespin = [10]int{0, 0, 0, 10, 15, 20, 25, 30} // 2 scatter

// Bet lines
var BetLines = []slot.Linex{
	{2, 2, 2, 2, 2}, // 1
	{1, 1, 1, 1, 1}, // 2
	{3, 3, 3, 3, 3}, // 3
	{1, 2, 3, 2, 1}, // 4
	{3, 2, 1, 2, 3}, // 5
	{1, 1, 2, 1, 1}, // 6
	{3, 3, 2, 3, 3}, // 7
	{2, 3, 3, 3, 2}, // 8
	{2, 1, 1, 1, 2}, // 9
	{2, 1, 2, 1, 2}, // 10
	{2, 3, 2, 3, 2}, // 11
	{1, 2, 1, 2, 1}, // 12
	{3, 2, 3, 2, 3}, // 13
	{2, 2, 1, 2, 2}, // 14
	{2, 2, 3, 2, 2}, // 15
	{1, 2, 2, 2, 1}, // 16
	{3, 2, 2, 2, 3}, // 17
	{1, 3, 1, 3, 1}, // 18
	{3, 1, 3, 1, 3}, // 19
	{1, 3, 3, 3, 1}, // 20
	{3, 1, 1, 1, 3}, // 21
	{1, 1, 3, 1, 1}, // 22
	{3, 3, 1, 3, 3}, // 23
	{1, 3, 2, 1, 3}, // 24
	{3, 1, 2, 3, 1}, // 25
	{2, 1, 3, 1, 2}, // 26
	{2, 3, 1, 3, 2}, // 27
	{1, 2, 3, 3, 3}, // 28
	{3, 2, 1, 1, 1}, // 29
	{2, 1, 2, 3, 2}, // 30
}

type Game struct {
	slot.Screen5x3 `yaml:",inline"`
	slot.Slotx     `yaml:",inline"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slotx: slot.Slotx{
			Sel: len(BetLines),
			Bet: 1,
		},
	}
}

func (g *Game) Clone() slot.SlotGame {
	var clone = *g
	return &clone
}

const wild, scat, scat2 = 1, 2, 3

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	for li, line := range BetLines[:g.Sel] {
		var numw, numl slot.Pos
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			if sx := g.LY(x, line); sx == wild {
				if syml == 0 {
					numw++
				}
			} else if sd := DoubleSym[sx-1]; syml == 0 {
				if sd > 0 {
					syml = sd
					numl++
				} else {
					syml = sx
				}
			} else if sd == syml {
				numl++
			} else if sx != syml {
				break
			}
			numl++
		}

		var payw, payl float64
		if numw >= 3 {
			payw = LinePay[wild-1][numw-1]
		}
		if numl >= 3 && syml > 0 {
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
				Line: li + 1,
				XY:   line.CopyL(x - 1),
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
				Line: li + 1,
				XY:   line.CopyL(numw),
			})
		}
	}
}

func (g *Game) ScatNumDbl() (n slot.Pos) {
	for x := range 5 {
		var r = g.Scr[x]
		if r[0] == scat || r[1] == scat || r[2] == scat {
			n++
		} else if r[0] == scat2 || r[1] == scat2 || r[2] == scat2 {
			n += 2
		}
	}
	return
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.ScatNumDbl(); count >= 3 {
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

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
