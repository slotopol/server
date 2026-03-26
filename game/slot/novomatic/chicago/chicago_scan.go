package chicago

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
	const L = 12
	var cost = g.Cost()
	var calc = func(w io.Writer) (float64, float64) {
		var N, S, Q = s.NSQ(cost)
		var µ = S / N
		var Dsym = Q/N - µ*µ
		var q = s.FSQ()
		var sq = 1 / (1 - q)
		var Pfg = s.FGQ()
		var Em1, Em2 float64
		for _, m := range MultChoose {
			Em1 += m
			Em2 += m * m
		}
		Em1 /= float64(len(MultChoose))
		Em2 /= float64(len(MultChoose))
		var rtpfs = Em1 * sq * µ
		var rtp = µ + q*rtpfs
		var Eser, Dser = L * sq, L * q * sq * sq * sq
		var D = Dsym + Pfg*(Em2*Eser*Dsym+Em1*Em1*µ*µ*Dser)
		if sp.IsMain() {
			fmt.Fprintf(w, "symbols: µ = %.8g%%, sigma(sym) = %.6g\n", µ*100, math.Sqrt(Dsym))
			fmt.Fprintf(w, "free spins: q = %.5g, sq = 1/(1-q) = %.6f\n", q, sq)
			fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", 1/Pfg)
			fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.8g%%\n", µ*100, q, rtpfs*100, rtp*100)
		}
		slot.Print_all(w, sp, s, rtp, D)
		return rtp, D
	}

	return slot.ScanReelsCommon(ctx, sp, s, g, reels, calc)
}
