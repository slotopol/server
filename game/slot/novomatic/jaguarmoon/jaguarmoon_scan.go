package jaguarmoon

import (
	"context"
	"fmt"
	"io"

	"go.uber.org/atomic"

	"github.com/slotopol/server/game/slot"
)

type Stat struct {
	slot.Stat
	FreeCount [6]atomic.Uint64
}

// Declare conformity with Stater interface.
var _ slot.Stater = (*Stat)(nil)

func (s *Stat) FreeCountF(n int) float64 {
	return float64(s.FreeCount[n-1].Load())
}

// Returns (q, sq), where q = free spins quantifier, sq = 1/(1-q)
// sum of a decreasing geometric progression for retriggered free spins.
func (s *Stat) FSQ(n int) (q float64) {
	q = s.FreeCountF(n) / s.Count()
	return
}

func (s *Stat) Update(wins slot.Wins, cfn int) {
	var lpay, spay float64
	for _, wi := range wins {
		if wi.Pay != 0 {
			if wi.LI != 0 { // line win
				lpay += wi.Pay * wi.MP
			} else { // scatter win
				spay += wi.Pay * wi.MP
			}
		}
		if wi.FS != 0 {
			s.FreeCount[wi.Num-1].Add(uint64(wi.FS))
			s.FreeHits.Inc()
		}
		if wi.BID != 0 {
			s.BonCount[wi.BID].Inc()
		}
		if wi.JID != 0 {
			s.JackCount[wi.JID].Inc()
		}
	}
	if lpay != 0 {
		s.LinePay.Add(lpay)
	}
	if spay != 0 {
		s.ScatPay.Add(spay)
	}
	if cfn <= slot.FallLimit {
		s.Falls[cfn-1].Inc()
	}
}

func CalcStatBon(ctx context.Context) float64 {
	var reels = ReelsBon
	var g = NewGame()
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var lrtp, srtp = s.SymRTP(g.Cost())
		var rtpsym = lrtp + srtp
		fmt.Fprintf(w, "RTP = %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		return rtpsym
	}

	return slot.ScanReelsCommon(ctx, &s, g, reels, calc)
}

func CalcStatReg(ctx context.Context, mrtp float64) float64 {
	fmt.Printf("*bonus reels calculations*\n")
	var rtpfs = CalcStatBon(ctx)
	if ctx.Err() != nil {
		return 0
	}
	fmt.Printf("*regular reels calculations*\n")
	var reels, _ = ReelsMap.FindClosest(mrtp)
	var g = NewGame()
	var s Stat

	var calc = func(w io.Writer) float64 {
		var lrtp, srtp = s.SymRTP(g.Cost())
		var rtpsym = lrtp + srtp
		var q3 = s.FSQ(3)
		var q4 = s.FSQ(4)
		var q5 = s.FSQ(5)
		var q6 = s.FSQ(6)
		var rtpqfs = q3*rtpfs*FreeMult[2] + q4*rtpfs*FreeMult[3] + q5*rtpfs*FreeMult[4] + q6*rtpfs*FreeMult[5]
		var rtp = rtpsym + rtpqfs
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "free spins %d, q3 = %.5g, q4 = %.5g, q5 = %.5g, q6 = %.5g\n",
			s.FreeCount[2].Load()+s.FreeCount[3].Load()+s.FreeCount[4].Load()+s.FreeCount[5].Load(),
			q3, q4, q5, q6)
		fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", s.FGF())
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g(fg) = %.6f%%\n", rtpsym, rtpqfs, rtp)
		return rtp
	}

	return slot.ScanReelsCommon(ctx, &s, g, reels, calc)
}
