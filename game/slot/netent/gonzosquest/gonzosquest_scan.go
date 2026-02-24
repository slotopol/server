package gonzosquest

import (
	"context"
	"fmt"
	"io"

	"github.com/slotopol/server/game/slot"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) float64 {
	var reels, _ = ReelsMap.FindClosest(sp.MRTP)
	var g = NewGame()
	var s = slot.NewStatCascade(sn, 5)

	var calc = func(w io.Writer) float64 {
		var N, S, _ = s.NSQ(g.Cost())
		var N2 = float64(s.Casc[1].N.Load())
		var N3 = float64(s.Casc[2].N.Load())
		var N4 = float64(s.Casc[3].N.Load())
		var N5 = float64(s.Casc[4].N.Load())
		var µ = S / N
		var q, sq = s.FSQ()
		var rtpfs = 3 * sq * µ
		var rtp = µ + q*rtpfs
		fmt.Fprintf(w, "symbols: rtp(sym) = %.6f%%\n", µ*100)
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.SumFSC(), q, sq)
		fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", s.FGF())
		fmt.Fprintf(w, "fall[2] = %.10g, Ec2 = Kf2 = 1/%.5g\n", N2, N/N2)
		fmt.Fprintf(w, "fall[3] = %.10g, Ec3 = 1/%.5g, Kf3 = 1/%.5g\n", N3, N/N3, N2/N3)
		fmt.Fprintf(w, "fall[4] = %.10g, Ec4 = 1/%.5g, Kf4 = 1/%.5g\n", N4, N/N4, N3/N4)
		fmt.Fprintf(w, "fall[5] = %.10g, Ec5 = 1/%.5g, Kf5 = 1/%.5g\n", N5, N/N5, N4/N5)
		fmt.Fprintf(w, "Mcascade = %.5g, ACL = %.5g, Kfading = 1/%.5g, Ncascmax = %d\n", s.Mcascade(), s.ACL(), s.Kfading(), s.Ncascmax())
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%\n", µ*100, q, rtpfs*100, rtp*100)
		return rtp
	}

	return slot.ScanReelsCommon(ctx, sp, s, g, reels, calc)
}
