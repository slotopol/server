package jaguarmoon

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStatBon(ctx context.Context, sp *slot.ScanPar) float64 {
	var reels = ReelsBon
	var g = NewGame()
	var s = slot.NewStatGeneric(sn, 5)

	var calc = func(w io.Writer) float64 {
		var N, S, _ = s.NSQ(g.Cost())
		var µ = S / N
		fmt.Fprintf(w, "RTP = %.6f%%\n", µ*100)
		return µ
	}

	return slot.ScanReelsCommon(ctx, sp, s, g, reels, calc)
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
	var s = slot.NewStatGeneric(sn, 5)
	s.SymDim(scat, 6)

	var calc = func(w io.Writer) float64 {
		var N, S, _ = s.NSQ(g.Cost())
		var µ = S / N
		var q3 = float64(s.C[scat][3].Load()) * float64(ScatFreespin[2]) / N
		var q4 = float64(s.C[scat][4].Load()) * float64(ScatFreespin[3]) / N
		var q5 = float64(s.C[scat][5].Load()) * float64(ScatFreespin[4]) / N
		var q6 = float64(s.C[scat][6].Load()) * float64(ScatFreespin[5]) / N
		var fgh = s.C[scat][3].Load() + s.C[scat][4].Load() + s.C[scat][5].Load() + s.C[scat][6].Load()
		var rtpqfs = q3*rtpfs*FreeMult[2] +
			q4*rtpfs*FreeMult[3] +
			q5*rtpfs*FreeMult[4] +
			q6*rtpfs*FreeMult[5]
		var rtp = µ + rtpqfs
		fmt.Fprintf(w, "symbols: rtp(sym) = %.6f%%\n", µ*100)
		fmt.Fprintf(w, "free games %d, q3 = %.5g, q4 = %.5g, q5 = %.5g, q6 = %.5g\n",
			fgh, q3, q4, q5, q6)
		fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", N/float64(fgh))
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g(fg) = %.6f%%\n", µ*100, rtpqfs*100, rtp*100)
		return rtp
	}

	return slot.ScanReelsCommon(ctx, sp, s, g, reels, calc)
}
