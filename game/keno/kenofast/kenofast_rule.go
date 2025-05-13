package kenofast

// See: https://demo.agtsoftware.com/games/agt/keno

import (
	"github.com/slotopol/server/game/keno"
)

// RTP[ 1] = 95.000000%
// RTP[ 2] = 98.101266%
// RTP[ 3] = 95.740019%
// RTP[ 4] = 95.685960%
// RTP[ 5] = 95.123231%
// RTP[ 6] = 95.083714%
// RTP[ 7] = 95.540780%
// RTP[ 8] = 95.576691%
// RTP[ 9] = 95.086210%
// RTP[10] = 95.231797%
// RTP[game] = 95.616967%
var Paytable = keno.Paytable{
	{0},                                      //  0 sel
	{0, 3.8},                                 //  1 sel
	{0, 1, 10},                               //  2 sel
	{0, 0, 2, 49},                            //  3 sel
	{0, 0, 1, 8, 130},                        //  4 sel
	{0, 0, 1, 4, 20, 160},                    //  5 sel
	{0, 0, 0, 2, 15, 60, 600},                //  6 sel
	{0, 0, 0, 2, 5, 29, 95, 1000},            //  7 sel
	{1, 0, 0, 0, 5, 16, 50, 250, 2000},       //  8 sel
	{2, 0, 0, 0, 2, 10, 30, 100, 1000, 8000}, //  9 sel
	{2, 0, 0, 0, 0, 8, 20, 100, 330, 2000, 15000}, // 10 sel
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
