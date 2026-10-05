package wildhorses

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func Efs(s *slot.StatGeneric) (HRfg, EK, EK2 float64) {
	var N = s.Count()
	var c = float64(s.Cfg.Load())
	const p = 0.5 // probability of getting or not getting free spins is equal
	const L = 10  // number of free spins
	HRfg = p * c / N
	EK = p * c * L / N
	EK2 = p * c * L * L / N
	return
}

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	fmt.Printf("\n(1/2) bonus reels calculations\n")
	var sb = slot.NewStatGeneric(sn, 5)
	{
		var g = NewGame(sp.Sel)
		g.FSR = 10 // set free spins mode
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

	fmt.Printf("\n(2/2) regular reels calculations\n")
	var sr = slot.NewStatGeneric(sn, 5)
	{
		var g = NewGame(sp.Sel)
		var calc = func(w io.Writer) (float64, float64) {
			var qr = sr.FSQ()
			var qb = 0.5 // probability of getting or not getting free spins is equal
			var ΣPL = sb.FGQ() * 10
			return slot.Parsheet_fgretrig_split_custom(w, sp, sr, sb, g.Cost(), 1, qr, qb, ΣPL)
		}
		return slot.ScanReelsCommon(ctx, sp, sr, g, calc)
	}
}
