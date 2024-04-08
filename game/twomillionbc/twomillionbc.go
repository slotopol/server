package twomillionbc

import (
	"github.com/slotopol/server/game"
)

// reels lengths [19, 20, 20, 20, 101], total reshuffles 15352000
// symbols: 54.944(lined) + 0(scatter) = 54.943884%
// free spins 1964628, q = 0.12797, sq = 1/(1-q) = 1.146752
// acorn bonuses: count 152000, rtp = 10.891089%
// diamond lion bonuses: count 8080, rtp = 9.210526%
// RTP = 54.944(sym) + 10.891(acorn) + 9.2105(dl) + 0.12797*138.6(fg) = 92.783055%
var ReelsReg93 = game.Reels5x{
	{9, 10, 9, 3, 7, 8, 5, 8, 6, 4, 6, 2, 10, 1, 13, 11, 4, 5, 7},
	{13, 6, 2, 6, 5, 4, 8, 8, 5, 3, 4, 10, 9, 11, 13, 7, 9, 10, 7, 1},
	{8, 4, 5, 6, 5, 4, 6, 13, 11, 8, 2, 9, 13, 10, 9, 7, 7, 10, 1, 3},
	{3, 9, 10, 7, 5, 6, 8, 10, 9, 6, 8, 2, 13, 5, 13, 4, 11, 7, 4, 1},
	{9, 2, 9, 7, 8, 1, 9, 4, 13, 13, 9, 13, 5, 12, 5, 8, 4, 6, 7, 8, 1, 10, 11, 3, 5, 7, 13, 7, 6, 7, 1, 7, 7, 4, 5, 6, 6, 9, 6, 5, 10, 11, 9, 5, 13, 1, 4, 7, 4, 10, 9, 3, 8, 4, 10, 5, 1, 10, 8, 7, 13, 11, 6, 2, 2, 5, 6, 11, 6, 8, 4, 5, 3, 10, 8, 13, 5, 13, 9, 4, 10, 2, 10, 4, 6, 7, 8, 2, 13, 11, 6, 8, 13, 3, 10, 4, 10, 8, 3, 9, 9},
}

// reels lengths [20, 21, 21, 21, 93], total reshuffles 17225460
// symbols: 60.224(lined) + 0(scatter) = 60.223855%
// free spins 1811808, q = 0.10518, sq = 1/(1-q) = 1.117546
// acorn bonuses: count 185220, rtp = 11.827957%
// diamond lion bonuses: count 7812, rtp = 7.936508%
// RTP = 60.224(sym) + 11.828(acorn) + 7.9365(dl) + 0.10518*138.6(fg) = 94.567051%
var ReelsReg95 = game.Reels5x{
	{8, 6, 8, 7, 4, 10, 6, 1, 2, 9, 4, 10, 9, 3, 13, 11, 7, 5, 3, 5},
	{13, 7, 1, 10, 7, 8, 6, 5, 10, 2, 9, 3, 4, 9, 5, 3, 4, 11, 6, 13, 8},
	{4, 10, 7, 13, 10, 5, 6, 1, 8, 3, 2, 8, 6, 9, 5, 9, 13, 11, 7, 3, 4},
	{2, 9, 5, 7, 9, 1, 4, 13, 6, 8, 10, 5, 6, 7, 10, 3, 11, 13, 3, 8, 4},
	{4, 5, 7, 9, 11, 13, 4, 8, 2, 6, 2, 10, 5, 9, 3, 10, 3, 9, 3, 4, 8, 6, 1, 3, 10, 1, 8, 7, 4, 5, 8, 1, 13, 6, 7, 6, 6, 7, 13, 1, 8, 6, 1, 3, 13, 13, 2, 5, 1, 10, 13, 2, 2, 5, 8, 2, 8, 4, 13, 12, 11, 6, 10, 2, 1, 6, 11, 9, 5, 3, 4, 3, 5, 4, 7, 9, 7, 9, 11, 2, 10, 8, 10, 9, 13, 7, 10, 4, 3, 9, 1, 5, 7},
}

// reels lengths [43, 43, 43, 43, 43], total reshuffles 147008443
// symbols: 93.09(lined) + 0(scatter) = 93.090022%
// free spins 48274380, q = 0.32838, sq = 1/(1-q) = 1.488933
// RTP = sq*rtp(sym) = 1.4889*93.09 = 138.604842%
var ReelsBon = game.Reels5x{
	{4, 10, 3, 7, 1, 9, 5, 7, 5, 9, 11, 2, 5, 1, 6, 1, 10, 4, 4, 6, 11, 5, 3, 7, 6, 2, 3, 6, 8, 4, 10, 1, 10, 7, 8, 2, 11, 8, 2, 9, 3, 9, 8},
	{8, 2, 1, 3, 10, 11, 1, 4, 9, 7, 7, 6, 5, 10, 5, 2, 5, 6, 9, 4, 5, 11, 7, 9, 8, 11, 3, 7, 8, 10, 2, 3, 1, 2, 1, 4, 8, 9, 10, 6, 3, 6, 4},
	{2, 7, 9, 2, 8, 4, 3, 2, 6, 10, 11, 10, 3, 7, 5, 10, 2, 6, 3, 9, 4, 4, 4, 8, 5, 1, 1, 11, 7, 9, 8, 6, 9, 5, 5, 10, 1, 8, 1, 6, 7, 3, 11},
	{8, 2, 3, 3, 5, 11, 8, 9, 5, 1, 3, 6, 1, 9, 1, 6, 4, 10, 7, 3, 8, 7, 5, 7, 4, 2, 10, 6, 10, 10, 11, 9, 6, 2, 7, 8, 4, 5, 11, 1, 9, 2, 4},
	{11, 3, 1, 2, 7, 4, 3, 10, 3, 9, 6, 2, 9, 1, 1, 7, 4, 10, 6, 8, 11, 8, 5, 9, 7, 10, 5, 11, 8, 2, 4, 6, 8, 1, 5, 10, 3, 2, 7, 5, 9, 4, 6},
}

// Map with available reels.
var ReelsMap = map[string]*game.Reels5x{
	"93":  &ReelsReg93,
	"95":  &ReelsReg95,
	"bon": &ReelsBon,
}

// Lined payment.
var LinePay = [13][5]float64{
	{0, 30, 100, 300, 500}, //  1 girl
	{0, 15, 75, 200, 400},  //  2 lion
	{0, 10, 60, 150, 300},  //  3 bee
	{0, 5, 50, 125, 250},   //  4 stone
	{0, 5, 40, 100, 200},   //  5 wheel
	{0, 2, 30, 90, 150},    //  6 club
	{0, 0, 25, 75, 125},    //  7 chaplet
	{0, 0, 20, 60, 100},    //  8 gold
	{0, 0, 15, 50, 75},     //  9 vase
	{0, 0, 10, 25, 50},     // 10 ruby
	{0, 0, 0, 0, 0},        // 11 fire
	{0, 0, 0, 0, 0},        // 12 acorn
	{0, 0, 40, 100, 200},   // 13 diamond
}

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 4, 12, 20} // 11 fire

const (
	acbn = 1 // acorn bonus
	dlbn = 2 // diamond lion bonus
	jid  = 1 // jackpot ID
)

// Jackpot win combinations.
var Jackpot = [13][5]int{
	{0, 0, 0, 0, 0}, //  1 girl
	{0, 0, 0, 0, 0}, //  2 lion
	{0, 0, 0, 0, 0}, //  3 bee
	{0, 0, 0, 0, 0}, //  4 stone
	{0, 0, 0, 0, 0}, //  5 wheel
	{0, 0, 0, 0, 0}, //  6 club
	{0, 0, 0, 0, 0}, //  7 chaplet
	{0, 0, 0, 0, 0}, //  8 gold
	{0, 0, 0, 0, 0}, //  9 vase
	{0, 0, 0, 0, 0}, // 10 ruby
	{0, 0, 0, 0, 0}, // 11 fire
	{0, 0, 0, 0, 0}, // 12 acorn
	{0, 0, 0, 0, 0}, // 13 diamond
}

type Game struct {
	game.Slot5x3 `yaml:",inline"`
	FS           int     `json:"fs" yaml:"fs" xml:"fs"` // free spin number
	AN           int     `json:"an" yaml:"an" xml:"an"` // acorns number
	AB           float64 `json:"ab" yaml:"ab" xml:"ab"` // acorns bet
}

func NewGame(rd string) *Game {
	return &Game{
		Slot5x3: game.Slot5x3{
			RD:  rd,
			BLI: "bs30",
			SBL: game.MakeSblNum(30),
			Bet: 1,
		},
		FS: 0,
	}
}

const scat, acorn, diamond = 11, 12, 13

func (g *Game) Scanner(screen game.Screen, ws *game.WinScan) {
	g.ScanLined(screen, ws)
	g.ScanScatters(screen, ws)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen game.Screen, ws *game.WinScan) {
	var bl = game.BetLines5x[g.BLI]
	for li := g.SBL.Next(0); li != 0; li = g.SBL.Next(li) {
		var line = bl.Line(li)

		var syml, numl = screen.At(1, line.At(1)), 1
		for x := 2; x <= 5; x++ {
			var sx = screen.At(x, line.At(x))
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
				XY:   line.CopyN(numl),
			})
		}
		if syml == diamond && numl >= 3 {
			ws.Wins = append(ws.Wins, game.WinItem{
				Mult: 1,
				Sym:  diamond,
				Num:  numl,
				Line: li,
				XY:   line.CopyN(numl),
				BID:  dlbn,
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen game.Screen, ws *game.WinScan) {
	if count := screen.ScatNum(scat); count >= 3 {
		var fs = ScatFreespin[count-1]
		ws.Wins = append(ws.Wins, game.WinItem{
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
			Free: fs,
		})
	}

	if screen.At(5, 1) == acorn || screen.At(5, 2) == acorn || screen.At(5, 3) == acorn {
		if (g.AN+1)%3 == 0 {
			ws.Wins = append(ws.Wins, game.WinItem{
				Mult: 1,
				Sym:  acorn,
				Num:  1,
				BID:  acbn,
			})
		}
	}
}

func (g *Game) Spin(screen game.Screen) {
	if g.FS == 0 {
		screen.Spin(ReelsMap[g.RD])
	} else {
		screen.Spin(&ReelsBon)
	}
}

func (g *Game) Spawn(screen game.Screen, sw *game.WinScan) {
	for i, wi := range sw.Wins {
		switch wi.BID {
		case acbn:
			sw.Wins[i].Pay = AcornSpawn(g.AB + g.Bet*float64(g.SBL.Num()))
		case dlbn:
			sw.Wins[i].Pay = DiamondLionSpawn(g.Bet)
		}
	}
}

func (g *Game) Apply(screen game.Screen, sw *game.WinScan) {
	if screen.At(5, 1) == acorn || screen.At(5, 2) == acorn || screen.At(5, 3) == acorn {
		g.AN++
		g.AN %= 3
		if g.AN > 0 {
			g.AB += g.Bet * float64(g.SBL.Num())
		} else {
			g.AB = 0
		}
	}

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
