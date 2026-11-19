package hypercuber

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	fmt.Printf("\n(1/2) bonus reels calculations\n")
	var sb = slot.NewSnapCascade(sn, 15)
	{
		var g = NewGame()
		g.FSR = 15                                     // set free spins mode
		g.M = [5]float64{Mavr, Mavr, Mavr, Mavr, Mavr} // set multipliers to average value for RTP calculation
		var calc = func(w io.Writer) (float64, float64) {
			return slot.Parsheet_affixrecur(w, sp, sb, g.Cost())
		}
		slot.ScanReelsCommon2(ctx, sp, sb, g, calc)
	}

	if ctx.Err() != nil {
		return 0, 0
	}

	fmt.Printf("\n(2/2) regular reels calculations\n")
	var sr = slot.NewSnapCascade(sn, 15)
	{
		var g = NewGame()
		var calc = func(w io.Writer) (float64, float64) {
			return slot.Parsheet_dualreelrecur(w, sp, sr, sb, g.Cost())
		}
		return slot.ScanReelsCommon2(ctx, sp, sr, g, calc)
	}
}
