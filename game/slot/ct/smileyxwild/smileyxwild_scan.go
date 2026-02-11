package smileyxwild

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) float64 {
	var reels, _ = ReelsMap.FindClosest(sp.MRTP)
	var g = NewGame(sp.Sel)
	g.M2 = 3 // average wild multiplier on 2 reel
	g.M4 = 3 // average wild multiplier on 4 reel
	var s slot.StatGeneric

	var calc = func(w io.Writer) float64 {
		var lrtp, srtp = s.RTPsym(g.Cost(), scat)
		var rtpsym = lrtp + srtp
		fmt.Fprintf(w, "RTP = %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp*100, srtp*100, rtpsym*100)
		return rtpsym
	}

	return slot.ScanReelsCommon(ctx, sp, &s, g, reels, calc)
}
