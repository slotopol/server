package atthemovies

import (
	slot "github.com/slotopol/server/game/slot"
)

// Lined payment.
var LinePay = [10][5]float64{
	{0, 20, 200, 500, 1000}, //  1 oscar
	{0, 10, 100, 250, 500},  //  2 popcorn
	{0, 5, 50, 100, 200},    //  3 poster
	{0, 2, 25, 50, 100},     //  4 a
	{0, 0, 20, 40, 80},      //  5 dummy
	{0, 0, 15, 30, 60},      //  6 maw
	{0, 0, 10, 20, 40},      //  7 starship
	{0, 0, 5, 10, 20},       //  8 heart
	{0, 0, 0, 0, 0},         //  9 masks
	{0, 0, 0, 0, 0},         // 10 projector
}

// Scatters payment.
var ScatPay = [5]float64{0, 2, 0, 0, 0} // 10 projector

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 8, 12, 20} // 10 projector

const (
	jid = 1 // jackpot ID
)

// Jackpot win combinations.
var Jackpot = [10][5]int{
	{0, 0, 0, 0, 0}, //  1 oscar
	{0, 0, 0, 0, 0}, //  2 popcorn
	{0, 0, 0, 0, 0}, //  3 poster
	{0, 0, 0, 0, 0}, //  4 a
	{0, 0, 0, 0, 0}, //  5 dummy
	{0, 0, 0, 0, 0}, //  6 maw
	{0, 0, 0, 0, 0}, //  7 starship
	{0, 0, 0, 0, 0}, //  8 heart
	{0, 0, 0, 0, 0}, //  9 masks
	{0, 0, 0, 0, 0}, // 10 projector
}

// Bet lines
var bl = slot.BetLinesBetSoft25

type Game struct {
	slot.Slot5x3 `yaml:",inline"`
	// free spin number
	FS int `json:"fs,omitempty" yaml:"fs,omitempty" xml:"fs,omitempty"`
}

func NewGame() *Game {
	return &Game{
		Slot5x3: slot.Slot5x3{
			Sel: slot.MakeBitNum(len(bl), 1),
			Bet: 1,
		},
		FS: 0,
	}
}

const wild, scat = 9, 10

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	g.ScanLined(screen, wins)
	g.ScanScatters(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen slot.Screen, wins *slot.Wins) {
	for li := g.Sel.Next(0); li != -1; li = g.Sel.Next(li) {
		var line = bl[li-1]

		var mw float64 = 1 // mult wild
		var numl slot.Pos = 1
		var syml = screen.Pos(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = screen.Pos(x, line)
			if sx == wild {
				mw = 2
			} else if sx != syml {
				break
			}
			numl++
		}

		if pay := LinePay[syml-1][numl-1]; pay > 0 {
			var mm float64 = 1 // mult mode
			if g.FS > 0 {
				mm = 2
			}
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * pay,
				Mult: mw * mm,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyL(numl),
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen slot.Screen, wins *slot.Wins) {
	if count := screen.ScatNum(scat); count >= 2 {
		var mm float64 = 1 // mult mode
		if g.FS > 0 {
			mm = 2
		}
		var pay, fs = ScatPay[count-1], ScatFreespin[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel.Num()) * pay,
			Mult: mm,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
			Free: fs,
		})
	}
}

func (g *Game) Spin(screen slot.Screen, mrtp float64) {
	var _, reels = slot.FindReels(ReelsMap, mrtp)
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
	return g.SetSelNum(sel, len(bl))
}
