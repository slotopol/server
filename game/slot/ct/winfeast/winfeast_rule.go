package winfeast

// See: https://www.livebet.com/casino/slots/ct-interactive/win-feast

import (
	"github.com/slotopol/server/game/slot"
)

const (
	sn         = 10   // number of symbols
	wild, scat = 1, 2 // wild & scatter symbol IDs
	cfn8, cfn9 = 1, 2 // bonus ID for CFN = 8 and CFN = 9
)

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [sn][5]float64{
	{},                    //  1 wild (2, 3, 4 reels only)
	{},                    //  2 scatter
	{0, 0, 30, 100, 1000}, //  3 man
	{0, 0, 5, 20, 400},    //  4 banana
	{0, 0, 5, 20, 100},    //  5 apple
	{0, 0, 5, 20, 100},    //  6 melon
	{0, 0, 5, 10, 100},    //  7 orange
	{0, 0, 5, 10, 100},    //  8 lemon
	{0, 0, 5, 10, 100},    //  9 plum
	{0, 0, 5, 10, 100},    // 10 cherry
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 20, 50, 500} // 2 scatter

// Bet lines
var BetLines = slot.BetLinesCT5x3[:]

type Game struct {
	slot.Cascade5x3 `yaml:",inline"`
	slot.Slotx      `yaml:",inline"`
}

// Declare conformity with SlotCascade interface.
var _ slot.SlotCascade = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slotx: slot.Slotx{
			Sel: 20,
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
	if len(*wins) == 0 {
		if g.CFN == 8 {
			*wins = append(*wins, slot.WinItem{
				Pay: 100,
				MP:  1,
				BID: cfn8,
			})
		} else if g.CFN >= 9 {
			*wins = append(*wins, slot.WinItem{
				Pay: 1000,
				MP:  1,
				BID: cfn9,
			})
		}
	}
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
