package secretelixir

// See: https://www.slotsmate.com/software/novomatic/secret-elixir

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed secretelixir_bon.yaml
var rbon []byte

var ReelsBon = slot.ReadObj[*slot.Reels5x](rbon)

//go:embed secretelixir_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [13][4]float64{
	{0, 0, 100, 1000}, //  1 lover
	{0, 0, 40, 200},   //  2 wife
	{0, 0, 40, 200},   //  3 husband
	{0, 0, 20, 80},    //  4 owl
	{0, 0, 15, 60},    //  5 gargoyle1
	{0, 0, 15, 60},    //  6 gargoyle2
	{0, 0, 10, 40},    //  7 ace
	{0, 0, 10, 40},    //  8 king
	{0, 0, 5, 20},     //  9 queen
	{0, 0, 5, 20},     // 10 jack
	{0, 0, 5, 20},     // 11 ten
	{},                // 12 scatter
}

// Scatters payment.
var ScatPay = [4]float64{0, 0, 5, 25} // 12 scatter

// Bet lines
var BetLines = slot.BetLinesNvm10

type Game struct {
	slot.Screen5x3 `yaml:",inline"`
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

const wild, scat = 1, 12

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	for li, line := range BetLines[:g.Sel] {
		var numw, numl slot.Pos = 0, 4
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 4; x++ {
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
		if numl >= 3 && syml > 0 {
			payl = LinePay[syml-1][numl-1]
		}
		if payl > payw {
			var ml = float64(g.LY(5, line))
			var xy = line.CopyL(numl)
			if ml > 1 {
				xy[4] = line[4]
			}
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payl,
				Mult: ml,
				Sym:  syml,
				Num:  numl,
				Line: li + 1,
				XY:   xy,
			})
		} else if payw > 0 {
			var ml = float64(g.LY(5, line))
			var xy = line.CopyL(numl)
			if ml > 1 {
				xy[4] = line[4]
			}
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payw,
				Mult: ml,
				Sym:  wild,
				Num:  numw,
				Line: li + 1,
				XY:   xy,
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
			Free: 12,
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	if g.FSR == 0 {
		var reels, _ = slot.FindClosest(ReelsMap, mrtp)
		g.ReelSpin(reels)
	} else {
		g.ReelSpin(ReelsBon)
	}
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
