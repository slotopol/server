package panda

import (
	"context"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	var reels, _ = ReelsMap.FindClosest(sp.MRTP)
	var g = NewGame(sp.Sel)
	var s = slot.NewStatGeneric(sn, 3)

	var calc = func(w io.Writer) (float64, float64) {
		return slot.Parsheet_generic_fgretrig_series(w, sp, s, g.Cost(), 1, []int{1, 2, 3}, scat)
	}

	return slot.ScanReelsCommon(ctx, sp, s, g, reels, calc)
}
