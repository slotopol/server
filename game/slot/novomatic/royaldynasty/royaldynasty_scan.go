package royaldynasty

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

type Stat struct {
	*slot.StatSnap
}

func NewStat() *Stat {
	return &Stat{
		StatSnap: slot.NewSnapGeneric(sn, 5),
	}
}

func (s *Stat) FGC() (Cfg, Cfs, Qfs float64) {
	var sum, sum2 float64
	for _, Li := range Freegames {
		sum += float64(Li)
		sum2 += float64(Li * Li)
	}
	Cfg = float64(s.A[0].Cfg)
	Cfs = Cfg * sum / float64(len(Freegames))
	Qfs = Cfg * sum2 / float64(len(Freegames))
	return
}

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	fmt.Printf("\n(1/2) bonus reels calculations\n")
	var sb = NewStat()
	{
		var g = NewGame(sp.Sel)
		g.FSR = 35 // set free spins mode
		g.TS = scat1
		var calc = func(w io.Writer) (float64, float64) {
			return slot.Parsheet_affixrecur(w, sp, sb, g.Cost())
		}
		slot.ScanReelsCommon2(ctx, sp, sb, g, calc)
	}

	if ctx.Err() != nil {
		return 0, 0
	}

	fmt.Printf("\n(2/2) regular reels calculations\n")
	var sr = NewStat()
	{
		var g = NewGame(sp.Sel)
		var calc = func(w io.Writer) (float64, float64) {
			return slot.Parsheet_dualreelrecur(w, sp, sr, sb, g.Cost())
		}
		return slot.ScanReelsCommon2(ctx, sp, sr, g, calc)
	}
}
