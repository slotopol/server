package firekeno

import (
	"github.com/slotopol/server/game/keno"
)

// RTP[ 2] = 92.088608%, sigma = 1.369910
// RTP[ 3] = 91.601753%, sigma = 3.015944
// RTP[ 4] = 92.009889%, sigma = 6.444535
// RTP[ 5] = 92.202418%, sigma = 16.977691
// RTP[ 6] = 91.995220%, sigma = 28.428062
// RTP[ 7] = 92.125355%, sigma = 49.481730
// RTP[ 8] = 91.999084%, sigma = 42.641870
// RTP[ 9] = 92.045254%, sigma = 43.076607
// RTP[10] = 92.053609%, sigma = 34.625466
// RTP[game] = 92.013465%
var Paytable = keno.Paytable{
	{0},                                //  0 sel
	{0, 0},                             //  1 sel
	{1, 0, 6},                          //  2 sel
	{1, 0, 1, 26},                      //  3 sel
	{1, 0, 0, 6, 115},                  //  4 sel
	{1, 0, 0, 1, 15, 666},              //  5 sel
	{2, 0, 0, 1, 2, 25, 2500},          //  6 sel
	{3, 0, 0, 0, 3, 10, 100, 10000},    //  7 sel
	{5, 0, 0, 0, 1, 4, 55, 666, 20000}, //  8 sel
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
	return Paytable.Scanner(&g.Grid, wins, g.Bet)
}

func (g *Game) SetSel(sel keno.Bitset) error {
	return g.CheckSel(sel, &Paytable)
}
