package kenocenturion

import (
	"github.com/slotopol/server/game/keno"
)

// RTP[ 2] = 98.101266%
// RTP[ 3] = 97.127556%
// RTP[ 4] = 98.028554%
// RTP[ 5] = 98.034877%
// RTP[ 6] = 98.364294%
// RTP[ 7] = 98.139903%
// RTP[ 8] = 97.988575%
// RTP[ 9] = 98.085519%
// RTP[10] = 97.950344%
// RTP[game] = 97.980099%
var Paytable = keno.Paytable{
	{0},                                   //  0 sel
	{0, 0},                                //  1 sel
	{0, 1, 10},                            //  2 sel
	{0, 0, 5, 20},                         //  3 sel
	{0, 0, 2, 10, 40},                     //  4 sel
	{0, 0, 1, 5, 20, 75},                  //  5 sel
	{0, 0, 1, 3, 5, 40, 150},              //  6 sel
	{0, 0, 0, 3, 5, 15, 80, 300},          //  7 sel
	{0, 0, 0, 1, 5, 15, 24, 150, 500},     //  8 sel
	{0, 0, 0, 1, 3, 5, 30, 80, 300, 1000}, //  9 sel
	{0, 0, 0, 0, 3, 5, 15, 50, 180, 500, 2000}, // 10 sel
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
