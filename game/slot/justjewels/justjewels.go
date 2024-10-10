package justjewels

import (
	slot "github.com/slotopol/server/game/slot"
)

// Lined payment.
var LinePay = [8][5]float64{
	{0, 0, 50, 500, 5000}, // 1 crown
	{0, 0, 30, 150, 500},  // 2 gold
	{0, 0, 30, 150, 500},  // 3 money
	{0, 0, 15, 50, 200},   // 4 ruby
	{0, 0, 15, 50, 200},   // 5 sapphire
	{0, 0, 10, 25, 150},   // 6 emerald
	{0, 0, 10, 25, 150},   // 7 amethyst
	{0, 0, 0, 0, 0},       // 8 euro
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 2, 10, 50} // 8 euro

const (
	jid = 1 // jackpot ID
)

// Jackpot win combinations.
var Jackpot = [8][5]int{
	{0, 0, 0, 0, 0}, //  1 crown
	{0, 0, 0, 0, 0}, //  2 gold
	{0, 0, 0, 0, 0}, //  3 money
	{0, 0, 0, 0, 0}, //  4 ruby
	{0, 0, 0, 0, 0}, //  5 sapphire
	{0, 0, 0, 0, 0}, //  6 emerald
	{0, 0, 0, 0, 0}, //  7 amethyst
	{0, 0, 0, 0, 0}, //  8 euro
}

// Bet lines
var bl = slot.BetLinesNvm10

type Game struct {
	slot.Slot5x3 `yaml:",inline"`
}

func NewGame() *Game {
	return &Game{
		Slot5x3: slot.Slot5x3{
			Sel: slot.MakeBitNum(5, 1),
			Bet: 1,
		},
	}
}

const scat = 8

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	g.ScanLined(screen, wins)
	g.ScanScatters(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen slot.Screen, wins *slot.Wins) {
	for li := g.Sel.Next(0); li != -1; li = g.Sel.Next(li) {
		var line = bl[li-1]

		var numl slot.Pos = 1
		var syml = screen.Pos(3, line)
		var xy slot.Linex
		xy.Set(3, line.At(3))
		if screen.Pos(2, line) == syml {
			xy.Set(2, line.At(2))
			numl++
			if screen.Pos(1, line) == syml {
				xy.Set(1, line.At(1))
				numl++
			}
		}
		if screen.Pos(4, line) == syml {
			xy.Set(4, line.At(4))
			numl++
			if screen.Pos(5, line) == syml {
				xy.Set(5, line.At(5))
				numl++
			}
		}

		if numl >= 3 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * LinePay[syml-1][numl-1],
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   xy,
			})
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
		})
	}
}

func (g *Game) Spin(screen slot.Screen, mrtp float64) {
	var reels, _ = slot.FindReels(ReelsMap, mrtp)
	screen.Spin(reels)
}

func (g *Game) SetSel(sel slot.Bitset) error {
	return g.SetSelNum(sel, len(bl))
}
