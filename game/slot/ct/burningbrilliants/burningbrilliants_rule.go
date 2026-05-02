package burningbrilliants

// See: https://www.slotsmate.com/software/ct-interactive/100-burning-brilliants

import (
	"github.com/slotopol/server/game/slot"
)

const (
	sn         = 10   // number of symbols
	wild, scat = 1, 2 // wild & scatter symbol IDs
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [sn][5]float64{
	{},                    //  1 wild
	{},                    //  2 scatter
	{0, 0, 35, 100, 2000}, //  3 star
	{0, 0, 10, 25, 100},   //  4 ruby
	{0, 0, 10, 25, 100},   //  5 emerald
	{0, 0, 10, 25, 100},   //  6 topaz
	{0, 0, 7, 10, 100},    //  7 spader
	{0, 0, 7, 10, 100},    //  8 heart
	{0, 0, 5, 10, 100},    //  9 diamond
	{0, 0, 5, 10, 100},    // 10 club
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 5, 20, 500} // 2 scatter

// Bet lines
var BetLines = slot.BetLinesCT5x4v3[:]

type Game struct {
	slot.Cascade5x4 `yaml:",inline"`
	slot.Slotx      `yaml:",inline"`
}

// Declare conformity with SlotCascade interface.
var _ slot.SlotCascade = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slotx: slot.Slotx{
			Sel: 100,
			Bet: 1,
		},
	}
}

func (g *Game) Clone() slot.SlotGeneric {
	var clone = *g
	return &clone
}

func (g *Game) FreeMode() bool {
	return g.FSR != 0 || g.Cascade()
}

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	for li, line := range BetLines[:g.Sel] {
		var numl slot.Pos = 5
		var syml = g.LX(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = g.LX(x, line)
			if sx != syml && sx != wild {
				numl = x - 1
				break
			}
		}
		if pay := LinePay[syml-1][numl-1]; pay > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * pay,
				MP:  1,
				Sym: syml,
				Num: numl,
				LI:  li + 1,
				XY:  line.HitxL(numl),
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
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.SpinReels(reels)
}

func (g *Game) Prepare() {
	g.UntoFall()
}

func (g *Game) Apply(wins slot.Wins) {
	g.Slotx.Apply(wins)
	g.Strike(wins)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
