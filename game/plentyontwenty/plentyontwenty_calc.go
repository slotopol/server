package plentyontwenty

import (
	"context"
	"fmt"
	"time"

	"github.com/slotopol/server/game"
)

func CalcStat(ctx context.Context, rn string) float64 {
	var reels *game.Reels5x
	if rn != "" {
		var ok bool
		if reels, ok = ReelsMap[rn]; !ok {
			return 0
		}
	} else {
		rn, reels = "93", &Reels93
	}
	var g = NewGame(rn)
	var sbl = float64(g.SBL.Num())
	var s game.Stat

	var total = float64(reels.Reshuffles())
	var dur = func() time.Duration {
		var t0 = time.Now()
		var ctx2, cancel2 = context.WithCancel(ctx)
		defer cancel2()
		go s.Progress(ctx2, time.NewTicker(2*time.Second), sbl, total)
		game.BruteForce5x(ctx2, &s, g, reels)
		return time.Since(t0)
	}()

	var reshuf = float64(s.Reshuffles)
	var lrtp, srtp = float64(s.LinePay) / reshuf / sbl * 100, float64(s.ScatPay) / reshuf * 100
	var rtpsym = lrtp + srtp
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", reshuf/total*100, g.SBL.Num(), dur)
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
	if s.JackCount[jid] > 0 {
		fmt.Printf("jackpots: count %d, frequency 1/%d\n", s.JackCount[jid], int(reshuf/float64(s.JackCount[jid])))
	}
	fmt.Printf("RTP = %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
	return rtpsym
}
