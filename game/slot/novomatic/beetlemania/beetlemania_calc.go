package beetlemania

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/slotopol/server/game/slot"
)

// Attention! On freespins can be calculated median only, not expectation.

func CalcStatBon(ctx context.Context) float64 {
	var reels = ReelsBon
	var g = NewGame()
	var sln float64 = 1
	g.Sel = int(sln)
	g.FSR = 10 // set free spins mode
	var s slot.Stat

	slot.ScanReels5x(ctx, &s, g, reels,
		time.Tick(2*time.Second), time.Tick(2*time.Second))

	var reshuf = float64(s.Reshuffles)
	var lrtp = s.LinePay / reshuf / sln * 100
	var qjazz = float64(s.BonusCount[jbonus]) / reshuf
	var jpow = math.Pow(2, 10*qjazz) // jazz power
	var rtpjazz = lrtp*jpow - lrtp
	var rtp = lrtp * jpow
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
	fmt.Printf("symbols: %.5g(lined) + 0(scatter) = %.6f%%\n", lrtp, lrtp)
	fmt.Printf("jazzbee bonuses: frequency 1/%.5g, pow = %.5g, rtp = %.6f%%\n", reshuf/float64(s.BonusCount[jbonus]), jpow, rtpjazz)
	fmt.Printf("RTP = rtp(sym) + rtp(jazz) = %.5g + %.5g = %.6f%%\n", lrtp, rtpjazz, rtp)
	return rtp
}

func CalcStatReg(ctx context.Context, mrtp float64) float64 {
	fmt.Printf("*bonus reels calculations*\n")
	var rtpfs = CalcStatBon(ctx)
	if ctx.Err() != nil {
		return 0
	}
	fmt.Printf("*regular reels calculations*\n")
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
	var q = float64(s.FreeCount) / reshuf
	var rtp = rtpsym + q*rtpfs
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
	fmt.Printf("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
	fmt.Printf("free spins %d, q = %.5g\n", s.FreeCount, q)
	fmt.Printf("free games frequency: 1/%.5g\n", reshuf/float64(s.FreeHits))
	fmt.Printf("RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%\n", rtpsym, q, rtpfs, rtp)
	return rtp
}
