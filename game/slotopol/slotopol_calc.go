package slotopol

import (
	"context"
	"fmt"
	"time"

	"github.com/schwarzlichtbezirk/slot-srv/game"
)

var Emje float64 // Eldorado game 1 spin calculated expectation

// Eldorado expectation.
func ExpEldorado() float64 {
	var sum = 0
	for _, v := range Eldorado {
		sum += v
	}
	var E = float64(sum) / float64(len(Eldorado))
	fmt.Printf("eldorado 1 spin: count = %d, E = %g\n", len(Eldorado), E)
	return E
}

func CalcStat(ctx context.Context, rn string) float64 {
	fmt.Printf("*bonus games calculations*\n")
	Emje = ExpEldorado()
	fmt.Printf("*reels calculations*\n")
	var reels *game.Reels5x
	if rn != "" {
		var ok bool
		if reels, ok = ReelsMap[rn]; !ok {
			return 0
		}
	} else {
		reels = &Reels100
	}
	var g = NewGame(reels)
	var sbl = float64(len(g.SBL))
	var s game.Stat
	var t0 = time.Now()
	go s.Progress(ctx, time.NewTicker(2*time.Second), sbl, float64(g.Reels.Reshuffles()))
	s.BruteForce5x(ctx, g, g.Reels)
	var dur = time.Since(t0)
	var n = float64(s.Reshuffles)
	var lrtp, srtp = float64(s.LinePay) / n / sbl * 100, float64(s.ScatPay) / n * 100
	var rtpsym = lrtp + srtp
	var qmje9 = float64(s.BonusCount[mje9]) / n / sbl
	var rtpmje9 = Emje * 9 * qmje9 * 100
	var Mmjm, qmjm = 286.60597422268, float64(s.BonusCount[mjm]) / n / sbl
	var rtpmjm = Mmjm * qmjm * 100
	var rtp = rtpsym + rtpmje9 + rtpmjm
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", float64(s.Reshuffles)/float64(g.Reels.Reshuffles())*100, len(g.SBL), dur)
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
		len(g.Reels.Reel(1)), len(g.Reels.Reel(2)), len(g.Reels.Reel(3)), len(g.Reels.Reel(4)), len(g.Reels.Reel(5)), g.Reels.Reshuffles())
	fmt.Printf("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
	fmt.Printf("spin9 bonuses: count %d, rtp = %.6f%%\n", s.BonusCount[mje9], rtpmje9)
	fmt.Printf("monopoly bonuses: count %d, rtp = %.6f%%\n", s.BonusCount[mjm], rtpmjm)
	fmt.Printf("jackpots: count %d, frequency 1/%d\n", s.JackCount[jid], int(n/float64(s.JackCount[jid])))
	fmt.Printf("RTP = %.5g(sym) + %.5g(mje9) + %.5g(mjm) = %.6f%%\n", rtpsym, rtpmje9, rtpmjm, rtp)
	return rtp
}
