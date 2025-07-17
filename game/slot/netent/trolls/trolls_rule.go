package trolls

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed trolls_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [14][5]float64{
	{0, 3, 25, 100, 750},      //  1 troll1
	{0, 0, 25, 100, 500},      //  2 troll2
	{0, 0, 15, 100, 500},      //  3 troll3
	{0, 0, 10, 75, 250},       //  4 troll4
	{0, 0, 10, 75, 250},       //  5 troll5
	{0, 0, 10, 50, 200},       //  6 troll6
	{0, 0, 5, 50, 150},        //  7 ace
	{0, 0, 5, 25, 125},        //  8 king
	{0, 0, 5, 25, 125},        //  9 queen
	{0, 0, 5, 25, 125},        // 10 jack
	{0, 2, 5, 25, 100},        // 11 ten
	{0, 10, 250, 2500, 10000}, // 12 wild
	{},                        // 13 golden
	{},                        // 14 scatter
}

// Scatters payment.
var ScatPay = [5]float64{0, 2, 5, 25, 500} // 14 scatter

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 10, 20, 30} // 14 scatter

// Bet lines
var BetLines = slot.BetLinesNetEnt5x3[:20]

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

const wild1, wild2, scat = 12, 13, 14

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
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
			var sx = g.LY(x, line)
			if sx == wild1 {
				if syml == 0 {
					numw = x
				}
				if mw < 4 {
					mw = 2
				}
			} else if sx == wild2 {
				if syml == 0 {
					numw = x
				}
				mw = 4
			} else if syml == 0 && sx != scat {
				syml = sx
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		var payw, payl float64
		if numw >= 2 {
			payw = LinePay[wild1-1][numw-1]
		}
		if numl >= 2 && syml > 0 {
			payl = LinePay[syml-1][numl-1]
		}
		if payl*mw > payw {
			var mm float64 = 1 // mult mode
			if g.FSR > 0 {
				mm = 3
			}
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payl,
				Mult: mw * mm,
				Sym:  syml,
				Num:  numl,
				Line: li + 1,
				XY:   line.CopyL(numl),
			})
		} else if payw > 0 {
			var mm float64 = 1 // mult mode
			if g.FSR > 0 {
				mm = 3
			}
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payw,
				Mult: mm, // no multiplyer on line by double symbol
				Sym:  wild1,
				Num:  numw,
				Line: li + 1,
				XY:   line.CopyL(numw),
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.ScatNum(scat); count >= 2 {
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
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	g.ReelSpin(reels)
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
