package goldentour

import (
	"context"
	"fmt"
	"time"

	slot "github.com/slotopol/server/game/slot"
)

var Egolfbn float64

func ExpGolf() float64 {
	var sum float64
	for _, v := range Golf {
		sum += float64(v)
	}
	var E = sum / float64(len(Golf))
	fmt.Printf("len = %d, E = %g\n", len(Golf), E)
	return E
}

func CalcStat(ctx context.Context, mrtp float64) float64 {
	fmt.Printf("*bonus games calculations*\n")
	Egolfbn = ExpGolf()
	fmt.Printf("*reels calculations*\n")
	var reels, _ = slot.FindReels(ReelsMap, mrtp)
	var g = NewGame()
	var sln float64 = 1
	g.Sel = int(sln)
	var s slot.Stat

	var dur = slot.ScanReels5x(ctx, &s, g, reels,
		time.Tick(2*time.Second), time.Tick(2*time.Second))

	var reshuf = float64(s.Reshuffles)
	var lrtp, srtp = s.LinePay / reshuf / sln * 100, s.ScatPay / reshuf / sln * 100
	var rtpsym = lrtp + srtp
	var qgolfbn = float64(s.BonusCount[golfbon]) / reshuf / sln
	var rtpgolfbn = Egolfbn * qgolfbn * 100
	var rtp = rtpsym + rtpgolfbn
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", reshuf/float64(s.Planned())*100, g.Sel, dur)
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
	fmt.Printf("golf bonuses: count %d, rtp = %.6f%%\n", s.BonusCount[golfbon], rtpgolfbn)
	fmt.Printf("golf bonuses frequency: 1/%.5g\n", reshuf/float64(s.BonusCount[golfbon]))
	fmt.Printf("RTP = %.5g(lined) + %.5g(scatter) + %.5g(golf) = %.6f%%\n", lrtp, srtp, rtpgolfbn, rtp)
	return rtp
}
