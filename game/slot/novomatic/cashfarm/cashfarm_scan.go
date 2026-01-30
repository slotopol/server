package cashfarm

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

// Bonus game expectation
// calculation see at generator/prov/novomatic/cashfarmbon.lua
const Ebon = 50

func CalcStat(ctx context.Context, mrtp float64) float64 {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	var g = NewGame()
	var s slot.StatCascade

	var calc = func(w io.Writer) float64 {
		var reshuf1 = float64(s.Casc[0].Reshuf.Load())
		var reshuf2 = float64(s.Casc[1].Reshuf.Load())
		var reshuf3 = float64(s.Casc[2].Reshuf.Load())
		var reshuf4 = float64(s.Casc[3].Reshuf.Load())
		var reshuf5 = float64(s.Casc[4].Reshuf.Load())
		var lrtp, srtp = s.SymRTP(g.Cost())
		var rtpsym = lrtp + srtp
		var qfarm = float64(s.SumBonusHits(farmbn)) / reshuf1
		var rtpbon = Ebon * qfarm * 100
		var rtp = rtpsym + rtpbon
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "fall[2] = %.10g, Ec2 = Kf2 = 1/%.5g\n", reshuf2, reshuf1/reshuf2)
		fmt.Fprintf(w, "fall[3] = %.10g, Ec3 = 1/%.5g, Kf3 = 1/%.5g\n", reshuf3, reshuf1/reshuf3, reshuf2/reshuf3)
		fmt.Fprintf(w, "fall[4] = %.10g, Ec4 = 1/%.5g, Kf4 = 1/%.5g\n", reshuf4, reshuf1/reshuf4, reshuf3/reshuf4)
		fmt.Fprintf(w, "fall[5] = %.10g, Ec5 = 1/%.5g, Kf5 = 1/%.5g\n", reshuf5, reshuf1/reshuf5, reshuf4/reshuf5)
		fmt.Fprintf(w, "Mcascade = %.5g, ACL = %.5g, Kfading = 1/%.5g, Ncascmax = %d\n", s.Mcascade(), s.ACL(), s.Kfading(), s.Ncascmax())
		fmt.Fprintf(w, "farm bonuses: hit rate 1/%.5g, rtp = %.6f%%\n", reshuf1/float64(s.SumBonusHits(farmbn)), rtpbon)
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g(farm) = %.6f%%\n", rtpsym, rtpbon, rtp)
		return rtp
	}

	return slot.ScanReelsCommon(ctx, &s, g, reels, calc)
}
