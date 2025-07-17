package lovelymermaid

// See: https://www.slotsmate.com/software/novomatic/lovely-mermaid

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed lovelymermaid_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

//go:embed lovelymermaid_jack.yaml
var jack []byte

var JackMap = slot.ReadMap[float64](jack)

// Lined payment.
var LinePay = [13][5]float64{
	{0, 0, 20, 200, 2000}, //  1 mermaid
	{0, 2, 10, 40, 400},   //  2 lobster
	{0, 2, 10, 40, 400},   //  3 turtle
	{0, 0, 8, 30, 300},    //  4 blowfish
	{0, 0, 6, 20, 200},    //  5 seahorse
	{0, 0, 6, 20, 200},    //  6 parrotfish
	{0, 0, 4, 10, 100},    //  7 ace
	{0, 0, 4, 10, 100},    //  8 king
	{0, 0, 4, 10, 80},     //  9 queen
	{0, 0, 4, 10, 80},     // 10 jack
	{0, 0, 4, 10, 80},     // 11 ten
	{0, 0, 4, 10, 80},     // 12 nine
	{},                    // 13 scatter
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 3, 20, 400} // 13 scatter

// Bet lines
var BetLines = slot.BetLinesNvm5x4[:40]

type Game struct {
	slot.Screen5x4 `yaml:",inline"`
	slot.Slotx     `yaml:",inline"`
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

const (
	lmj        = 1     // jackpot ID
	wild, scat = 1, 13 // symbols
)

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

func (g *Game) Filled() slot.Sym {
	var sym = g.Scr[4][3]
	for x := range 5 {
		for y := range 4 {
			if g.Scr[x][y] != sym {
				return 0
			}
		}
	}
	return sym
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	if sym := g.Filled(); sym != 0 {
		*wins = append(*wins, slot.WinItem{
			Sym: sym,
			JID: lmj,
		})
		return
	}
	for li, line := range BetLines[:g.Sel] {
		var numw, numl slot.Pos = 0, 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = g.LY(x, line)
			if sx == wild {
				if syml == 0 {
					numw = x
				}
			} else if syml == 0 {
				syml = sx
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		var payw, payl float64
		if numw >= 3 {
			payw = LinePay[wild-1][numw-1]
		}
		if numl >= 2 && syml > 0 {
			payl = LinePay[syml-1][numl-1]
		}
		if payl > payw {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payl,
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li + 1,
				XY:   line.CopyL(numl),
			})
		} else if payw > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payw,
				Mult: 1,
				Sym:  wild,
				Num:  numw,
				Line: li + 1,
				XY:   line.CopyL(numw),
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.ScatNum(scat); count >= 3 {
		var pay = ScatPay[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   g.ScatPos(scat),
			Free: 25,
		})
	}
}

func (g *Game) Cost() (float64, bool) {
	return g.Bet * float64(g.Sel), true
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	g.ReelSpin(reels)
}

func (g *Game) Spawn(wins slot.Wins, fund, mrtp float64) {
	for i, wi := range wins {
		if wi.JID != 0 {
			var bulk, _ = slot.FindClosest(JackMap, mrtp)
			var jf = bulk * g.Bet / slot.JackBasis
			if jf > 1 {
				jf = 1
			}
			wins[i].Jack = jf * fund
		}
	}
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
