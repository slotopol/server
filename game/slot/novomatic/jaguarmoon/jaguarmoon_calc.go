package jaguarmoon

import (
	"context"
	"fmt"
	"io"
	"sync"
	"sync/atomic"

	"github.com/slotopol/server/game/slot"
)

type Stat struct {
	planned    uint64
	reshuffles [slot.FallLimit]uint64
	errcount   uint64
	linepay    float64
	scatpay    float64
	freecount  uint64
	freehits   uint64
	bonuscount [8]uint64
	jackcount  [4]uint64
	lpm, spm   sync.Mutex
}

// Declare conformity with Stater interface.
var _ slot.Stater = (*Stat)(nil)

func (s *Stat) SetPlan(n uint64) {
	atomic.StoreUint64(&s.planned, n)
}

func (s *Stat) Planned() float64 {
	return float64(atomic.LoadUint64(&s.planned))
}

func (s *Stat) Count() float64 {
	var n uint64
	for i := range slot.FallLimit {
		n += atomic.LoadUint64(&s.reshuffles[i])
	}
	return float64(n)
}

func (s *Stat) Reshuf(cfn int) float64 {
	var n uint64
	for i := cfn - 1; i < slot.FallLimit; i++ {
		n += atomic.LoadUint64(&s.reshuffles[i])
	}
	return float64(n)
}

func (s *Stat) Errors() float64 {
	return float64(atomic.LoadUint64(&s.errcount))
}

func (s *Stat) IncErr() {
	atomic.AddUint64(&s.errcount, 1)
}

func (s *Stat) LineRTP(cost float64) float64 {
	s.lpm.Lock()
	var lp = s.linepay
	s.lpm.Unlock()
	return lp / s.Count() / cost * 100
}

func (s *Stat) ScatRTP(cost float64) float64 {
	s.spm.Lock()
	var sp = s.scatpay
	s.spm.Unlock()
	return sp / s.Count() / cost * 100
}

func (s *Stat) SymRTP(cost float64) (lrtp, srtp float64) {
	s.lpm.Lock()
	var lp = s.linepay
	s.lpm.Unlock()
	s.spm.Lock()
	var sp = s.scatpay
	s.spm.Unlock()
	var reshuf = s.Count()
	lrtp = lp / reshuf / cost * 100
	srtp = sp / reshuf / cost * 100
	return
}

func (s *Stat) FreeCountU() uint64 {
	return atomic.LoadUint64(&s.freecount)
}

func (s *Stat) FreeCount() float64 {
	return float64(atomic.LoadUint64(&s.freecount))
}

// Returns (q, sq), where q = free spins quantifier, sq = 1/(1-q)
// sum of a decreasing geometric progression for retriggered free spins.
func (s *Stat) FSQ() (q float64, sq float64) {
	q = s.FreeCount() / s.Count()
	sq = 1 / (1 - q)
	return
}

func (s *Stat) FreeHits() float64 {
	return float64(atomic.LoadUint64(&s.freehits))
}

// Quantifier of free games per reshuffles.
func (s *Stat) FGQ() float64 {
	return s.FreeHits() / s.Count()
}

// Free Games Frequency: average number of reshuffles per free games hit.
func (s *Stat) FGF() float64 {
	return s.Count() / s.FreeHits()
}

func (s *Stat) BonusCount(bid int) float64 {
	return float64(atomic.LoadUint64(&s.bonuscount[bid]))
}

func (s *Stat) JackCount(jid int) float64 {
	return float64(atomic.LoadUint64(&s.jackcount[jid]))
}

func (s *Stat) Update(wins slot.Wins, cfn int) {
	for _, wi := range wins {
		if wi.Pay != 0 {
			if wi.LI != 0 {
				s.lpm.Lock()
				s.linepay += wi.Pay * wi.MP
				s.lpm.Unlock()
			} else {
				s.spm.Lock()
				s.scatpay += wi.Pay * wi.MP
				s.spm.Unlock()
			}
		}
		if wi.FS != 0 {
			atomic.AddUint64(&s.freecount, uint64(wi.FS*int(FreeMult[wi.Num-1])))
			atomic.AddUint64(&s.freehits, 1)
		}
		if wi.BID != 0 {
			atomic.AddUint64(&s.bonuscount[wi.BID], 1)
		}
		if wi.JID != 0 {
			atomic.AddUint64(&s.jackcount[wi.JID], 1)
		}
	}
	if cfn <= slot.FallLimit {
		atomic.AddUint64(&s.reshuffles[cfn-1], 1)
	}
}

func CalcStatBon(ctx context.Context) float64 {
	var reels = ReelsBon
	var g = NewGame()
	g.Sel = 10 // bet on 243 ways
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var cost, _ = g.Cost()
		var lrtp, srtp = s.SymRTP(cost)
		var rtpsym = lrtp + srtp
		fmt.Fprintf(w, "RTP = %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		return rtpsym
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc)
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
	g.Sel = 10 // bet on 243 ways
	var s Stat

	var calc = func(w io.Writer) float64 {
		var cost, _ = g.Cost()
		var lrtp, srtp = s.SymRTP(cost)
		var rtpsym = lrtp + srtp
		var q, sq = s.FSQ()
		var rtp = rtpsym + q*rtpfs
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCountU(), q, sq)
		fmt.Fprintf(w, "free games frequency: 1/%.5g\n", s.FGF())
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%\n", rtpsym, q, rtpfs, rtp)
		return rtp
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc)
}
