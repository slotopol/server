package katana

import (
	"github.com/slotopol/server/game"
)

// reels lengths [36, 36, 36, 36, 36], total reshuffles 60466176
// symbols: 61.542(lined) + 1.495(scatter) = 63.036564%
// free spins 3076380, q = 0.050878, sq = 1/(1-q) = 1.053605
// RTP = 63.037(sym) + 1.0536*477.78(fg) = 87.344837%
var ReelsReg87 = game.Reels5x{
	{7, 11, 8, 9, 11, 1, 8, 4, 12, 3, 6, 7, 2, 9, 8, 10, 11, 9, 6, 5, 2, 9, 4, 10, 7, 2, 3, 10, 5, 8, 3, 10, 11, 5, 6, 4},
	{10, 2, 6, 8, 3, 9, 4, 1, 11, 7, 5, 2, 4, 11, 5, 4, 10, 6, 9, 3, 7, 8, 11, 5, 9, 11, 8, 9, 6, 10, 7, 2, 8, 10, 12, 3},
	{3, 4, 6, 11, 5, 2, 7, 8, 3, 10, 12, 7, 9, 5, 11, 7, 2, 10, 6, 4, 1, 8, 10, 9, 8, 6, 9, 4, 11, 8, 3, 10, 2, 11, 9, 5},
	{2, 5, 9, 11, 10, 6, 2, 5, 8, 10, 6, 5, 7, 6, 11, 7, 9, 8, 4, 3, 10, 7, 9, 8, 11, 9, 4, 3, 8, 12, 11, 2, 10, 1, 4, 3},
	{7, 3, 10, 11, 6, 7, 1, 6, 8, 10, 4, 6, 3, 9, 2, 5, 4, 8, 11, 5, 10, 9, 4, 8, 7, 12, 2, 3, 8, 9, 5, 11, 2, 9, 10, 11},
}

// reels lengths [35, 36, 36, 36, 35], total reshuffles 57153600
// symbols: 62.53(lined) + 1.5511(scatter) = 64.080772%
// free spins 3003750, q = 0.052556, sq = 1/(1-q) = 1.055471
// RTP = 64.081(sym) + 1.0555*477.78(fg) = 89.190777%
var ReelsReg89 = game.Reels5x{
	{2, 7, 10, 9, 3, 6, 9, 10, 5, 4, 7, 2, 10, 4, 6, 9, 10, 7, 8, 11, 2, 5, 12, 8, 6, 3, 11, 1, 4, 11, 8, 9, 5, 11, 3},
	{10, 2, 6, 8, 3, 9, 4, 1, 11, 7, 5, 2, 4, 11, 5, 4, 10, 6, 9, 3, 7, 8, 11, 5, 9, 11, 8, 9, 6, 10, 7, 2, 8, 10, 12, 3},
	{3, 4, 6, 11, 5, 2, 7, 8, 3, 10, 12, 7, 9, 5, 11, 7, 2, 10, 6, 4, 1, 8, 10, 9, 8, 6, 9, 4, 11, 8, 3, 10, 2, 11, 9, 5},
	{2, 5, 9, 11, 10, 6, 2, 5, 8, 10, 6, 5, 7, 6, 11, 7, 9, 8, 4, 3, 10, 7, 9, 8, 11, 9, 4, 3, 8, 12, 11, 2, 10, 1, 4, 3},
	{6, 3, 2, 8, 6, 10, 9, 7, 4, 2, 11, 10, 4, 11, 8, 7, 9, 4, 3, 8, 5, 12, 6, 3, 11, 10, 9, 2, 5, 1, 11, 10, 5, 9, 7},
}

// reels lengths [35, 36, 35, 36, 35], total reshuffles 55566000
// symbols: 63.65(lined) + 1.5799(scatter) = 65.230227%
// free spins 2967840, q = 0.053411, sq = 1/(1-q) = 1.056425
// RTP = 65.23(sym) + 1.0564*477.78(fg) = 90.748893%
var ReelsReg91 = game.Reels5x{
	{2, 7, 10, 9, 3, 6, 9, 10, 5, 4, 7, 2, 10, 4, 6, 9, 10, 7, 8, 11, 2, 5, 12, 8, 6, 3, 11, 1, 4, 11, 8, 9, 5, 11, 3},
	{10, 2, 6, 8, 3, 9, 4, 1, 11, 7, 5, 2, 4, 11, 5, 4, 10, 6, 9, 3, 7, 8, 11, 5, 9, 11, 8, 9, 6, 10, 7, 2, 8, 10, 12, 3},
	{9, 12, 6, 7, 4, 6, 3, 8, 9, 3, 1, 10, 4, 2, 9, 3, 2, 8, 5, 10, 6, 11, 9, 7, 11, 5, 4, 11, 7, 5, 10, 8, 11, 2, 10},
	{2, 5, 9, 11, 10, 6, 2, 5, 8, 10, 6, 5, 7, 6, 11, 7, 9, 8, 4, 3, 10, 7, 9, 8, 11, 9, 4, 3, 8, 12, 11, 2, 10, 1, 4, 3},
	{6, 3, 2, 8, 6, 10, 9, 7, 4, 2, 11, 10, 4, 11, 8, 7, 9, 4, 3, 8, 5, 12, 6, 3, 11, 10, 9, 2, 5, 1, 11, 10, 5, 9, 7},
}

// reels lengths [35, 35, 35, 35, 35], total reshuffles 52521875
// symbols: 65.372(lined) + 1.6389(scatter) = 67.010688%
// free spins 2896830, q = 0.055155, sq = 1/(1-q) = 1.058374
// RTP = 67.011(sym) + 1.0584*477.78(fg) = 93.362435%
var ReelsReg93 = game.Reels5x{
	{2, 7, 10, 9, 3, 6, 9, 10, 5, 4, 7, 2, 10, 4, 6, 9, 10, 7, 8, 11, 2, 5, 12, 8, 6, 3, 11, 1, 4, 11, 8, 9, 5, 11, 3},
	{9, 7, 2, 5, 11, 7, 3, 9, 5, 2, 9, 10, 6, 3, 2, 9, 12, 10, 11, 1, 5, 11, 4, 8, 10, 4, 8, 7, 11, 4, 6, 8, 10, 6, 3},
	{9, 12, 6, 7, 4, 6, 3, 8, 9, 3, 1, 10, 4, 2, 9, 3, 2, 8, 5, 10, 6, 11, 9, 7, 11, 5, 4, 11, 7, 5, 10, 8, 11, 2, 10},
	{11, 9, 4, 7, 5, 2, 3, 6, 4, 11, 8, 6, 11, 7, 8, 4, 12, 6, 7, 9, 3, 11, 5, 10, 9, 8, 2, 10, 3, 1, 10, 9, 2, 10, 5},
	{6, 3, 2, 8, 6, 10, 9, 7, 4, 2, 11, 10, 4, 11, 8, 7, 9, 4, 3, 8, 5, 12, 6, 3, 11, 10, 9, 2, 5, 1, 11, 10, 5, 9, 7},
}

// reels lengths [34, 35, 35, 35, 34], total reshuffles 49563500
// symbols: 66.438(lined) + 1.7024(scatter) = 68.140123%
// free spins 2826360, q = 0.057025, sq = 1/(1-q) = 1.060474
// RTP = 68.14(sym) + 1.0605*477.78(fg) = 95.385457%
var ReelsReg95 = game.Reels5x{
	{11, 3, 6, 1, 5, 11, 3, 4, 6, 7, 9, 11, 2, 6, 5, 10, 8, 9, 10, 7, 11, 4, 10, 9, 5, 12, 8, 2, 4, 3, 10, 7, 8, 2},
	{9, 7, 2, 5, 11, 7, 3, 9, 5, 2, 9, 10, 6, 3, 2, 9, 12, 10, 11, 1, 5, 11, 4, 8, 10, 4, 8, 7, 11, 4, 6, 8, 10, 6, 3},
	{9, 12, 6, 7, 4, 6, 3, 8, 9, 3, 1, 10, 4, 2, 9, 3, 2, 8, 5, 10, 6, 11, 9, 7, 11, 5, 4, 11, 7, 5, 10, 8, 11, 2, 10},
	{11, 9, 4, 7, 5, 2, 3, 6, 4, 11, 8, 6, 11, 7, 8, 4, 12, 6, 7, 9, 3, 11, 5, 10, 9, 8, 2, 10, 3, 1, 10, 9, 2, 10, 5},
	{4, 3, 10, 2, 6, 8, 7, 2, 9, 11, 8, 4, 2, 6, 5, 11, 10, 8, 5, 3, 7, 9, 11, 5, 7, 1, 9, 11, 10, 4, 12, 3, 10, 6},
}

// reels lengths [34, 35, 34, 35, 34], total reshuffles 48147400
// symbols: 67.657(lined) + 1.7349(scatter) = 69.391506%
// free spins 2791530, q = 0.057979, sq = 1/(1-q) = 1.061547
// RTP = 69.392(sym) + 1.0615*477.78(fg) = 97.092546%
var ReelsReg97 = game.Reels5x{
	{11, 3, 6, 1, 5, 11, 3, 4, 6, 7, 9, 11, 2, 6, 5, 10, 8, 9, 10, 7, 11, 4, 10, 9, 5, 12, 8, 2, 4, 3, 10, 7, 8, 2},
	{9, 7, 2, 5, 11, 7, 3, 9, 5, 2, 9, 10, 6, 3, 2, 9, 12, 10, 11, 1, 5, 11, 4, 8, 10, 4, 8, 7, 11, 4, 6, 8, 10, 6, 3},
	{7, 1, 2, 6, 10, 7, 11, 9, 10, 4, 2, 9, 11, 10, 3, 8, 5, 12, 7, 6, 3, 11, 5, 3, 6, 9, 8, 4, 5, 2, 10, 8, 11, 4},
	{11, 9, 4, 7, 5, 2, 3, 6, 4, 11, 8, 6, 11, 7, 8, 4, 12, 6, 7, 9, 3, 11, 5, 10, 9, 8, 2, 10, 3, 1, 10, 9, 2, 10, 5},
	{4, 3, 10, 2, 6, 8, 7, 2, 9, 11, 8, 4, 2, 6, 5, 11, 10, 8, 5, 3, 7, 9, 11, 5, 7, 1, 9, 11, 10, 4, 12, 3, 10, 6},
}

// reels lengths [34, 34, 34, 34, 34], total reshuffles 45435424
// symbols: 69.539(lined) + 1.8018(scatter) = 71.340415%
// free spins 2722680, q = 0.059924, sq = 1/(1-q) = 1.063744
// RTP = 71.34(sym) + 1.0637*477.78(fg) = 99.970895%
var ReelsReg100 = game.Reels5x{
	{11, 3, 6, 1, 5, 11, 3, 4, 6, 7, 9, 11, 2, 6, 5, 10, 8, 9, 10, 7, 11, 4, 10, 9, 5, 12, 8, 2, 4, 3, 10, 7, 8, 2},
	{7, 10, 2, 3, 11, 4, 9, 3, 6, 5, 1, 11, 8, 7, 2, 10, 4, 6, 9, 10, 4, 11, 2, 8, 6, 5, 10, 7, 9, 11, 5, 3, 12, 8},
	{7, 1, 2, 6, 10, 7, 11, 9, 10, 4, 2, 9, 11, 10, 3, 8, 5, 12, 7, 6, 3, 11, 5, 3, 6, 9, 8, 4, 5, 2, 10, 8, 11, 4},
	{4, 10, 7, 4, 10, 9, 5, 11, 6, 10, 2, 8, 3, 11, 6, 2, 7, 11, 10, 3, 9, 8, 6, 1, 4, 5, 12, 3, 2, 5, 7, 8, 11, 9},
	{4, 3, 10, 2, 6, 8, 7, 2, 9, 11, 8, 4, 2, 6, 5, 11, 10, 8, 5, 3, 7, 9, 11, 5, 7, 1, 9, 11, 10, 4, 12, 3, 10, 6},
}

// reels lengths [33, 33, 33, 33, 33], total reshuffles 39135393
// symbols: 242.27(lined) + 20.962(scatter) = 263.231929%
// free spins 17573760, q = 0.44905, sq = 1/(1-q) = 1.815048
// RTP = sq*rtp(sym) = 1.815*263.23 = 477.778515%
var ReelsBon = game.Reels5x{
	{7, 3, 1, 4, 5, 6, 3, 7, 5, 4, 7, 11, 4, 8, 10, 3, 9, 12, 8, 10, 6, 9, 11, 12, 5, 2, 10, 9, 2, 11, 6, 8, 2},
	{11, 2, 6, 9, 8, 7, 4, 10, 6, 4, 1, 8, 6, 3, 5, 9, 2, 5, 10, 7, 3, 8, 10, 9, 11, 7, 5, 4, 12, 11, 3, 2, 12},
	{12, 11, 5, 7, 1, 4, 5, 7, 6, 4, 8, 10, 6, 12, 9, 4, 11, 2, 8, 10, 3, 5, 10, 3, 11, 2, 9, 8, 7, 3, 9, 2, 6},
	{11, 7, 2, 6, 7, 1, 2, 4, 9, 10, 8, 9, 12, 5, 4, 8, 7, 11, 8, 3, 4, 11, 5, 10, 3, 9, 6, 3, 2, 6, 12, 5, 10},
	{11, 12, 5, 6, 9, 8, 5, 10, 3, 9, 12, 2, 4, 3, 2, 7, 6, 4, 9, 10, 4, 1, 2, 7, 8, 11, 3, 6, 7, 11, 10, 5, 8},
}

// Map with available reels.
var ReelsMap = map[string]*game.Reels5x{
	"87":  &ReelsReg87,
	"89":  &ReelsReg89,
	"91":  &ReelsReg91,
	"93":  &ReelsReg93,
	"95":  &ReelsReg95,
	"97":  &ReelsReg97,
	"100": &ReelsReg100,
	"bon": &ReelsBon,
}

// Lined payment.
var LinePay = [12][5]int{
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
	{0, 0, 0, 0, 0},        // 12 katana
}

// Scatters payment.
var ScatPay = [5]int{0, 0, 2, 20, 200} // 12 katana

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

type Game struct {
	game.Slot5x3 `yaml:",inline"`
	FS           int `json:"fs" yaml:"fs" xml:"fs"` // free spin number
}

func NewGame(ri string) *Game {
	return &Game{
		Slot5x3: game.Slot5x3{
			RI:  ri,
			BLI: "nvm10",
			SBL: game.MakeSBL(1),
			Bet: 1,
		},
		FS: 0,
	}
}

const wild, scat = 1, 12

func (g *Game) Scanner(screen game.Screen, ws *game.WinScan) {
	g.ScanLined(screen, ws)
	g.ScanScatters(screen, ws)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen game.Screen, ws *game.WinScan) {
	var reelwild [5]bool
	if g.FS > 0 {
		for x := 1; x <= 5; x++ {
			for y := 1; y <= 3; y++ {
				if screen.At(x, y) == wild {
					reelwild[x-1] = true
				}
			}
		}
	}

	var bl = game.BetLines5x[g.BLI]
	for li := g.SBL.Next(0); li != 0; li = g.SBL.Next(li) {
		var line = bl.Line(li)

		var numw, numl = 0, 5
		var syml game.Sym
		for x := 1; x <= 5; x++ {
			var sx = screen.At(x, line.At(x))
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

		var payw, payl int
		if numw > 0 {
			payw = LinePay[wild-1][numw-1]
		}
		if numl > 0 && syml > 0 {
			payl = LinePay[syml-1][numl-1]
		}
		if payl > payw {
			ws.Wins = append(ws.Wins, game.WinItem{
				Pay:  g.Bet * payl,
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyN(numl),
			})
		} else if payw > 0 {
			ws.Wins = append(ws.Wins, game.WinItem{
				Pay:  g.Bet * payw,
				Mult: 1,
				Sym:  wild,
				Num:  numw,
				Line: li,
				XY:   line.CopyN(numw),
				Jack: Jackpot[wild-1][numw-1],
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen game.Screen, ws *game.WinScan) {
	if count := screen.ScatNum(scat); count >= 3 {
		var pay, fs = ScatPay[count-1], ScatFreespin[count-1]
		var xy = game.NewLine5x()
		for x := 1; x <= 5; x++ {
			for y := 1; y <= 3; y++ {
				if screen.At(x, y) == scat {
					xy.Set(x, y)
					break
				}
			}
		}
		ws.Wins = append(ws.Wins, game.WinItem{
			Pay:  g.Bet * pay, // independent from selected lines
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   xy,
			Free: fs,
		})
	}
}

func (g *Game) Spin(screen game.Screen) {
	if g.FS == 0 {
		screen.Spin(ReelsMap[g.RI])
	} else {
		screen.Spin(&ReelsBon)
	}
}

func (g *Game) Apply(screen game.Screen, sw *game.WinScan) {
	g.Gain = sw.Gain()
	if g.FS > 0 {
		g.FS--
	}
	for _, wi := range sw.Wins {
		if wi.Free > 0 {
			g.FS += wi.Free
		}
	}
}

func (g *Game) FreeSpins() int {
	return g.FS
}
