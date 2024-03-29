package atthemovies

import "github.com/slotopol/server/game"

// reels lengths [36, 37, 36, 37, 36], total reshuffles 63872064
// symbols: 75.457(lined) + 10.499(scatter) = 85.955963%
// free spins 2576664, q = 0.040341, sq = 1/(1-q) = 1.042037
// RTP = rtp(sym) + q*sq*2*rtp(sym) = 85.956 + 0.040341*179.14 = 93.182595%
var Reels93 = game.Reels5x{
	{4, 8, 7, 6, 5, 7, 4, 5, 6, 4, 8, 3, 6, 4, 7, 5, 6, 8, 3, 7, 10, 8, 7, 6, 8, 5, 6, 7, 8, 6, 5, 1, 8, 5, 2, 7},
	{4, 8, 6, 4, 5, 6, 9, 8, 2, 7, 8, 4, 7, 5, 8, 7, 5, 8, 6, 4, 5, 8, 6, 3, 7, 6, 8, 10, 7, 6, 3, 5, 7, 1, 5, 6, 7},
	{8, 4, 6, 7, 4, 6, 2, 7, 3, 8, 7, 5, 8, 6, 10, 5, 6, 8, 7, 4, 6, 7, 3, 8, 5, 1, 6, 5, 4, 8, 7, 5, 8, 6, 5, 7},
	{4, 7, 8, 6, 7, 8, 6, 5, 3, 7, 2, 5, 1, 8, 4, 7, 5, 4, 7, 9, 6, 5, 8, 6, 5, 4, 6, 8, 10, 3, 7, 6, 8, 7, 5, 6, 8},
	{4, 6, 8, 4, 7, 6, 4, 7, 8, 3, 7, 5, 2, 8, 5, 7, 8, 1, 7, 6, 8, 7, 5, 6, 10, 7, 6, 3, 5, 4, 6, 5, 8, 6, 5, 8},
}

// reels lengths [31, 35, 34, 35, 31], total reshuffles 40025650
// symbols: 72.319(lined) + 12.352(scatter) = 84.671307%
// free spins 2119824, q = 0.052962, sq = 1/(1-q) = 1.055923
// RTP = rtp(sym) + q*sq*2*rtp(sym) = 84.671 + 0.052962*178.81 = 94.141528%
var Reels94 = game.Reels5x{
	{5, 2, 7, 8, 5, 6, 8, 7, 3, 5, 6, 7, 8, 3, 6, 7, 5, 6, 10, 8, 7, 4, 8, 7, 5, 8, 6, 7, 1, 8, 4},
	{5, 6, 8, 1, 6, 8, 7, 3, 9, 7, 4, 6, 8, 4, 6, 7, 5, 10, 7, 8, 5, 7, 4, 8, 5, 3, 7, 5, 4, 6, 2, 8, 7, 6, 8},
	{8, 6, 3, 5, 8, 4, 6, 1, 7, 6, 8, 7, 5, 8, 7, 6, 8, 7, 4, 8, 5, 2, 7, 6, 4, 7, 5, 6, 3, 5, 8, 7, 4, 10},
	{2, 7, 6, 10, 4, 8, 3, 7, 5, 6, 8, 1, 5, 7, 6, 8, 4, 5, 8, 7, 9, 3, 7, 4, 8, 6, 7, 5, 8, 7, 6, 5, 4, 6, 8},
	{7, 8, 3, 7, 8, 5, 6, 8, 4, 7, 8, 6, 5, 7, 8, 5, 7, 10, 6, 3, 8, 4, 7, 6, 8, 5, 2, 7, 6, 5, 1},
}

// reels lengths [34, 34, 34, 34, 34], total reshuffles 45435424
// symbols: 74.5(lined) + 11.802(scatter) = 86.301737%
// free spins 2231280, q = 0.049109, sq = 1/(1-q) = 1.051645
// RTP = rtp(sym) + q*sq*2*rtp(sym) = 86.302 + 0.049109*181.52 = 95.215851%
var Reels95 = game.Reels5x{
	{7, 8, 5, 6, 7, 5, 3, 8, 6, 5, 8, 7, 5, 1, 7, 4, 8, 6, 5, 7, 4, 6, 10, 8, 7, 4, 6, 3, 8, 7, 2, 6, 8, 4},
	{5, 2, 7, 4, 6, 3, 8, 7, 5, 6, 10, 8, 3, 6, 8, 5, 7, 6, 8, 7, 6, 8, 7, 4, 6, 8, 5, 4, 6, 1, 7, 9, 8, 7},
	{8, 6, 3, 5, 8, 4, 6, 1, 7, 6, 8, 7, 5, 8, 7, 6, 8, 7, 4, 8, 5, 2, 7, 6, 4, 7, 5, 6, 3, 5, 8, 7, 4, 10},
	{7, 6, 8, 5, 6, 7, 4, 8, 7, 6, 5, 7, 3, 8, 4, 6, 8, 7, 5, 8, 2, 5, 7, 8, 6, 1, 8, 6, 7, 9, 3, 6, 10, 4},
	{5, 8, 6, 1, 5, 7, 8, 6, 10, 5, 7, 6, 5, 7, 6, 4, 7, 5, 8, 4, 7, 8, 4, 7, 3, 8, 6, 7, 3, 8, 2, 6, 8, 4},
}

// reels lengths [34, 35, 34, 35, 34], total reshuffles 48147400
// symbols: 76.744(lined) + 11.572(scatter) = 88.316046%
// free spins 2287008, q = 0.0475, sq = 1/(1-q) = 1.049869
// RTP = rtp(sym) + q*sq*2*rtp(sym) = 88.316 + 0.0475*185.44 = 97.124497%
var Reels97 = game.Reels5x{
	{8, 6, 7, 4, 8, 6, 1, 8, 5, 2, 7, 8, 5, 7, 6, 5, 3, 6, 5, 7, 4, 8, 6, 7, 8, 6, 4, 8, 7, 10, 5, 6, 7, 3},
	{4, 8, 6, 7, 2, 8, 6, 5, 7, 8, 4, 6, 5, 7, 6, 8, 4, 10, 8, 7, 6, 8, 1, 5, 3, 7, 6, 5, 8, 7, 9, 3, 6, 5, 7},
	{7, 8, 5, 6, 8, 7, 5, 6, 7, 5, 3, 7, 5, 6, 8, 4, 7, 8, 6, 3, 8, 6, 4, 8, 7, 6, 1, 8, 6, 7, 2, 5, 10, 4},
	{5, 7, 6, 8, 7, 6, 8, 2, 6, 7, 1, 5, 4, 8, 10, 5, 6, 7, 5, 8, 3, 6, 4, 8, 7, 5, 6, 7, 4, 6, 8, 9, 3, 7, 8},
	{1, 6, 5, 8, 7, 2, 8, 6, 5, 4, 7, 10, 8, 3, 5, 7, 4, 8, 6, 7, 8, 6, 5, 7, 8, 6, 5, 8, 6, 7, 3, 6, 4, 7},
}

// reels lengths [33, 34, 33, 34, 33], total reshuffles 41543172
// symbols: 78.344(lined) + 12.171(scatter) = 90.514056%
// free spins 2148660, q = 0.051721, sq = 1/(1-q) = 1.054542
// RTP = rtp(sym) + q*sq*2*rtp(sym) = 90.514 + 0.051721*190.9 = 100.387712%
var Reels100 = game.Reels5x{
	{8, 7, 3, 6, 5, 7, 8, 1, 6, 4, 7, 8, 2, 6, 5, 8, 7, 6, 8, 5, 6, 7, 10, 8, 4, 7, 6, 8, 3, 7, 6, 4, 5},
	{5, 2, 7, 4, 6, 3, 8, 7, 5, 6, 10, 8, 3, 6, 8, 5, 7, 6, 8, 7, 6, 8, 7, 4, 6, 8, 5, 4, 6, 1, 7, 9, 8, 7},
	{7, 6, 4, 5, 8, 7, 6, 8, 5, 6, 3, 8, 6, 1, 7, 8, 4, 7, 8, 2, 7, 6, 8, 7, 5, 8, 6, 7, 10, 4, 6, 3, 5},
	{7, 6, 8, 5, 6, 7, 4, 8, 7, 6, 5, 7, 3, 8, 4, 6, 8, 7, 5, 8, 2, 5, 7, 8, 6, 1, 8, 6, 7, 9, 3, 6, 10, 4},
	{5, 6, 7, 8, 6, 7, 8, 6, 7, 5, 8, 3, 6, 2, 5, 7, 3, 8, 6, 4, 7, 6, 4, 8, 5, 10, 6, 8, 7, 4, 8, 1, 7},
}

// Map with available reels.
var ReelsMap = map[string]*game.Reels5x{
	"93":  &Reels93,
	"94":  &Reels94,
	"95":  &Reels95,
	"97":  &Reels97,
	"100": &Reels100,
}

// Lined payment.
var LinePay = [10][5]float64{
	{0, 20, 200, 500, 1000}, //  1 oscar
	{0, 10, 100, 250, 500},  //  2 popcorn
	{0, 5, 50, 100, 200},    //  3 poster
	{0, 2, 25, 50, 100},     //  4 a
	{0, 0, 20, 40, 80},      //  5 dummy
	{0, 0, 15, 30, 60},      //  6 maw
	{0, 0, 10, 20, 40},      //  7 starship
	{0, 0, 5, 10, 20},       //  8 heart
	{0, 0, 0, 0, 0},         //  9 masks
	{0, 0, 0, 0, 0},         // 10 projector
}

// Scatters payment.
var ScatPay = [5]float64{0, 2, 0, 0, 0} // 10 projector

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 8, 12, 20} // 10 projector

const (
	jid = 1 // jackpot ID
)

// Jackpot win combinations.
var Jackpot = [10][5]int{
	{0, 0, 0, 0, 0}, //  1 oscar
	{0, 0, 0, 0, 0}, //  2 popcorn
	{0, 0, 0, 0, 0}, //  3 poster
	{0, 0, 0, 0, 0}, //  4 a
	{0, 0, 0, 0, 0}, //  5 dummy
	{0, 0, 0, 0, 0}, //  6 maw
	{0, 0, 0, 0, 0}, //  7 starship
	{0, 0, 0, 0, 0}, //  8 heart
	{0, 0, 0, 0, 0}, //  9 masks
	{0, 0, 0, 0, 0}, // 10 projector
}

type Game struct {
	game.Slot5x3 `yaml:",inline"`
	FS           int `json:"fs" yaml:"fs" xml:"fs"` // free spin number
}

func NewGame(rd string) *Game {
	return &Game{
		Slot5x3: game.Slot5x3{
			RD:  rd,
			BLI: "bs25",
			SBL: game.MakeSBL(1),
			Bet: 1,
		},
		FS: 0,
	}
}

const wild, scat = 9, 10

func (g *Game) Scanner(screen game.Screen, ws *game.WinScan) {
	g.ScanLined(screen, ws)
	g.ScanScatters(screen, ws)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen game.Screen, ws *game.WinScan) {
	var mm float64 = 1 // mult mode
	if g.FS > 0 {
		mm = 2
	}

	var bl = game.BetLines5x[g.BLI]
	for li := g.SBL.Next(0); li != 0; li = g.SBL.Next(li) {
		var line = bl.Line(li)

		var syml, numl = screen.At(1, line.At(1)), 1
		var mw = mm // mult wild
		for x := 2; x <= 5; x++ {
			var sx = screen.At(x, line.At(x))
			if sx == wild {
				mw = 2 * mm
			} else if sx != syml {
				break
			}
			numl++
		}

		if pay := LinePay[syml-1][numl-1]; pay > 0 {
			ws.Wins = append(ws.Wins, game.WinItem{
				Pay:  g.Bet * pay,
				Mult: mw,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyN(numl),
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen game.Screen, ws *game.WinScan) {
	if count := screen.ScatNum(scat); count >= 2 {
		var mm float64 = 1 // mult mode
		if g.FS > 0 {
			mm = 2
		}
		var pay, fs = ScatPay[count-1], ScatFreespin[count-1]
		ws.Wins = append(ws.Wins, game.WinItem{
			Pay:  g.Bet * pay, // independent from selected lines
			Mult: mm,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
			Free: fs,
		})
	}
}

func (g *Game) Spin(screen game.Screen) {
	screen.Spin(ReelsMap[g.RD])
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
