package winstorm

import (
	"context"
	"fmt"
	"io"
	"math"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) float64 {
	var reels, _ = ReelsMap.FindClosest(sp.MRTP)
	var g = NewGame()
	var s = slot.NewStatCascade(sn, 5)

	var calc = func(w io.Writer) float64 {
		var N, S, Q = s.NSQ(g.Cost())
		var N2 = float64(s.Casc[1].N.Load())
		var N3 = float64(s.Casc[2].N.Load())
		var N4 = float64(s.Casc[3].N.Load())
		var N5 = float64(s.Casc[4].N.Load())
		var rtp = S / N
		var sigma = math.Sqrt(Q/N - rtp*rtp)
		var vi = slot.GetZ(sp.Conf) * sigma
		fmt.Fprintf(w, "fall[2] = %.10g, Ec2 = Kf2 = 1/%.5g\n", N2, N/N2)
		fmt.Fprintf(w, "fall[3] = %.10g, Ec3 = 1/%.5g, Kf3 = 1/%.5g\n", N3, N/N3, N2/N3)
		fmt.Fprintf(w, "fall[4] = %.10g, Ec4 = 1/%.5g, Kf4 = 1/%.5g\n", N4, N/N4, N3/N4)
		fmt.Fprintf(w, "fall[5] = %.10g, Ec5 = 1/%.5g, Kf5 = 1/%.5g\n", N5, N/N5, N4/N5)
		fmt.Fprintf(w, "Mcascade = %.5g, ACL = %.5g, Kfading = 1/%.5g, Ncascmax = %d\n", s.Mcascade(), s.ACL(), s.Kfading(), s.Ncascmax())
		fmt.Fprintf(w, "RTP = %.6f%%\n", rtp*100)
		fmt.Fprintf(w, "sigma = %.6g, VI[%.4g%%] = %.6g (%s)\n", sigma, sp.Conf*100, vi, slot.VIname6[slot.VIclass6(vi)])
		fmt.Fprintf(w, "CI[90%%] = %d, CI[95%%] = %d, CI[99%%] = %d\n",
			int(slot.CI(0.90, rtp, sigma)), int(slot.CI(0.95, rtp, sigma)), int(slot.CI(0.99, rtp, sigma)))
		return rtp
	}

	return slot.ScanReelsCommon(ctx, sp, s, g, reels, calc)
}
