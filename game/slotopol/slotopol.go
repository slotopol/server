package slotopol

import (
	"github.com/slotopol/server/game"
)

// Original reels.
// reels lengths [32, 32, 32, 32, 32], total reshuffles 33554432
// symbols: 48.848(lined) + 39.546(scatter) = 88.394135%
// spin9 bonuses: count 2700, rtp = 7.676482%
// monopoly bonuses: count 4320, rtp = 3.689938%
// jackpots: count 32, frequency 1/1048576
// RTP = 88.394(sym) + 7.6765(mje9) + 3.6899(mjm) = 99.760556%
var Reels100 = game.Reels5x{
	{13, 1, 5, 12, 13, 11, 12, 11, 13, 8, 2, 12, 13, 3, 4, 6, 13, 2, 5, 10, 13, 9, 7, 8, 13, 10, 7, 9, 13, 3, 4, 6},
	{9, 5, 10, 13, 9, 6, 3, 4, 13, 2, 12, 8, 12, 13, 11, 12, 11, 13, 5, 7, 10, 6, 3, 4, 13, 2, 12, 8, 13, 7, 1, 12},
	{12, 13, 11, 12, 11, 13, 5, 10, 9, 7, 1, 12, 13, 3, 8, 6, 12, 13, 8, 4, 12, 2, 5, 10, 13, 7, 2, 13, 6, 3, 4, 9},
	{12, 1, 2, 13, 6, 5, 12, 4, 8, 12, 13, 3, 10, 9, 7, 13, 11, 11, 11, 11, 13, 5, 12, 9, 8, 6, 13, 3, 10, 2, 7, 4},
	{13, 11, 13, 12, 6, 4, 12, 3, 2, 5, 12, 10, 7, 12, 8, 1, 9, 12, 8, 9, 12, 4, 3, 12, 2, 5, 12, 10, 7, 13, 12, 6},
}

// Map with available reels.
var ReelsMap = map[string]*game.Reels5x{
	"100": &Reels100,
}

// Lined payment.
var LinePay = [13][5]float64{
	{0, 0, 0, 0, 0},           //  1 dollar
	{0, 2, 5, 15, 100},        //  2 cherry
	{0, 2, 5, 15, 100},        //  3 plum
	{0, 0, 5, 15, 100},        //  4 wmelon
	{0, 0, 5, 15, 100},        //  5 grapes
	{0, 0, 5, 15, 100},        //  6 ananas
	{0, 0, 5, 15, 100},        //  7 lemon
	{0, 0, 5, 15, 100},        //  8 drink
	{0, 2, 5, 15, 100},        //  9 palm
	{0, 2, 5, 15, 100},        // 10 yacht
	{0, 10, 100, 2000, 10000}, // 11 eldorado
	{0, 0, 0, 0, 0},           // 12 spin
	{0, 0, 0, 0, 0},           // 13 dice
}

// Scatters payment.
var ScatPay = [5]float64{0, 5, 8, 20, 1000} // 1 dollar

const (
	mje1 = 1 // Eldorado9
	mje3 = 2 // Eldorado9
	mje6 = 3 // Eldorado9
	mje9 = 4 // Eldorado9
	mjm  = 5 // Monopoly
	mjc  = 6 // Champagne
)

// Lined bonus games
var LineBonus = [13][5]int{
	{0, 0, 0, 0, 0},    //  1
	{0, 0, 0, 0, 0},    //  2
	{0, 0, 0, 0, 0},    //  3
	{0, 0, 0, 0, 0},    //  4
	{0, 0, 0, 0, 0},    //  5
	{0, 0, 0, 0, 0},    //  6
	{0, 0, 0, 0, 0},    //  7
	{0, 0, 0, 0, 0},    //  8
	{0, 0, 0, 0, 0},    //  9
	{0, 0, 0, 0, 0},    // 10
	{0, 0, 0, 0, 0},    // 11
	{0, 0, 0, 0, mje9}, // 12 Eldorado9
	{0, 0, 0, 0, mjm},  // 13 Monopoly
}

const (
	jid = 1 // jackpot ID
)

// Jackpot win combinations.
var Jackpot = [13][5]int{
	{0, 0, 0, 0, 0},   //  1 dollar
	{0, 0, 0, 0, 0},   //  2 cherry
	{0, 0, 0, 0, 0},   //  3 plum
	{0, 0, 0, 0, 0},   //  4 wmelon
	{0, 0, 0, 0, 0},   //  5 grapes
	{0, 0, 0, 0, 0},   //  6 ananas
	{0, 0, 0, 0, 0},   //  7 lemon
	{0, 0, 0, 0, 0},   //  8 drink
	{0, 0, 0, 0, 0},   //  9 palm
	{0, 0, 0, 0, 0},   // 10 yacht
	{0, 0, 0, 0, jid}, // 11 eldorado
	{0, 0, 0, 0, 0},   // 12 spin
	{0, 0, 0, 0, 0},   // 13 dice
}

type Game struct {
	game.Slot5x3 `yaml:",inline"`
}

func NewGame(rd string) *Game {
	return &Game{
		Slot5x3: game.Slot5x3{
			RD:  rd,
			BLI: "mgj",
			SBL: game.MakeSBL(1),
			Bet: 1,
		},
	}
}

// Not from lined paytable.
var special = [13]bool{
	true,  //  1
	false, //  2
	false, //  3
	false, //  4
	false, //  5
	false, //  6
	false, //  7
	false, //  8
	false, //  9
	false, // 10
	false, // 11
	true,  // 12
	true,  // 13
}

const wild, scat = 11, 1

func (g *Game) Scanner(screen game.Screen, ws *game.WinScan) {
	g.ScanLined(screen, ws)
	g.ScanScatters(screen, ws)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen game.Screen, ws *game.WinScan) {
	var bl = game.BetLines5x[g.BLI]
	for li := g.SBL.Next(0); li != 0; li = g.SBL.Next(li) {
		var line = bl.Line(li)

		var numw, numl = 0, 5
		var syml game.Sym
		var mw float64 = 1 // mult wild
		for x := 1; x <= 5; x++ {
			var sx = screen.At(x, line.At(x))
			if sx == wild {
				if syml == 0 {
					numw = x
				} else if special[syml-1] {
					numl = x - 1
					break
				}
				mw = 2
			} else if numw > 0 && special[sx-1] {
				numl = x - 1
				break
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
			ws.Wins = append(ws.Wins, game.WinItem{
				Pay:  g.Bet * payl,
				Mult: mw,
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
		} else if syml > 0 && numl > 0 && LineBonus[syml-1][numl-1] > 0 {
			ws.Wins = append(ws.Wins, game.WinItem{
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyN(numl),
				BID:  LineBonus[syml-1][numl-1],
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen game.Screen, ws *game.WinScan) {
	if count := screen.ScatNum(scat); count >= 2 {
		var pay = ScatPay[count-1]
		ws.Wins = append(ws.Wins, game.WinItem{
			Pay:  g.Bet * pay, // independent from selected lines
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
		})
	}
}

func (g *Game) Spin(screen game.Screen) {
	screen.Spin(ReelsMap[g.RD])
}

func (g *Game) Spawn(screen game.Screen, sw *game.WinScan) {
	for i, wi := range sw.Wins {
		switch wi.BID {
		case mje1:
			sw.Wins[i].Bon, sw.Wins[i].Pay = EldoradoSpawn(g.Bet, 1)
		case mje3:
			sw.Wins[i].Bon, sw.Wins[i].Pay = EldoradoSpawn(g.Bet, 3)
		case mje6:
			sw.Wins[i].Bon, sw.Wins[i].Pay = EldoradoSpawn(g.Bet, 6)
		case mje9:
			sw.Wins[i].Bon, sw.Wins[i].Pay = EldoradoSpawn(g.Bet, 9)
		case mjm:
			sw.Wins[i].Bon, sw.Wins[i].Pay = MonopolySpawn(g.Bet)
		}
	}
}

func (g *Game) SetReels(rd string) error {
	if _, ok := ReelsMap[rd]; !ok {
		return game.ErrNoReels
	}
	g.RD = rd
	return nil
}
