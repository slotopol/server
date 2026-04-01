package beetlemania

import (
	"context"
	"fmt"
	"io"
	"math"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	var s = slot.NewStatGeneric(sn, 5)
	var reels, _ = ReelsMap.FindClosest(sp.MRTP)
	var g = NewGame(sp.Sel)
	// custom parsheet
	var njb float64
	for _, sym := range reels[2] {
		if sym == jazz {
			njb++
		}
	}
	var Pjb = njb * 3 / float64(len(reels[2]))
	var cost = g.Cost()
	const L = 10.0
	var calc = func(w io.Writer) (float64, float64) {
		var N, S, Q = s.NSQ(cost)
		var µ = S / N
		var Dsym = Q/N - µ*µ
		var Pfg = s.FGQ()
		var p = Pjb
		var mjb = (math.Pow(1+p, L) - 1) / p
		var md = math.Pow(1+3*p, L)
		var Wavg = µ * mjb
		var rtp = µ + Pfg*Wavg
		var Vbon = L*L*µ*µ*(md-math.Pow(1+p, 2*L)) + md*L*Dsym
		var D = Dsym + Pfg*(Vbon+(Wavg-µ)*(Wavg-µ))
		if sp.IsMain() {
			fmt.Fprintf(w, "symbols: µ = %.8g%%, sigma(sym) = %.6g\n", µ*100, math.Sqrt(Dsym))
			fmt.Fprintf(w, "jazzbee: HRjb = 1/%.5g, mjb = %.5g\n", 1/p, mjb)
			fmt.Fprintf(w, "free games: HRfg = 1/%.5g\n", 1/Pfg)
			fmt.Fprintf(w, "average bonus win: Wavg = %.5g\n", Wavg)
			fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.8g%%\n", µ*100, Pfg, Wavg*100, rtp*100)
		}
		slot.Print_all(w, sp, s, rtp, D)
		return rtp, D
	}
	return slot.ScanReelsCommon(ctx, sp, s, g, reels, calc)
}
