package jewels

import (
	"context"
	"fmt"
	"strconv"
	"time"

	slot "github.com/slotopol/server/game/slot"
	"github.com/slotopol/server/util"
)

func CalcStat(ctx context.Context, rn string) float64 {
	var reels *slot.Reels5x
	if mrtp, _ := strconv.ParseFloat(rn, 64); mrtp != 0 {
		var _, r = FindReels(mrtp)
		reels = r.(*slot.Reels5x)
	} else {
		reels = &Reels93
	}
	var g = NewGame()
	g.Sel = util.MakeBitNum(1, 1)
	var sln = float64(g.Sel.Num())
	var s slot.Stat

	var total = float64(reels.Reshuffles())
	var dur = func() time.Duration {
		var t0 = time.Now()
		var ctx2, cancel2 = context.WithCancel(ctx)
		defer cancel2()
		go s.Progress(ctx2, time.Tick(2*time.Second), sln, total)
		slot.BruteForce5x(ctx2, &s, g, reels)
		return time.Since(t0)
	}()

	var reshuf = float64(s.Reshuffles)
	var lrtp = s.LinePay / reshuf / sln * 100
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", reshuf/total*100, g.Sel.Num(), dur)
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
	if s.JackCount[jid] > 0 {
		fmt.Printf("jackpots: count %d, frequency 1/%d\n", s.JackCount[jid], int(reshuf/float64(s.JackCount[jid])))
	}
	fmt.Printf("RTP = %g%%\n", lrtp)
	return lrtp
}
