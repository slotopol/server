package firejoker

// See: https://freeslotshub.com/playngo/fire-joker/

import (
	"math/rand/v2"

	"github.com/slotopol/server/game"
)

// *bonus reels calculations*
// RTP[1] = 816.33(lined) + 0.36735(scatter) = 816.693878%
// RTP[2] = 318.37(lined) + 0.36735(scatter) = 318.734694%
// RTP[3] = 318.37(lined) + 0.36735(scatter) = 318.734694%
// RTP[4] = 163.27(lined) + 0.36735(scatter) = 163.632653%
// RTP[5] = 163.27(lined) + 0.36735(scatter) = 163.632653%
// RTP[6] = 163.27(lined) + 0.36735(scatter) = 163.632653%
// RTP[7] = 163.27(lined) + 0.36735(scatter) = 163.632653%
// average freespins RTP = 301.241983%
// *regular reels calculations*
// reels lengths [35, 35, 35, 35, 35], total reshuffles 52521875
// symbols: 16.875(lined) + 50.653(scatter) = 67.528054%
// free spins 2764800, q = 0.052641
// free games frequency: 1/189.97
// RTP = 67.528(sym) + 0.052641*301.24(fg) = 83.385710%
var Reels83 = game.Reels5x{
	{5, 8, 7, 2, 4, 7, 3, 9, 4, 5, 2, 4, 6, 1, 3, 6, 4, 2, 5, 6, 1, 5, 4, 3, 7, 2, 3, 1, 6, 7, 1, 6, 5, 7, 1},
	{4, 2, 1, 8, 7, 1, 9, 7, 5, 3, 4, 6, 1, 2, 4, 5, 1, 4, 7, 3, 1, 6, 4, 5, 7, 2, 3, 6, 7, 5, 6, 2, 5, 3, 6},
	{1, 8, 5, 2, 4, 6, 3, 7, 4, 5, 2, 7, 6, 3, 1, 2, 4, 6, 7, 3, 5, 7, 1, 4, 9, 6, 5, 1, 7, 4, 5, 1, 3, 2, 6},
	{7, 4, 3, 5, 1, 6, 2, 4, 1, 6, 5, 9, 4, 3, 2, 6, 1, 5, 4, 6, 7, 8, 6, 7, 3, 4, 7, 1, 5, 2, 7, 5, 2, 1, 3},
	{5, 2, 3, 6, 4, 3, 7, 4, 3, 2, 1, 4, 5, 7, 1, 6, 7, 5, 3, 4, 2, 6, 8, 5, 6, 1, 7, 2, 9, 1, 5, 6, 1, 7, 4},
}

// *bonus reels calculations*
// RTP[1] = 657.44(lined) + 0.38927(scatter) = 657.828720%
// RTP[2] = 328.72(lined) + 0.38927(scatter) = 329.108997%
// RTP[3] = 328.72(lined) + 0.38927(scatter) = 329.108997%
// RTP[4] = 168.69(lined) + 0.38927(scatter) = 169.074394%
// RTP[5] = 168.69(lined) + 0.38927(scatter) = 169.074394%
// RTP[6] = 168.69(lined) + 0.38927(scatter) = 169.074394%
// RTP[7] = 168.69(lined) + 0.38927(scatter) = 169.074394%
// average freespins RTP = 284.620613%
// *regular reels calculations*
// reels lengths [34, 34, 34, 34, 34], total reshuffles 45435424
// symbols: 14.374(lined) + 58.146(scatter) = 72.520518%
// free spins 2594700, q = 0.057107
// free games frequency: 1/175.11
// RTP = 72.521(sym) + 0.057107*284.62(fg) = 88.774468%
var Reels89 = game.Reels5x{
	{2, 6, 4, 5, 7, 3, 6, 7, 5, 1, 3, 6, 4, 5, 7, 1, 4, 2, 8, 6, 3, 7, 5, 6, 2, 9, 4, 1, 3, 4, 2, 5, 7, 1},
	{5, 4, 9, 6, 7, 5, 2, 3, 1, 6, 5, 1, 6, 3, 5, 1, 8, 4, 3, 5, 7, 4, 6, 2, 4, 7, 6, 4, 2, 1, 7, 3, 2, 7},
	{2, 3, 4, 2, 5, 4, 7, 5, 6, 1, 5, 6, 1, 3, 4, 2, 7, 5, 1, 6, 7, 3, 5, 7, 3, 8, 4, 7, 6, 2, 1, 6, 9, 4},
	{6, 7, 5, 6, 8, 7, 5, 2, 7, 4, 6, 1, 5, 2, 3, 1, 4, 3, 9, 4, 1, 6, 5, 7, 3, 2, 6, 5, 2, 4, 3, 7, 1, 4},
	{6, 7, 4, 1, 5, 3, 6, 4, 1, 5, 4, 7, 5, 2, 4, 7, 2, 8, 6, 3, 5, 6, 2, 1, 3, 2, 7, 1, 5, 6, 4, 7, 3, 9},
}

// *bonus reels calculations*
// RTP[1] = 816.33(lined) + 0.36735(scatter) = 816.693878%
// RTP[2] = 318.37(lined) + 0.36735(scatter) = 318.734694%
// RTP[3] = 318.37(lined) + 0.36735(scatter) = 318.734694%
// RTP[4] = 163.27(lined) + 0.36735(scatter) = 163.632653%
// RTP[5] = 163.27(lined) + 0.36735(scatter) = 163.632653%
// RTP[6] = 163.27(lined) + 0.36735(scatter) = 163.632653%
// RTP[7] = 163.27(lined) + 0.36735(scatter) = 163.632653%
// average freespins RTP = 301.241983%
// *regular reels calculations*
// reels lengths [35, 35, 36, 35, 35], total reshuffles 54022500
// symbols: 16.406(lined) + 51.118(scatter) = 67.524561%
// free spins 4354560, q = 0.080606
// free games frequency: 1/124.06
// RTP = 67.525(sym) + 0.080606*301.24(fg) = 91.806597%
var Reels92 = game.Reels5x{
	{5, 8, 7, 2, 4, 7, 3, 9, 4, 5, 2, 4, 6, 1, 3, 6, 4, 2, 5, 6, 1, 5, 4, 3, 7, 2, 3, 1, 6, 7, 1, 6, 5, 7, 1},
	{4, 2, 1, 8, 7, 1, 9, 7, 5, 3, 4, 6, 1, 2, 4, 5, 1, 4, 7, 3, 1, 6, 4, 5, 7, 2, 3, 6, 7, 5, 6, 2, 5, 3, 6},
	{2, 1, 3, 2, 7, 5, 8, 7, 6, 1, 3, 4, 5, 6, 1, 5, 6, 3, 7, 6, 3, 7, 4, 6, 8, 1, 5, 4, 2, 7, 4, 1, 5, 2, 4, 9},
	{7, 4, 3, 5, 1, 6, 2, 4, 1, 6, 5, 9, 4, 3, 2, 6, 1, 5, 4, 6, 7, 8, 6, 7, 3, 4, 7, 1, 5, 2, 7, 5, 2, 1, 3},
	{5, 2, 3, 6, 4, 3, 7, 4, 3, 2, 1, 4, 5, 7, 1, 6, 7, 5, 3, 4, 2, 6, 8, 5, 6, 1, 7, 2, 9, 1, 5, 6, 1, 7, 4},
}

// *bonus reels calculations*
// RTP[1] = 816.33(lined) + 0.36735(scatter) = 816.693878%
// RTP[2] = 318.37(lined) + 0.36735(scatter) = 318.734694%
// RTP[3] = 318.37(lined) + 0.36735(scatter) = 318.734694%
// RTP[4] = 163.27(lined) + 0.36735(scatter) = 163.632653%
// RTP[5] = 163.27(lined) + 0.36735(scatter) = 163.632653%
// RTP[6] = 163.27(lined) + 0.36735(scatter) = 163.632653%
// RTP[7] = 163.27(lined) + 0.36735(scatter) = 163.632653%
// average freespins RTP = 301.241983%
// *regular reels calculations*
// reels lengths [35, 34, 30, 34, 35], total reshuffles 42483000
// symbols: 15.505(lined) + 62.01(scatter) = 77.515026%
// free spins 2525850, q = 0.059456
// free games frequency: 1/168.19
// RTP = 77.515(sym) + 0.059456*301.24(fg) = 95.425533%
var Reels95 = game.Reels5x{
	{5, 8, 7, 2, 4, 7, 3, 9, 4, 5, 2, 4, 6, 1, 3, 6, 4, 2, 5, 6, 1, 5, 4, 3, 7, 2, 3, 1, 6, 7, 1, 6, 5, 7, 1},
	{5, 6, 4, 5, 6, 1, 4, 3, 6, 2, 4, 1, 7, 9, 4, 3, 7, 2, 6, 7, 3, 2, 5, 6, 7, 5, 1, 4, 8, 2, 7, 5, 1, 3},
	{7, 5, 4, 3, 5, 6, 1, 4, 6, 3, 1, 8, 6, 2, 7, 1, 2, 7, 5, 4, 2, 5, 3, 7, 2, 3, 1, 9, 4, 6},
	{4, 5, 2, 1, 3, 7, 1, 4, 6, 2, 1, 7, 4, 5, 2, 9, 6, 3, 7, 5, 1, 6, 4, 3, 7, 4, 5, 6, 3, 2, 6, 8, 5, 7},
	{5, 2, 3, 6, 4, 3, 7, 4, 3, 2, 1, 4, 5, 7, 1, 6, 7, 5, 3, 4, 2, 6, 8, 5, 6, 1, 7, 2, 9, 1, 5, 6, 1, 7, 4},
}

// *bonus reels calculations*
// RTP[1] = 657.44(lined) + 0.38927(scatter) = 657.828720%
// RTP[2] = 328.72(lined) + 0.38927(scatter) = 329.108997%
// RTP[3] = 328.72(lined) + 0.38927(scatter) = 329.108997%
// RTP[4] = 168.69(lined) + 0.38927(scatter) = 169.074394%
// RTP[5] = 168.69(lined) + 0.38927(scatter) = 169.074394%
// RTP[6] = 168.69(lined) + 0.38927(scatter) = 169.074394%
// RTP[7] = 168.69(lined) + 0.38927(scatter) = 169.074394%
// average freespins RTP = 284.620613%
// *regular reels calculations*
// reels lengths [34, 34, 35, 34, 34], total reshuffles 46771760
// symbols: 13.964(lined) + 58.476(scatter) = 72.439667%
// free spins 4084560, q = 0.087330
// free games frequency: 1/114.51
// RTP = 72.44(sym) + 0.08733*284.62(fg) = 97.295476%
var Reels97 = game.Reels5x{
	{2, 6, 4, 5, 7, 3, 6, 7, 5, 1, 3, 6, 4, 5, 7, 1, 4, 2, 8, 6, 3, 7, 5, 6, 2, 9, 4, 1, 3, 4, 2, 5, 7, 1},
	{5, 4, 9, 6, 7, 5, 2, 3, 1, 6, 5, 1, 6, 3, 5, 1, 8, 4, 3, 5, 7, 4, 6, 2, 4, 7, 6, 4, 2, 1, 7, 3, 2, 7},
	{5, 3, 6, 1, 5, 3, 7, 8, 4, 5, 2, 6, 7, 4, 3, 1, 8, 6, 5, 7, 6, 2, 1, 4, 2, 3, 7, 5, 4, 2, 9, 4, 6, 1, 7},
	{6, 7, 5, 6, 8, 7, 5, 2, 7, 4, 6, 1, 5, 2, 3, 1, 4, 3, 9, 4, 1, 6, 5, 7, 3, 2, 6, 5, 2, 4, 3, 7, 1, 4},
	{6, 7, 4, 1, 5, 3, 6, 4, 1, 5, 4, 7, 5, 2, 4, 7, 2, 8, 6, 3, 5, 6, 2, 1, 3, 2, 7, 1, 5, 6, 4, 7, 3, 9},
}

// *bonus reels calculations*
// RTP[1] = 679.52(lined) + 0.41322(scatter) = 679.935721%
// RTP[2] = 339.76(lined) + 0.41322(scatter) = 340.174472%
// RTP[3] = 339.76(lined) + 0.41322(scatter) = 340.174472%
// RTP[4] = 135.9(lined) + 0.41322(scatter) = 136.317723%
// RTP[5] = 174.47(lined) + 0.41322(scatter) = 174.885216%
// RTP[6] = 174.47(lined) + 0.41322(scatter) = 174.885216%
// RTP[7] = 174.47(lined) + 0.41322(scatter) = 174.885216%
// average freespins RTP = 288.751148%
// *regular reels calculations*
// reels lengths [33, 33, 33, 33, 33], total reshuffles 39135393
// symbols: 14.911(lined) + 67.06(scatter) = 81.970129%
// free spins 2430000, q = 0.062092
// free games frequency: 1/161.05
// RTP = 81.97(sym) + 0.062092*288.75(fg) = 99.899303%
var Reels100 = game.Reels5x{
	{7, 2, 5, 4, 6, 5, 7, 4, 6, 1, 3, 5, 2, 6, 1, 4, 2, 1, 7, 6, 4, 3, 1, 6, 7, 2, 3, 8, 5, 7, 3, 9, 5},
	{6, 8, 5, 1, 2, 7, 1, 6, 2, 3, 1, 6, 4, 7, 5, 4, 3, 6, 5, 2, 4, 5, 1, 3, 7, 4, 2, 9, 5, 7, 6, 3, 7},
	{5, 2, 3, 6, 4, 5, 2, 7, 5, 6, 7, 3, 5, 2, 8, 4, 7, 5, 2, 1, 4, 3, 7, 4, 6, 3, 1, 6, 7, 1, 9, 6, 1},
	{7, 6, 1, 7, 9, 5, 3, 7, 6, 5, 1, 3, 2, 5, 6, 4, 1, 7, 2, 4, 5, 7, 6, 4, 1, 3, 2, 6, 8, 3, 2, 5, 4},
	{3, 1, 2, 6, 5, 2, 6, 7, 5, 4, 6, 1, 4, 3, 7, 1, 5, 4, 1, 3, 5, 9, 2, 7, 6, 8, 5, 7, 4, 2, 3, 6, 7},
}

// *bonus reels calculations*
// RTP[1] = 872.36(lined) + 0.41322(scatter) = 872.773186%
// RTP[2] = 436.18(lined) + 0.41322(scatter) = 436.593205%
// RTP[3] = 436.18(lined) + 0.41322(scatter) = 436.593205%
// RTP[4] = 135.9(lined) + 0.41322(scatter) = 136.317723%
// RTP[5] = 135.9(lined) + 0.41322(scatter) = 136.317723%
// RTP[6] = 135.9(lined) + 0.41322(scatter) = 136.317723%
// RTP[7] = 135.9(lined) + 0.41322(scatter) = 136.317723%
// average freespins RTP = 327.318641%
// *regular reels calculations*
// reels lengths [33, 33, 33, 33, 33], total reshuffles 39135393
// symbols: 21.346(lined) + 67.06(scatter) = 88.405521%
// free spins 2430000, q = 0.062092
// free games frequency: 1/161.05
// RTP = 88.406(sym) + 0.062092*327.32(fg) = 108.729433%
var Reels109 = game.Reels5x{
	{7, 4, 2, 6, 1, 2, 4, 1, 3, 2, 5, 7, 9, 1, 5, 6, 7, 3, 4, 1, 5, 3, 8, 1, 6, 2, 7, 4, 3, 6, 5, 2, 3},
	{1, 3, 9, 6, 4, 5, 1, 8, 2, 5, 1, 7, 4, 6, 2, 1, 4, 2, 3, 5, 4, 2, 1, 7, 3, 6, 7, 3, 2, 6, 3, 7, 5},
	{6, 1, 5, 9, 3, 7, 1, 3, 6, 4, 5, 2, 4, 7, 1, 5, 2, 3, 1, 8, 3, 2, 7, 4, 2, 1, 6, 2, 3, 5, 6, 7, 4},
	{4, 2, 1, 3, 4, 2, 6, 5, 2, 7, 1, 3, 6, 5, 4, 1, 2, 6, 8, 3, 1, 7, 3, 5, 7, 4, 2, 3, 5, 9, 6, 7, 1},
	{1, 2, 7, 1, 5, 2, 4, 5, 2, 6, 1, 3, 9, 4, 2, 3, 7, 6, 4, 2, 5, 3, 8, 4, 5, 1, 3, 6, 1, 7, 3, 6, 7},
}

// Map with available reels.
var ReelsMap = map[string]*game.Reels5x{
	"83":  &Reels83,
	"89":  &Reels89,
	"92":  &Reels92,
	"95":  &Reels95,
	"97":  &Reels97,
	"100": &Reels100,
	"109": &Reels109,
}

// Lined payment.
var LinePay = [13][5]float64{
	{0, 0, 20, 50, 100}, // 1 seven
	{0, 0, 10, 25, 50},  // 2 bell
	{0, 0, 10, 25, 50},  // 3 melon
	{0, 0, 4, 10, 20},   // 4 plum
	{0, 0, 4, 10, 20},   // 5 orange
	{0, 0, 4, 10, 20},   // 6 lemon
	{0, 0, 4, 10, 20},   // 7 cherry
	{0, 0, 0, 0, 0},     // 8 bonus
	{0, 0, 0, 0, 0},     // 9 joker
}

// Scatters payment.
var ScatPay = [5]float64{0, 0.5, 3} // 8 bonus

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 10} // 8 bonus

type Game struct {
	game.Slot5x3 `yaml:",inline"`
	FS           int `json:"fs" yaml:"fs" xml:"fs"` // free spin number
}

func NewGame(rd string) *Game {
	return &Game{
		Slot5x3: game.Slot5x3{
			RD:  rd,
			SBL: game.MakeSblNum(5),
			Bet: 1,
		},
		FS: 0,
	}
}

const scat, jack = 8, 9

var bl = game.Lineset5x{
	{2, 2, 2, 2, 2}, // 1
	{1, 1, 1, 1, 1}, // 2
	{3, 3, 3, 3, 3}, // 3
	{1, 2, 3, 2, 1}, // 4
	{3, 2, 1, 2, 3}, // 5
}

func (g *Game) Scanner(screen game.Screen, ws *game.WinScan) {
	g.ScanLined(screen, ws)
	g.ScanScatters(screen, ws)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen game.Screen, ws *game.WinScan) {
	for li := g.SBL.Next(0); li != 0; li = g.SBL.Next(li) {
		var line = bl.Line(li)

		var syml, numl = screen.Pos(1, line), 1
		for x := 2; x <= 5; x++ {
			var sx = screen.Pos(x, line)
			if sx != syml {
				break
			}
			numl++
		}

		if pay := LinePay[syml-1][numl-1]; pay > 0 {
			ws.Wins = append(ws.Wins, game.WinItem{
				Pay:  g.Bet * pay,
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyL(numl),
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen game.Screen, ws *game.WinScan) {
	if count := screen.ScatNum(scat); count >= 2 {
		var pay, fs = ScatPay[count-1], ScatFreespin[count-1]
		ws.Wins = append(ws.Wins, game.WinItem{
			Pay:  g.Bet * float64(g.SBL.Num()) * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
			Free: fs,
		})
	}
	if count := screen.ScatNum(jack); count == 5 {
		ws.Wins = append(ws.Wins, game.WinItem{
			Pay:  g.Bet * float64(g.SBL.Num()) * 100000,
			Mult: 1,
			Sym:  jack,
			Num:  5,
			XY:   screen.ScatPos(scat),
		})
	}
}

func (g *Game) Spin(screen game.Screen) {
	var reels = ReelsMap[g.RD]
	if g.FS == 0 {
		screen.Spin(reels)
	} else {
		var reel []game.Sym
		var hit int
		reel = reels.Reel(1)
		hit = rand.N(len(reel))
		screen.SetCol(1, reel, hit)
		var gs = game.Sym(rand.N(7) + 1)
		for x := 2; x <= 4; x++ {
			screen.Set(x, 1, gs)
			screen.Set(x, 2, gs)
			screen.Set(x, 3, gs)
		}
		reel = reels.Reel(5)
		hit = rand.N(len(reel))
		screen.SetCol(5, reel, hit)
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

func (g *Game) SetLines(sbl game.SBL) error {
	return game.ErrNoFeature
}

func (g *Game) SetReels(rd string) error {
	if _, ok := ReelsMap[rd]; !ok {
		return game.ErrNoReels
	}
	g.RD = rd
	return nil
}
