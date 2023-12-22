package slotopol

import (
	"context"
	"fmt"
	"time"

	"github.com/schwarzlichtbezirk/slot-srv/game"
)

// Lined payment.
var LinePay = [13][5]int{
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

// Reels sets.
var Reels = game.Reels5x{
	{13, 1, 5, 12, 13, 11, 12, 11, 13, 8, 2, 12, 13, 3, 4, 6, 13, 2, 5, 10, 13, 9, 7, 8, 13, 10, 7, 9, 13, 3, 4, 6}, // 1 reel
	{9, 5, 10, 13, 9, 6, 3, 4, 13, 2, 12, 8, 12, 13, 11, 12, 11, 13, 5, 7, 10, 6, 3, 4, 13, 2, 12, 8, 13, 7, 1, 12}, // 2 reel
	{12, 13, 11, 12, 11, 13, 5, 10, 9, 7, 1, 12, 13, 3, 8, 6, 12, 13, 8, 4, 12, 2, 5, 10, 13, 7, 2, 13, 6, 3, 4, 9}, // 3 reel
	{12, 1, 2, 13, 6, 5, 12, 4, 8, 12, 13, 3, 10, 9, 7, 13, 11, 11, 11, 11, 13, 5, 12, 9, 8, 6, 13, 3, 10, 2, 7, 4}, // 4 reel
	{13, 11, 13, 12, 6, 4, 12, 3, 2, 5, 12, 10, 7, 12, 8, 1, 9, 12, 8, 9, 12, 4, 3, 12, 2, 5, 12, 10, 7, 13, 12, 6}, // 5 reel
}

// Scatters payment.
var ScatPay = [5]int{0, 5, 8, 20, 1000} // 1 dollar

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 0, 0, 0} // 1 dollar

const (
	mje9 = 1 // Eldorado9
	mjm  = 2 // Monopoly
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
	bet  int          // bet value
	free int          // free spin number
	sbl  []int        // selected bet lines
	bls  game.Lineset // bet lines set
}

func NewGame() *Game {
	return &Game{
		bet:  1,
		free: 0,
		sbl:  []int{1},
		bls:  &game.BetLinesMgj,
	}
}

func (g *Game) NewScreen() game.Screen {
	return &game.Screen5x3{}
}

func (g *Game) Bet() int {
	return g.bet
}

func (g *Game) SetBet(bet int) error {
	g.bet = bet
	return nil
}

func (g *Game) Line() []int {
	return g.sbl
}

func (g *Game) SetLine(sbl []int) error {
	g.sbl = sbl
	return nil
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
	var mm = 1
	if g.free > 0 {
		mm = 2
	}

	for _, li := range g.sbl {
		var line = g.bls.Line(li)

		var xy = []int{0, 0, 0, 0, 0}
		var cntw, cntl = 0, 5
		var sl, m = 0, mm
		for x := 1; x <= 5; x++ {
			var sx = screen.At(x, line[x-1])
			if sx == wild {
				if sl == 0 {
					cntw = x
				} else if special[sl-1] {
					cntl = x - 1
					break
				}
				m = 2 * mm
			} else if cntw > 0 && special[sx-1] {
				cntl = x - 1
				break
			} else if sl == 0 && sx != scat {
				sl = sx
			} else if sx != sl {
				cntl = x - 1
				break
			}
			xy[x-1] = line[x-1]
		}

		var payw, payl int
		if cntw > 0 {
			payw = LinePay[wild-1][cntw-1]
		}
		if cntl > 0 && sl > 0 {
			payl = LinePay[sl-1][cntl-1]
		}
		if payw > 0 && payl > 0 {
			if payw*mm < payl*m {
				payw = 0
			} else {
				payl = 0
				// delete non-wild line
				for x := cntw; x < cntl; x++ {
					xy[x] = 0
				}
			}
		}
		if payl > 0 {
			ws.Wins = append(ws.Wins, game.WinItem{
				Pay:  g.bet * payl,
				Mult: m,
				Sym:  sl,
				Num:  cntl,
				Line: li,
				XY:   xy,
			})
		} else if payw > 0 {
			ws.Wins = append(ws.Wins, game.WinItem{
				Pay:  g.bet * payw,
				Mult: mm,
				Sym:  wild,
				Num:  cntw,
				Line: li,
				XY:   xy,
				Jack: Jackpot[wild-1][cntw-1],
			})
		} else if sl > 0 && cntl > 0 && LineBonus[sl-1][cntl-1] > 0 {
			ws.Wins = append(ws.Wins, game.WinItem{
				Mult: 1,
				Sym:  sl,
				Num:  cntl,
				Line: li,
				XY:   xy,
				BID:  LineBonus[sl-1][cntl-1],
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen game.Screen, ws *game.WinScan) {
	var xy = []int{0, 0, 0, 0, 0}
	var count = 0
	for x := 1; x <= 5; x++ {
		for y := 1; y <= 3; y++ {
			if screen.At(x, y) == scat {
				xy[x-1] = y
				count++
				break
			}
		}
	}

	if count > 1 {
		if ScatPay[count-1] > 0 || ScatFreespin[count-1] > 0 {
			ws.Wins = append(ws.Wins, game.WinItem{
				Pay:  g.bet * ScatPay[count-1],
				Mult: 1,
				Sym:  scat,
				Num:  count,
				XY:   xy,
				Free: ScatFreespin[count-1],
			})
		}
	}
}

func (g *Game) Spawn(screen game.Screen, sw *game.WinScan) {
	for i, wi := range sw.Wins {
		switch wi.BID {
		case mje9:
			sw.Wins[i].Bon, sw.Wins[i].Pay = Eldorado9Spawn(g.bet)
		case mjm:
			sw.Wins[i].Bon, sw.Wins[i].Pay = MonopolySpawn(g.bet)
		}
	}
}

func CalcStat() {
	var g = NewGame()
	var sbl = float64(len(g.sbl))
	var s game.Stat
	var t0 = time.Now()
	func() {
		var ctx, cancel = context.WithCancel(context.Background())
		defer cancel()
		go s.Progress(ctx, time.NewTicker(2*time.Second), sbl, float64(Reels.Reshuffles()))
		s.Rotator5x(ctx, g, &Reels)
	}()
	var dur = time.Since(t0)
	var n = float64(s.Reshuffles)
	var lp, sp = float64(s.LinePay) / n / sbl * 100, float64(s.ScatPay) / n * 100
	var rtpsym = lp + sp
	var Mmjs9, qmjs9 = 106.0 * 9, float64(s.BonusCount[mje9]) / n / sbl
	var rtpmjs9 = Mmjs9 * qmjs9 * 100
	var Mmjm, qmjm = 286.60597422268, float64(s.BonusCount[mjm]) / n / sbl
	var rtpmjm = Mmjm * qmjm * 100
	fmt.Printf("selected %d lines, reshuffles %d, time spent %v\n", len(g.sbl), s.Reshuffles, dur)
	fmt.Printf("symbols: %2.5g(lined) + %2.5g(scatter) = %g%%\n", lp, sp, rtpsym)
	fmt.Printf("mje9 bonuses: count %d, q = %2.5g, rtp = %g%%\n", s.BonusCount[mje9], qmjs9, rtpmjs9)
	fmt.Printf("mjm bonuses: count %d, q = %2.5g, rtp = %g%%\n", s.BonusCount[mjm], qmjm, rtpmjm)
	fmt.Printf("jackpots: count %d, frequency 1/%d\n", s.JackCount[jid], int(n/float64(s.JackCount[jid])))
	fmt.Printf("RTP = %2.5g(sym) + %2.5g(mje9) + %2.5g(mjm) = %g%%\n", rtpsym, rtpmjs9, rtpmjm, rtpsym+rtpmjs9+rtpmjm)
}
