package jaguarmoon

import (
	"context"
	"fmt"
	"io"
	"sync/atomic"
	"time"

	"github.com/slotopol/server/game/slot"
)

type Stat struct {
	slot.Stat
}

func (s *Stat) Update(wins slot.Wins) {
	for _, wi := range wins {
		if wi.Pay != 0 {
			if wi.Line != 0 {
				s.LPM.Lock()
				s.LinePay += wi.Pay * wi.Mult
				s.LPM.Unlock()
			} else {
				s.SPM.Lock()
				s.ScatPay += wi.Pay * wi.Mult
				s.SPM.Unlock()
			}
		}
		if wi.Free != 0 {
			atomic.AddUint64(&s.FreeCount, uint64(wi.Free*int(FreeMult[wi.Num-1])))
			atomic.AddUint64(&s.FreeHits, 1)
		}
		if wi.BID != 0 {
			atomic.AddUint64(&s.BonusCount[wi.BID], 1)
		}
		if wi.JID != 0 {
			atomic.AddUint64(&s.JackCount[wi.JID], 1)
		}
	}
	atomic.AddUint64(&s.Reshuffles, 1)
}

func CalcStatBon(ctx context.Context) float64 {
	var reels = ReelsBon
	var g = NewGame()
	var sln float64 = 10
	g.Sel = int(sln)
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var reshuf = float64(s.Reshuffles)
		var lrtp, srtp = s.LinePay / reshuf / sln * 100, s.ScatPay / reshuf / sln * 100
		var rtpsym = lrtp + srtp
		fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
			len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
		fmt.Printf("RTP = %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		return rtpsym
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc,
		time.Tick(2*time.Second), time.Tick(2*time.Second))
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
	var sln float64 = 10
	g.Sel = int(sln)
	var s Stat

	var calc = func(w io.Writer) float64 {
		var reshuf = float64(s.Reshuffles)
		var lrtp, srtp = s.LinePay / reshuf / sln * 100, s.ScatPay / reshuf / sln * 100
		var rtpsym = lrtp + srtp
		var q = float64(s.FreeCount) / reshuf
		var sq = 1 / (1 - q)
		var rtp = rtpsym + q*rtpfs
		fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
			len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
		fmt.Printf("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Printf("free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCount, q, sq)
		fmt.Printf("free games frequency: 1/%.5g\n", reshuf/float64(s.FreeHits))
		fmt.Printf("RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%\n", rtpsym, q, rtpfs, rtp)
		return rtp
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc,
		time.Tick(2*time.Second), time.Tick(2*time.Second))
}
