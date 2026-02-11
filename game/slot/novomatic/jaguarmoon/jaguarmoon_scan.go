package jaguarmoon

import (
	"context"
	"fmt"
	"io"

	"go.uber.org/atomic"

	"github.com/slotopol/server/game/slot"
)

type Stat struct {
	slot.StatGeneric
	FreeCount [6]atomic.Uint64
}

// Declare conformity with Stater interface.
var _ slot.Simulator = (*Stat)(nil)

func (s *Stat) FreeCountF(n int) float64 {
	return float64(s.FreeCount[n-1].Load())
}

// Returns free spins quantifier, sum of a decreasing
// geometric progression for retriggered free spins.
func (s *Stat) FSQ(n int) (q float64) {
	q = s.FreeCountF(n) / s.Count()
	return
}

func (s *Stat) Update(wins slot.Wins) (pay float64) {
	for _, wi := range wins {
		if wi.Pay != 0 {
			var p = wi.Pay * wi.MP
			s.S[wi.Sym].Add(p)
			pay += p
		}
		if wi.FS != 0 {
			s.FreeCount[wi.Num-1].Add(uint64(wi.FS))
			s.FHC.Inc()
		}
		if wi.BID != 0 {
			s.BHC[wi.BID].Inc()
		}
		if wi.JID != 0 {
			s.JHC[wi.JID].Inc()
		}
	}
	s.N.Inc()
	return
}

func (s *Stat) Simulate(g slot.SlotGame, reels slot.Reelx, wins *slot.Wins) {
	if g.Scanner(wins) != nil {
		s.EC.Inc()
		return
	}
	if pay := s.Update(*wins); pay != 0 {
		s.Q.Add(pay * pay)
	}
}

func CalcStatBon(ctx context.Context, sp *slot.ScanPar) float64 {
	var reels = ReelsBon
	var g = NewGame()
	var s slot.StatGeneric

	var calc = func(w io.Writer) float64 {
		var N, S, _ = s.NSQ(g.Cost())
		var rtpsym = S / N
		fmt.Fprintf(w, "RTP = %.6f%%\n", rtpsym*100)
		return rtpsym
	}

	return slot.ScanReelsCommon(ctx, sp, &s, g, reels, calc)
}

func CalcStatReg(ctx context.Context, sp *slot.ScanPar) float64 {
	fmt.Printf("*bonus reels calculations*\n")
	var rtpfs = CalcStatBon(ctx, sp)
	if ctx.Err() != nil {
		return 0
	}
	fmt.Printf("*regular reels calculations*\n")
	var reels, _ = ReelsMap.FindClosest(sp.MRTP)
	var g = NewGame()
	var s Stat

	var calc = func(w io.Writer) float64 {
		var N, S, _ = s.NSQ(g.Cost())
		var rtpsym = S / N
		var q3 = s.FSQ(3)
		var q4 = s.FSQ(4)
		var q5 = s.FSQ(5)
		var q6 = s.FSQ(6)
		var rtpqfs = q3*rtpfs*FreeMult[2] + q4*rtpfs*FreeMult[3] + q5*rtpfs*FreeMult[4] + q6*rtpfs*FreeMult[5]
		var rtp = rtpsym + rtpqfs
		fmt.Fprintf(w, "symbols: rtp(sym) = %.6f%%\n", rtpsym*100)
		fmt.Fprintf(w, "free spins %d, q3 = %.5g, q4 = %.5g, q5 = %.5g, q6 = %.5g\n",
			s.FreeCount[2].Load()+s.FreeCount[3].Load()+s.FreeCount[4].Load()+s.FreeCount[5].Load(),
			q3, q4, q5, q6)
		fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", s.FGF())
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g(fg) = %.6f%%\n", rtpsym*100, rtpqfs*100, rtp*100)
		return rtp
	}

	return slot.ScanReelsCommon(ctx, sp, &s, g, reels, calc)
}
