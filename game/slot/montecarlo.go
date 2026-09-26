package slot

import (
	"context"
	"fmt"
	"io"
	"math"
	"sync"
	"time"

	"github.com/slotopol/server/game"
)

const lolim = 1e6 // lower limit, let some space to get approximate sigma

// Function to report about progress of Monte Carlo calculation
func ProgressMC(ctx context.Context, sp *ScanPar, s Simulator, calc func(io.Writer) (float64, float64), cost float64) {
	const stepdur = 1000 * time.Millisecond
	var t0 = time.Now()
	var steps = time.Tick(stepdur)
	fmt.Printf("calculation started...\r")
	var (
		dur    time.Duration
		N      float64
		RTP, D float64
		VI     float64
		ΔRTP   float64
		total  float64
	)
	var param = func() {
		dur = time.Since(t0)
		if RTP, D = calc(io.Discard); math.IsNaN(D) {
			_, D = s.EvD(cost)
		}
		VI = game.GetZ(sp.Conf) * math.Sqrt(D)
		N = s.Count()
		ΔRTP = VI / math.Sqrt(N)
		var tc, tp float64
		tc = max(float64(sp.Total), lolim)
		if sp.Prec > 0 {
			var t2 = VI / sp.Prec
			tp = t2 * t2
		}
		total = max(tc, tp)
	}
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-steps:
			param()
			var exp = time.Duration(float64(dur) * total / N)
			fmt.Printf("processed %.1fm/%.1fm, ready %2.2f%% (%v / %v), RTP = %2.2f%%, Δ[%.4g%%] = %.4g%%  \r",
				N/1e6, total/1e6, N/total*100,
				dur.Truncate(stepdur), exp.Truncate(stepdur),
				RTP*100, sp.Conf*100, ΔRTP*100)
		}
	}

	// report results
	param()
	fmt.Printf("completed %.5g%% (%d), time spent %v, precision Δ[%.4g%%] = %.4g%%                \n",
		N/total*100, int(N), dur, sp.Conf*100, ΔRTP*100)
}

func MonteCarlo(ctx context.Context, sp *ScanPar, s Simulator, g SlotGeneric) {
	var reels = g.GetReels(sp.MRTP)
	var tn64 = uint64(sp.TN)
	var total = max(sp.Total, lolim)
	var wg sync.WaitGroup
	wg.Add(sp.TN)
	for range tn64 {
		var gt = g.Clone()
		go func() {
			defer wg.Done()

			var wins Wins
			var i uint64
			for i = 0; i < total; i += tn64 {
				if i%CtxGranulation == 0 {
					// check on break
					select {
					case <-ctx.Done():
						return
					default:
					}
					// recalculate total iterations number
					var tc, tp uint64
					tc = max(sp.Total, lolim)
					if sp.Prec > 0 {
						var _, D = s.EvD(gt.Cost())
						var VI = game.GetZ(sp.Conf) * math.Sqrt(D)
						var t2 = VI / sp.Prec
						tp = uint64(t2 * t2)
					}
					total = max(tc, tp)
				}
				gt.SpinReels(reels)
				s.Simulate(gt, &wins, sp.MRTP)
				wins.Reset()
			}
		}()
	}
	wg.Wait()
}
