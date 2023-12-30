package slotopoldeluxe

import (
	"context"
	"fmt"
	"time"

	"github.com/schwarzlichtbezirk/slot-srv/game"
)

func CalcStat(ctx context.Context, rn string) float64 {
	var reels *game.Reels5x
	if rn != "" {
		var ok bool
		if reels, ok = ReelsMap[rn]; !ok {
			return 0
		}
	} else {
		reels = &Reels103
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
	var Mmje1, qmje1 = 106.0 * 1, float64(s.BonusCount[mje1]) / n / sbl
	var rtpmje1 = Mmje1 * qmje1 * 100
	var Mmje3, qmje3 = 106.0 * 3, float64(s.BonusCount[mje3]) / n / sbl
	var rtpmje3 = Mmje3 * qmje3 * 100
	var Mmje6, qmje6 = 106.0 * 6, float64(s.BonusCount[mje6]) / n / sbl
	var rtpmje6 = Mmje6 * qmje6 * 100
	var Mmjm, qmjm = 286.60597422268, float64(s.BonusCount[mjm]) / n / sbl
	var rtpmjm = Mmjm * qmjm * 100
	var rtp = rtpsym + rtpmje1 + rtpmje3 + rtpmje6 + rtpmjm
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", float64(s.Reshuffles)/float64(g.Reels.Reshuffles())*100, len(g.SBL), dur)
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
