package gonzosquest

// See: https://www.slotsmate.com/software/netent/gonzos-quest

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed gonzosquest_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [9][5]float64{
	{},                    // 1 wild
	{},                    // 2 freefall
	{0, 0, 50, 250, 2500}, // 3 mask1
	{0, 0, 20, 100, 1000}, // 4 mask2
	{0, 0, 15, 50, 500},   // 5 mask3
	{0, 0, 10, 25, 200},   // 6 mask4
	{0, 0, 5, 20, 100},    // 7 mask5
	{0, 0, 4, 15, 75},     // 8 mask6
	{0, 0, 3, 10, 50},     // 9 mask7
}

// Bet lines
var BetLines = slot.BetLinesNetEnt5x3[:20]

type Game struct {
	slot.Cascade5x3 `yaml:",inline"`
	slot.Slotx      `yaml:",inline"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slotx: slot.Slotx{
			Sel: len(BetLines),
			Bet: 1,
		},
	}
}

func (g *Game) Clone() slot.SlotGame {
	var clone = *g
	return &clone
}

func (g *Game) Free() bool {
	return g.FSR != 0 || g.Cascade()
}

const wild, scat = 1, 2

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	for li, line := range BetLines[:g.Sel] {
		var numl slot.Pos = 5
		var syml = g.LY(1, line)
		var x slot.Pos
		for x = 2; x <= 5; x++ {
			var sx = g.LY(x, line)
			if sx == wild {
				continue
			} else if sx != syml {
				numl = x - 1
				break
			}
		}
		if numl >= 3 && syml > scat {
			var fm = float64(min(g.CFN, 5))
			var mm float64 = 1 // mult mode
			if g.FSR > 0 {
				mm = 3
			}
			var pay = LinePay[syml-1][numl-1]
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * pay,
				Mult: fm * mm,
				Sym:  syml,
				Num:  numl,
				Line: li + 1,
				XY:   line.CopyL(numl),
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.ScatNum(scat); count >= 3 {
		*wins = append(*wins, slot.WinItem{
			Sym:  scat,
			Num:  count,
			XY:   g.ScatPos(scat),
			Free: 10,
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	g.ReelSpin(reels)
}

func (g *Game) Prepare() {
	g.NewFall()
}

func (g *Game) Apply(wins slot.Wins) {
	g.Slotx.Apply(wins)
	g.Strike(wins)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
