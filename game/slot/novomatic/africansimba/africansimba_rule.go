package africansimba

// See: https://www.slotsmate.com/software/novomatic/african-simba

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
	{},                     //  1 wild (2, 3, 4 reels only)
	{},                     //  2 scatter (1, 3, 5 reels only)
	{0, 0, 100, 500, 2500}, //  3 giraffe
	{0, 0, 50, 150, 750},   //  4 buffalo
	{0, 0, 25, 75, 250},    //  5 lemur
	{0, 0, 25, 75, 250},    //  6 flamingo
	{0, 0, 10, 25, 125},    //  7 ace
	{0, 0, 10, 25, 125},    //  8 king
	{0, 0, 10, 25, 125},    //  9 queen
	{0, 0, 5, 20, 100},     // 10 jack
	{0, 0, 5, 20, 100},     // 11 ten
	{0, 0, 5, 20, 100},     // 12 nine
}

type Game struct {
	slot.Grid5x3 `yaml:",inline"`
	slot.Slotx   `yaml:",inline"`
}

// Declare conformity with SlotGeneric interface.
var _ slot.SlotGeneric = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slotx: slot.Slotx{
			Bet: 1,
		},
	}
}

func (g *Game) Clone() slot.SlotGeneric {
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
	// Ways calculation
	var combs = [sn + 1]int{0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	for x, cx := range counts {
		for sym, c := range combs {
			var n = cx[sym] + cx[wild]
			combs[sym] = c * n
			if x >= linemin && c > 0 && n == 0 {
				var pay = LinePay[sym-1][x-1]
				var mm float64 = 1 // mult mode
				if g.FSR > 0 {
					mm = 3
				}
				*wins = append(*wins, slot.WinItem{
					Pay: g.Bet * pay,
					MP:  mm * float64(c),
					Sym: slot.Sym(sym),
					Num: slot.Pos(x),
					LI:  243,
					XY:  g.SymPosL2(slot.Pos(x), slot.Sym(sym), wild),
				})
			}
		}
	}
	// Scatters calculation
	if count := counts[0][scat] + counts[2][scat] + counts[4][scat]; count >= 3 {
		*wins = append(*wins, slot.WinItem{
			Sym: scat,
			Num: slot.Pos(count),
			XY:  g.SymPos(scat),
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
