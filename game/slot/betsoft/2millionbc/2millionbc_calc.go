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
	var g = NewGame()
	g.Sel = 1
	g.FSR = 4 // set free spins mode
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var reshuf = s.Count()
		var cost, _ = g.Cost()
		var lrtp, srtp = s.LineRTP(cost), s.ScatRTP(cost)
		var rtpsym = lrtp + srtp
		var q = s.FreeCount() / reshuf
		var sq = 1 / (1 - q)
		var rtp = sq * rtpsym
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCountU(), q, sq)
		fmt.Fprintf(w, "free games frequency: 1/%.5g\n", reshuf/s.FreeHits())
		fmt.Fprintf(w, "RTP = sq*rtp(sym) = %.5g*%.5g = %.6f%%\n", sq, rtpsym, rtp)
		return rtp
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc)
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
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	var g = NewGame()
	g.Sel = 1
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var reshuf = s.Count()
		var cost, _ = g.Cost()
		var lrtp, srtp = s.LineRTP(cost), s.ScatRTP(cost)
		var rtpsym = lrtp + srtp
		var q = s.FreeCount() / reshuf
		var sq = 1 / (1 - q)
		var qacbn = 1 / float64(len(reels.Reel(5)))
		var rtpacbn = Eacbn * qacbn * 100
		var qdlbn = s.BonusCount(dlbn) / reshuf / float64(g.Sel)
		var rtpdlbn = Edlbn * qdlbn * 100
		var rtp = rtpsym + rtpacbn + rtpdlbn + q*rtpfs
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCountU(), q, sq)
		fmt.Fprintf(w, "free games frequency: 1/%.5g\n", reshuf/s.FreeHits())
		fmt.Fprintf(w, "acorn bonuses: frequency 1/%d, rtp = %.6f%%\n", len(reels.Reel(5)), rtpacbn)
		fmt.Fprintf(w, "diamond lion bonuses: frequency 1/%.5g, rtp = %.6f%%\n", reshuf/s.BonusCount(dlbn), rtpdlbn)
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g(acorn) + %.5g(dl) + %.5g*%.5g(fg) = %.6f%%\n", rtpsym, rtpacbn, rtpdlbn, q, rtpfs, rtp)
		return rtp
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc)
}
