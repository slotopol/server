package egypt

import (
	"context"
	"fmt"
	"time"

	"github.com/slotopol/server/game/slot"
)

// Minislot expectation calculation:
// total combinations: 3*3*3 = 27
//
//	x1 = 27-3=24
//	x3 = 1
//	x6 = 1
//	x9 = 1
//
// Em = (24*1 + 1*3 + 1*6 + 1*9)/27 = 42/27 = 1.5555555556
const Em = 42. / 27.

func CalcStat(ctx context.Context, mrtp float64) float64 {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	var g = NewGame()
	var sln float64 = 1
	g.Sel = int(sln)
	var s slot.Stat

	slot.ScanReels5x(ctx, &s, g, reels,
		time.Tick(2*time.Second), time.Tick(2*time.Second))

	var reshuf = float64(s.Reshuffles)
	var lrtp, srtp = s.LinePay / reshuf / sln * 100, s.ScatPay / reshuf / sln * 100
	var rtpsym = lrtp + srtp
	var rtp = rtpsym * Em
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
	fmt.Printf("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
	fmt.Printf("RTP = %.5g(sym) * %.5g(Em) = %.6f%%\n", rtpsym, Em, rtp)
	return rtpsym
}
