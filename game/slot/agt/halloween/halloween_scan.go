package halloween

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
		var lrtp = s.LineRTP(g.Cost())
		fmt.Fprintf(w, "RTP = %.6f%%\n", lrtp*100)
		return lrtp
	}

	return slot.ScanReelsCommon(ctx, &s, g, reels, calc)
}
