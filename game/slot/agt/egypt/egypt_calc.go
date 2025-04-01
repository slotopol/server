package egypt

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

// Minislot expectation calculation:
// total combinations: 3*3*3 = 27
//
//	x1 = 27-3=24
//	x3 = 1
//	x6 = 1
//	x9 = 1
//
// Em = (24*1 + 1*3 + 1*6 + 1*9)/27 = 42/27 = 1.5555555556
const Em = 42. / 27.

func CalcStat(ctx context.Context, mrtp float64) float64 {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	var g = NewGame()
	g.Sel = 1
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var cost, _ = g.Cost()
		var lrtp, srtp = s.LineRTP(cost), s.ScatRTP(cost)
		var rtpsym = lrtp + srtp
		var rtp = rtpsym * Em
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "RTP = %.5g(sym) * %.5g(Em) = %.6f%%\n", rtpsym, Em, rtp)
		return rtpsym
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc)
}
