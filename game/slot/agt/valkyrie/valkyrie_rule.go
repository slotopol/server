package valkyrie

// See: https://demo.agtsoftware.com/games/agt/valkyrie

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed valkyrie_bon.yaml
var rbon []byte

var BonusReel = slot.ReadObj[[]slot.Sym](rbon)

//go:embed valkyrie_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [13][5]float64{
	{0, 2, 50, 500, 1000}, //  1 wild
	{},                    //  2 scatter
	{0, 2, 25, 250, 500},  //  3 warrior
	{0, 2, 25, 100, 200},  //  4 helmet
	{0, 0, 20, 100, 200},  //  5 shield
	{0, 0, 15, 50, 100},   //  6 axe
	{0, 0, 15, 50, 100},   //  7 mug
	{0, 0, 10, 25, 50},    //  8 ace
	{0, 0, 10, 25, 50},    //  9 king
	{0, 0, 5, 10, 25},     // 10 queen
	{0, 0, 5, 10, 25},     // 11 jack
	{0, 0, 5, 10, 25},     // 12 ten
	{0, 2, 5, 10, 25},     // 13 nine
}

// Bet lines
var BetLines = slot.BetLinesAgt5x3[:30]

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

const wild, scat = 1, 2

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
		if numw >= 2 {
			payw = LinePay[wild-1][numw-1]
		}
		if numl >= 2 && syml > 0 {
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
		*wins = append(*wins, slot.WinItem{
			Sym:  scat,
			Num:  count,
			XY:   g.ScatPos(scat),
			Free: 15,
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	if g.FSR == 0 {
		g.ReelSpin(reels)
	} else {
		g.Screen5x3.SpinBig(reels.Reel(1), BonusReel, reels.Reel(5))
	}
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
