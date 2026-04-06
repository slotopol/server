package greatblue

import (
	"context"
	"fmt"
	"io"
	"math"

	"github.com/slotopol/server/game/slot"
)

// combinations of multiplier & freespins number
// of two shells from set [x5, x8, 7, 10, 15]
var combs = []struct {
	m, L float64
}{
	{2 + 5 + 8, 8},
	{2 + 5, 8 + 7},
	{2 + 5, 8 + 10},
	{2 + 5, 8 + 15},
	{2 + 8, 8 + 7},
	{2 + 8, 8 + 10},
	{2 + 8, 8 + 15},
	{2, 8 + 7 + 10},
	{2, 8 + 7 + 15},
	{2, 8 + 10 + 15},
}

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	var s = slot.NewStatGeneric(sn, 5)
	var reels, _ = ReelsMap.FindClosest(sp.MRTP)
	var g = NewGame(sp.Sel)

	// custom parsheet
	var cost = g.Cost()
	var calc = func(w io.Writer) (float64, float64) {
		var µ, Dsym = slot.EvD(s, cost)
		var Pfg = s.FGQ()
		var q = Pfg * 15
		var sq = 1 / (1 - q)
		var ΣmN, Σm2N2, ΣV float64
		for i := range 10 {
			var mi, Li = combs[i].m, combs[i].L
			var Ni = Li + 15*q*sq
			ΣmN += mi * Ni
			Σm2N2 += mi * Ni * mi * Ni
			var Vi = mi * mi * (Ni*Dsym + 15*µ*15*µ*q*sq*sq*sq)
			ΣV += Vi
		}
		var EWbon = µ * ΣmN / 10
		var Vbon = ΣV/10 + µ*µ*Σm2N2/10 - EWbon*EWbon
		var rtp = µ + Pfg*EWbon
		var D = Dsym + Pfg*(Vbon+(EWbon-µ)*(EWbon-µ))
		if sp.IsMain() {
			fmt.Fprintf(w, "symbols: µ = %.8g%%, sigma(sym) = %.6g\n", µ*100, math.Sqrt(Dsym))
			fmt.Fprintf(w, "free: HRfg = 1/%.5g, q = %.5g, sq = 1/(1-q) = %.5g\n", 1/Pfg, q, sq)
			fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g = %.8g%%\n", µ*100, Pfg, EWbon*100, rtp*100)
		}
		slot.Print_all(w, sp, s, rtp, D)
		return rtp, D
	}

	return slot.ScanReelsCommon(ctx, sp, s, g, reels, calc)
}
