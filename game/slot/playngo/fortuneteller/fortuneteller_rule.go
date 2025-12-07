package fortuneteller

// See: https://freeslotshub.com/playngo/fortune-teller/
// See: https://www.slotsmate.com/software/play-n-go/play-n-go-fortune-teller
// See: https://www.youtube.com/watch?v=bFQq3cCz9XY

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsMap slot.ReelsMap[*slot.Reels5x]

// Lined payment.
var LinePay = [12][5]float64{
	{0, 0, 50, 500, 5000}, //  1 wild
	{},                    //  2 cat
	{},                    //  3 bonus
	{0, 0, 25, 250, 1000}, //  4 girl
	{0, 0, 15, 100, 500},  //  5 hand
	{0, 0, 15, 100, 500},  //  6 zodiac
	{0, 0, 15, 75, 250},   //  7 candle
	{0, 0, 10, 50, 150},   //  8 ace
	{0, 0, 10, 50, 150},   //  9 king
	{0, 0, 5, 25, 100},    // 10 queen
	{0, 0, 5, 25, 100},    // 11 jack
	{0, 0, 5, 25, 100},    // 12 ten
}

// Scatters payment.
var ScatPay = [5]float64{0, 2, 3, 30, 300} // 2 cat

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 15, 20, 25} // 2 cat

// Bet lines
var BetLines = slot.BetLinesPlt5x3[:]

const (
	cbn = 1 // cards bonus
)

type Game struct {
	slot.Screen5x3 `yaml:",inline"`
	slot.Slotx     `yaml:",inline"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame(sel int) *Game {
	return &Game{
		Slotx: slot.Slotx{
			Sel: sel,
			Bet: 1,
		},
	}
}

func (g *Game) Clone() slot.SlotGame {
	var clone = *g
	return &clone
}

const wild, scat, bon = 1, 2, 3

func (g *Game) Scanner(wins *slot.Wins) error {
	if g.FSR == 0 {
		g.ScanLinedReg(wins)
		g.ScanScattersReg(wins)
	} else {
		g.ScanLinedBon(wins)
		g.ScanScattersBon(wins)
	}
	return nil
}

// Lined symbols calculation on regular games.
func (g *Game) ScanLinedReg(wins *slot.Wins) {
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
				Pay: g.Bet * payl,
				MP:  1,
				Sym: syml,
				Num: numl,
				LI:  li + 1,
				XY:  line.HitxL(numl),
			})
		} else if payw > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * payw,
				MP:  1,
				Sym: wild,
				Num: numw,
				LI:  li + 1,
				XY:  line.HitxL(numw),
			})
		}
	}
}

// Lined symbols calculation on free spins.
func (g *Game) ScanLinedBon(wins *slot.Wins) {
	for li, line := range BetLines[:g.Sel] {
		var numw, numl slot.Pos = 0, 5
		var syml slot.Sym
		var x slot.Pos
		var cw = true // continues wilds
		for x = 1; x <= 5; x++ {
			var sx = g.LY(x, line)
			if sx == wild {
				if cw && syml == 0 {
					numw = x
				}
			} else if sx == scat {
				cw = false
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
				Pay: g.Bet * payl,
				MP:  1,
				Sym: syml,
				Num: numl,
				LI:  li + 1,
				XY:  line.HitxL(numl),
			})
		} else if payw > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * payw,
				MP:  1,
				Sym: wild,
				Num: numw,
				LI:  li + 1,
				XY:  line.HitxL(numw),
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScattersReg(wins *slot.Wins) {
	if count := g.SymNum(scat); count >= 2 {
		var pay, fs = ScatPay[count-1], ScatFreespin[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * float64(g.Sel) * pay,
			MP:  1,
			Sym: scat,
			Num: count,
			XY:  g.SymPos(scat),
			FS:  fs,
		})
	}
	if count := g.SymNum(bon); count >= 3 {
		*wins = append(*wins, slot.WinItem{
			MP:  1,
			Sym: bon,
			Num: count,
			XY:  g.SymPos(bon),
			BID: cbn,
		})
	}
}

// Scatters calculation.
func (g *Game) ScanScattersBon(wins *slot.Wins) {
	if count := g.SymNum(bon); count >= 3 {
		*wins = append(*wins, slot.WinItem{
			MP:  1,
			Sym: bon,
			Num: count,
			XY:  g.SymPos(bon),
			BID: cbn,
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.ReelSpin(reels)
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
