package beetlemania

import (
	slot "github.com/slotopol/server/game/slot"
)

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
	{0, 0, 0, 0, 0},         // 10 note
	{0, 0, 0, 0, 0},         // 11 jazzbee
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 2, 15, 50} // 10 note

const (
	jbonus = 1 // jazzbee freespins bonus
)

const (
	jid = 1 // jackpot ID
)

// Jackpot win combinations.
var Jackpot = [11][5]int{
	{0, 0, 0, 0, 0}, //  1 bee
	{0, 0, 0, 0, 0}, //  2 snail
	{0, 0, 0, 0, 0}, //  3 fly
	{0, 0, 0, 0, 0}, //  4 worm
	{0, 0, 0, 0, 0}, //  5 ace
	{0, 0, 0, 0, 0}, //  6 king
	{0, 0, 0, 0, 0}, //  7 queen
	{0, 0, 0, 0, 0}, //  8 jack
	{0, 0, 0, 0, 0}, //  9 ten
	{0, 0, 0, 0, 0}, // 10 note
	{0, 0, 0, 0, 0}, // 11 jazzbee
}

// Bet lines
var bl = slot.BetLinesNvm10

type Game struct {
	slot.Slot5x3 `yaml:",inline"`
	// free spin number
	FS int `json:"fs,omitempty" yaml:"fs,omitempty" xml:"fs,omitempty"`
}

func NewGame() *Game {
	return &Game{
		Slot5x3: slot.Slot5x3{
			Sel: slot.MakeBitNum(5, 1),
			Bet: 1,
		},
		FS: 0,
	}
}

const wild, scat = 1, 10
const jazz = 11

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	g.ScanLined(screen, wins)
	g.ScanScatters(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen slot.Screen, wins *slot.Wins) {
	for li := g.Sel.Next(0); li != -1; li = g.Sel.Next(li) {
		var line = bl[li-1]

		/*var numw, numl int
		var syml slot.Sym
		for x := 1; x <= 5; x++ {
			var symx = screen.Pos(x, line)
			if symx == wild {
				if syml == 0 {
					numw = x
				} else {
					numl = x
				}
			} else if symx == scat || symx == jazz {
				break
			} else if syml == 0 {
				syml = symx
				numl = x
			} else if symx == syml {
				numl = x
			} else {
				break
			}
		}*/

		var numw, numl slot.Pos
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = screen.Pos(x, line)
			if sx != wild {
				if sx != scat && sx != jazz {
					syml = sx
					numl = numw + 1
					for x := numl + 1; x <= 5; x++ {
						var sx = screen.Pos(x, line)
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
		}

		var payw, payl float64
		if numw > 0 {
			payw = LinePay[wild-1][numw-1]
		}
		if numl > 0 {
			payl = LinePay[syml-1][numl-1]
		}
		if payl > payw {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payl,
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyL(numl),
			})
		} else if payw > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payw,
				Mult: 1,
				Sym:  wild,
				Num:  numw,
				Line: li,
				XY:   line.CopyL(numw),
				Jack: Jackpot[wild-1][numw-1],
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen slot.Screen, wins *slot.Wins) {
	if g.FS > 0 {
		var y slot.Pos
		if screen.At(3, 1) == jazz {
			y = 1
		} else if screen.At(3, 1) == jazz {
			y = 2
		} else if screen.At(3, 3) == jazz {
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

	if count := screen.ScatNum(scat); count >= 3 {
		var pay = ScatPay[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel.Num()) * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPosCont(scat),
			Free: 10,
		})
	}
}

func (g *Game) Spin(screen slot.Screen, mrtp float64) {
	if g.FS == 0 {
		var _, reels = slot.FindReels(ReelsMap, mrtp)
		screen.Spin(reels)
	} else {
		screen.Spin(&ReelsBonu)
	}
}

func (g *Game) Spawn(screen slot.Screen, wins slot.Wins) {
	for i, wi := range wins {
		switch wi.BID {
		case jbonus:
			wins[i].Pay = min(g.Gain, 100_000*g.Bet)
		}
	}
}

func (g *Game) Apply(screen slot.Screen, wins slot.Wins) {
	if g.FS > 0 {
		g.Gain += wins.Gain()
	} else {
		g.Gain = wins.Gain()
	}

	if g.FS > 0 {
		g.FS--
	}
	for _, wi := range wins {
		if wi.Free > 0 {
			g.FS += wi.Free
		}
	}
}

func (g *Game) FreeSpins() int {
	return g.FS
}

func (g *Game) SetSel(sel slot.Bitset) error {
	return g.SetSelNum(sel, len(bl))
}
