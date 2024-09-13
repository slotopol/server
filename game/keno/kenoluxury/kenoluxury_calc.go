package kenoluxury

import (
	"context"
	"fmt"

	keno "github.com/slotopol/server/game/keno"
)

func CalcStat(ctx context.Context) float64 {
	var rtp float64
	for n := 2; n <= 10; n++ {
		var nrtp float64
		for r := 0; r <= n; r++ {
			var pay = Paytable[n][r]
			nrtp += pay * keno.Prob(n, r)
		}
		fmt.Printf("RTP[%2d] = %.6f%%\n", n, nrtp*100)
		rtp += nrtp
	}
	rtp *= 100e0 / 9e0
	fmt.Printf("RTP[game] = %.6f%%", rtp)
	return rtp
}
