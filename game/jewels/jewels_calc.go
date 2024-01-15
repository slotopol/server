package jewels

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

	var dur = func() time.Duration {
		var t0 = time.Now()
		var ctx2, cancel2 = context.WithCancel(ctx)
		defer cancel2()
		go s.Progress(ctx2, time.NewTicker(2*time.Second), sbl, float64(reels.Reshuffles()))
		s.BruteForce5x(ctx2, g, reels)
		return time.Since(t0)
	}()

	var n = float64(s.Reshuffles)
	var lrtp = float64(s.LinePay) / n / sbl * 100
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", float64(s.Reshuffles)/float64(reels.Reshuffles())*100, g.SBL.Num(), dur)
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
	fmt.Printf("RTP = %g%%\n", lrtp)
	return lrtp
}
