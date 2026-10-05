package royaldynasty

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func Efs(s *slot.StatGeneric) (HRfg, EK, EK2 float64) {
	var N = s.Count()
	var sum, sum2 float64
	for _, Li := range Freegames {
		sum += float64(Li)
		sum2 += float64(Li * Li)
	}
	HRfg = float64(s.Cfg.Load()) / N
	EK = HRfg * sum / float64(len(Freegames))
	EK2 = HRfg * sum2 / float64(len(Freegames))
	return
}

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	fmt.Printf("\n(1/2) bonus reels calculations\n")
	var sb = slot.NewStatGeneric(sn, 5)
	{
		var g = NewGame(sp.Sel)
		g.FSR = 35 // set free spins mode
		g.TS = scat1
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
			// bonus reels parameters
			var Nb = sb.Count()
			var qb = float64(sb.Cfg.Load()*35) / Nb
			// regular reels parameters
			var Nr = sr.Count()
			var qr = float64(sr.Cfg.Load()*35) / Nr
			var Pfgi = float64(sr.Cfg.Load()) / float64(len(Freegames)) / Nr
			var ΣPL float64
			for _, Li := range Freegames {
				ΣPL += Pfgi * float64(Li)
			}
			return slot.Parsheet_fgretrig_split_custom(w, sp, sr, sb, g.Cost(), 1, qr, qb, ΣPL)
		}
		return slot.ScanReelsCommon(ctx, sp, sr, g, calc)
	}
}
