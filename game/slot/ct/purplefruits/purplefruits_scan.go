package purplefruits

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, mrtp float64) float64 {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	var g = NewGame(1)
	var s slot.StatGeneric

	var calc = func(w io.Writer) float64 {
		var lrtp, srtp = s.SymRTP(g.Cost())
		var rtpsym = lrtp + srtp
		fmt.Fprintf(w, "RTP = %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp*100, srtp*100, rtpsym*100)
		return rtpsym
	}

	return slot.ScanReelsCommon(ctx, &s, g, reels, calc)
}
