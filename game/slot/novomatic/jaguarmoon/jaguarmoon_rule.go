package jaguarmoon

// See: https://www.slotsmate.com/software/novomatic/jaguar-moon

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed jaguarmoon_bon.yaml
var rbon []byte

var ReelsBon = slot.ReadObj[*slot.Reels5x](rbon)

//go:embed jaguarmoon_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [12][5]float64{
	{},                   //  1 wild
	{},                   //  2 scatter
	{0, 0, 40, 200, 800}, //  3 wooman
	{0, 0, 20, 60, 200},  //  4 panther
	{0, 0, 10, 50, 100},  //  5 footprint
	{0, 0, 10, 30, 100},  //  6 rings
	{0, 0, 4, 10, 50},    //  7 ace
	{0, 0, 4, 10, 50},    //  8 king
	{0, 0, 4, 10, 50},    //  9 queen
	{0, 0, 2, 8, 40},     // 10 jack
	{0, 0, 2, 8, 40},     // 11 ten
	{0, 0, 2, 8, 40},     // 12 nine
}

// Scatter freespins table
var ScatFreespin = [6]int{0, 0, 8, 12, 15, 20} // 2 scatter

// Free games multipliers
var FreeMult = [6]float64{0, 0, 2, 3, 4, 5}

type Game struct {
	slot.Screen5x3 `yaml:",inline"`
	slot.Slotx     `yaml:",inline"`
	// multiplier on freespins
	M float64 `json:"m,omitempty" yaml:"m,omitempty" xml:"m,omitempty"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slotx: slot.Slotx{
			Sel: 10,
			Bet: 1,
		},
		M: 0,
	}
}

func (g *Game) Clone() slot.SlotGame {
	var clone = *g
	return &clone
}

const wild, scat = 1, 2

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
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
							if sx == wild {
								continue
							} else if syml == 0 {
								syml = sx
							} else if sx != syml {
								numl = x - 1
								break
							}
						}

						if numl >= 3 && syml > scat {
							var mm float64 = 1 // mult mode
							if g.FSR > 0 {
								mm = g.M
							}
							// var li = (int(line[0])-1)*81 + (int(line[1])-1)*27 + (int(line[2])-1)*9 + (int(line[line[4]])-1)*3 + int(line[5])
							*wins = append(*wins, slot.WinItem{
								Pay:  g.Bet * LinePay[syml-1][numl-1],
								Mult: mm,
								Sym:  syml,
								Num:  numl,
								Line: 243,
								XY:   line.CopyL(numl),
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

func (g *Game) BonNum() (n, count slot.Pos) {
	for x := range 3 { // only on reels 1, 2, 3
		for y := range 3 {
			if g.Scr[x][y] == scat {
				n++
				count++
				if y < 2 && g.Scr[x][y+1] == scat {
					count++
				}
			}
		}
	}
	return
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if n, count := g.BonNum(); n == 3 {
		var fs = ScatFreespin[count]
		*wins = append(*wins, slot.WinItem{
			Sym:  scat,
			Num:  count,
			XY:   g.ScatPos(scat),
			Free: fs,
		})
	}
}

func (g *Game) Cost() (float64, bool) {
	return g.Bet * 10, false
}

func (g *Game) Spin(mrtp float64) {
	if g.FSR == 0 {
		var reels, _ = slot.FindClosest(ReelsMap, mrtp)
		g.ReelSpin(reels)
	} else {
		g.ReelSpin(ReelsBon)
	}
}

func (g *Game) Prepare() {
	if g.FSR == 0 {
		g.M = 0
	}
}

func (g *Game) Apply(wins slot.Wins) {
	g.Slotx.Apply(wins)
	for _, wi := range wins {
		if wi.Free > 0 {
			g.M = FreeMult[wi.Num-1]
		}
	}
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
