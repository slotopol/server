package twomillionbc

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

var Eacbn float64

func ExpAcorn() {
	var sum float64
	for _, v := range Acorn {
		sum += float64(v)
	}
	Eacbn = sum / float64(len(Acorn))
}

var Edlbn float64

func ExpDiamondLion() {
	var sum float64
	for _, v := range DiamondLion {
		sum += float64(v)
	}
	Edlbn = sum / float64(len(DiamondLion))
}

func CalcStatBon(ctx context.Context) float64 {
	var reels = ReelsBon
	var g = NewGame(1)
	g.FSR = 4 // set free spins mode
	var s slot.StatGeneric

	var calc = func(w io.Writer) float64 {
		var lrtp, srtp = s.SymRTP(g.Cost())
		var rtpsym = lrtp + srtp
		var q, sq = s.FSQ()
		var rtp = sq * rtpsym
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp*100, srtp*100, rtpsym*100)
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCount.Load(), q, sq)
		fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", s.FGF())
		fmt.Fprintf(w, "RTP = sq*rtp(sym) = %.5g*%.5g = %.6f%%\n", sq, rtpsym*100, rtp*100)
		return rtp
	}

	return slot.ScanReelsCommon(ctx, &s, g, reels, calc)
}

func CalcStatReg(ctx context.Context, mrtp float64) float64 {
	fmt.Printf("*bonus games calculations*\n")
	ExpAcorn()
	fmt.Printf("len = %d, E = %g\n", len(Acorn), Eacbn)
	ExpDiamondLion()
	fmt.Printf("len = %d, E = %g\n", len(DiamondLion), Edlbn)
	fmt.Printf("*bonus reels calculations*\n")
	var rtpfs = CalcStatBon(ctx)
	if ctx.Err() != nil {
		return 0
	}
	fmt.Printf("*regular reels calculations*\n")
	var reels, _ = ReelsMap.FindClosest(mrtp)
	var g = NewGame(1)
	var s slot.StatGeneric

	var calc = func(w io.Writer) float64 {
		var N = s.Count()
		var lrtp, srtp = s.SymRTP(g.Cost())
		var rtpsym = lrtp + srtp
		var q, sq = s.FSQ()
		var qacbn = 1 / float64(len(reels.Reel(5)))
		var rtpacbn = Eacbn * qacbn
		var qdlbn = s.BonusHitsF(dlbn) / N / float64(g.Sel)
		var rtpdlbn = Edlbn * qdlbn
		var rtp = rtpsym + rtpacbn + rtpdlbn + q*rtpfs
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp*100, srtp*100, rtpsym*100)
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCount.Load(), q, sq)
		fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", s.FGF())
		fmt.Fprintf(w, "acorn bonuses: hit rate 1/%d, rtp = %.6f%%\n", len(reels.Reel(5)), rtpacbn*100)
		fmt.Fprintf(w, "diamond lion bonuses: hit rate 1/%.5g, rtp = %.6f%%\n", N/s.BonusHitsF(dlbn), rtpdlbn*100)
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g(acorn) + %.5g(dl) + %.5g*%.5g(fg) = %.6f%%\n", rtpsym*100, rtpacbn*100, rtpdlbn*100, q, rtpfs*100, rtp*100)
		return rtp
	}

	return slot.ScanReelsCommon(ctx, &s, g, reels, calc)
}
