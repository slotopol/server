package africansimba

import (
	"context"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	var reels, _ = ReelsMap.FindClosest(sp.MRTP)
	var g = NewGame()
	var s = slot.NewStatGeneric(sn, 5)

	var calc = func(w io.Writer) (float64, float64) {
		return slot.Parsheet_generic_fgretrig(w, sp, s, g.Cost(), 3, 12)
	}

	return slot.ScanReelsCommon(ctx, sp, s, g, reels, calc)
}
