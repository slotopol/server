package indiandreaming

// See: https://freeslotshub.com/aristocrat/indian-dreaming/

import (
	slot "github.com/slotopol/server/game/slot"
)

// Lined payment.
var LinePay = [13][5]float64{
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
	slot.Slot5x3 `yaml:",inline"`
	// free spin number
	FS int `json:"fs,omitempty" yaml:"fs,omitempty" xml:"fs,omitempty"`
}

func NewGame() *Game {
	return &Game{
		Slot5x3: slot.Slot5x3{
			Sel: slot.MakeBitNum(25, 1),
			Bet: 1,
		},
		FS: 0,
	}
}

const wild, scat = 1, 2

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	if screen.ScatNum(wild) < 5 {
		g.ScanLined(screen, wins)
	}
	g.ScanScatters(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen slot.Screen, wins *slot.Wins) {
	var syml slot.Sym
	var x, y slot.Pos
	var pool [243]slot.WinItem
	for syml = 3; syml <= 12; syml++ {
		clear(pool[:])
		var wn = 0
		for y = 1; y <= 3; y++ {
			var sx = screen.At(1, y)
			if sx == syml || sx == wild {
				pool[wn].Num++
				if sx == syml {
					pool[wn].Sym = syml
				}
				pool[wn].XY.Set(1, y)
				wn++
			}
		}
		for x = 2; x <= 5; x++ {
			var wx = wn
			var yi = 0
			for y = 1; y <= 3; y++ {
				var sx = screen.At(x, y)
				if sx == syml || sx == wild {
					if yi > 0 {
						copy(pool[wx*yi:], pool[:wx])
					}
					for i := range wx {
						pool[i+wx*yi].XY.Set(x, y)
					}
					yi++
					wn += wx
				}
			}
			if yi == 0 {
				break
			}
			var pay = g.Bet * LinePay[syml-1][x-1]
			for i := range wn {
				var sx = screen.At(x, pool[i].XY.At(x))
				if sx == syml {
					pool[i].Sym = syml
				} else if g.FS > 0 {
					pool[i].Mult = 5
				}
				pool[i].Pay = pay
				pool[i].Num = x
			}
		}
		if wn > 0 {
			for i := range wn {
				if pool[i].Sym == syml {
					*wins = append(*wins, pool[i])
				}
			}
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen slot.Screen, wins *slot.Wins) {
	if count := screen.ScatNum(wild) + screen.ScatNum(scat); count >= 3 {
		var pay, fs = ScatPay[count-1], 12
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
			Free: fs,
		})
	}
}

func (g *Game) Spin(screen slot.Screen, mrtp float64) {
	var _, reels = FindReels(mrtp)
	screen.Spin(reels)
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
	return slot.ErrNoFeature
}
