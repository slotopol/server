package ultrahot

import (
	"context"
	"fmt"
	"strconv"
	"time"

	slot "github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, rn string) float64 {
	var reels *slot.Reels3x
	if mrtp, _ := strconv.ParseFloat(rn, 64); mrtp != 0 {
		_, reels = slot.FindReels(ReelsMap, mrtp)
	} else {
		reels = &Reels93
	}
	var g = NewGame()
	var sln float64 = 1
	g.Sel.SetNum(int(sln), 1)
	var s slot.Stat

	var dur = slot.ScanReels(ctx, &s, g, reels,
		time.Tick(2*time.Second), time.Tick(2*time.Second))

	var reshuf = float64(s.Reshuffles)
	var lrtp = s.LinePay / reshuf / sln * 100
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", reshuf/float64(s.Planned())*100, g.Sel.Num(), dur)
	fmt.Printf("reels lengths [%d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), reels.Reshuffles())
	fmt.Printf("RTP = %.6f%%\n", lrtp)
	return lrtp
}
