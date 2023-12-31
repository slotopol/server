package champagne

import (
	"context"
	"fmt"
	"time"

	"github.com/slotopol/server/game"
)

var Emjc float64 // Bottle game calculated expectation

func ExpBottle() float64 {
	// avr 1 bottle gain
	var m float64
	for _, v := range Bottles {
		m += float64(v)
	}
	m /= float64(len(Bottles))

	// expectation
	var E float64
	var n = 0
	for i := 0; i < len(Bottles); i++ {
		for j := i + 1; j < len(Bottles); j++ {
			if Bottles[i] == Bottles[j] {
				E += float64(Bottles[i]) * 4
			} else {
				E += float64(Bottles[i] + Bottles[j])
			}
			n++
		}
	}
	E /= float64(n)
	fmt.Printf("len = %d, avr bottle gain = %.5g, E = %g\n", len(Bottles), m, E)
	return E
}

func CalcStatBon(ctx context.Context, rn string) float64 {
	var reels *game.Reels5x
	if rn != "" {
		var ok bool
		if reels, ok = ReelsMap[rn]; !ok {
			return 0
		}
	} else {
		reels = &Reels964
	}
	var g = NewGame(reels)
	g.FS = 15 // set free spins mode
	var sbl = float64(len(g.SBL))
	var s game.Stat

	var dur = func() time.Duration {
		var t0 = time.Now()
		var ctx2, cancel2 = context.WithCancel(ctx)
		defer cancel2()
		go s.Progress(ctx2, time.NewTicker(2*time.Second), sbl, float64(g.Reels.Reshuffles()))
		s.BruteForce5x(ctx2, g, g.Reels)
		return time.Since(t0)
	}()

	var n = float64(s.Reshuffles)
	var lrtp, srtp = float64(s.LinePay) / n / sbl * 100, float64(s.ScatPay) / n * 100
	var rtpsym = lrtp + srtp
	var q = float64(s.FreeCount) / n
	var sq = 1 / (1 - q)
	var qmjc = float64(s.BonusCount[mjc]) / n / sbl
	var rtpmjc = Emjc * qmjc * 100
	var rtp = sq * (rtpsym + rtpmjc)
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", float64(s.Reshuffles)/float64(g.Reels.Reshuffles())*100, len(g.SBL), dur)
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
		len(g.Reels.Reel(1)), len(g.Reels.Reel(2)), len(g.Reels.Reel(3)), len(g.Reels.Reel(4)), len(g.Reels.Reel(5)), g.Reels.Reshuffles())
	fmt.Printf("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
	fmt.Printf("free games %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCount, q, sq)
	fmt.Printf("champagne bonuses: count %d, rtp = %.6f%%\n", s.BonusCount[mjc], rtpmjc)
	fmt.Printf("jackpots: count %d, frequency 1/%d\n", s.JackCount[jid], int(n/float64(s.JackCount[jid])))
	fmt.Printf("RTP = sq*(rtp(sym)+rtp(mjc)) = %.5g*(%.5g+%.5g) = %.6f%%\n", sq, rtpsym, rtpmjc, rtp)
	return rtp
}

func CalcStatReg(ctx context.Context, rn string) float64 {
	fmt.Printf("*bonus games calculations*\n")
	Emjc = ExpBottle()
	fmt.Printf("*bonus reels calculations*\n")
	var rtpfs = CalcStatBon(ctx, rn)
	fmt.Printf("*regular reels calculations*\n")
	var reels *game.Reels5x
	if rn != "" {
		var ok bool
		if reels, ok = ReelsMap[rn]; !ok {
			return 0
		}
	} else {
		reels = &Reels964
	}
	var g = NewGame(reels)
	g.FS = 0 // no free spins
	var sbl = float64(len(g.SBL))
	var s game.Stat

	var dur = func() time.Duration {
		var t0 = time.Now()
		var ctx2, cancel2 = context.WithCancel(ctx)
		defer cancel2()
		go s.Progress(ctx2, time.NewTicker(2*time.Second), sbl, float64(g.Reels.Reshuffles()))
		s.BruteForce5x(ctx2, g, g.Reels)
		return time.Since(t0)
	}()

	var n = float64(s.Reshuffles)
	var lrtp, srtp = float64(s.LinePay) / n / sbl * 100, float64(s.ScatPay) / n * 100
	var rtpsym = lrtp + srtp
	var q = float64(s.FreeCount) / n
	var qmjc = float64(s.BonusCount[mjc]) / n / sbl
	var rtpmjc = Emjc * qmjc * 100
	var rtp = rtpsym + rtpmjc + q*rtpfs
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", float64(s.Reshuffles)/float64(g.Reels.Reshuffles())*100, len(g.SBL), dur)
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
		len(g.Reels.Reel(1)), len(g.Reels.Reel(2)), len(g.Reels.Reel(3)), len(g.Reels.Reel(4)), len(g.Reels.Reel(5)), g.Reels.Reshuffles())
	fmt.Printf("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
	fmt.Printf("free games %d, q = %.6f\n", s.FreeCount, q)
	fmt.Printf("champagne bonuses: count %d, rtp = %.6f%%\n", s.BonusCount[mjc], rtpmjc)
	fmt.Printf("jackpots: count %d, frequency 1/%d\n", s.JackCount[jid], int(n/float64(s.JackCount[jid])))
	fmt.Printf("RTP = rtp(sym) + rtp(mjc) + q*rtp(fg) = %.5g + %.5g + %.5g*%.5g = %.6f%%\n", rtpsym, rtpmjc, q, rtpfs, rtp)
	return rtp
}
