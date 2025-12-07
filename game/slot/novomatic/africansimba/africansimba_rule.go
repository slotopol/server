package africansimba

// See: https://www.slotsmate.com/software/novomatic/african-simba

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsMap slot.ReelsMap[*slot.Reels5x]

// Lined payment.
var LinePay = [12][5]float64{
	{},                     //  1 wild
	{},                     //  2 scatter
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
						var numl slot.Pos = 5
						var syml = g.LY(1, line)
						var x slot.Pos
						for x = 2; x <= 5; x++ {
							var sx = g.LY(x, line)
							if sx != syml && sx != wild {
								numl = x - 1
								break
							}
						}

						if numl >= 3 && syml > scat {
							var mm float64 = 1 // mult mode
							if g.FSR > 0 {
								mm = 3
							}
							// var li = (int(line[0])-1)*81 + (int(line[1])-1)*27 + (int(line[2])-1)*9 + (int(line[line[4]])-1)*3 + int(line[5])
							*wins = append(*wins, slot.WinItem{
								Pay: g.Bet * LinePay[syml-1][numl-1],
								MP:  mm,
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
		*wins = append(*wins, slot.WinItem{
			Sym: scat,
			Num: count,
			XY:  g.SymPos(scat),
			FS:  12,
		})
	}
}

func (g *Game) Cost() float64 {
	return g.Bet * 25
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.ReelSpin(reels)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
