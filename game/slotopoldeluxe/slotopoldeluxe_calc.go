package slotopoldeluxe

import (
	"context"
	"fmt"
	"time"

	"github.com/slotopol/server/game"
	"github.com/slotopol/server/game/slotopol"
)

func CalcStat(ctx context.Context, rn string) float64 {
	fmt.Printf("*bonus games calculations*\n")
	slotopol.Emje = slotopol.ExpEldorado()
	slotopol.Emjm = slotopol.ExpMonopoly()
	fmt.Printf("*reels calculations*\n")
	var reels *game.Reels5x
	if rn != "" {
		var ok bool
		if reels, ok = ReelsMap[rn]; !ok {
			return 0
		}
	} else {
		reels = &Reels104
	}
	var g = NewGame(reels)
	var sbl = float64(g.SBL.Num())
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
	var qmje1 = float64(s.BonusCount[mje1]) / n / sbl
	var rtpmje1 = slotopol.Emje * 1 * qmje1 * 100
	var qmje3 = float64(s.BonusCount[mje3]) / n / sbl
	var rtpmje3 = slotopol.Emje * 3 * qmje3 * 100
	var qmje6 = float64(s.BonusCount[mje6]) / n / sbl
	var rtpmje6 = slotopol.Emje * 6 * qmje6 * 100
	var qmjm = float64(s.BonusCount[mjm]) / n / sbl
	var rtpmjm = slotopol.Emjm * qmjm * 100
	var rtp = rtpsym + rtpmje1 + rtpmje3 + rtpmje6 + rtpmjm
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", float64(s.Reshuffles)/float64(g.Reels.Reshuffles())*100, g.SBL.Num(), dur)
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
		len(g.Reels.Reel(1)), len(g.Reels.Reel(2)), len(g.Reels.Reel(3)), len(g.Reels.Reel(4)), len(g.Reels.Reel(5)), g.Reels.Reshuffles())
	fmt.Printf("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
	fmt.Printf("spin1 bonuses: count1 %d, rtp = %.6f%%\n", s.BonusCount[mje1], rtpmje1)
	fmt.Printf("spin3 bonuses: count3 %d, rtp = %.6f%%\n", s.BonusCount[mje3], rtpmje3)
	fmt.Printf("spin6 bonuses: count6 %d, rtp = %.6f%%\n", s.BonusCount[mje6], rtpmje6)
	fmt.Printf("monopoly bonuses: count %d, rtp = %.6f%%\n", s.BonusCount[mjm], rtpmjm)
	fmt.Printf("jackpots: count %d, frequency 1/%d\n", s.JackCount[jid], int(n/float64(s.JackCount[jid])))
	fmt.Printf("RTP = %.5g(sym) + %.5g(mje) + %.5g(mjm) = %.6f%%\n", rtpsym, rtpmje1+rtpmje3+rtpmje6, rtpmjm, rtp)
	return rtp
}
