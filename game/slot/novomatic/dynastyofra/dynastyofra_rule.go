package dynastyofra

// See: https://www.slotsmate.com/software/novomatic/dynasty-of-ra

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed dynastyofra_bon.yaml
var rbon []byte

var ReelsBon = slot.ReadObj[*slot.Reels5x](rbon)

//go:embed dynastyofra_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [11][5]float64{
	{0, 10, 50, 200, 1000}, //  1 Cleopatra
	{0, 5, 25, 100, 500},   //  2 sphinx
	{0, 5, 20, 50, 250},    //  3 mask
	{0, 5, 20, 50, 250},    //  4 necklace
	{0, 5, 20, 50, 250},    //  5 beads
	{0, 0, 10, 30, 125},    //  6 ace
	{0, 0, 10, 30, 125},    //  7 king
	{0, 0, 5, 20, 100},     //  8 queen
	{0, 0, 5, 20, 100},     //  9 jack
	{0, 0, 5, 20, 100},     // 10 ten
	{},                     // 11 scatter
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 2, 20, 200} // 11 scatter

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 10, 10, 10} // 11 scatter

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

const book = 11

func (g *Game) Scanner(wins *slot.Wins) error {
	if g.FSR == 0 {
		g.ScanLinedReg(wins)
	} else {
		g.ScanLinedBon(wins)
	}
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation on regular games.
func (g *Game) ScanLinedReg(wins *slot.Wins) {
	for li, line := range BetLines[:g.Sel] {
		var numl slot.Pos = 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = g.LY(x, line)
			if sx == book {
				continue
			} else if syml == 0 {
				syml = sx
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		if syml > 0 {
			if payl := LinePay[syml-1][numl-1]; payl > 0 {
				*wins = append(*wins, slot.WinItem{
					Pay:  g.Bet * payl,
					Mult: 1,
					Sym:  syml,
					Num:  numl,
					Line: li + 1,
					XY:   line.CopyL(numl),
				})
			}
		}
	}
}

// Lined symbols calculation on bonus games.
func (g *Game) ScanLinedBon(wins *slot.Wins) {
	var line slot.Linex
	var pays float64
	var wi slot.WinItem
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
							if sx != syml {
								numl = x - 1
								break
							}
						}

						if syml > 0 {
							if payl := LinePay[syml-1][numl-1]; payl > pays {
								wi.Sym = syml
								wi.Num = numl
								wi.XY = line.CopyL(numl)
								pays = payl
								switch numl {
								case 3:
									continue loop3
								case 4:
									continue loop4
								case 5:
									continue loop5
								}
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
	if wi.Sym == 0 {
		return
	}

	wi.Line = 243
	*wins = append(*wins, wi)

	wi.Pay = pays
	wi.Mult = 1
	for li, line := range BetLines[:g.Sel] {
		wi.Line = li + 1
		wi.XY = line.CopyL(wi.Num)
		*wins = append(*wins, wi)
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.ScatNum(book); count >= 3 {
		var pay, fs = ScatPay[count-1], ScatFreespin[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * pay,
			Mult: 1,
			Sym:  book,
			Num:  count,
			XY:   g.ScatPos(book),
			Free: fs,
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
