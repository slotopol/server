package kenoluxury

import (
	"github.com/slotopol/server/game/keno"
)

// RTP[ 2] = 92.088608%
// RTP[ 3] = 91.577410%
// RTP[ 4] = 92.027909%
// RTP[ 5] = 92.575146%
// RTP[ 6] = 92.048964%
// RTP[ 7] = 92.443234%
// RTP[ 8] = 91.949466%
// RTP[ 9] = 92.001419%
// RTP[10] = 92.228834%
// RTP[game] = 92.104554%
var Paytable = keno.Paytable{
	{0},                                      //  0 sel
	{0, 0},                                   //  1 sel
	{0, 1, 9},                                //  2 sel
	{0, 0, 2, 46},                            //  3 sel
	{0, 0, 2, 5, 91},                         //  4 sel
	{0, 0, 0, 3, 12, 820},                    //  5 sel
	{0, 0, 0, 3, 4, 68, 1600},                //  6 sel
	{0, 0, 0, 1, 2, 21, 400, 7000},           //  7 sel
	{0, 0, 0, 0, 2, 12, 100, 1600, 10000},    //  8 sel
	{0, 0, 0, 0, 1, 6, 44, 335, 4700, 10000}, //  9 sel
	{0, 0, 0, 0, 0, 5, 24, 140, 1000, 4500, 10000}, // 10 sel
}

type Game struct {
	keno.Keno80 `yaml:",inline"`
}

// Declare conformity with KenoGame interface.
var _ keno.KenoGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Keno80: keno.Keno80{
			Bet: 1,
		},
	}
}

func (g *Game) Scanner(wins *keno.Wins) error {
	return Paytable.Scanner(&g.Scr, wins, g.Bet)
}

func (g *Game) SetSel(sel keno.Bitset) error {
	return g.CheckSel(sel, &Paytable)
}
