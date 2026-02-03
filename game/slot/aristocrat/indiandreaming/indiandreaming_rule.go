package indiandreaming

// See: https://freeslotshub.com/aristocrat/indian-dreaming/

import (
	"github.com/slotopol/server/game/slot"
)

const (
	sn         = 12   // number of symbols
	wild, scat = 1, 2 // wild & scatter symbol IDs
	linemin    = 3    // minimum line symbols to win
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [sn][5]float64{
	{},                     //  1 wild
	{},                     //  2 scatter
	{0, 0, 100, 200, 5000}, //  3 catcher
	{0, 0, 50, 100, 2500},  //  4 man
	{0, 0, 50, 100, 1000},  //  5 woman
	{0, 0, 10, 40, 250},    //  6 guy
	{0, 0, 6, 25, 150},     //  7 bull
	{0, 0, 6, 25, 150},     //  8 hatchet
	{0, 0, 6, 15, 80},      //  9 ace
	{0, 0, 4, 10, 80},      // 10 king
	{0, 0, 3, 10, 70},      // 11 queen
	{0, 0, 3, 10, 60},      // 12 jack
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 2, 15, 100} //  2 scatter

type Game struct {
	slot.Grid5x3 `yaml:",inline"`
	slot.Slotx     `yaml:",inline"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slotx: slot.Slotx{
			Bet: 1,
		},
	}
}

func (g *Game) Clone() slot.SlotGame {
	var clone = *g
	return &clone
}

func (g *Game) Scanner(wins *slot.Wins) error {
	// Count symbols
	var counts [5 + 1][sn + 1]int // symbol counts per reel
	for x := range 5 {
		var cx = &counts[x]
		for _, sym := range g.Grid[x] {
			cx[sym]++
		}
	}
	var nw = counts[0][wild] + counts[1][wild] + counts[2][wild] + counts[3][wild] + counts[4][wild]
	// Ways calculation
	var mm = 1 // mult mode
	if g.FSR > 0 && nw > 0 {
		mm = 5
	}
	if nw < 5 {
		var combs1 = [sn + 1]int{0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1} // pure symbols
		var combs2 = [sn + 1]int{0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1} // symbols + wilds
		for x, cx := range counts {
			var cw = combs1[wild]
			for sym := range combs2 {
				var c1, c2 = combs1[sym], combs2[sym]
				var n = cx[sym] + cx[wild]
				combs1[sym] = c1 * cx[sym]
				combs2[sym] = c2 * n
				var c = (c2-c1-cw)*mm + c1
				if x >= linemin && c > 0 && n == 0 {
					var pay = LinePay[sym-1][x-1]
					*wins = append(*wins, slot.WinItem{
						Pay: g.Bet * pay,
						MP:  float64(c),
						Sym: slot.Sym(sym),
						Num: slot.Pos(x),
						LI:  243,
						XY:  g.SymPosL2(slot.Pos(x), slot.Sym(sym), wild),
					})
				}
			}
		}
		/*
			// this algorithm has lower performance due to cache misses
			var sym slot.Sym
			for sym = 3; sym <= sn; sym++ { // ignore wild & scatter
				var cw = 1 // current ways on wilds only
				var c1 = 1 // current ways on pure symbols
				var c2 = 1 // current ways on symbols + wilds
				for x, cx := range counts {
					var n1, nw = cx[sym], cx[wild]
					var n2 = n1 + nw
					if n2 == 0 {
						if x >= linemin && c2 > cw {
							var pay = LinePay[sym-1][x-1]
							*wins = append(*wins, slot.WinItem{
								Pay: g.Bet * pay,
								MP:  float64(c1 + (c2-c1-cw)*mm),
								Sym: sym,
								Num: slot.Pos(x),
								LI:  243,
								XY:  g.SymPosL2(slot.Pos(x), sym, wild),
							})
						}
						break
					}
					cw *= nw
					c1 *= n1
					c2 *= n2
				}
			}
		*/
	}
	// Scatters calculation
	var ns = counts[0][scat] + counts[1][scat] + counts[2][scat] + counts[3][scat] + counts[4][scat]
	if ns+nw >= 3 {
		var pay = ScatPay[ns+nw-1]
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * pay,
			MP:  float64(mm),
			Sym: scat,
			Num: slot.Pos(ns + nw),
			XY:  g.SymPos2(scat, wild),
			FS:  12,
		})
	}
	return nil
}

func (g *Game) Cost() float64 {
	return g.Bet * 25
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.SpinReels(reels)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
