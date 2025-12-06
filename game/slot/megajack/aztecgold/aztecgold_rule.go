package aztecgold

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsMap slot.ReelsMap[*slot.Reels5x]

var JackMap = slot.ReelsMap[float64]{}

// Lined payment.
var LinePay = [12][5]float64{
	{0, 0, 5, 10, 100},     //  1 tomat
	{0, 0, 5, 10, 100},     //  2 corn
	{0, 0, 5, 10, 100},     //  3 lama
	{0, 0, 5, 10, 100},     //  4 frog
	{0, 0, 10, 20, 100},    //  5 jaguar
	{0, 0, 20, 100, 500},   //  6 condor
	{0, 0, 20, 100, 750},   //  7 queen
	{0, 0, 20, 100, 1000},  //  8 king
	{0, 2, 25, 200, 10000}, //  9 dragon
	{},                     // 10 scatter
	{},                     // 11 idol
	{},                     // 12 pyramid
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 2, 10, 50} // 10 scatter

const (
	mje1 = 1 // Eldorado9
	mje3 = 2 // Eldorado9
	mje6 = 3 // Eldorado9
	mje9 = 4 // Eldorado9
	mjm  = 5 // Monopoly
	mjc  = 6 // Champagne
	mjap = 7 // AztecPyramid
)

// Bet lines
var BetLines = slot.BetLinesMgj[:]

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

const (
	mjj             = 1 // jackpot ID
	wild, scat, bon = 11, 10, 12
)

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	var reelwild [5]bool
	for x := 1; x < 4; x++ { // 2, 3, 4 reels only
		for y := range 3 {
			if g.Scr[x][y] == wild {
				reelwild[x] = true
				break
			}
		}
	}

	for li, line := range BetLines[:g.Sel] {
		var numl slot.Pos = 5
		var syml = g.LY(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = g.LY(x, line)
			if reelwild[x-1] {
				continue
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		if pay := LinePay[syml-1][numl-1]; pay > 0 {
			var jid int
			if numl == 5 {
				jid = mjj
			}
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * pay,
				MP:  1,
				Sym: syml,
				Num: numl,
				LI:  li + 1,
				XY:  line.HitxL(numl),
				JID: jid,
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if sn, wn := g.SymNum2(scat, wild); sn+wn >= 3 {
		var pay = ScatPay[sn+wn-1]
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * float64(g.Sel) * pay,
			MP:  1,
			Sym: scat,
			Num: sn + wn,
			XY:  g.SymPos2(scat, wild),
		})
	}
	if count := g.SymNum(bon); count >= 3 {
		*wins = append(*wins, slot.WinItem{
			MP:  1,
			Sym: bon,
			Num: count,
			XY:  g.SymPos(bon),
			BID: mjap,
		})
	}
}

func (g *Game) JackFreq(mrtp float64) []float64 {
	var bulk, _ = JackMap.FindClosest(mrtp)
	return []float64{bulk}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.ReelSpin(reels)
}

func (g *Game) Spawn(wins slot.Wins, fund, mrtp float64) {
	for i, wi := range wins {
		switch wi.BID {
		case mjap:
			wins[i].Bon, wins[i].Pay = AztecPyramidSpawn(g.Bet * float64(g.Sel))
		}
		if wi.JID != 0 {
			var bulk, _ = JackMap.FindClosest(mrtp)
			var jf = min(bulk*g.Bet/slot.JackBasis, 1)
			wins[i].JR = jf * fund
		}
	}
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
