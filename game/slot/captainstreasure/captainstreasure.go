package captainstreasure

// See: https://freeslotshub.com/playtech/captains-treasure/

import (
	slot "github.com/slotopol/server/game/slot"
)

// Lined payment.
var LinePay = [11][5]float64{
	{0, 0, 0, 0, 0},         //  1 wild
	{0, 0, 0, 0, 0},         //  2 scatter
	{2, 10, 100, 500, 5000}, //  3 sabers
	{0, 5, 50, 250, 2500},   //  4 map
	{0, 3, 20, 100, 1000},   //  5 anchor
	{0, 0, 10, 30, 500},     //  6 ace
	{0, 0, 5, 25, 300},      //  7 king
	{0, 0, 5, 20, 200},      //  8 queen
	{0, 0, 5, 20, 200},      //  9 jack
	{0, 0, 5, 15, 100},      // 10 ten
	{0, 0, 5, 15, 100},      // 11 nine
}

// Scatters payment.
var ScatPay = [5]float64{0, 1, 5, 10, 100} //  2 suitcase

type Game struct {
	slot.Slot5x3 `yaml:",inline"`
}

func NewGame() *Game {
	return &Game{
		Slot5x3: slot.Slot5x3{
			Sel: slot.MakeBitNum(9, 1),
			Bet: 1,
		},
	}
}

const wild, scat = 1, 2

var bl = slot.Lineset5x{
	{2, 2, 2, 2, 2}, // 1
	{1, 1, 1, 1, 1}, // 2
	{3, 3, 3, 3, 3}, // 3
	{1, 2, 3, 2, 1}, // 4
	{3, 2, 1, 2, 3}, // 5
	{1, 1, 2, 1, 1}, // 6
	{3, 3, 2, 3, 3}, // 7
	{2, 1, 1, 1, 2}, // 8
	{2, 3, 3, 3, 2}, // 9
}

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	g.ScanLined(screen, wins)
	g.ScanScatters(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen slot.Screen, wins *slot.Wins) {
	for li := range g.Sel.Bits() {
		var line = bl.Line(li)

		var syml, numl = screen.Pos(1, line), 1
		var mw float64 = 1 // mult wild
		for x := 2; x <= 5; x++ {
			var sx = screen.Pos(x, line)
			if sx == wild {
				mw = 2
			} else if sx != syml {
				break
			}
			numl++
		}

		if pay := LinePay[syml-1][numl-1]; pay > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * pay,
				Mult: mw,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyL(numl),
			})
		}

		if numl < 5 {
			var symr, numr = screen.Pos(5, line), 1
			var mw float64 = 1 // mult wild
			for x := 4; x > numl; x-- {
				var sx = screen.Pos(x, line)
				if sx == wild {
					mw = 2
				} else if sx != symr {
					break
				}
				numr++
			}

			if pay := LinePay[symr-1][numr-1]; pay > 0 {
				*wins = append(*wins, slot.WinItem{
					Pay:  g.Bet * pay,
					Mult: mw,
					Sym:  symr,
					Num:  numr,
					Line: li,
					XY:   line.CopyR(numr),
				})
			}
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen slot.Screen, wins *slot.Wins) {
	if count := screen.ScatNum(scat); count >= 2 {
		var pay = ScatPay[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel.Num()) * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
		})
	}
}

func (g *Game) Spin(screen slot.Screen, mrtp float64) {
	var _, reels = FindReels(mrtp)
	screen.Spin(reels)
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
