package rodeopower

// See: https://www.slotsmate.com/software/ct-interactive/rodeo-power

import (
	"github.com/slotopol/server/game/slot"
)

const (
	sn         = 13   // number of symbols
	wild, scat = 1, 2 // wild & scatter symbol IDs
	linemin    = 3    // minimum line symbols to win
)

var ReelsBon slot.Reelx

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [sn][5]float64{
	{},                    //  1 wild (2, 4 reels only)
	{},                    //  2 scatter
	{0, 0, 50, 300, 1000}, //  3 shoe
	{0, 0, 35, 300, 500},  //  4 woman
	{0, 0, 25, 250, 400},  //  5 spurs
	{0, 0, 25, 250, 400},  //  6 belt
	{0, 0, 10, 20, 120},   //  7 saddle
	{0, 0, 10, 20, 120},   //  8 hat
	{0, 0, 10, 20, 120},   //  9 boots
	{0, 0, 5, 10, 100},    // 10 ace
	{0, 0, 5, 10, 100},    // 11 king
	{0, 0, 5, 10, 100},    // 12 queen
	{0, 0, 5, 10, 100},    // 13 jack
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 5, 20, 100} // 2 scatter

type Game struct {
	slot.Screen5x3 `yaml:",inline"`
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
	var active [sn + 1]bool       // symbols present on 1st reel
	for x := range 5 {
		var cx = &counts[x]
		for _, sym := range g.Scr[x] {
			cx[sym]++
		}
	}
	for _, sym := range g.Scr[0] {
		active[sym] = true
	}
	// Ways calculation
	var sym slot.Sym
	for sym = 3; sym <= sn; sym++ { // ignore wild & scatter
		if !active[sym] {
			continue
		}
		var c = 1 // current ways
		for x, cx := range counts {
			var n = cx[sym]
			if x == 1 {
				n += cx[wild] * 2
			} else if x == 3 {
				n += cx[wild] * 5
			}
			if n == 0 {
				if x >= linemin {
					var pay = LinePay[sym-1][x-1]
					*wins = append(*wins, slot.WinItem{
						Pay: g.Bet * pay,
						MP:  float64(c),
						Sym: sym,
						Num: slot.Pos(x),
						LI:  243,
						XY:  g.SymPosL2(slot.Pos(x), sym, wild),
					})
				}
				break
			}
			c *= n
		}
	}
	// Scatters calculation
	if count := counts[0][scat] + counts[1][scat] + counts[2][scat] + counts[3][scat] + counts[4][scat]; count >= 3 {
		var pay = ScatPay[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * pay,
			MP:  1,
			Sym: scat,
			Num: slot.Pos(count),
			XY:  g.SymPos(scat),
			FS:  15,
		})
	}
	return nil
}

func (g *Game) Cost() float64 {
	return g.Bet * 25
}

func (g *Game) Spin(mrtp float64) {
	if g.FSR == 0 {
		var reels, _ = ReelsMap.FindClosest(mrtp)
		g.SpinReels(reels)
	} else {
		g.SpinReels(ReelsBon)
	}
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
