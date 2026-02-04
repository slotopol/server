package greatblue

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func FirstSreespins() (fsavr1 float64, multavr float64) {
	// combinations of multiplier & freespins number
	// of two shells from set [x5, x8, 7, 10, 15]
	var combs = []struct {
		mult, fsnum float64
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
	for _, c := range combs {
		fsavr1 += c.mult * c.fsnum
		multavr += c.mult
	}
	fsavr1 /= float64(len(combs))
	multavr /= float64(len(combs))
	return
}

func CalcStat(ctx context.Context, mrtp float64) float64 {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	var g = NewGame(5)
	var s slot.StatGeneric

	var calc = func(w io.Writer) float64 {
		var N = s.Count()
		var lrtp, srtp = s.SymRTP(g.Cost())
		var rtpsym = lrtp + srtp
		var fghits = float64(s.FreeHits.Load())
		var fsavr1, multavr = FirstSreespins()
		var q = fghits * fsavr1 / N
		var sq = 1 / (1 - fghits*multavr*15/N)
		var rtp = rtpsym + q*sq*rtpsym
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp*100, srtp*100, rtpsym*100)
		fmt.Fprintf(w, "average plain freespins at 1st iteration: %g\n", fsavr1)
		fmt.Fprintf(w, "average multiplier at free games: %g\n", multavr)
		fmt.Fprintf(w, "free games %g, q = %.5g, sq = %.5g\n", fghits, q, sq)
		fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", N/fghits)
		fmt.Fprintf(w, "RTP = rtpsym + q*sq*rtpsym = %.5g + %.5g = %.6f%%\n", rtpsym*100, q*sq*rtpsym*100, rtp*100)
		return rtp
	}

	return slot.ScanReelsCommon(ctx, &s, g, reels, calc)
}
