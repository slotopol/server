package bookofra

import (
	"context"
	"fmt"
	"io"
	"math"

	"github.com/slotopol/server/game/slot"
)

func CalcStatBon(ctx context.Context, sp *slot.ScanPar, s *slot.StatGeneric, es slot.Sym) (float64, float64) {
	var reels = ReelsBon
	var g = NewGame(sp.Sel)
	g.FSR = 10 // set free spins mode
	g.ES = es
	var calc = func(w io.Writer) (float64, float64) {
		return slot.Parsheet_fgretrig(w, sp, s, g.Cost(), 1, 10)
	}
	return slot.ScanReelsCommon(ctx, sp, s, g, reels, calc)
}

// custom parsheet
func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	const L = 10
	var g = NewGame(sp.Sel)
	var cost = g.Cost()
	var Σrtp, Σµ, Σµ2, ΣD float64
	var q, sq float64
	var es slot.Sym
	for es = 1; es < wsc; es++ {
		fmt.Printf("\n(%d/16) calculations for expanding symbol [%d]\n", es, es)
		var s = slot.NewStatGeneric(sn, 5)
		var rtpi, _ = CalcStatBon(ctx, sp, s, es)
		if ctx.Err() != nil {
			return 0, 0
		}
		if es == 1 {
			q = s.FSQ()
			sq = 1 / (1 - q)
		}
		var µ, D = slot.EvD(s, cost)
		Σrtp += rtpi
		Σµ += µ
		Σµ2 += µ * µ
		ΣD += D
	}
	var rtpfs = Σrtp / 9
	var µB = Σµ / 9
	var DB = ΣD/9 + Σµ2/9 - µB
	var Vbon = (L*sq)*DB + L*µB*L*µB*q*sq*sq*sq
	var Ebon = (L * sq) * µB

	fmt.Printf("\n(10/10) regular reels calculations\n")
	var reels, _ = ReelsMap.FindClosest(sp.MRTP)
	var sr = slot.NewStatGeneric(sn, 5)
	var calc = func(w io.Writer) (float64, float64) {
		var µr, Dsymr = slot.EvD(sr, cost)
		var qr = sr.FSQ()
		var sqr = 1 / (1 - qr)
		var Pfg = sr.FGQ()
		var rtp = µr + qr*rtpfs
		var D = Dsymr + Pfg*(Vbon+(Ebon-µr)*(Ebon-µr))
		if sp.IsMain() {
			fmt.Fprintf(w, "RTP(fg) = %.6f%%\n", rtpfs*100)
			fmt.Fprintf(w, "symbols: µ = %.8g%%, sigma(sym) = %.6g\n", µr*100, math.Sqrt(Dsymr))
			fmt.Fprintf(w, "free: HRfg = 1/%.5g, q = %.5g, sq = 1/(1-q) = %.5g\n", 1/Pfg, qr, sqr)
			fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.8g%%\n", µr*100, qr, rtpfs*100, rtp*100)
		}
		slot.Print_all(w, sp, sr, rtp, D)
		return rtp, D
	}
	return slot.ScanReelsCommon(ctx, sp, sr, g, reels, calc)
}
