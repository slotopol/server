package iceiceice

import (
	"context"
	"fmt"
	"time"

	"github.com/slotopol/server/game/slot"
)

func CalcStatBon(ctx context.Context, mrtp float64) float64 {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	var g = NewGame()
	var sln float64 = 1
	g.Sel = int(sln)
	g.FSR = 20 // set free spins mode
	var s slot.Stat

	var dur = slot.ScanReels3x(ctx, &s, g, reels,
		time.Tick(2*time.Second), time.Tick(2*time.Second))

	var reshuf = float64(s.Reshuffles)
	var lrtp, srtp = s.LinePay / reshuf / sln * 100, s.ScatPay / reshuf / sln * 100
	var rtpsym = lrtp + srtp
	var q = float64(s.FreeCount) / reshuf
	var sq = 1 / (1 - q)
	var rtp = sq * rtpsym
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", reshuf/float64(s.Planned())*100, g.Sel, dur)
	fmt.Printf("reels lengths [%d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), reels.Reshuffles())
	fmt.Printf("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
	fmt.Printf("free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCount, q, sq)
	fmt.Printf("free games frequency: 1/%.5g\n", reshuf/float64(s.FreeHits))
	fmt.Printf("RTP = sq*rtp(sym) = %.5g*%.5g = %.6f%%\n", sq, rtpsym, rtp)
	return rtp
}

func CalcStatReg(ctx context.Context, mrtp float64) float64 {
	fmt.Printf("*bonus reels calculations*\n")
	var rtpfs = CalcStatBon(ctx, mrtp)
	if ctx.Err() != nil {
		return 0
	}
	fmt.Printf("*regular reels calculations*\n")
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	var g = NewGame()
	var sln float64 = 1
	g.Sel = int(sln)
	var s slot.Stat

	var dur = slot.ScanReels3x(ctx, &s, g, reels,
		time.Tick(2*time.Second), time.Tick(2*time.Second))

	var reshuf = float64(s.Reshuffles)
	var lrtp, srtp = s.LinePay / reshuf / sln * 100, s.ScatPay / reshuf / sln * 100
	var rtpsym = lrtp + srtp
	var q = float64(s.FreeCount) / reshuf
	var sq = 1 / (1 - q)
	var rtp = rtpsym + q*rtpfs
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", reshuf/float64(s.Planned())*100, g.Sel, dur)
	fmt.Printf("reels lengths [%d, %d, %d], total reshuffles %d\n",
		len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), reels.Reshuffles())
	fmt.Printf("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
	fmt.Printf("free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCount, q, sq)
	fmt.Printf("free games frequency: 1/%.5g\n", reshuf/float64(s.FreeHits))
	fmt.Printf("RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%\n", rtpsym, q, rtpfs, rtp)
	return rtp
}
