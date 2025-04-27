package indiandreaming

// See: https://freeslotshub.com/aristocrat/indian-dreaming/

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed indiandreaming_reel.yaml
var reels []byte

var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [12][5]float64{
	{},                     //  1 wild
	{},                     //  2 scatter
	{0, 0, 100, 200, 5000}, //  3 cash catcher
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
			Sel: 25,
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
	if g.ScatNum(wild) < 5 {
		g.ScanLined(wins)
	}
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
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
							var sx = g.LY(x, line)
							if sx == wild {
								mw = mwm
							} else if syml == 0 && sx != scat {
								syml = sx
							} else if sx != syml {
								numl = x - 1
								break
							}
						}

						if numl >= 3 && syml > 0 {
							// var li = (int(line[0])-1)*81 + (int(line[1])-1)*27 + (int(line[2])-1)*9 + (int(line[line[4]])-1)*3 + int(line[5])
							*wins = append(*wins, slot.WinItem{
								Pay:  g.Bet * LinePay[syml-1][numl-1],
								Mult: mw,
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

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	var sn, wn = g.ScatNum(scat), g.ScatNum(wild)
	if count := sn + wn; count >= 3 {
		var mw float64 = 1 // mult wild
		if g.FSR > 0 && wn > 0 {
			mw = 5
		}
		var pay = ScatPay[count-1]
		var line = g.ScatPos(scat)
		line.Cover(g.ScatPos(wild))
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * pay,
			Mult: mw,
			Sym:  scat,
			Num:  count,
			XY:   line,
			Free: 12,
		})
	}
}

func (g *Game) Cost() (float64, bool) {
	return g.Bet * 25, false
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	g.ReelSpin(reels)
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
