package aislot

// See: https://demo.agtsoftware.com/games/agt/aislot

import (
	slot "github.com/slotopol/server/game/slot"
)

// Lined payment.
var LinePay = [10][5]float64{
	{},                       //  1 scatter
	{0, 10, 100, 1000, 5000}, //  2 man
	{0, 5, 40, 400, 2000},    //  3 mind
	{0, 5, 30, 100, 750},     //  4 internet
	{0, 5, 30, 100, 750},     //  5 eye
	{0, 0, 5, 40, 150},       //  6 ace
	{0, 0, 5, 40, 150},       //  7 king
	{0, 0, 5, 25, 100},       //  8 queen
	{0, 0, 5, 25, 100},       //  9 jack
	{0, 0, 5, 25, 100},       // 10 ten
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 2, 25, 250} // 1 scatter

// Bet lines
var BetLines = slot.BetLinesAgt15

type Game struct {
	slot.Slot5x3 `yaml:",inline"`
	// free spin number
	FS int `json:"fs,omitempty" yaml:"fs,omitempty" xml:"fs,omitempty"`
}

func NewGame() *Game {
	return &Game{
		Slot5x3: slot.Slot5x3{
			Sel: slot.MakeBitNum(len(BetLines), 1),
			Bet: 1,
		},
		FS: 0,
	}
}

const wild, scat = 1, 1

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	g.ScanLined(screen, wins)
	g.ScanScatters(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen slot.Screen, wins *slot.Wins) {
	for li := g.Sel.Next(0); li != -1; li = g.Sel.Next(li) {
		var line = BetLines[li-1]

		var numl slot.Pos = 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = screen.Pos(x, line)
			if sx == wild {
				continue
			} else if syml == 0 {
				syml = sx
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		if syml > 0 {
			if pay := LinePay[syml-1][numl-1]; pay > 0 {
				*wins = append(*wins, slot.WinItem{
					Pay:  g.Bet * pay,
					Mult: 1,
					Sym:  syml,
					Num:  numl,
					Line: li,
					XY:   line.CopyL(numl),
				})
			}
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen slot.Screen, wins *slot.Wins) {
	if count := screen.ScatNum(scat); count >= 3 {
		var pay = ScatPay[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel.Num()) * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
			Free: 12,
		})
	}
}

func (g *Game) Spin(screen slot.Screen, mrtp float64) {
	var reels, _ = slot.FindReels(ReelsMap, mrtp)
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
	return g.SetSelNum(sel, len(BetLines))
}
