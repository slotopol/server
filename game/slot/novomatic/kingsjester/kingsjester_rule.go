package kingsjester

// See: https://www.slotsmate.com/software/novomatic/kings-jester

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsMap slot.ReelsMap[slot.Reelx]

var JackMap slot.ReelsMap[[2]float64]

// Lined payment.
var LinePay = [14][5]float64{
	{0, 0, 20, 100, 1000}, //  1 double
	{0, 0, 10, 50, 500},   //  2 jester
	{0, 2, 20, 100, 500},  //  3 funny king
	{0, 2, 20, 100, 500},  //  4 funny queen
	{0, 0, 10, 75, 350},   //  5 cards
	{0, 0, 10, 50, 250},   //  6 bandura
	{0, 0, 10, 50, 250},   //  7 pan flute
	{0, 0, 5, 35, 125},    //  8 ace
	{0, 0, 5, 35, 125},    //  9 king
	{0, 0, 5, 25, 100},    // 10 queen
	{0, 0, 5, 25, 100},    // 11 jack
	{0, 0, 5, 25, 100},    // 12 ten
	{0, 2, 5, 25, 100},    // 13 nine
	{},                    // 14 scatter
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 5, 50, 500} // 12 scatter

// Bet lines
var BetLines = slot.BetLinesNvm10[:]

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

const (
	kjj1, kjj2       = 1, 2     // jackpot ID
	jest, wild, scat = 1, 2, 14 // symbols
)

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	for li, line := range BetLines[:g.Sel] {
		var mw float64 = 1 // mult wild
		var numj, numw, numl slot.Pos = 0, 0, 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = g.LY(x, line)
			if sx == jest {
				if syml == 0 {
					numw = x
				}
				numj++
				mw = 2
			} else if sx == wild {
				if syml == 0 {
					numw = x
				}
				numj = 0
			} else if syml == 0 {
				syml = sx
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		var payj, payw, payl float64
		if numj >= 2 && numj == numw {
			payj = LinePay[jest-1][numj-1]
		}
		if numw >= 2 {
			payw = LinePay[wild-1][numw-1]
		}
		if numl >= 2 && syml > 0 {
			payl = LinePay[syml-1][numl-1]
		}
		if payl*mw > payj && payl*mw > payw {
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * payl,
				MP:  mw,
				Sym: syml,
				Num: numl,
				LI:  li + 1,
				XY:  line.HitxL(numl),
			})
		} else if payw > payj {
			var jid int
			if numw == 5 {
				jid = kjj2
			}
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * payw,
				MP:  1,
				Sym: wild,
				Num: numw,
				LI:  li + 1,
				XY:  line.HitxL(numw),
				JID: jid,
			})
		} else if payj > 0 {
			var jid int
			if numw == 5 {
				jid = kjj1
			}
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * payj,
				MP:  1,
				Sym: jest,
				Num: numj,
				LI:  li + 1,
				XY:  line.HitxL(numj),
				JID: jid,
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.SymNum(scat); count >= 3 {
		var pay = ScatPay[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * float64(g.Sel) * pay,
			MP:  1,
			Sym: scat,
			Num: count,
			XY:  g.SymPos(scat),
			FS:  15,
		})
	}
}

func (g *Game) JackFreq(mrtp float64) []float64 {
	var bulk, _ = JackMap.FindClosest(mrtp)
	return bulk[:]
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.ReelSpin(reels)
}

func (g *Game) Spawn(wins slot.Wins, fund, mrtp float64) {
	for i, wi := range wins {
		if wi.JID != 0 {
			var bulk, _ = JackMap.FindClosest(mrtp)
			var jf = min(bulk[wi.JID-1]*g.Bet/slot.JackBasis, 1)
			wins[i].JR = jf * fund
		}
	}
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
