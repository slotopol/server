package indiandreaming

// See: https://freeslotshub.com/aristocrat/indian-dreaming/

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsMap slot.ReelsMap[slot.Reelx]

const (
	sn         = 12   // number of symbols
	wild, scat = 1, 2 // wild & scatter symbol IDs
	linemin    = 3    // minimum line symbols to win
)

// Lined payment.
var LinePay = [sn][5]float64{
	{},                     //  1 wild
	{},                     //  2 scatter
	{0, 0, 100, 200, 5000}, //  3 catcher
	{0, 0, 50, 100, 2500},  //  4 man
	{0, 0, 50, 100, 1000},  //  5 woman
	{0, 0, 10, 40, 250},    //  6 guy
	{0, 0, 6, 25, 150},     //  7 bull
	{0, 0, 6, 25, 150},     //  8 hatchet
	{0, 0, 6, 15, 80},      //  9 ace
	{0, 0, 4, 10, 80},      // 10 king
	{0, 0, 3, 10, 70},      // 11 queen
	{0, 0, 3, 10, 60},      // 12 jack
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 2, 15, 100} //  2 scatter

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

func (g *Game) Scanner(wins *slot.Wins) error {
	// Count symbols
	var counts [5 + 1][sn + 1]int
	for x := range 5 {
		var r = g.Scr[x]
		counts[x][r[0]]++
		counts[x][r[1]]++
		counts[x][r[2]]++
	}
	var wn = counts[0][wild] + counts[1][wild] + counts[2][wild] + counts[3][wild] + counts[4][wild]
	// Ways calculation
	if wn < 5 {
		var mwm = 1 // mult wild mode
		if g.FSR > 0 {
			mwm = 5
		}
		var combs1 = [sn + 1]int{0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1} // pure symbols
		var combs2 = [sn + 1]int{0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1} // symbols + wilds
		for x, cx := range counts {
			var cw = combs1[wild]
			for sym := range combs2 {
				var c1, c2 = combs1[sym], combs2[sym]
				var n = cx[sym] + cx[wild]
				combs1[sym] = c1 * cx[sym]
				combs2[sym] = c2 * n
				var c = (c2-c1-cw)*mwm + c1
				if x >= linemin && c > 0 && n == 0 {
					var pay = LinePay[sym-1][x-1]
					*wins = append(*wins, slot.WinItem{
						Pay: g.Bet * pay,
						MP:  float64(c),
						Sym: slot.Sym(sym),
						Num: slot.Pos(x),
						LI:  243,
						XY:  g.SymPosL2(slot.Pos(x), slot.Sym(sym), wild),
					})
				}
			}
		}
	}
	// Scatters calculation
	var sn = counts[0][scat] + counts[1][scat] + counts[2][scat] + counts[3][scat] + counts[4][scat]
	if sn+wn >= 3 {
		var mw float64 = 1 // mult wild
		if g.FSR > 0 && wn > 0 {
			mw = 5
		}
		var pay = ScatPay[sn+wn-1]
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * pay,
			MP:  mw,
			Sym: scat,
			Num: slot.Pos(sn + wn),
			XY:  g.SymPos2(scat, wild),
			FS:  12,
		})
	}
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanWays(wins *slot.Wins) {
	var mwm float64 = 1 // mult wild mode
	if g.FSR > 0 {
		mwm = 5
	}
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
						var syml slot.Sym
						var x slot.Pos
						for x = 1; x <= 5; x++ {
							var sx = g.LX(x, line)
							if sx == wild {
								mw = mwm
							} else if syml == 0 {
								syml = sx
							} else if sx != syml {
								numl = x - 1
								break
							}
						}

						if numl >= 3 && syml > 0 {
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
	if sn, wn := g.SymNum2(scat, wild); sn+wn >= 3 {
		var mw float64 = 1 // mult wild
		if g.FSR > 0 && wn > 0 {
			mw = 5
		}
		var pay = ScatPay[sn+wn-1]
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * pay,
			MP:  mw,
			Sym: scat,
			Num: sn + wn,
			XY:  g.SymPos2(scat, wild),
			FS:  12,
		})
	}
}

func (g *Game) Cost() float64 {
	return g.Bet * 25
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.SpinReels(reels)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
