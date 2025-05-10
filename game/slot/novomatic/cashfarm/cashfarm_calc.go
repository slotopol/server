package cashfarm

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

// Bonus game expectation
// calculation see at helper/prov/novomatic/cashfarmbon.lua
const Ebon = 50

func CalcStat(ctx context.Context, mrtp float64) float64 {
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	var g = NewGame()
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var reshuf1 = s.Reshuf(1)
		var reshuf2 = s.Reshuf(2)
		var reshuf3 = s.Reshuf(3)
		var reshuf4 = s.Reshuf(4)
		var reshuf5 = s.Reshuf(5)
		var cost, _ = g.Cost()
		var lrtp, srtp = s.LineRTP(cost), s.ScatRTP(cost)
		var rtpsym = lrtp + srtp
		var qfarm = s.BonusCount(farmbn) / reshuf1
		var rtpbon = Ebon * qfarm * 100
		var rtp = rtpsym + rtpbon
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "fall[2] = %.10g, freq = 1/%.5g\n", reshuf2, reshuf1/reshuf2)
		fmt.Fprintf(w, "fall[3] = %.10g, freq = 1/%.5g, freq2 = 1/%.5g\n", reshuf3, reshuf1/reshuf3, reshuf2/reshuf3)
		fmt.Fprintf(w, "fall[4] = %.10g, freq = 1/%.5g, freq2 = 1/%.5g\n", reshuf4, reshuf1/reshuf4, reshuf3/reshuf4)
		fmt.Fprintf(w, "fall[5] = %.10g, freq = 1/%.5g, freq2 = 1/%.5g\n", reshuf5, reshuf1/reshuf5, reshuf4/reshuf5)
		fmt.Fprintf(w, "farm bonuses: frequency 1/%.5g, rtp = %.6f%%\n", reshuf1/s.BonusCount(farmbn), rtpbon)
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g(farm) = %.6f%%\n", rtpsym, rtpbon, rtp)
		return rtp
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc)
}
