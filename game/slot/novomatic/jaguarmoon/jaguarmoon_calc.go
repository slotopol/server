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
	reshuffles uint64
	linepay    float64
	scatpay    float64
	freecount  uint64
	freehits   uint64
	bonuscount [8]uint64
	jackcount  [4]uint64
	lpm, spm   sync.Mutex
}

func (s *Stat) SetPlan(n uint64) {
	atomic.StoreUint64(&s.planned, n)
}

func (s *Stat) Planned() uint64 {
	return atomic.LoadUint64(&s.planned)
}

func (s *Stat) Count() uint64 {
	return atomic.LoadUint64(&s.reshuffles)
}

func (s *Stat) LineRTP(sel int) float64 {
	var reshuf = float64(atomic.LoadUint64(&s.reshuffles))
	s.lpm.Lock()
	var lp = s.linepay
	s.lpm.Unlock()
	return lp / reshuf / float64(sel) * 100
}

func (s *Stat) ScatRTP(sel int) float64 {
	var reshuf = float64(atomic.LoadUint64(&s.reshuffles))
	s.spm.Lock()
	var sp = s.scatpay
	s.spm.Unlock()
	return sp / reshuf / float64(sel) * 100
}

func (s *Stat) FreeCount() uint64 {
	return atomic.LoadUint64(&s.freecount)
}

func (s *Stat) FreeHits() uint64 {
	return atomic.LoadUint64(&s.freehits)
}

func (s *Stat) BonusCount(bid int) uint64 {
	return atomic.LoadUint64(&s.bonuscount[bid])
}

func (s *Stat) JackCount(jid int) uint64 {
	return atomic.LoadUint64(&s.jackcount[jid])
}

func (s *Stat) Update(wins slot.Wins) {
	for _, wi := range wins {
		if wi.Pay != 0 {
			if wi.Line != 0 {
				s.lpm.Lock()
				s.linepay += wi.Pay * wi.Mult
				s.lpm.Unlock()
			} else {
				s.spm.Lock()
				s.scatpay += wi.Pay * wi.Mult
				s.spm.Unlock()
			}
		}
		if wi.Free != 0 {
			atomic.AddUint64(&s.freecount, uint64(wi.Free*int(FreeMult[wi.Num-1])))
			atomic.AddUint64(&s.freehits, 1)
		}
		if wi.BID != 0 {
			atomic.AddUint64(&s.bonuscount[wi.BID], 1)
		}
		if wi.JID != 0 {
			atomic.AddUint64(&s.jackcount[wi.JID], 1)
		}
	}
	atomic.AddUint64(&s.reshuffles, 1)
}

func CalcStatBon(ctx context.Context) float64 {
	var reels = ReelsBon
	var g = NewGame()
	g.Sel = 10 // bet on 243 ways
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var lrtp, srtp = s.LineRTP(g.Sel), s.ScatRTP(g.Sel)
		var rtpsym = lrtp + srtp
		fmt.Fprintf(w, "reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
			len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
		fmt.Fprintf(w, "RTP = %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
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
		var q = float64(s.FreeCount()) / reshuf
		var sq = 1 / (1 - q)
		var rtp = rtpsym + q*rtpfs
		fmt.Fprintf(w, "reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
			len(reels.Reel(1)), len(reels.Reel(2)), len(reels.Reel(3)), len(reels.Reel(4)), len(reels.Reel(5)), reels.Reshuffles())
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCount(), q, sq)
		fmt.Fprintf(w, "free games frequency: 1/%.5g\n", reshuf/float64(s.FreeHits()))
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%\n", rtpsym, q, rtpfs, rtp)
		return rtp
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc,
		time.Tick(2*time.Second), time.Tick(2*time.Second))
}
