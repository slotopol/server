package jaguarmoon

import (
	"context"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/slotopol/server/game/slot"
)

type Stat struct {
	planned    uint64
	Reshuffles uint64
	LinePay    float64
	ScatPay    float64
	FreeCount  uint64
	FreeHits   uint64
	BonusCount [8]uint64
	JackCount  [4]uint64
	LPM, SPM   sync.Mutex
}

func (s *Stat) SetPlan(n uint64) {
	atomic.StoreUint64(&s.planned, n)
}

func (s *Stat) Planned() uint64 {
	return atomic.LoadUint64(&s.planned)
}

func (s *Stat) Count() uint64 {
	return atomic.LoadUint64(&s.Reshuffles)
}

func (s *Stat) LineRTP(sel int) float64 {
	var reshuf = float64(atomic.LoadUint64(&s.Reshuffles))
	s.LPM.Lock()
	var lp = s.LinePay
	s.LPM.Unlock()
	return lp / reshuf / float64(sel) * 100
}

func (s *Stat) ScatRTP(sel int) float64 {
	var reshuf = float64(atomic.LoadUint64(&s.Reshuffles))
	s.SPM.Lock()
	var sp = s.ScatPay
	s.SPM.Unlock()
	return sp / reshuf / float64(sel) * 100
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
	g.Sel = 10 // bet on 243 ways
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var lrtp, srtp = s.LineRTP(g.Sel), s.ScatRTP(g.Sel)
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
	g.Sel = 10 // bet on 243 ways
	var s Stat

	var calc = func(w io.Writer) float64 {
		var reshuf = float64(s.Count())
		var lrtp, srtp = s.LineRTP(g.Sel), s.ScatRTP(g.Sel)
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
