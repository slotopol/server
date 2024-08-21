package greatblue

// See: https://freeslotshub.com/playtech/great-blue/

import (
	"math"
	"math/rand/v2"

	"github.com/slotopol/server/game"
)

// reels lengths [44, 44, 44, 44, 44], total reshuffles 164916224
// symbols: 48.269(lined) + 9.1736(scatter) = 57.442790%
// average plain freespins at 1st iteration: 124.8
// average multiplier at free games: 7.2
// free games 470718, q = 0.35621, sq = 1.4456
// free games frequency: 1/350.35
// RTP = rtpsym + q*sq*rtpsym = 57.443 + 29.581 = 87.023352%
var Reels87 = game.Reels5x{
	{10, 4, 8, 13, 6, 11, 10, 2, 11, 4, 9, 6, 10, 5, 12, 10, 6, 9, 1, 7, 3, 12, 6, 7, 3, 12, 8, 5, 7, 11, 8, 5, 9, 2, 7, 4, 12, 8, 3, 12, 2, 11, 9, 5},
	{12, 2, 11, 5, 7, 2, 10, 3, 7, 2, 9, 1, 11, 5, 10, 8, 6, 11, 10, 3, 8, 5, 12, 4, 7, 6, 8, 3, 9, 6, 12, 9, 4, 13, 12, 5, 11, 10, 6, 12, 4, 9, 7, 8},
	{7, 6, 12, 11, 5, 10, 3, 12, 6, 8, 9, 5, 8, 9, 11, 6, 12, 8, 13, 4, 7, 1, 9, 10, 2, 12, 5, 7, 2, 12, 3, 10, 5, 9, 7, 4, 11, 3, 10, 2, 11, 6, 8, 4},
	{5, 9, 3, 8, 2, 13, 4, 11, 3, 7, 12, 6, 9, 3, 8, 5, 7, 4, 12, 7, 6, 10, 5, 12, 2, 8, 4, 7, 11, 12, 9, 1, 10, 5, 11, 6, 8, 2, 10, 6, 9, 12, 10, 11},
	{6, 11, 5, 13, 10, 6, 9, 8, 1, 12, 9, 4, 12, 3, 10, 6, 8, 9, 5, 8, 4, 7, 10, 2, 12, 5, 11, 6, 7, 3, 10, 2, 12, 3, 8, 4, 11, 7, 5, 11, 2, 7, 9, 12},
}

// reels lengths [44, 44, 44, 44, 44], total reshuffles 164916224
// symbols: 49.461(lined) + 9.1736(scatter) = 58.634659%
// average plain freespins at 1st iteration: 124.8
// average multiplier at free games: 7.2
// free games 470718, q = 0.35621, sq = 1.4456
// free games frequency: 1/350.35
// RTP = rtpsym + q*sq*rtpsym = 58.635 + 30.194 = 88.828983%
var Reels89 = game.Reels5x{
	{9, 7, 4, 9, 3, 8, 5, 9, 10, 6, 7, 4, 12, 2, 8, 11, 2, 7, 5, 11, 6, 10, 2, 12, 4, 7, 3, 11, 5, 8, 1, 12, 10, 3, 9, 6, 10, 4, 8, 5, 13, 12, 11, 6},
	{8, 10, 7, 11, 4, 1, 5, 12, 9, 8, 6, 11, 4, 9, 3, 13, 4, 4, 5, 6, 2, 9, 3, 8, 7, 5, 10, 6, 8, 10, 2, 12, 5, 7, 12, 2, 11, 7, 6, 12, 10, 11, 9, 3},
	{6, 11, 5, 9, 3, 12, 10, 4, 12, 1, 11, 3, 8, 6, 10, 4, 9, 8, 4, 12, 6, 8, 5, 7, 2, 10, 12, 9, 5, 8, 2, 7, 10, 2, 7, 4, 9, 6, 11, 5, 13, 11, 3, 7},
	{13, 4, 10, 1, 8, 4, 11, 2, 7, 8, 4, 9, 12, 5, 9, 6, 10, 5, 7, 8, 6, 11, 4, 9, 11, 3, 12, 6, 7, 3, 9, 5, 10, 8, 5, 12, 11, 6, 10, 2, 7, 3, 12, 2},
	{3, 8, 6, 7, 4, 10, 3, 11, 6, 7, 2, 10, 5, 12, 4, 10, 2, 9, 1, 7, 8, 5, 9, 3, 12, 6, 11, 5, 9, 8, 13, 4, 11, 12, 6, 8, 2, 9, 10, 5, 12, 4, 11, 7},
}

// reels lengths [43, 43, 43, 43, 43], total reshuffles 147008443
// symbols: 49.178(lined) + 9.6086(scatter) = 58.786938%
// average plain freespins at 1st iteration: 124.8
// average multiplier at free games: 7.2
// free games 448443, q = 0.3807, sq = 1.4913
// free games frequency: 1/327.82
// RTP = rtpsym + q*sq*rtpsym = 58.787 + 33.376 = 92.162519%
var Reels92 = game.Reels5x{
	{5, 9, 6, 12, 3, 10, 5, 8, 6, 12, 2, 7, 3, 8, 4, 9, 8, 1, 11, 4, 7, 5, 12, 6, 11, 9, 2, 7, 8, 10, 3, 13, 10, 4, 11, 9, 5, 11, 12, 2, 7, 6, 10},
	{9, 2, 7, 5, 10, 6, 9, 12, 5, 8, 11, 6, 10, 3, 11, 1, 7, 5, 9, 2, 7, 4, 12, 7, 8, 4, 9, 5, 10, 12, 2, 8, 6, 11, 12, 3, 13, 6, 10, 4, 11, 8, 3},
	{10, 5, 13, 4, 11, 3, 12, 4, 10, 8, 2, 7, 3, 9, 6, 12, 9, 5, 8, 7, 6, 11, 2, 8, 4, 12, 5, 11, 6, 7, 10, 1, 9, 6, 7, 2, 8, 3, 10, 5, 12, 11, 9},
	{2, 7, 6, 9, 3, 8, 2, 13, 7, 2, 11, 6, 9, 5, 7, 1, 12, 9, 5, 10, 8, 6, 10, 8, 3, 11, 4, 10, 12, 3, 8, 5, 9, 4, 11, 6, 7, 12, 4, 11, 12, 5, 10},
	{3, 11, 6, 10, 9, 4, 10, 7, 5, 11, 4, 10, 12, 6, 13, 5, 11, 2, 7, 3, 8, 1, 12, 2, 8, 5, 9, 6, 10, 11, 4, 9, 2, 12, 5, 7, 3, 8, 6, 9, 7, 8, 12},
}

// reels lengths [42, 42, 42, 42, 42], total reshuffles 130691232
// symbols: 47.535(lined) + 10.076(scatter) = 57.610673%
// average plain freespins at 1st iteration: 124.8
// average multiplier at free games: 7.2
// free games 426708, q = 0.40747, sq = 1.5447
// free games frequency: 1/306.28
// RTP = rtpsym + q*sq*rtpsym = 57.611 + 36.261 = 93.871960%
var Reels94 = game.Reels5x{
	{4, 12, 5, 9, 4, 7, 6, 13, 5, 8, 1, 12, 10, 6, 8, 9, 12, 10, 3, 12, 11, 3, 12, 7, 8, 2, 11, 6, 7, 8, 5, 9, 11, 2, 10, 4, 7, 6, 10, 11, 5, 9},
	{5, 12, 9, 4, 8, 2, 11, 6, 9, 1, 7, 6, 10, 4, 11, 7, 12, 11, 6, 12, 8, 13, 6, 7, 10, 5, 9, 3, 10, 5, 12, 10, 4, 8, 2, 12, 3, 9, 5, 8, 11, 7},
	{2, 10, 3, 12, 5, 11, 1, 7, 2, 11, 12, 5, 7, 6, 10, 4, 9, 13, 6, 12, 9, 6, 8, 10, 4, 11, 6, 9, 3, 8, 12, 5, 8, 4, 11, 7, 10, 5, 8, 9, 12, 7},
	{12, 5, 7, 9, 12, 1, 11, 7, 6, 11, 5, 12, 6, 11, 2, 12, 10, 5, 12, 3, 9, 4, 10, 6, 8, 2, 7, 4, 8, 9, 4, 10, 3, 8, 11, 5, 10, 13, 8, 7, 9, 6},
	{4, 12, 6, 9, 4, 13, 7, 11, 6, 10, 8, 12, 5, 10, 12, 1, 11, 5, 10, 9, 5, 11, 8, 4, 7, 6, 11, 2, 8, 6, 7, 3, 8, 2, 12, 3, 10, 5, 9, 12, 7, 9},
}

// reels lengths [42, 42, 42, 42, 42], total reshuffles 130691232
// symbols: 49.047(lined) + 10.076(scatter) = 59.122719%
// average plain freespins at 1st iteration: 124.8
// average multiplier at free games: 7.2
// free games 426708, q = 0.40747, sq = 1.5447
// free games frequency: 1/306.28
// RTP = rtpsym + q*sq*rtpsym = 59.123 + 37.213 = 96.335718%
var Reels96 = game.Reels5x{
	{1, 9, 4, 10, 3, 8, 5, 9, 6, 11, 4, 9, 12, 5, 7, 12, 4, 10, 11, 6, 12, 7, 5, 8, 3, 11, 8, 7, 4, 11, 6, 13, 2, 10, 5, 7, 2, 8, 6, 12, 9, 10},
	{4, 9, 11, 4, 8, 3, 11, 1, 7, 5, 8, 6, 11, 8, 7, 5, 9, 2, 12, 6, 9, 7, 10, 6, 8, 13, 2, 9, 6, 12, 4, 10, 11, 5, 12, 4, 7, 10, 3, 12, 5, 10},
	{5, 8, 4, 12, 2, 9, 7, 4, 8, 1, 11, 6, 10, 9, 7, 5, 10, 6, 9, 3, 8, 12, 4, 11, 6, 9, 2, 12, 3, 10, 5, 13, 4, 7, 8, 5, 11, 12, 6, 10, 7, 11},
	{6, 10, 4, 8, 3, 10, 2, 8, 5, 9, 6, 11, 8, 4, 12, 6, 7, 12, 3, 10, 5, 8, 1, 7, 2, 11, 6, 9, 7, 4, 10, 12, 5, 9, 12, 7, 5, 11, 4, 13, 9, 11},
	{9, 7, 1, 8, 4, 9, 2, 11, 7, 6, 12, 4, 10, 8, 6, 13, 12, 5, 9, 3, 8, 10, 6, 8, 9, 4, 11, 5, 7, 3, 12, 4, 10, 6, 11, 2, 10, 5, 12, 11, 7, 5},
}

// reels lengths [42, 42, 42, 42, 42], total reshuffles 130691232
// symbols: 49.605(lined) + 10.076(scatter) = 59.680835%
// average plain freespins at 1st iteration: 124.8
// average multiplier at free games: 7.2
// free games 426708, q = 0.40747, sq = 1.5447
// free games frequency: 1/306.28
// RTP = rtpsym + q*sq*rtpsym = 59.681 + 37.564 = 97.245122%
var Reels97 = game.Reels5x{
	{7, 5, 10, 8, 4, 9, 2, 7, 12, 9, 3, 11, 6, 9, 10, 2, 8, 1, 12, 5, 7, 4, 11, 5, 12, 6, 13, 7, 6, 11, 10, 12, 9, 4, 8, 3, 11, 2, 10, 12, 8, 3},
	{8, 2, 10, 11, 6, 10, 12, 8, 9, 7, 5, 12, 4, 11, 9, 3, 7, 6, 12, 4, 8, 3, 10, 7, 1, 11, 3, 12, 5, 9, 2, 10, 8, 4, 9, 11, 6, 13, 2, 7, 12, 5},
	{5, 13, 4, 12, 6, 7, 3, 9, 5, 11, 1, 10, 8, 9, 10, 2, 7, 12, 2, 11, 12, 2, 9, 4, 11, 7, 4, 8, 11, 3, 9, 12, 6, 10, 3, 12, 6, 8, 5, 7, 10, 8},
	{9, 5, 11, 9, 10, 11, 4, 9, 2, 12, 6, 7, 5, 8, 12, 6, 7, 2, 10, 12, 11, 10, 3, 12, 5, 11, 2, 8, 3, 7, 13, 10, 4, 7, 3, 8, 4, 12, 6, 9, 1, 8},
	{3, 12, 11, 6, 12, 5, 11, 3, 12, 6, 8, 5, 13, 2, 7, 4, 9, 6, 10, 5, 9, 10, 11, 4, 12, 2, 7, 9, 1, 8, 10, 4, 7, 10, 2, 8, 12, 3, 11, 7, 8, 9},
}

// reels lengths [41, 41, 41, 41, 41], total reshuffles 115856201
// symbols: 48.18(lined) + 10.578(scatter) = 58.758424%
// average plain freespins at 1st iteration: 124.8
// average multiplier at free games: 7.2
// free games 405513, q = 0.43682, sq = 1.6078
// free games frequency: 1/285.7
// RTP = rtpsym + q*sq*rtpsym = 58.758 + 41.266 = 100.024241%
var Reels100 = game.Reels5x{
	{5, 9, 1, 8, 4, 12, 3, 8, 6, 11, 9, 2, 7, 8, 2, 10, 11, 4, 12, 8, 4, 13, 7, 11, 6, 9, 12, 5, 10, 7, 6, 10, 9, 5, 12, 6, 11, 5, 10, 3, 7},
	{7, 5, 8, 6, 13, 10, 7, 6, 11, 2, 9, 11, 5, 12, 8, 5, 10, 4, 12, 9, 11, 2, 7, 5, 12, 6, 10, 4, 11, 3, 8, 6, 9, 4, 8, 9, 1, 12, 7, 3, 10},
	{6, 9, 12, 1, 8, 6, 10, 7, 4, 9, 6, 10, 8, 5, 10, 3, 7, 11, 6, 7, 4, 11, 5, 9, 8, 2, 12, 5, 8, 4, 10, 5, 11, 12, 2, 9, 13, 11, 12, 3, 7},
	{12, 8, 11, 4, 10, 3, 7, 5, 8, 6, 10, 7, 4, 11, 7, 4, 8, 9, 6, 10, 9, 5, 10, 2, 12, 6, 9, 12, 8, 5, 11, 6, 7, 5, 9, 3, 11, 1, 12, 2, 13},
	{6, 7, 3, 10, 6, 8, 5, 11, 9, 2, 10, 7, 6, 8, 10, 2, 9, 3, 8, 5, 9, 12, 6, 7, 5, 13, 4, 12, 5, 11, 12, 4, 8, 7, 11, 4, 9, 1, 12, 11, 10},
}

// reels lengths [41, 41, 41, 41, 41], total reshuffles 115856201
// symbols: 52.855(lined) + 10.578(scatter) = 63.433344%
// average plain freespins at 1st iteration: 124.8
// average multiplier at free games: 7.2
// free games 405513, q = 0.43682, sq = 1.6078
// free games frequency: 1/285.7
// RTP = rtpsym + q*sq*rtpsym = 63.433 + 44.549 = 107.982338%
var Reels108 = game.Reels5x{
	{2, 12, 6, 8, 5, 12, 4, 11, 12, 3, 9, 4, 13, 2, 8, 9, 6, 10, 2, 7, 12, 5, 10, 4, 8, 3, 11, 6, 9, 10, 5, 11, 9, 5, 11, 6, 7, 3, 10, 1, 7},
	{2, 10, 6, 9, 5, 12, 4, 11, 6, 10, 8, 1, 11, 5, 9, 2, 12, 4, 10, 7, 4, 9, 6, 12, 3, 8, 6, 10, 11, 5, 8, 3, 9, 5, 12, 2, 13, 7, 11, 3, 7},
	{1, 12, 2, 10, 4, 11, 6, 12, 5, 11, 3, 12, 10, 5, 9, 4, 7, 3, 12, 6, 10, 8, 9, 4, 10, 2, 7, 6, 13, 11, 3, 9, 6, 7, 5, 11, 2, 8, 9, 5, 8},
	{5, 9, 2, 12, 8, 1, 11, 3, 10, 4, 7, 12, 11, 4, 8, 2, 13, 10, 5, 9, 6, 10, 5, 9, 4, 12, 2, 9, 6, 10, 3, 7, 6, 11, 5, 8, 11, 6, 7, 3, 12},
	{3, 8, 6, 12, 2, 11, 4, 12, 10, 9, 4, 10, 3, 11, 9, 1, 7, 5, 8, 4, 12, 2, 10, 9, 6, 10, 5, 12, 6, 11, 5, 13, 6, 7, 5, 11, 3, 9, 2, 8, 7},
}

// Map with available reels.
var ReelsMap = map[float64]*game.Reels5x{
	87.023352:  &Reels87,
	88.828983:  &Reels89,
	92.162519:  &Reels92,
	93.871960:  &Reels94,
	96.335718:  &Reels96,
	97.245122:  &Reels97,
	100.024241: &Reels100,
	107.982338: &Reels108,
}

func FindReels(mrtp float64) (rtp float64, reels game.Reels) {
	for p, r := range ReelsMap {
		if math.Abs(mrtp-p) < math.Abs(mrtp-rtp) {
			rtp, reels = p, r
		}
	}
	return
}

// Lined payment.
var LinePay = [13][5]float64{
	{0, 10, 250, 2500, 10000}, //  1 wild
	{0, 2, 25, 125, 750},      //  2 dolphin
	{0, 2, 25, 125, 750},      //  3 turtle
	{0, 0, 20, 100, 400},      //  4 fish
	{0, 0, 15, 75, 250},       //  5 seahorse
	{0, 0, 15, 75, 250},       //  6 starfish
	{0, 0, 10, 50, 150},       //  7 ace
	{0, 0, 10, 50, 150},       //  8 king
	{0, 0, 5, 25, 100},        //  9 queen
	{0, 0, 5, 25, 100},        // 10 jack
	{0, 0, 5, 25, 100},        // 11 ten
	{0, 2, 5, 25, 100},        // 12 nine
	{0, 0, 0, 0, 0},           // 13 scatter
}

// Scatters payment.
var ScatPay = [5]float64{0, 2, 5, 20, 500} // 13 scatter

type Seashells struct {
	Sel1 string  `json:"sel1" yaml:"sel1" xml:"sel1"`
	Sel2 string  `json:"sel2" yaml:"sel2" xml:"sel2"`
	Mult float64 `json:"mult" yaml:"mult" xml:"mult"`
	Free int     `json:"free" yaml:"free" xml:"free"`
}

func (s *Seashells) SetupShell(shell string) {
	switch shell {
	case "x5":
		s.Mult += 5
	case "x8":
		s.Mult += 8
	case "7":
		s.Free += 7
	case "10":
		s.Free += 10
	case "15":
		s.Free += 15
	}
}

type Game struct {
	game.Slot5x3 `yaml:",inline"`
	// free spin number
	FS int `json:"fs,omitempty" yaml:"fs,omitempty" xml:"fs,omitempty"`
	// multiplier on freespins
	M float64 `json:"m,omitempty" yaml:"m,omitempty" xml:"m,omitempty"`
}

func NewGame(rtp float64) *Game {
	return &Game{
		Slot5x3: game.Slot5x3{
			RTP: rtp,
			SBL: game.MakeBitNum(25),
			Bet: 1,
		},
		FS: 0,
		M:  0,
	}
}

const wild, scat = 1, 13

var bl = game.BetLinesPlt30

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
			var mm float64 = 1 // mult mode
			if g.FS > 0 {
				mm = g.M
			}
			*wins = append(*wins, game.WinItem{
				Pay:  g.Bet * payl,
				Mult: mw * mm,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyL(numl),
			})
		} else if payw > 0 {
			var mm float64 = 1 // mult mode
			if g.FS > 0 {
				mm = g.M
			}
			*wins = append(*wins, game.WinItem{
				Pay:  g.Bet * payw,
				Mult: mm,
				Sym:  wild,
				Num:  numw,
				Line: li,
				XY:   line.CopyL(numw),
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen game.Screen, wins *game.Wins) {
	if count := screen.ScatNum(scat); count >= 2 {
		var pay, fs = ScatPay[count-1], 0
		var mm float64 = 1 // mult mode
		if g.FS > 0 {
			mm, fs = g.M, 15
		} else if count >= 3 {
			fs = 8
		}
		*wins = append(*wins, game.WinItem{
			Pay:  g.Bet * float64(g.SBL.Num()) * pay,
			Mult: mm,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
			Free: fs,
		})
	}
}

func (g *Game) Spin(screen game.Screen) {
	var _, reels = FindReels(g.RTP)
	screen.Spin(reels)
}

func (g *Game) Spawn(screen game.Screen, wins game.Wins) {
	if g.FS > 0 {
		return
	}
	for i := range wins {
		if wi := &wins[i]; wi.Sym == scat {
			var idx = []string{"x5", "x8", "7", "10", "15"}
			rand.Shuffle(len(idx), func(i, j int) {
				idx[i], idx[j] = idx[j], idx[i]
			})
			var bon = Seashells{
				Sel1: idx[0],
				Sel2: idx[1],
				Mult: 2,
				Free: 8,
			}
			bon.SetupShell(idx[0])
			bon.SetupShell(idx[1])
			wi.Mult = 1
			wi.Free = bon.Free
			wi.Bon = bon
		}
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
		if wi.Sym == scat {
			if g.FS > 0 {
				g.FS += wi.Free
			} else {
				var bon = wi.Bon.(Seashells)
				g.FS = bon.Free
				g.M = bon.Mult
			}
		}
	}
	if g.FS == 0 {
		g.M = 0
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
