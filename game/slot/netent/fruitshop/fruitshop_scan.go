package fruitshop

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func ΣPL(s *slot.StatGeneric) (sum float64) {
	var N = s.Count()
	for sym, L := range LineFreespinReg {
		for i, Li := range L {
			var Pfgi = float64(s.C[sym][i].Load()) / N
			sum += Pfgi * float64(Li)
		}
	}
	return
}

func Efs(s *slot.StatGeneric) (HRfg, EK, EK2 float64) {
	var N = s.Count()
	for sym, L := range LineFreespinReg {
		for i, Li := range L {
			var c = float64(s.C[sym][i].Load())
			HRfg += c
			EK += c * float64(Li)
			EK2 += c * float64(Li*Li)
		}
	}
	HRfg /= N
	EK /= N
	EK2 /= N
	return
}

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	fmt.Printf("\n(1/2) free games calculations\n")
	var sb = slot.NewStatGeneric(sn, 5)
	{
		var g = NewGame()
		g.FSR = 5 // set free spins mode
		var calc = func(w io.Writer) (float64, float64) {
			return slot.Parsheet_fgrecur(w, sp, sb, g.Cost(), 1, func() (float64, float64, float64) {
				return Efs(sb)
			})
		}
		slot.ScanReelsCommon(ctx, sp, sb, g, calc)
	}

	if ctx.Err() != nil {
		return 0, 0
	}

	fmt.Printf("\n(2/2) regular games calculations\n")
	var sr = slot.NewStatGeneric(sn, 5)
	{
		var g = NewGame()
		var calc = func(w io.Writer) (float64, float64) {
			return slot.Parsheet_fgretrig_split_custom(w, sp, sr, sb, g.Cost(), 1, sr.FSQ(), sb.FSQ(), ΣPL(sr))
		}
		return slot.ScanReelsCommon(ctx, sp, sr, g, calc)
	}
}
