package smileyxwild

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, mrtp float64) float64 {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	var g = NewGame(1)
	g.M2 = 3 // average wild multiplier on 2 reel
	g.M4 = 3 // average wild multiplier on 4 reel
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var lrtp, srtp = s.SymRTP(g.Cost())
		var rtpsym = lrtp + srtp
		fmt.Fprintf(w, "RTP = %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		return rtpsym
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc)
}
