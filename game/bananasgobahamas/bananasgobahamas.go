package bananasgobahamas

import (
	"github.com/slotopol/server/game"
)

// reels lengths [30, 43, 43, 43, 43], total reshuffles 102564030
// symbols: 43.094(lined) + 10.921(scatter) = 54.015807%
// free games 17605350, q = 0.17165, sq = 1/(1-q) = 1.207222
// free games frequency: 1/262.16
// RTP = 54.016(sym) + 0.17165*210.46(fg) = 90.142004%
var ReelsReg90 = game.Reels5x{
	{6, 11, 1, 9, 12, 8, 11, 7, 2, 11, 6, 8, 11, 4, 7, 9, 13, 10, 9, 3, 7, 5, 11, 9, 6, 11, 9, 10, 4, 7},
	{5, 8, 1, 7, 8, 11, 2, 9, 7, 12, 6, 10, 8, 3, 10, 9, 2, 8, 10, 5, 8, 11, 13, 12, 4, 10, 5, 8, 6, 3, 11, 12, 4, 8, 10, 5, 7, 10, 5, 12, 4, 9, 7},
	{11, 9, 1, 11, 4, 10, 9, 2, 8, 6, 10, 4, 11, 12, 2, 7, 8, 9, 12, 13, 9, 11, 5, 11, 3, 10, 9, 4, 7, 5, 9, 11, 2, 7, 5, 12, 7, 3, 11, 10, 5, 7, 3},
	{8, 11, 1, 9, 3, 7, 10, 4, 9, 2, 7, 6, 10, 5, 12, 11, 6, 12, 8, 5, 11, 13, 8, 6, 10, 2, 9, 6, 12, 5, 8, 7, 6, 10, 3, 6, 12, 4, 10, 5, 12, 10, 3},
	{8, 7, 1, 8, 4, 9, 5, 12, 10, 6, 11, 3, 7, 8, 2, 9, 5, 10, 11, 2, 12, 13, 5, 10, 4, 12, 6, 7, 2, 10, 8, 2, 12, 10, 3, 8, 10, 4, 12, 5, 8, 12, 2},
}

// reels lengths [30, 43, 43, 43, 43], total reshuffles 102564030
// symbols: 45.156(lined) + 10.921(scatter) = 56.077317%
// free games 17605350, q = 0.17165, sq = 1/(1-q) = 1.207222
// free games frequency: 1/262.16
// RTP = 56.077(sym) + 0.17165*210.46(fg) = 92.203515%
var ReelsReg92 = game.Reels5x{
	{6, 11, 1, 9, 12, 8, 11, 7, 2, 11, 5, 8, 11, 4, 7, 9, 13, 10, 9, 3, 7, 5, 10, 8, 6, 11, 9, 10, 4, 7},
	{5, 8, 1, 7, 8, 11, 2, 9, 10, 12, 6, 10, 8, 3, 10, 9, 2, 8, 10, 5, 8, 11, 13, 12, 4, 10, 5, 8, 6, 3, 10, 12, 4, 7, 10, 6, 8, 7, 5, 12, 4, 9, 11},
	{11, 9, 1, 11, 4, 7, 9, 2, 8, 6, 10, 4, 10, 12, 2, 7, 8, 9, 12, 13, 9, 10, 5, 11, 3, 10, 9, 4, 7, 5, 10, 11, 2, 7, 6, 12, 8, 3, 11, 9, 6, 7, 3},
	{8, 11, 1, 9, 3, 7, 11, 4, 9, 2, 7, 6, 10, 5, 12, 11, 6, 12, 8, 5, 11, 13, 8, 6, 10, 2, 9, 6, 12, 5, 8, 7, 6, 10, 3, 6, 12, 4, 10, 5, 12, 10, 3},
	{8, 7, 1, 8, 4, 9, 5, 12, 10, 6, 11, 3, 7, 8, 2, 9, 5, 10, 11, 2, 12, 13, 5, 10, 4, 12, 6, 7, 2, 10, 8, 2, 12, 10, 3, 8, 10, 4, 12, 5, 8, 12, 2},
}

// reels lengths [30, 43, 43, 43, 43], total reshuffles 102564030
// symbols: 47.432(lined) + 10.921(scatter) = 58.353031%
// free games 17605350, q = 0.17165, sq = 1/(1-q) = 1.207222
// free games frequency: 1/262.16
// RTP = 58.353(sym) + 0.17165*210.46(fg) = 94.479229%
var ReelsReg94 = game.Reels5x{
	{6, 11, 1, 9, 12, 8, 11, 7, 2, 11, 6, 8, 9, 4, 7, 9, 13, 10, 9, 3, 7, 5, 11, 9, 10, 11, 9, 10, 4, 7},
	{5, 8, 1, 7, 8, 11, 2, 9, 10, 12, 6, 11, 7, 3, 10, 9, 2, 8, 10, 5, 8, 11, 13, 12, 4, 10, 5, 8, 6, 3, 10, 12, 4, 8, 10, 6, 8, 9, 5, 12, 4, 9, 7},
	{11, 10, 1, 11, 4, 8, 9, 2, 8, 6, 10, 4, 11, 12, 2, 10, 8, 9, 12, 13, 9, 11, 5, 11, 3, 10, 9, 4, 7, 5, 9, 11, 2, 7, 6, 12, 7, 3, 11, 9, 6, 7, 3},
	{8, 11, 1, 9, 3, 7, 10, 4, 9, 2, 7, 6, 10, 5, 12, 10, 6, 12, 8, 5, 11, 13, 8, 6, 10, 2, 7, 6, 12, 5, 8, 11, 1, 9, 3, 6, 12, 4, 10, 5, 12, 10, 3},
	{8, 7, 1, 8, 4, 9, 5, 12, 10, 6, 11, 3, 7, 8, 2, 9, 5, 10, 11, 2, 12, 13, 5, 10, 4, 12, 6, 7, 2, 10, 8, 2, 12, 10, 3, 8, 10, 4, 12, 5, 8, 12, 2},
}

// reels lengths [30, 43, 43, 43, 43], total reshuffles 102564030
// symbols: 48.317(lined) + 10.921(scatter) = 59.238329%
// free games 17605350, q = 0.17165, sq = 1/(1-q) = 1.207222
// free games frequency: 1/262.16
// RTP = 59.238(sym) + 0.17165*210.46(fg) = 95.364527%
var ReelsReg95 = game.Reels5x{
	{6, 11, 1, 9, 12, 8, 11, 7, 2, 11, 6, 8, 9, 4, 7, 9, 13, 10, 9, 3, 7, 5, 11, 8, 10, 11, 5, 10, 4, 7},
	{5, 8, 1, 7, 8, 11, 2, 9, 10, 12, 6, 11, 7, 3, 10, 9, 2, 8, 10, 5, 8, 11, 13, 12, 4, 10, 5, 8, 6, 3, 9, 12, 4, 8, 10, 6, 8, 10, 5, 12, 4, 9, 7},
	{11, 9, 1, 11, 4, 8, 9, 2, 8, 6, 10, 4, 7, 12, 2, 10, 11, 9, 12, 13, 9, 11, 5, 11, 3, 10, 9, 4, 7, 5, 10, 11, 2, 7, 6, 12, 7, 3, 11, 10, 6, 7, 3},
	{8, 11, 1, 9, 3, 7, 10, 4, 9, 2, 7, 6, 9, 5, 12, 10, 6, 12, 8, 5, 11, 13, 8, 6, 10, 2, 7, 6, 12, 5, 8, 11, 1, 9, 3, 6, 12, 4, 10, 5, 12, 10, 3},
	{8, 7, 1, 8, 4, 9, 5, 12, 10, 6, 11, 3, 7, 8, 2, 9, 5, 10, 11, 2, 12, 13, 5, 10, 4, 12, 6, 7, 2, 10, 8, 2, 12, 10, 3, 8, 10, 4, 12, 5, 8, 12, 2},
}

// reels lengths [30, 43, 43, 43, 43], total reshuffles 102564030
// symbols: 49.364(lined) + 10.921(scatter) = 60.285580%
// free games 17605350, q = 0.17165, sq = 1/(1-q) = 1.207222
// free games frequency: 1/262.16
// RTP = 60.286(sym) + 0.17165*210.46(fg) = 96.411777%
var ReelsReg96 = game.Reels5x{
	{6, 11, 1, 9, 12, 8, 11, 7, 2, 11, 6, 8, 9, 4, 7, 9, 13, 10, 9, 3, 7, 5, 11, 8, 10, 11, 5, 10, 4, 7},
	{5, 8, 1, 7, 8, 11, 2, 9, 10, 12, 6, 11, 7, 3, 10, 9, 2, 8, 10, 5, 8, 11, 13, 12, 4, 10, 5, 8, 6, 3, 9, 12, 4, 8, 10, 6, 8, 10, 5, 12, 4, 9, 7},
	{11, 9, 1, 11, 4, 8, 9, 2, 8, 6, 10, 4, 7, 12, 2, 10, 8, 9, 12, 13, 9, 11, 5, 11, 3, 10, 9, 4, 7, 5, 10, 11, 2, 7, 6, 12, 7, 3, 11, 10, 6, 7, 3},
	{8, 11, 1, 9, 3, 7, 10, 4, 9, 2, 7, 6, 10, 5, 12, 10, 6, 12, 8, 5, 11, 13, 8, 6, 10, 2, 7, 6, 12, 5, 8, 11, 1, 9, 3, 6, 12, 4, 10, 5, 12, 10, 3},
	{8, 7, 1, 8, 4, 9, 5, 12, 10, 6, 11, 3, 7, 8, 2, 9, 5, 10, 11, 2, 12, 13, 5, 10, 4, 12, 6, 7, 2, 10, 8, 2, 12, 10, 3, 8, 10, 4, 12, 5, 8, 12, 2},
}

// reels lengths [30, 43, 43, 43, 43], total reshuffles 102564030
// symbols: 50.334(lined) + 10.921(scatter) = 61.255553%
// free games 17605350, q = 0.17165, sq = 1/(1-q) = 1.207222
// free games frequency: 1/262.16
// RTP = 61.256(sym) + 0.17165*210.46(fg) = 97.381751%
var ReelsReg97 = game.Reels5x{
	{6, 11, 1, 9, 12, 8, 11, 7, 2, 11, 6, 8, 9, 4, 7, 9, 13, 10, 9, 3, 7, 5, 11, 8, 10, 11, 5, 10, 4, 7},
	{5, 8, 1, 7, 8, 11, 2, 9, 10, 12, 6, 11, 7, 3, 10, 9, 2, 8, 10, 5, 8, 11, 13, 12, 4, 10, 5, 8, 6, 3, 9, 12, 4, 8, 10, 6, 8, 10, 5, 12, 4, 9, 7},
	{11, 9, 1, 11, 4, 8, 9, 2, 8, 6, 10, 4, 7, 12, 2, 10, 11, 9, 12, 13, 9, 11, 5, 11, 3, 10, 9, 4, 7, 5, 10, 8, 2, 7, 6, 12, 7, 3, 11, 10, 6, 7, 3},
	{8, 11, 1, 9, 3, 7, 10, 4, 9, 2, 7, 6, 9, 5, 12, 10, 6, 12, 10, 5, 11, 13, 8, 6, 10, 2, 7, 6, 12, 5, 8, 11, 1, 9, 3, 6, 12, 4, 10, 5, 12, 10, 3},
	{8, 7, 1, 8, 4, 9, 5, 12, 9, 6, 11, 3, 7, 12, 2, 9, 5, 10, 11, 2, 12, 13, 5, 10, 2, 12, 6, 9, 3, 12, 4, 8, 7, 1, 8, 4, 10, 2, 12, 5, 8, 9, 2},
}

// reels lengths [30, 40, 40, 40, 40], total reshuffles 76800000
// symbols: 52.632(lined) + 22.694(scatter) = 75.325924%
// free games 36469440, q = 0.47486, sq = 1/(1-q) = 1.904263
// free games frequency: 1/94.764
// RTP = 75.326(sym) + 0.47486*210.46(fg) = 175.266200%
var ReelsReg175 = game.Reels5x{
	{6, 11, 1, 9, 12, 8, 11, 7, 2, 11, 6, 8, 9, 4, 7, 9, 13, 10, 9, 3, 7, 5, 11, 8, 10, 11, 5, 10, 4, 7},
	{5, 8, 1, 7, 13, 11, 2, 9, 4, 12, 6, 11, 7, 3, 10, 9, 2, 8, 10, 5, 9, 11, 13, 12, 4, 10, 5, 8, 6, 3, 9, 12, 4, 8, 12, 6, 8, 10, 5, 7},
	{11, 9, 1, 11, 4, 8, 9, 2, 8, 6, 10, 4, 7, 12, 2, 10, 11, 6, 12, 13, 9, 11, 5, 10, 3, 10, 9, 4, 7, 5, 10, 8, 2, 7, 6, 12, 3, 11, 7, 3},
	{8, 11, 1, 9, 3, 7, 10, 5, 9, 2, 7, 6, 9, 5, 12, 10, 13, 12, 10, 5, 11, 13, 8, 6, 10, 2, 7, 6, 12, 5, 8, 11, 1, 12, 9, 6, 12, 4, 10, 3},
	{8, 7, 1, 8, 4, 9, 5, 12, 9, 6, 11, 3, 7, 12, 2, 9, 5, 10, 11, 2, 9, 13, 5, 10, 2, 12, 6, 9, 3, 12, 4, 8, 7, 1, 8, 4, 10, 2, 12, 5},
}

// reels lengths [35, 35, 36, 36, 36], total reshuffles 57153600
// symbols: 79.548(lined) + 24.435(scatter) = 103.982479%
// free games 28915785, q = 0.50593, sq = 1/(1-q) = 2.024009
// free games frequency: 1/88.945
// RTP = sq*rtp(sym) = 2.024*103.98 = 210.461503%
var ReelsBon = game.Reels5x{
	{6, 11, 1, 9, 12, 5, 11, 8, 2, 7, 6, 10, 12, 4, 7, 9, 13, 10, 9, 3, 7, 5, 10, 9, 2, 8, 11, 10, 3, 11, 6, 8, 4, 11, 9},
	{5, 8, 1, 7, 8, 10, 2, 9, 4, 11, 5, 8, 1, 7, 8, 11, 10, 12, 4, 10, 7, 3, 11, 6, 10, 4, 6, 11, 13, 12, 5, 9, 6, 11, 9},
	{11, 9, 1, 11, 4, 7, 8, 2, 6, 11, 9, 1, 11, 4, 7, 12, 10, 9, 8, 2, 12, 10, 5, 12, 3, 6, 7, 12, 13, 9, 10, 5, 3, 7, 5, 10},
	{8, 11, 1, 9, 3, 7, 10, 4, 7, 2, 8, 6, 10, 4, 5, 11, 10, 8, 6, 12, 11, 5, 7, 9, 10, 2, 5, 11, 13, 8, 6, 9, 3, 12, 7, 4},
	{8, 7, 1, 8, 4, 9, 12, 13, 5, 10, 11, 3, 7, 8, 2, 12, 13, 5, 10, 6, 7, 9, 6, 7, 4, 9, 2, 12, 13, 5, 10, 4, 11, 6, 3, 11},
}

// Map with available reels.
var ReelsMap = map[string]*game.Reels5x{
	"90":  &ReelsReg90,
	"92":  &ReelsReg92,
	"94":  &ReelsReg94,
	"95":  &ReelsReg95,
	"96":  &ReelsReg96,
	"97":  &ReelsReg97,
	"175": &ReelsReg175,
	"bon": &ReelsBon,
}

// Lined payment.
var LinePay = [13][5]float64{
	{0, 10, 250, 2500, 9000}, //  1 banana
	{0, 2, 30, 120, 800},     //  2 strawberry
	{0, 2, 30, 120, 800},     //  3 water
	{0, 0, 20, 100, 400},     //  4 pineapple
	{0, 0, 20, 70, 250},      //  5 mango
	{0, 0, 20, 70, 250},      //  6 coconu
	{0, 0, 10, 50, 120},      //  7 ace
	{0, 0, 10, 50, 120},      //  8 king
	{0, 0, 4, 30, 100},       //  9 queen
	{0, 0, 4, 30, 100},       // 10 jack
	{0, 0, 4, 30, 100},       // 11 ten
	{0, 2, 4, 30, 100},       // 12 nine
	{0, 0, 0, 0, 0},          // 13 suitcase
}

// Scatters payment.
var ScatPay = [5]float64{0, 2, 4, 20, 500} // 13 suitcase

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 45, 45, 45} // 13 suitcase

const (
	jid = 1 // jackpot ID
)

// Jackpot win combinations.
var Jackpot = [13][5]int{
	{0, 0, 0, 0, 0}, //  1 banana
	{0, 0, 0, 0, 0}, //  2 strawberry
	{0, 0, 0, 0, 0}, //  3 water
	{0, 0, 0, 0, 0}, //  4 pineapple
	{0, 0, 0, 0, 0}, //  5 mango
	{0, 0, 0, 0, 0}, //  6 coconu
	{0, 0, 0, 0, 0}, //  7 ace
	{0, 0, 0, 0, 0}, //  8 king
	{0, 0, 0, 0, 0}, //  9 queen
	{0, 0, 0, 0, 0}, // 10 jack
	{0, 0, 0, 0, 0}, // 11 ten
	{0, 0, 0, 0, 0}, // 12 nine
	{0, 0, 0, 0, 0}, // 13 suitcase
}

type Game struct {
	game.Slot5x3 `yaml:",inline"`
	// free spin number
	FS int `json:"fs,omitempty" yaml:"fs,omitempty" xml:"fs,omitempty"`
}

func NewGame(rd string) *Game {
	return &Game{
		Slot5x3: game.Slot5x3{
			RD:  rd,
			SBL: game.MakeBitNum(5),
			Bet: 1,
		},
		FS: 0,
	}
}

const wild, scat = 1, 13

var bl = game.BetLinesNvm10

func (g *Game) Scanner(screen game.Screen, wins *game.Wins) {
	g.ScanLined(screen, wins)
	g.ScanScatters(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen game.Screen, wins *game.Wins) {
	for li := g.SBL.Next(0); li != 0; li = g.SBL.Next(li) {
		var line = bl.Line(li)

		var numw, numl = 0, 5
		var syml game.Sym
		var mw float64 = 1 // mult wild
		for x := 1; x <= 5; x++ {
			var sx = screen.Pos(x, line)
			if sx == wild {
				if syml == 0 {
					numw = x
				}
				mw = 2
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
		if payl*mw > payw {
			*wins = append(*wins, game.WinItem{
				Pay:  g.Bet * payl,
				Mult: mw,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyL(numl),
			})
		} else if payw > 0 {
			*wins = append(*wins, game.WinItem{
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
func (g *Game) ScanScatters(screen game.Screen, wins *game.Wins) {
	if count := screen.ScatNum(scat); count >= 2 {
		var pay, fs = ScatPay[count-1], ScatFreespin[count-1]
		*wins = append(*wins, game.WinItem{
			Pay:  g.Bet * float64(g.SBL.Num()) * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
			Free: fs,
		})
	}
}

func (g *Game) Spin(screen game.Screen) {
	if g.FS == 0 {
		screen.Spin(ReelsMap[g.RD])
	} else {
		screen.Spin(&ReelsBon)
	}
}

func (g *Game) Apply(screen game.Screen, wins game.Wins) {
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

func (g *Game) SetLines(sbl game.Bitset) error {
	var mask game.Bitset = (1<<len(bl) - 1) << 1
	if sbl == 0 {
		return game.ErrNoLineset
	}
	if sbl&^mask != 0 {
		return game.ErrLinesetOut
	}
	if g.FreeSpins() > 0 {
		return game.ErrNoFeature
	}
	g.SBL = sbl
	return nil
}

func (g *Game) SetReels(rd string) error {
	if _, ok := ReelsMap[rd]; !ok {
		return game.ErrNoReels
	}
	g.RD = rd
	return nil
}
