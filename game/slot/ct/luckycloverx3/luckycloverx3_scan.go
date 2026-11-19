package luckycloverx3

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

// custom parsheet
func Parsheet(w io.Writer, sp *slot.ScanPar, ss *slot.StatSnap, cost float64) (float64, float64) {
	var µ, D = slot.EvD(ss, cost)
	var N = ss.Count()
	var cbx2, cbx3 = float64(ss.A[0].Cb[bonx2-1]), float64(ss.A[0].Cb[bonx3-1])
	var HRx2 = cbx2 / N
	var HRx3 = cbx3 / N
	var HRxx = (cbx2 + cbx3) / N
	if sp.IsMain() {
		fmt.Fprintf(w, "bonus mult: HRx2 = 1/%.5g, HRx3 = 1/%.5g, HRxx = 1/%.5g\n", 1/HRx2, 1/HRx3, 1/HRxx)
		fmt.Fprintf(w, "RTP = %.8g%%\n", µ*100)
	}
	slot.Print_all(w, sp, ss, µ, D)
	return µ, D
}

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	var g = NewGame(sp.Sel)
	var s = slot.NewStatSnap(1, sn, 5, bonx3, 0)

	var calc = func(w io.Writer) (float64, float64) {
		return Parsheet(w, sp, s, g.Cost())
	}

	return slot.ScanReelsCommon2(ctx, sp, s, g, calc)
}
