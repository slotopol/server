package suncity

import (
	"context"
	"fmt"
	"time"

	"github.com/slotopol/server/game/slot"
)

func CalcStatBon(ctx context.Context) (rtp, num float64) {
	var reels = ReelsBon
	var g = NewGame()
	var sln float64 = 1
	g.Sel = int(sln)
	g.FSR = -1 // set free spins mode
	var s slot.Stat

	var dur = slot.ScanReels5x(ctx, &s, g, reels,
		time.Tick(2*time.Second), time.Tick(2*time.Second))

	var reshuf = float64(s.Reshuffles)
	var lrtp, srtp = s.LinePay / reshuf / sln * 100, s.ScatPay / reshuf / sln * 100
	if srtp > 0 {
		panic("scatters have no pays")
	}
	var rtpsym = lrtp + srtp
	var fgf = reshuf / float64(s.FreeHits)
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", reshuf/float64(s.Planned())*100, g.Sel, dur)
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
	fmt.Printf("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
	fmt.Printf("free games frequency: 1/%.5g\n", reshuf/float64(s.FreeHits))
	fmt.Printf("RTP = rtp(sym) = %.6f%%\n", rtpsym)
	return rtpsym, fgf
}

func CalcStatReg(ctx context.Context, mrtp float64) float64 {
	fmt.Printf("*bonus reels calculations*\n")
	var rtpfs, numfs = CalcStatBon(ctx)
	if ctx.Err() != nil {
		return 0
	}
	fmt.Printf("*regular reels calculations*\n")
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	var g = NewGame()
	var sln float64 = 1
	g.Sel = int(sln)
	var s slot.Stat

	var dur = slot.ScanReels5x(ctx, &s, g, reels,
		time.Tick(2*time.Second), time.Tick(2*time.Second))

	var reshuf = float64(s.Reshuffles)
	var lrtp, srtp = s.LinePay / reshuf / sln * 100, s.ScatPay / reshuf / sln * 100
	if srtp > 0 {
		panic("scatters have no pays")
	}
	var rtpsym = lrtp + srtp
	var fgf = reshuf / float64(s.FreeHits)
	var rtp = rtpsym + rtpfs*numfs/fgf
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", reshuf/float64(s.Planned())*100, g.Sel, dur)
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
	fmt.Printf("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
	fmt.Printf("free games frequency: 1/%.5g\n", fgf)
	fmt.Printf("RTP = %.5g(sym) + %.5g(fg)*%.5g/%.5g = %.6f%%\n", rtpsym, rtpfs, numfs, fgf, rtp)
	return rtp
}
