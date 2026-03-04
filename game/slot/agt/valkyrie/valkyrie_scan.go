package valkyrie

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/slotopol/server/game/slot"
)

func CalcStatReg(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	fmt.Printf("\n(1/2) bonus reels calculations\n")
	var sb = slot.NewStatGeneric(sn, 5)
	{
		var reels, _ = ReelsMap.FindClosest(sp.MRTP)
		var g = NewGame(sp.Sel)
		g.FSR = 15 // set free spins mode
		var calc = func(w io.Writer) (float64, float64) {
			return slot.Parsheet_generic_simple(w, sp, sb, g.Cost())
		}
		func() {
			var ctx2, cancel2 = context.WithCancel(ctx)
			defer cancel2()
			slot.BruteForce5x3Big(ctx2, sp, sb, g, reels.Reel(1), ReelBig, reels.Reel(5))
		}()
		calc(os.Stdout)
	}

	if ctx.Err() != nil {
		return 0, 0
	}

	fmt.Printf("\n(2/2) regular reels calculations\n")
	var sr = slot.NewStatGeneric(sn, 5)
	{
		var reels, _ = ReelsMap.FindClosest(sp.MRTP)
		var g = NewGame(sp.Sel)
		var calc = func(w io.Writer) (float64, float64) {
			return slot.Parsheet_generic_fgonce_split(w, sp, sr, sb, g.Cost(), 1, 15)
		}
		return slot.ScanReelsCommon(ctx, sp, sr, g, reels, calc)
	}
}
