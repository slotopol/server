package firekeno

import (
	"github.com/slotopol/server/game/keno"
)

// RTP[ 2] = 92.088608%
// RTP[ 3] = 91.601753%
// RTP[ 4] = 91.739590%
// RTP[ 5] = 92.202418%
// RTP[ 6] = 91.995220%
// RTP[ 7] = 92.491393%
// RTP[ 8] = 92.041872%
// RTP[ 9] = 92.045254%
// RTP[10] = 92.053609%
// RTP[game] = 92.028857%
var Paytable = keno.Paytable{
	{0},                                //  0 sel
	{0, 0},                             //  1 sel
	{1, 0, 6},                          //  2 sel
	{1, 0, 1, 26},                      //  3 sel
	{1, 0, 0, 7, 100},                  //  4 sel
	{1, 0, 0, 1, 15, 666},              //  5 sel
	{2, 0, 0, 1, 2, 25, 2500},          //  6 sel
	{3, 0, 0, 0, 3, 10, 100, 10000},    //  7 sel
	{5, 0, 0, 0, 1, 4, 46, 666, 25000}, //  8 sel
	{10, 0, 0, 0, 0, 3, 10, 100, 1000, 50000},     //  9 sel
	{15, 0, 0, 0, 0, 1, 3, 25, 666, 1000, 100000}, // 10 sel
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
