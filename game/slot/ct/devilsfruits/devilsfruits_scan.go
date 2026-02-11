package devilsfruits

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) float64 {
	var reels, _ = ReelsMap.FindClosest(sp.MRTP)
	var g = NewGame(sp.Sel)
	var s slot.StatGeneric

	var calc = func(w io.Writer) float64 {
		var N, S, _ = s.NSQ(g.Cost())
		var rtp = S / N
		fmt.Fprintf(w, "RTP = %.6f%%\n", rtp*100)
		return rtp
	}

	return slot.ScanReelsCommon(ctx, sp, &s, g, reels, calc)
}
