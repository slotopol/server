package fortuneteller

// See: https://freeslotshub.com/playngo/fortune-teller/
// See: https://www.slotsmate.com/software/play-n-go/play-n-go-fortune-teller
// See: https://www.youtube.com/watch?v=bFQq3cCz9XY

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed fortuneteller_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

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
var BetLines = []slot.Linex{
	{2, 2, 2, 2, 2}, // 1
	{1, 1, 1, 1, 1}, // 2
	{3, 3, 3, 3, 3}, // 3
	{1, 2, 3, 2, 1}, // 4
	{3, 2, 1, 2, 3}, // 5
	{2, 1, 1, 1, 2}, // 6
	{2, 3, 3, 3, 2}, // 7
	{1, 1, 2, 3, 3}, // 8
	{3, 3, 2, 1, 1}, // 9
	{2, 3, 2, 1, 2}, // 10
	{2, 1, 2, 3, 2}, // 11
	{1, 2, 2, 2, 1}, // 12
	{3, 2, 2, 2, 3}, // 13
	{1, 2, 1, 2, 1}, // 14
	{3, 2, 3, 2, 3}, // 15
	{2, 2, 1, 2, 2}, // 16
	{2, 2, 3, 2, 2}, // 17
	{1, 1, 3, 1, 1}, // 18
	{3, 3, 1, 3, 3}, // 19
	{1, 3, 3, 3, 1}, // 20
}

const (
	cbn = 1 // cards bonus
)

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
func (g *Game) ScanScattersReg(wins *slot.Wins) {
	if count := g.ScatNum(scat); count >= 2 {
		var pay, fs = ScatPay[count-1], ScatFreespin[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   g.ScatPos(scat),
			Free: fs,
		})
	}
	if count := g.ScatNum(bon); count >= 3 {
		*wins = append(*wins, slot.WinItem{
			Mult: 1,
			Sym:  bon,
			Num:  count,
			XY:   g.ScatPos(bon),
			BID:  cbn,
		})
	}
}

// Scatters calculation.
func (g *Game) ScanScattersBon(wins *slot.Wins) {
	if count := g.ScatNum(bon); count >= 3 {
		*wins = append(*wins, slot.WinItem{
			Mult: 1,
			Sym:  bon,
			Num:  count,
			XY:   g.ScatPos(bon),
			BID:  cbn,
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
