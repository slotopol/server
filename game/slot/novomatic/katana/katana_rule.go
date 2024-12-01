package katana

import (
	"github.com/slotopol/server/game/slot"
)

// Lined payment.
var LinePay = [12][5]float64{
	{0, 0, 100, 500, 1000}, // 1  shogun
	{0, 0, 80, 250, 800},   // 2  geisha
	{0, 0, 40, 200, 600},   // 3  general
	{0, 0, 40, 200, 600},   // 4  archers
	{0, 0, 30, 150, 500},   // 5  template
	{0, 0, 30, 150, 500},   // 6  gazebo
	{0, 0, 10, 50, 200},    // 7  ace
	{0, 0, 10, 30, 150},    // 8  king
	{0, 0, 10, 30, 150},    // 9  queen
	{0, 0, 10, 20, 100},    // 10 jack
	{0, 0, 10, 20, 100},    // 11 ten
	{},                     // 12 katana
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 2, 20, 200} // 12 katana

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 10, 10, 10} // 12 katana

const (
	jid = 1 // jackpot ID
)

// Jackpot win combinations.
var Jackpot = [12][5]int{
	{0, 0, 0, 0, 0}, //  1 shogun
	{0, 0, 0, 0, 0}, //  2 geisha
	{0, 0, 0, 0, 0}, //  3 general
	{0, 0, 0, 0, 0}, //  4 archers
	{0, 0, 0, 0, 0}, //  5 template
	{0, 0, 0, 0, 0}, //  6 gazebo
	{0, 0, 0, 0, 0}, //  7 ace
	{0, 0, 0, 0, 0}, //  8 king
	{0, 0, 0, 0, 0}, //  9 queen
	{0, 0, 0, 0, 0}, // 10 jack
	{0, 0, 0, 0, 0}, // 11 ten
	{0, 0, 0, 0, 0}, // 12 katana
}

// Bet lines
var BetLines = slot.BetLinesNvm20

type Game struct {
	slot.Slot5x3 `yaml:",inline"`
	// free spin number
	FS int `json:"fs,omitempty" yaml:"fs,omitempty" xml:"fs,omitempty"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slot5x3: slot.Slot5x3{
			Sel: len(BetLines),
			Bet: 1,
		},
		FS: 0,
	}
}

const wild, scat = 1, 12

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	g.ScanLined(screen, wins)
	g.ScanScatters(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen slot.Screen, wins *slot.Wins) {
	var reelwild [5]bool
	if g.FS > 0 {
		var x, y slot.Pos
		for x = 1; x <= 5; x++ {
			for y = 1; y <= 3; y++ {
				if screen.At(x, y) == wild {
					reelwild[x-1] = true
					break
				}
			}
		}
	}

	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]

		var numw, numl slot.Pos = 0, 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = screen.Pos(x, line)
			if sx == wild || reelwild[x-1] {
				if syml == 0 {
					numw = x
				}
			} else if syml == 0 && sx != scat {
				syml = sx
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		var payw, payl float64
		if numw > 0 {
			payw = LinePay[wild-1][numw-1]
		}
		if numl > 0 && syml > 0 {
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
	if count := screen.ScatNum(scat); count >= 3 {
		var pay, fs = ScatPay[count-1], ScatFreespin[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
			Free: fs,
		})
	}
}

func (g *Game) Spin(screen slot.Screen, mrtp float64) {
	if g.FS == 0 {
		var reels, _ = slot.FindReels(ReelsMap, mrtp)
		screen.Spin(reels)
	} else {
		screen.Spin(&ReelsBon)
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

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
