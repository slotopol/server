package flamedancer

// See: https://casino.ru/flame-dancer-novomatic/

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed flamedancer_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [12][5]float64{
	{},                     //  1 wild
	{0, 0, 100, 400, 2000}, //  2 volcano
	{0, 0, 40, 100, 750},   //  3 drums
	{0, 0, 25, 75, 400},    //  4 guitar
	{0, 0, 25, 75, 400},    //  5 coconut
	{0, 0, 10, 40, 150},    //  6 ace
	{0, 0, 10, 25, 125},    //  7 king
	{0, 0, 10, 25, 125},    //  8 queen
	{0, 0, 5, 20, 100},     //  9 jack
	{0, 0, 5, 20, 100},     // 10 ten
	{},                     // 11 scatter
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 2, 10, 25} // 11 scatter

// Bet lines
var BetLines = slot.BetLinesNvm20v2

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

const wild, scat = 1, 11
const wilds = 5

func (g *Game) Scanner(wins *slot.Wins) error {
	if g.FSR == 0 {
		g.ScanLinedReg(wins)
	} else {
		g.ScanLinedBon(wins)
	}
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation on regular games.
func (g *Game) ScanLinedReg(wins *slot.Wins) {
	var reelwild [5]bool
	for x := 1; x < 4; x += 2 { // 2, 4 reel only
		for y := 0; y < 3; y++ {
			if g.Scr[x][y] == wild {
				reelwild[x] = true
				break
			}
		}
	}

	for li, line := range BetLines[:g.Sel] {
		var numl slot.Pos = 5
		var syml = g.LY(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = g.LY(x, line)
			if reelwild[x-1] {
				continue
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		if numl >= 3 && syml != scat {
			var pay = LinePay[syml-1][numl-1]
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * pay,
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li + 1,
				XY:   line.CopyL(numl),
			})
		}
	}
}

// Lined symbols calculation on bonus games.
func (g *Game) ScanLinedBon(wins *slot.Wins) {
	var reelwild [5]bool
	for x := 1; x < 4; x += 2 { // 2, 4 reel only
		for y := 0; y < 3; y++ {
			if g.Scr[x][y] == wild {
				reelwild[x] = true
				break
			}
		}
	}

	for li, line := range BetLines[:g.Sel] {
		var numl slot.Pos = 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = g.LY(x, line)
			if sx <= wilds || reelwild[x-1] {
				continue
			} else if syml == 0 {
				syml = sx
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		if numl >= 3 && syml > 0 && syml != scat {
			var pay = LinePay[syml-1][numl-1]
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * pay,
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li + 1,
				XY:   line.CopyL(numl),
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.ScatNum(scat); count >= 3 {
		var pay = ScatPay[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   g.ScatPos(scat),
			Free: 7,
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
