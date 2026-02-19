package smileyxwild

import (
	"context"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) float64 {
	var reels, _ = ReelsMap.FindClosest(sp.MRTP)
	var g = NewGame(sp.Sel)
	g.M2 = 3 // average wild multiplier on 2 reel
	g.M4 = 3 // average wild multiplier on 4 reel
	var s = slot.NewStatGeneric(sn, 5)

	var calc = func(w io.Writer) (rtp float64) {
		rtp, _ = slot.Parsheet_generic_simple(w, sp, s, g.Cost())
		return
	}

	return slot.ScanReelsCommon(ctx, sp, s, g, reels, calc)
}
