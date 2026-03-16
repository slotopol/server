package hypercuber

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	fmt.Printf("\n(1/2) bonus reels calculations\n")
	var sb = slot.NewStatCascade(sn, 15)
	{
		var reels = ReelsBon
		var g = NewGame()
		g.FSR = 15                                     // set free spins mode
		g.M = [5]float64{Mavr, Mavr, Mavr, Mavr, Mavr} // set multipliers to average value for RTP calculation
		var calc = func(w io.Writer) (float64, float64) {
			return slot.Parsheet_fgretrig(w, sp, sb, g.Cost(), 1, 15)
		}
		slot.ScanReelsCommon(ctx, sp, sb, g, reels, calc)
	}

	if ctx.Err() != nil {
		return 0, 0
	}

	fmt.Printf("\n(2/2) regular reels calculations\n")
	var sr = slot.NewStatCascade(sn, 15)
	{
		var reels, _ = ReelsMap.FindClosest(sp.MRTP)
		var g = NewGame()
		var calc = func(w io.Writer) (float64, float64) {
			return slot.Parsheet_fgretrig_split(w, sp, sr, sb, g.Cost(), 1, 15)
		}
		return slot.ScanReelsCommon(ctx, sp, sr, g, reels, calc)
	}
}
