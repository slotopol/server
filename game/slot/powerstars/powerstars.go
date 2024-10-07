package powerstars

// See: https://freeslotshub.com/novomatic/power-stars/

import (
	"math/rand/v2"

	slot "github.com/slotopol/server/game/slot"
)

// Lined payment.
var LinePay = [9][5]float64{
	{0, 0, 100, 500, 1000}, // 1 seven
	{0, 0, 50, 200, 500},   // 2 bell
	{0, 0, 20, 50, 200},    // 3 melon
	{0, 0, 20, 50, 200},    // 4 grapes
	{0, 0, 10, 30, 150},    // 5 plum
	{0, 0, 10, 30, 150},    // 6 orange
	{0, 0, 10, 20, 100},    // 7 lemon
	{0, 0, 10, 20, 100},    // 8 cherry
	{0, 0, 0, 0, 0},        // 9 star
}

const (
	jid = 1 // jackpot ID
)

// Jackpot win combinations.
var Jackpot = [9][5]int{
	{0, 0, 0, 0, 0}, //  1 seven
	{0, 0, 0, 0, 0}, //  2 bell
	{0, 0, 0, 0, 0}, //  3 melon
	{0, 0, 0, 0, 0}, //  4 grapes
	{0, 0, 0, 0, 0}, //  5 plum
	{0, 0, 0, 0, 0}, //  6 orange
	{0, 0, 0, 0, 0}, //  7 lemon
	{0, 0, 0, 0, 0}, //  8 cherry
	{0, 0, 0, 0, 0}, // 9 star
}

// Bet lines
var bl = slot.BetLinesNvm10

type Game struct {
	slot.Slot5x3 `yaml:",inline"`
	PRW          [5]int `json:"prw" yaml:"prw" xml:"prw"` // pinned reel wild
}

func NewGame() *Game {
	return &Game{
		Slot5x3: slot.Slot5x3{
			Sel: slot.MakeBitNum(len(bl), 1),
			Bet: 1,
		},
	}
}

const wild = 9

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	g.ScanLined(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen slot.Screen, wins *slot.Wins) {
	var reelwild [5]bool
	var fs int
	var x, y slot.Pos
	for x = 2; x <= 4; x++ {
		if g.PRW[x-1] > 0 {
			reelwild[x-1] = true
		} else {
			for y = 1; y <= 3; y++ {
				if screen.At(x, y) == wild {
					reelwild[x-1] = true
					fs = 1
					break
				}
			}
		}
	}

	for li := g.Sel.Next(0); li != -1; li = g.Sel.Next(li) {
		var line = bl[li-1]
		var syml, symr slot.Sym
		var numl, numr slot.Pos
		var payl, payr float64

		syml, numl = screen.Pos(1, line), 1
		for x = 2; x <= 5; x++ {
			var sx = screen.Pos(x, line)
			if sx != syml && !reelwild[x-1] {
				break
			}
			numl++
		}
		payl = LinePay[syml-1][numl-1]

		if numl < 4 {
			symr, numr = screen.Pos(5, line), 1
			for x = 4; x >= 2; x-- {
				var sx = screen.Pos(x, line)
				if sx != symr && !reelwild[x-1] {
					break
				}
				numr++
			}
			payr = LinePay[symr-1][numr-1]
		}

		if payl > payr {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payl,
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyL(numl),
			})
		} else if payr > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payr,
				Mult: 1,
				Sym:  symr,
				Num:  numr,
				Line: li,
				XY:   line.CopyL(numr),
			})
		}
		if fs > 0 {
			*wins = append(*wins, slot.WinItem{
				Sym:  wild,
				Free: fs,
			})
		}
	}
}

func (g *Game) Spin(screen slot.Screen, mrtp float64) {
	screen.Spin(&Reels)
	if g.FreeSpins() == 0 {
		var _, wc = FindChance(mrtp) // wild chance
		var x, y slot.Pos
		for x = 2; x <= 4; x++ {
			if rand.Float64() < wc {
				y = rand.N[slot.Pos](3) + 1
				screen.Set(x, y, wild)
			}
		}
	}
}

func (g *Game) Apply(screen slot.Screen, wins slot.Wins) {
	if g.FreeSpins() > 0 {
		g.Gain += wins.Gain()
	} else {
		g.Gain = wins.Gain()
	}

	var x, y slot.Pos
	for x = 2; x <= 4; x++ {
		if g.PRW[x-1] > 0 {
			g.PRW[x-1]--
		} else {
			for y = 1; y <= 3; y++ {
				if screen.At(x, y) == wild {
					g.PRW[x-1] = 1
					break
				}
			}
		}
	}
}

func (g *Game) FreeSpins() int {
	return max(g.PRW[1], g.PRW[2], g.PRW[3])
}

func (g *Game) SetSel(sel slot.Bitset) error {
	return g.SetSelNum(sel, len(bl))
}
