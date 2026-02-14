package shiningstars

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) float64 {
	var reels, _ = ReelsMap.FindClosest(sp.MRTP)
	var g = NewGame(sp.Sel)
	var s = slot.NewStatGeneric(sn, 5)

	var calc = func(w io.Writer) float64 {
		var lrtp, srtp = s.RTPsym2(g.Cost(), scat1, scat2)
		var rtpsym = lrtp + srtp
		fmt.Fprintf(w, "RTP = %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp*100, srtp*100, rtpsym*100)
		return rtpsym
	}

	return slot.ScanReelsCommon(ctx, sp, s, g, reels, calc)
}
