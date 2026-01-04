package rodeopower

// See: https://www.slotsmate.com/software/ct-interactive/rodeo-power

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsBon slot.Reelx

var ReelsMap slot.ReelsMap[slot.Reelx]

// Lined payment.
var LinePay = [13][5]float64{
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

const wild, scat = 1, 2

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanWays(wins)
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanWays(wins *slot.Wins) {
	var line slot.Linex
loop1:
	for line[0] = 1; line[0] <= 3; line[0]++ {
	loop2:
		for line[1] = 1; line[1] <= 3; line[1]++ {
		loop3:
			for line[2] = 1; line[2] <= 3; line[2]++ {
			loop4:
				for line[3] = 1; line[3] <= 3; line[3]++ {
				loop5:
					for line[4] = 1; line[4] <= 3; line[4]++ {
						var mw float64 = 1 // mult wild
						var numl slot.Pos = 5
						var syml = g.LX(1, line)
						var x slot.Pos
						for x = 2; x <= 5; x++ {
							var sx = g.LX(x, line)
							if sx == wild {
								if x == 2 {
									mw *= 2
								} else { // x == 4
									mw *= 5
								}
							} else if sx != syml {
								numl = x - 1
								break
							}
						}

						if numl >= 3 && syml > scat {
							// var li = (int(line[0])-1)*81 + (int(line[1])-1)*27 + (int(line[2])-1)*9 + (int(line[line[4]])-1)*3 + int(line[5])
							*wins = append(*wins, slot.WinItem{
								Pay: g.Bet * LinePay[syml-1][numl-1],
								MP:  mw,
								Sym: syml,
								Num: numl,
								LI:  243,
								XY:  line.HitxL(numl),
							})
							switch numl {
							case 3:
								continue loop3
							case 4:
								continue loop4
							case 5:
								continue loop5
							}
						}
						switch numl + 1 {
						case 1:
							continue loop1
						case 2:
							continue loop2
						case 3:
							continue loop3
						case 4:
							continue loop4
						case 5:
							continue loop5
						}
					}
				}
			}
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.SymNum(scat); count >= 3 {
		var pay = ScatPay[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * pay,
			MP:  1,
			Sym: scat,
			Num: count,
			XY:  g.SymPos(scat),
			FS:  15,
		})
	}
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
