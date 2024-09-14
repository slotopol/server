package ultrahot

import (
	"context"
	"fmt"
	"strconv"
	"time"

	slot "github.com/slotopol/server/game/slot"
	"github.com/slotopol/server/util"
)

func CalcStat(ctx context.Context, rn string) float64 {
	var reels *slot.Reels3x
	if mrtp, _ := strconv.ParseFloat(rn, 64); mrtp != 0 {
		var _, r = FindReels(mrtp)
		reels = r.(*slot.Reels3x)
	} else {
		reels = &Reels93
	}
	var g = NewGame()
	g.SBL = util.MakeBitNum(1, 1)
	var sbl = float64(g.SBL.Num())
	var s slot.Stat

	var total = float64(reels.Reshuffles())
	var dur = func() time.Duration {
		var t0 = time.Now()
		var ctx2, cancel2 = context.WithCancel(ctx)
		defer cancel2()
		go s.Progress(ctx2, time.Tick(2*time.Second), sbl, total)
		slot.BruteForce3x(ctx2, &s, g, reels)
		return time.Since(t0)
	}()

	var reshuf = float64(s.Reshuffles)
	var lrtp = s.LinePay / reshuf / sbl * 100
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", reshuf/total*100, g.SBL.Num(), dur)
	fmt.Printf("reels lengths [%d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), reels.Reshuffles())
	fmt.Printf("RTP = %.6f%%\n", lrtp)
	return lrtp
}
