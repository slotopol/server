package slotopoldeluxe

import (
	"context"
	"fmt"
	"time"

	"github.com/schwarzlichtbezirk/slot-srv/game"
	"github.com/schwarzlichtbezirk/slot-srv/game/slotopol"
)

// Original reels.
// symbols: 81.213(lined) + 2.7777(scatter) = 83.990312%
// spin1 bonuses: count1 28672, rtp = 9.057617%
// spin3 bonuses: count3 3328, rtp = 3.153992%
// spin6 bonuses: count6 768, rtp = 1.455688%
// monopoly bonuses: count 6912, rtp = 5.903901%
// jackpots: count 32, frequency 1/1048576
// RTP = 83.99(sym) + 13.667(mje) + 5.9039(mjm) = 103.561510%
var ReelsOrig = game.Reels5x{
	{1, 2, 13, 4, 13, 3, 13, 5, 9, 7, 8, 13, 10, 13, 12, 11, 13, 12, 11, 13, 13, 2, 4, 5, 2, 6, 7, 9, 8, 3, 10, 6},
	{13, 2, 12, 9, 13, 2, 5, 6, 9, 7, 13, 10, 12, 13, 11, 12, 13, 11, 12, 13, 3, 4, 5, 2, 8, 7, 10, 4, 6, 8, 3, 1},
	{2, 1, 12, 3, 4, 5, 2, 6, 7, 10, 8, 4, 3, 13, 12, 11, 13, 12, 11, 13, 12, 3, 5, 13, 9, 6, 7, 10, 13, 8, 9, 13},
	{1, 2, 12, 4, 3, 6, 4, 2, 7, 8, 10, 5, 13, 12, 11, 13, 12, 11, 13, 12, 9, 4, 5, 3, 13, 6, 10, 7, 13, 8, 9, 13},
	{1, 2, 12, 4, 3, 5, 12, 6, 7, 10, 8, 12, 9, 13, 12, 11, 13, 12, 11, 13, 12, 3, 4, 5, 2, 6, 7, 10, 13, 8, 9, 5},
}

// reels lengths [32, 32, 32, 32, 32], total reshuffles 33554432
// symbols: 79.254(lined) + 2.7777(scatter) = 82.031429%
// spin1 bonuses: count1 20736, rtp = 6.550598%
// spin3 bonuses: count3 3120, rtp = 2.956867%
// spin6 bonuses: count6 720, rtp = 1.364708%
// monopoly bonuses: count 6720, rtp = 5.739904%
// jackpots: count 32, frequency 1/1048576
// RTP = 82.031(sym) + 10.872(mje) + 5.7399(mjm) = 98.643506%
var Reels986 = game.Reels5x{
	{1, 5, 13, 4, 13, 3, 13, 5, 9, 7, 8, 13, 10, 13, 12, 11, 13, 12, 11, 13, 13, 2, 4, 5, 2, 6, 7, 9, 8, 3, 10, 6},
	{13, 2, 12, 9, 13, 4, 5, 6, 9, 7, 13, 10, 12, 13, 11, 13, 13, 11, 12, 13, 3, 4, 5, 2, 8, 7, 10, 4, 6, 8, 3, 1},
	{2, 1, 12, 3, 4, 5, 2, 6, 7, 10, 8, 4, 5, 13, 12, 11, 13, 12, 11, 13, 12, 3, 5, 13, 9, 6, 7, 10, 13, 8, 9, 13},
	{1, 2, 12, 4, 3, 6, 4, 2, 7, 8, 10, 5, 13, 12, 11, 13, 12, 11, 13, 12, 9, 4, 5, 3, 12, 6, 10, 7, 13, 8, 9, 13},
	{1, 2, 12, 4, 3, 5, 12, 6, 7, 10, 8, 12, 9, 13, 12, 11, 13, 12, 11, 13, 12, 3, 4, 5, 2, 6, 7, 10, 13, 8, 9, 5},
}

// reels lengths [32, 32, 32, 32, 32], total reshuffles 33554432
// symbols: 77.605(lined) + 2.7777(scatter) = 80.382288%
// spin1 bonuses: count1 20736, rtp = 6.550598%
// spin3 bonuses: count3 3120, rtp = 2.956867%
// spin6 bonuses: count6 720, rtp = 1.364708%
// monopoly bonuses: count 6720, rtp = 5.739904%
// jackpots: count 32, frequency 1/1048576
// RTP = 80.382(sym) + 10.872(mje) + 5.7399(mjm) = 96.994365%
var Reels970 = game.Reels5x{
	{1, 5, 13, 4, 13, 3, 13, 5, 9, 7, 8, 13, 10, 13, 12, 11, 13, 12, 11, 13, 13, 2, 4, 5, 2, 6, 7, 9, 8, 3, 10, 6},
	{13, 2, 12, 9, 13, 4, 5, 6, 9, 7, 13, 10, 12, 13, 11, 13, 13, 11, 12, 13, 3, 4, 5, 2, 8, 7, 10, 4, 6, 8, 3, 1},
	{2, 1, 12, 3, 4, 5, 2, 6, 7, 10, 8, 4, 5, 13, 12, 11, 13, 12, 11, 13, 12, 3, 5, 13, 9, 6, 7, 10, 13, 8, 9, 13},
	{1, 2, 12, 4, 3, 6, 4, 2, 7, 8, 10, 5, 13, 12, 11, 13, 12, 11, 13, 12, 2, 4, 5, 3, 12, 6, 10, 7, 13, 8, 9, 13},
	{1, 2, 12, 4, 3, 5, 12, 6, 7, 3, 8, 12, 9, 13, 12, 11, 13, 12, 11, 13, 12, 3, 4, 5, 2, 6, 7, 10, 13, 8, 9, 5},
}

// reels lengths [32, 32, 32, 32, 32], total reshuffles 33554432
// symbols: 76.288(lined) + 2.7777(scatter) = 79.065740%
// spin1 bonuses: count1 20736, rtp = 6.550598%
// spin3 bonuses: count3 3120, rtp = 2.956867%
// spin6 bonuses: count6 720, rtp = 1.364708%
// monopoly bonuses: count 6720, rtp = 5.739904%
// jackpots: count 32, frequency 1/1048576
// RTP = 79.066(sym) + 10.872(mje) + 5.7399(mjm) = 95.677817%
var Reels957 = game.Reels5x{
	{1, 5, 13, 4, 13, 3, 13, 5, 9, 7, 8, 13, 10, 13, 12, 11, 13, 12, 11, 13, 13, 2, 4, 5, 2, 6, 7, 9, 8, 3, 10, 6},
	{13, 2, 12, 9, 13, 4, 5, 6, 9, 7, 13, 10, 12, 13, 11, 13, 13, 11, 12, 13, 3, 4, 5, 2, 8, 7, 10, 4, 6, 8, 3, 1},
	{2, 1, 12, 3, 4, 5, 2, 6, 7, 10, 8, 4, 5, 13, 12, 11, 13, 12, 11, 13, 12, 3, 5, 13, 9, 6, 7, 10, 13, 8, 9, 13},
	{1, 2, 12, 4, 3, 6, 4, 2, 7, 8, 3, 5, 13, 12, 11, 13, 12, 11, 13, 12, 2, 4, 5, 3, 12, 6, 10, 7, 13, 8, 9, 13},
	{1, 2, 12, 4, 3, 5, 12, 6, 7, 3, 8, 12, 2, 13, 12, 11, 13, 12, 11, 13, 12, 3, 4, 5, 2, 6, 7, 10, 13, 8, 9, 5},
}

// Map with available reels.
var ReelsMap = map[string]*game.Reels5x{
	"orig": &ReelsOrig,
	"98.6": &Reels986,
	"97.0": &Reels970,
	"95.7": &Reels957,
}

// Lined payment.
var LinePay = [13][5]int{
	{0, 0, 0, 0, 0},           //  1 dollar
	{0, 2, 5, 25, 100},        //  2 cherry
	{0, 2, 5, 25, 100},        //  3 plum
	{0, 0, 5, 25, 100},        //  4 wmelon
	{0, 0, 5, 25, 100},        //  5 grapes
	{0, 0, 10, 100, 250},      //  6 ananas
	{0, 0, 10, 100, 250},      //  7 lemon
	{0, 0, 10, 100, 250},      //  8 drink
	{0, 2, 10, 100, 500},      //  9 palm
	{0, 2, 10, 100, 500},      // 10 yacht
	{0, 10, 200, 2000, 10000}, // 11 eldorado
	{0, 0, 0, 0, 0},           // 12 spin
	{0, 0, 0, 0, 0},           // 13 dice
}

// Scatters payment.
var ScatPay = [5]int{0, 0, 2, 20, 1000} // 1 dollar

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
	{0, 0, 0, 0, 0},          //  1
	{0, 0, 0, 0, 0},          //  2
	{0, 0, 0, 0, 0},          //  3
	{0, 0, 0, 0, 0},          //  4
	{0, 0, 0, 0, 0},          //  5
	{0, 0, 0, 0, 0},          //  6
	{0, 0, 0, 0, 0},          //  7
	{0, 0, 0, 0, 0},          //  8
	{0, 0, 0, 0, 0},          //  9
	{0, 0, 0, 0, 0},          // 10
	{0, 0, 0, 0, 0},          // 11
	{0, 0, mje1, mje3, mje6}, // 12 Eldorado1, Eldorado3, Eldorado6
	{0, 0, 0, 0, mjm},        // 13 Monopoly
}

func NewGame(reels *game.Reels5x) *slotopol.Game {
	return &slotopol.Game{
		Slot5x3: game.Slot5x3{
			SBL:      []int{1},
			Bet:      1,
			FS:       0,
			Reels:    reels,
			BetLines: &game.BetLinesMgj,
		},
		LinePay:   &LinePay,
		ScatPay:   &ScatPay,
		ScatFree:  &slotopol.ScatFreespin,
		LineBonus: &LineBonus,
	}
}

const (
	jid = 1 // jackpot ID
)

func CalcStat(rn string) {
	var reels *game.Reels5x
	if rn != "" {
		var ok bool
		if reels, ok = ReelsMap[rn]; !ok {
			return
		}
	} else {
		reels = &ReelsOrig
	}
	var g = NewGame(reels)
	var sbl = float64(len(g.SBL))
	var s game.Stat
	var t0 = time.Now()
	func() {
		var ctx, cancel = context.WithCancel(context.Background())
		defer cancel()
		go s.Progress(ctx, time.NewTicker(2*time.Second), sbl, float64(g.Reels.Reshuffles()))
		s.BruteForce5x(ctx, g, g.Reels)
	}()
	var dur = time.Since(t0)
	var n = float64(s.Reshuffles)
	var lp, sp = float64(s.LinePay) / n / sbl * 100, float64(s.ScatPay) / n * 100
	var rtpsym = lp + sp
	var Mmje1, qmje1 = 106.0 * 1, float64(s.BonusCount[mje1]) / n / sbl
	var rtpmje1 = Mmje1 * qmje1 * 100
	var Mmje3, qmje3 = 106.0 * 3, float64(s.BonusCount[mje3]) / n / sbl
	var rtpmje3 = Mmje3 * qmje3 * 100
	var Mmje6, qmje6 = 106.0 * 6, float64(s.BonusCount[mje6]) / n / sbl
	var rtpmje6 = Mmje6 * qmje6 * 100
	var Mmjm, qmjm = 286.60597422268, float64(s.BonusCount[mjm]) / n / sbl
	var rtpmjm = Mmjm * qmjm * 100
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", float64(s.Reshuffles)/float64(g.Reels.Reshuffles())*100, len(g.SBL), dur)
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
		len(g.Reels.Reel(1)), len(g.Reels.Reel(2)), len(g.Reels.Reel(3)), len(g.Reels.Reel(4)), len(g.Reels.Reel(5)), g.Reels.Reshuffles())
	fmt.Printf("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lp, sp, rtpsym)
	fmt.Printf("spin1 bonuses: count1 %d, rtp = %.6f%%\n", s.BonusCount[mje1], rtpmje1)
	fmt.Printf("spin3 bonuses: count3 %d, rtp = %.6f%%\n", s.BonusCount[mje3], rtpmje3)
	fmt.Printf("spin6 bonuses: count6 %d, rtp = %.6f%%\n", s.BonusCount[mje6], rtpmje6)
	fmt.Printf("monopoly bonuses: count %d, rtp = %.6f%%\n", s.BonusCount[mjm], rtpmjm)
	fmt.Printf("jackpots: count %d, frequency 1/%d\n", s.JackCount[jid], int(n/float64(s.JackCount[jid])))
	fmt.Printf("RTP = %.5g(sym) + %.5g(mje) + %.5g(mjm) = %.6f%%\n", rtpsym, rtpmje1+rtpmje3+rtpmje6, rtpmjm, rtpsym+rtpmje1+rtpmje3+rtpmje6+rtpmjm)
}
