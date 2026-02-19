package rainbowcharm

import (
	"context"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) float64 {
	var reels, _ = ReelsMap.FindClosest(sp.MRTP)
	var g = NewGame()
	g.M = [5]float64{4, 4, 4, 4, 4} // set multipliers to average value for RTP calculation
	var s = slot.NewStatGeneric(sn, 15)

	var calc = func(w io.Writer) (rtp float64) {
		rtp, _ = slot.Parsheet_generic_simple(w, sp, s, g.Cost())
		return
	}

	return slot.ScanReelsCommon(ctx, sp, s, g, reels, calc)
}
