package beetlemania

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed beetlemania_bon.yaml
var rbon []byte

var ReelsBon = slot.ReadObj[*slot.Reels5x](rbon)

//go:embed beetlemania_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [11][5]float64{
	{0, 10, 80, 1000, 5000}, //  1 bee
	{0, 5, 30, 200, 1000},   //  2 snail
	{0, 5, 25, 100, 500},    //  3 fly
	{0, 5, 15, 65, 250},     //  4 worm
	{0, 0, 10, 40, 200},     //  5 ace
	{0, 0, 10, 40, 200},     //  6 king
	{0, 0, 5, 20, 100},      //  7 queen
	{0, 0, 5, 20, 100},      //  8 jack
	{0, 0, 5, 20, 100},      //  9 ten
	{},                      // 10 note
	{},                      // 11 jazzbee
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 2, 15, 50} // 10 note

const (
	jbonus = 1 // jazzbee freespins bonus
)

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

const wild, scat = 1, 10
const jazz = 11

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

func (g *Game) ScatNumCont() (n slot.Pos) {
	for x := 0; x < 5; x++ {
		var r = g.Scr[x]
		if r[0] == scat || r[1] == scat || r[2] == scat {
			n++
		} else {
			break // scatters should be continuous
		}
	}
	return
}

func (g *Game) ScatPosCont() (l slot.Linex) {
	var x int
	for x = 0; x < 5; x++ {
		var r = g.Scr[x]
		if r[0] == scat {
			l[x] = 1
		} else if r[1] == scat {
			l[x] = 2
		} else if r[2] == scat {
			l[x] = 3
		} else {
			break // scatters should be continuous
		}
	}
	return
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
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

		/*var numw, numl slot.Pos
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = g.LY(x, line)
			if sx != wild {
				if sx != scat && sx != jazz {
					syml = sx
					numl = numw + 1
					for x := numl + 1; x <= 5; x++ {
						var sx = g.LY(x, line)
						if sx == syml || sx == wild {
							numl++
						} else {
							break
						}
					}
				}
				break
			}
			numw++
		}*/

		var payw, payl float64
		if numw >= 2 {
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
	if g.FSR > 0 {
		var y slot.Pos
		if g.At(3, 1) == jazz {
			y = 1
		} else if g.At(3, 2) == jazz {
			y = 2
		} else if g.At(3, 3) == jazz {
			y = 3
		} else {
			return // ignore scatters on freespins
		}
		var xy slot.Linex
		xy.Set(3, y)
		*wins = append(*wins, slot.WinItem{
			Mult: 1,
			Sym:  jazz,
			Num:  1,
			XY:   xy,
			BID:  jbonus,
		})
		return
	}

	if count := g.ScatNumCont(); count >= 3 {
		var pay = ScatPay[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   g.ScatPosCont(),
			Free: 10,
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

func (g *Game) Spawn(wins slot.Wins, fund, mrtp float64) {
	for i, wi := range wins {
		switch wi.BID {
		case jbonus:
			wins[i].Pay = min(g.Gain, 100_000*g.Bet)
		}
	}
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
