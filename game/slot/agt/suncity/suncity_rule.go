package suncity

// See: https://demo.agtsoftware.com/games/agt/suncity

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed suncity_bon.yaml
var rbon []byte

var ReelsBon = slot.ReadObj[*slot.Reels5x](rbon)

//go:embed suncity_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [13][5]float64{
	{0, 10, 250, 2500, 9000}, //  1 wild
	{},                       //  2 scatter
	{0, 2, 25, 125, 750},     //  3 yacht
	{0, 2, 25, 125, 750},     //  4 bike
	{0, 0, 20, 100, 400},     //  5 lady
	{0, 0, 15, 75, 250},      //  6 salesgirl
	{0, 0, 15, 75, 250},      //  7 courier
	{0, 0, 10, 50, 125},      //  8 ace
	{0, 0, 10, 50, 125},      //  9 king
	{0, 0, 5, 25, 100},       // 10 queen
	{0, 0, 5, 25, 100},       // 11 jack
	{0, 0, 5, 25, 100},       // 12 ten
	{0, 2, 5, 25, 100},       // 13 nine
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
	if count := g.ScatNum(scat); count == 2 {
		*wins = append(*wins, slot.WinItem{
			Sym:  scat,
			Num:  count,
			XY:   g.ScatPos(scat),
			Free: -1,
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

func (g *Game) Apply(wins slot.Wins) {
	if g.FSR != 0 {
		g.Gain += wins.Gain()
		g.FSN++
	} else {
		g.Gain = wins.Gain()
		g.FSN = 0
	}

	for _, wi := range wins {
		if wi.Free != 0 {
			if g.FSR != 0 {
				g.FSR = 0 // stop free games
			} else {
				g.FSR = -1 // start free games
			}
		}
	}
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
