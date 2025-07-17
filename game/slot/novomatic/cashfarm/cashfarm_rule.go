package cashfarm

// See: https://casino.ru/cash-farm-novomatic/

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed cashfarm_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [13][5]float64{
	{0, 0, 100, 750, 5000}, //  1 wild
	{0, 0, 75, 300, 1000},  //  2 cow
	{0, 0, 50, 150, 500},   //  3 ram
	{0, 0, 50, 150, 500},   //  4 pig
	{0, 0, 30, 100, 300},   //  5 rabbit
	{0, 0, 30, 100, 300},   //  6 rooster
	{0, 0, 15, 75, 200},    //  7 ace
	{0, 0, 15, 75, 200},    //  8 king
	{0, 0, 10, 50, 150},    //  9 queen
	{0, 0, 10, 50, 150},    // 10 jack
	{0, 0, 5, 25, 75},      // 11 ten
	{0, 0, 5, 25, 75},      // 12 nine
	{},                     // 13 tractor
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 5, 20, 200} // 13 tractor

// Bet lines
var BetLines = slot.BetLinesNvm25

type Game struct {
	slot.Cascade5x3 `yaml:",inline"`
	slot.Slotx      `yaml:",inline"`
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

func (g *Game) Free() bool {
	return g.FSR != 0 || g.Cascade()
}

const wild, scat = 1, 13
const farmbn = 1

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	for li, line := range BetLines[:g.Sel] {
		var numw, numl slot.Pos = 0, 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = g.LY(x, line)
			if sx == wild {
				if syml == 0 {
					numw = x
				}
			} else if syml == 0 {
				syml = sx
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		var payw, payl float64
		if numw >= 3 {
			payw = LinePay[wild-1][numw-1]
		}
		if numl >= 3 && syml > 0 {
			payl = LinePay[syml-1][numl-1]
		}
		if payl > payw {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payl,
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li + 1,
				XY:   line.CopyL(numl),
			})
		} else if payw > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payw,
				Mult: 1,
				Sym:  wild,
				Num:  numw,
				Line: li + 1,
				XY:   line.CopyL(numw),
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
		})
		*wins = append(*wins, slot.WinItem{
			Sym: scat,
			Num: count,
			XY:  g.ScatPos(scat),
			BID: farmbn,
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	g.ReelSpin(reels)
}

func (g *Game) Spawn(wins slot.Wins, fund, mrtp float64) {
	for i, wi := range wins {
		switch wi.BID {
		case farmbn:
			wins[i].Bon, wins[i].Pay, wins[i].Mult = FarmSpawn(g.Bet)
		}
	}
}

func (g *Game) Prepare() {
	g.NewFall()
}

func (g *Game) Apply(wins slot.Wins) {
	g.Slotx.Apply(wins)
	g.Strike(wins)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
