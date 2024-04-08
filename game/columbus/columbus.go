package columbus

import (
	"github.com/slotopol/server/game"
)

// reels lengths [32, 28, 32, 28, 32], total reshuffles 25690112
// symbols: 60.236(lined) + 0(scatter) = 60.236366%
// free spins 2540160, q = 0.098877, sq = 1/(1-q) = 1.109726
// RTP = 60.236(sym) + 0.098877*251.43(fg) = 85.096797%
var ReelsReg85 = game.Reels5x{
	{1, 5, 7, 9, 8, 2, 10, 6, 4, 10, 6, 3, 4, 8, 9, 10, 7, 8, 9, 7, 8, 9, 3, 8, 9, 7, 5, 9, 7, 2, 5, 9},
	{1, 7, 6, 3, 7, 4, 6, 9, 3, 8, 5, 6, 8, 9, 6, 8, 5, 2, 8, 5, 8, 6, 9, 8, 4, 5, 8, 6},
	{1, 6, 8, 4, 7, 5, 9, 6, 7, 9, 5, 6, 8, 7, 5, 3, 9, 7, 8, 7, 6, 9, 7, 5, 8, 4, 10, 7, 9, 2, 10, 9},
	{1, 6, 9, 8, 6, 7, 9, 4, 6, 8, 9, 6, 8, 7, 3, 6, 9, 5, 7, 6, 8, 9, 2, 9, 7, 6, 5, 8},
	{10, 6, 5, 9, 6, 10, 5, 6, 4, 1, 8, 6, 1, 8, 9, 3, 7, 2, 4, 3, 9, 4, 2, 9, 3, 4, 5, 3, 9, 4, 7, 5},
}

// reels lengths [32, 28, 32, 28, 32], total reshuffles 25690112
// symbols: 63.094(lined) + 0(scatter) = 63.094159%
// free spins 2540160, q = 0.098877, sq = 1/(1-q) = 1.109726
// RTP = 63.094(sym) + 0.098877*251.43(fg) = 87.954589%
var ReelsReg88 = game.Reels5x{
	{1, 5, 7, 9, 8, 2, 10, 6, 4, 10, 6, 3, 4, 8, 9, 10, 7, 8, 9, 7, 8, 9, 3, 8, 9, 7, 5, 9, 7, 2, 5, 9},
	{1, 7, 6, 3, 7, 4, 6, 9, 3, 8, 5, 6, 8, 9, 6, 8, 6, 2, 8, 5, 8, 6, 9, 8, 4, 5, 8, 6},
	{1, 6, 8, 4, 7, 5, 9, 6, 7, 9, 5, 6, 8, 7, 5, 3, 9, 7, 8, 7, 5, 9, 7, 5, 8, 4, 10, 7, 9, 2, 10, 9},
	{1, 6, 9, 8, 6, 7, 9, 4, 6, 8, 1, 6, 8, 7, 3, 6, 9, 5, 7, 6, 8, 9, 2, 9, 7, 6, 5, 8},
	{10, 6, 5, 9, 6, 10, 5, 6, 4, 1, 8, 6, 1, 8, 9, 3, 7, 2, 4, 3, 9, 4, 2, 9, 3, 4, 5, 3, 9, 4, 7, 5},
}

// reels lengths [32, 28, 32, 28, 32], total reshuffles 25690112
// symbols: 65.101(lined) + 0(scatter) = 65.101001%
// free spins 2540160, q = 0.098877, sq = 1/(1-q) = 1.109726
// RTP = 65.101(sym) + 0.098877*251.43(fg) = 89.961431%
var ReelsReg90 = game.Reels5x{
	{1, 5, 7, 9, 8, 2, 10, 6, 4, 10, 6, 3, 4, 8, 9, 10, 7, 8, 9, 7, 8, 9, 3, 8, 9, 7, 5, 9, 7, 2, 5, 9},
	{1, 7, 6, 3, 7, 4, 6, 9, 3, 8, 5, 6, 8, 9, 6, 8, 6, 2, 8, 5, 8, 6, 9, 8, 4, 5, 8, 6},
	{1, 6, 8, 4, 7, 5, 9, 6, 7, 9, 5, 6, 8, 7, 5, 3, 9, 7, 8, 7, 5, 9, 7, 5, 8, 4, 10, 7, 9, 2, 10, 9},
	{1, 6, 9, 5, 6, 7, 9, 4, 6, 8, 1, 6, 8, 7, 3, 6, 1, 5, 7, 6, 8, 9, 2, 9, 7, 6, 5, 8},
	{10, 6, 5, 9, 6, 10, 5, 6, 4, 9, 8, 6, 1, 8, 9, 3, 7, 2, 4, 3, 9, 4, 2, 9, 3, 4, 8, 3, 9, 4, 7, 5},
}

// reels lengths [32, 28, 32, 28, 32], total reshuffles 25690112
// symbols: 67.18(lined) + 0(scatter) = 67.180478%
// free spins 2540160, q = 0.098877, sq = 1/(1-q) = 1.109726
// RTP = 67.18(sym) + 0.098877*251.43(fg) = 92.040908%
var ReelsReg92 = game.Reels5x{
	{1, 5, 7, 9, 8, 2, 10, 6, 4, 10, 6, 3, 4, 8, 9, 10, 7, 8, 9, 7, 8, 9, 3, 8, 9, 7, 5, 9, 7, 2, 5, 9},
	{1, 7, 6, 3, 7, 4, 6, 9, 3, 8, 5, 6, 8, 9, 6, 8, 5, 2, 8, 5, 8, 6, 9, 8, 4, 5, 8, 6},
	{1, 6, 8, 4, 7, 5, 9, 6, 7, 9, 5, 6, 8, 7, 5, 3, 9, 6, 8, 7, 6, 9, 7, 5, 8, 4, 10, 7, 9, 2, 10, 9},
	{1, 6, 9, 8, 6, 7, 9, 4, 6, 8, 9, 5, 8, 7, 3, 6, 9, 5, 7, 1, 8, 9, 2, 9, 7, 6, 5, 8},
	{10, 6, 1, 9, 6, 10, 5, 6, 4, 1, 8, 6, 1, 8, 9, 3, 7, 2, 6, 3, 9, 4, 2, 9, 3, 4, 5, 3, 9, 4, 7, 5},
}

// reels lengths [32, 28, 32, 28, 32], total reshuffles 25690112
// symbols: 69.142(lined) + 0(scatter) = 69.142322%
// free spins 2540160, q = 0.098877, sq = 1/(1-q) = 1.109726
// RTP = 69.142(sym) + 0.098877*251.43(fg) = 94.002752%
var ReelsReg94 = game.Reels5x{
	{1, 5, 7, 9, 8, 2, 10, 6, 4, 10, 6, 3, 4, 8, 9, 10, 7, 8, 9, 7, 8, 9, 3, 8, 9, 7, 5, 9, 7, 2, 5, 9},
	{1, 7, 6, 3, 7, 4, 6, 9, 3, 8, 5, 6, 8, 9, 6, 8, 5, 2, 8, 5, 8, 6, 9, 8, 4, 5, 8, 6},
	{1, 6, 8, 4, 7, 5, 9, 6, 7, 9, 5, 6, 8, 7, 5, 3, 9, 6, 8, 7, 6, 9, 7, 5, 8, 4, 10, 7, 9, 2, 10, 9},
	{1, 6, 9, 8, 1, 7, 9, 4, 6, 8, 9, 3, 5, 2, 6, 3, 9, 5, 7, 1, 8, 9, 2, 9, 7, 4, 5, 8},
	{10, 6, 8, 9, 6, 10, 5, 6, 4, 1, 8, 7, 1, 8, 9, 3, 7, 2, 6, 8, 9, 4, 7, 9, 3, 8, 5, 3, 9, 4, 7, 5},
}

// reels lengths [32, 28, 32, 28, 32], total reshuffles 25690112
// symbols: 70.162(lined) + 0(scatter) = 70.162364%
// free spins 2540160, q = 0.098877, sq = 1/(1-q) = 1.109726
// RTP = 70.162(sym) + 0.098877*251.43(fg) = 95.022795%
var ReelsReg95 = game.Reels5x{
	{1, 5, 7, 9, 8, 2, 10, 6, 4, 10, 6, 3, 4, 8, 9, 10, 7, 8, 9, 7, 8, 9, 3, 8, 9, 7, 5, 9, 7, 2, 5, 9},
	{1, 7, 6, 3, 7, 4, 6, 9, 3, 8, 5, 6, 2, 9, 6, 8, 5, 2, 8, 5, 8, 6, 9, 8, 4, 5, 8, 6},
	{1, 6, 8, 4, 7, 5, 9, 6, 7, 9, 5, 6, 8, 7, 5, 3, 9, 6, 8, 7, 6, 9, 7, 5, 8, 4, 10, 7, 9, 2, 10, 9},
	{1, 6, 4, 8, 1, 7, 9, 4, 6, 8, 9, 3, 5, 2, 6, 3, 9, 5, 7, 1, 8, 9, 2, 5, 7, 4, 5, 8},
	{10, 6, 8, 9, 6, 10, 5, 6, 4, 1, 8, 7, 1, 8, 9, 3, 7, 2, 6, 8, 9, 4, 7, 9, 3, 7, 5, 3, 9, 4, 7, 5},
}

// reels lengths [32, 28, 32, 28, 32], total reshuffles 25690112
// symbols: 71.226(lined) + 0(scatter) = 71.225614%
// free spins 2540160, q = 0.098877, sq = 1/(1-q) = 1.109726
// RTP = 71.226(sym) + 0.098877*251.43(fg) = 96.086044%
var ReelsReg96 = game.Reels5x{
	{1, 5, 7, 9, 8, 2, 10, 6, 4, 10, 6, 3, 4, 8, 9, 10, 7, 8, 9, 7, 8, 9, 3, 8, 9, 7, 5, 9, 7, 2, 5, 9},
	{1, 7, 6, 3, 7, 4, 6, 9, 3, 8, 5, 6, 2, 9, 6, 8, 5, 2, 8, 5, 8, 6, 9, 8, 4, 5, 8, 6},
	{1, 6, 8, 4, 7, 5, 9, 6, 7, 9, 5, 6, 8, 7, 5, 3, 9, 6, 8, 7, 6, 9, 7, 5, 8, 4, 10, 7, 9, 2, 10, 9},
	{1, 6, 9, 8, 1, 7, 9, 4, 6, 8, 9, 3, 5, 2, 6, 3, 9, 5, 7, 1, 8, 9, 2, 9, 7, 4, 5, 8},
	{10, 6, 8, 9, 6, 10, 5, 6, 4, 1, 8, 7, 1, 8, 9, 3, 7, 2, 6, 8, 9, 4, 7, 9, 3, 8, 5, 3, 9, 4, 7, 5},
}

// reels lengths [32, 28, 32, 28, 32], total reshuffles 25690112
// symbols: 72.169(lined) + 0(scatter) = 72.168623%
// free spins 2540160, q = 0.098877, sq = 1/(1-q) = 1.109726
// RTP = 72.169(sym) + 0.098877*251.43(fg) = 97.029053%
var ReelsReg97 = game.Reels5x{
	{1, 5, 7, 9, 8, 2, 10, 6, 4, 10, 6, 3, 4, 8, 9, 10, 7, 8, 9, 7, 8, 9, 3, 8, 9, 7, 5, 9, 7, 2, 5, 9},
	{1, 7, 6, 3, 7, 4, 6, 9, 3, 8, 5, 6, 2, 9, 6, 8, 5, 2, 8, 5, 8, 6, 9, 8, 4, 5, 8, 6},
	{1, 6, 8, 4, 7, 5, 9, 6, 7, 9, 5, 6, 8, 7, 5, 3, 9, 6, 8, 7, 6, 9, 7, 5, 8, 4, 10, 7, 9, 2, 10, 9},
	{1, 6, 9, 8, 1, 7, 9, 8, 6, 8, 9, 3, 5, 2, 6, 3, 9, 5, 7, 1, 8, 9, 2, 9, 7, 4, 5, 8},
	{10, 6, 8, 9, 6, 10, 5, 6, 4, 1, 8, 7, 1, 8, 9, 3, 7, 2, 6, 8, 9, 4, 7, 9, 8, 3, 5, 8, 9, 4, 7, 5},
}

// reels lengths [32, 28, 32, 28, 32], total reshuffles 25690112
// symbols: 68.866(lined) + 0(scatter) = 68.865562%
// free spins 7620480, q = 0.29663, sq = 1/(1-q) = 1.421729
// RTP = 68.866(sym) + 0.29663*251.43(fg) = 143.446853%
var ReelsReg143 = game.Reels5x{
	{1, 5, 7, 10, 8, 2, 10, 6, 4, 10, 6, 3, 4, 8, 9, 10, 7, 8, 9, 7, 8, 9, 3, 8, 9, 7, 5, 9, 7, 2, 5, 9},
	{1, 7, 6, 3, 7, 4, 6, 9, 3, 8, 5, 6, 2, 9, 6, 8, 5, 2, 8, 5, 8, 6, 9, 8, 4, 5, 8, 6},
	{1, 6, 8, 4, 7, 5, 9, 6, 10, 9, 5, 6, 8, 7, 5, 3, 9, 6, 8, 7, 6, 9, 7, 5, 8, 4, 10, 7, 9, 2, 10, 9},
	{1, 6, 9, 8, 1, 7, 9, 8, 6, 8, 9, 3, 5, 2, 6, 3, 9, 5, 7, 1, 8, 9, 2, 9, 7, 4, 5, 8},
	{10, 6, 8, 9, 6, 10, 5, 6, 4, 1, 8, 7, 1, 8, 9, 3, 7, 2, 6, 8, 9, 4, 7, 9, 8, 3, 5, 10, 9, 4, 7, 5},
}

// reels lengths [25, 25, 24, 25, 24], total reshuffles 9000000
// symbols: 194.86(lined) + 0(scatter) = 194.856667%
// free spins 2025000, q = 0.225, sq = 1/(1-q) = 1.290323
// RTP = sq*rtp(sym) = 1.2903*194.86 = 251.427957%
var ReelsBon = game.Reels5x{
	{5, 1, 9, 8, 2, 10, 6, 5, 10, 8, 3, 4, 8, 5, 10, 7, 8, 5, 9, 8, 5, 3, 6, 5, 8},
	{7, 9, 3, 7, 4, 6, 9, 7, 8, 9, 7, 8, 9, 7, 8, 5, 2, 8, 5, 4, 9, 6, 7, 9, 1},
	{6, 8, 4, 7, 5, 9, 6, 2, 9, 5, 6, 8, 10, 6, 3, 9, 6, 8, 7, 5, 10, 6, 8, 1},
	{6, 9, 1, 6, 7, 9, 4, 2, 8, 9, 3, 8, 7, 3, 4, 9, 5, 4, 6, 3, 5, 8, 3, 9, 1},
	{7, 4, 9, 6, 10, 9, 6, 4, 9, 8, 7, 1, 8, 9, 3, 7, 2, 9, 3, 8, 10, 4, 5, 2},
}

// Map with available reels.
var ReelsMap = map[string]*game.Reels5x{
	"85":  &ReelsReg85,
	"88":  &ReelsReg88,
	"90":  &ReelsReg90,
	"92":  &ReelsReg92,
	"94":  &ReelsReg94,
	"95":  &ReelsReg95,
	"96":  &ReelsReg96,
	"97":  &ReelsReg97,
	"143": &ReelsReg143,
	"bon": &ReelsBon,
}

// Lined payment.
var LinePay = [10][5]float64{
	{0, 10, 100, 1000, 5000}, //  1 columbus
	{0, 5, 50, 200, 1000},    //  2 spain
	{0, 5, 25, 100, 500},     //  3 necklace
	{0, 5, 15, 75, 250},      //  4 sextant
	{0, 0, 10, 40, 150},      //  5 ace
	{0, 0, 10, 40, 150},      //  6 king
	{0, 0, 10, 40, 150},      //  7 queen
	{0, 0, 5, 20, 100},       //  8 jack
	{0, 0, 5, 20, 100},       //  9 ten
	{0, 0, 0, 0, 0},          // 10 ship
}

const (
	jid = 1 // jackpot ID
)

// Jackpot win combinations.
var Jackpot = [10][5]int{
	{0, 0, 0, 0, 0}, //  1 columbus
	{0, 0, 0, 0, 0}, //  2 spain
	{0, 0, 0, 0, 0}, //  3 necklace
	{0, 0, 0, 0, 0}, //  4 sextant
	{0, 0, 0, 0, 0}, //  5 ace
	{0, 0, 0, 0, 0}, //  6 king
	{0, 0, 0, 0, 0}, //  7 queen
	{0, 0, 0, 0, 0}, //  8 jack
	{0, 0, 0, 0, 0}, //  9 ten
	{0, 0, 0, 0, 0}, // 10 ship
}

type Game struct {
	game.Slot5x3 `yaml:",inline"`
	FS           int `json:"fs" yaml:"fs" xml:"fs"` // free spin number
}

func NewGame(rd string) *Game {
	return &Game{
		Slot5x3: game.Slot5x3{
			RD:  rd,
			BLI: "nvm10",
			SBL: game.MakeSblNum(5),
			Bet: 1,
		},
		FS: 0,
	}
}

const wild, scat = 1, 10

func (g *Game) Scanner(screen game.Screen, ws *game.WinScan) {
	g.ScanLined(screen, ws)
	g.ScanScatters(screen, ws)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen game.Screen, ws *game.WinScan) {
	var wbon game.Sym
	if g.FS > 0 {
		wbon = scat
	}

	var bl = game.BetLines5x[g.BLI]
	for li := g.SBL.Next(0); li != 0; li = g.SBL.Next(li) {
		var line = bl.Line(li)

		var numw, numl = 0, 5
		var syml game.Sym
		for x := 1; x <= 5; x++ {
			var sx = screen.At(x, line.At(x))
			if sx == wild || sx == wbon {
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
	if count := screen.ScatNumOdd(scat); count >= 3 {
		ws.Wins = append(ws.Wins, game.WinItem{
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPosOdd(scat),
			Free: 10,
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

func (g *Game) Apply(screen game.Screen, sw *game.WinScan) {
	if g.FS > 0 {
		g.Gain += sw.Gain()
	} else {
		g.Gain = sw.Gain()
	}

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

func (g *Game) SetReels(rd string) error {
	if _, ok := ReelsMap[rd]; !ok {
		return game.ErrNoReels
	}
	g.RD = rd
	return nil
}
