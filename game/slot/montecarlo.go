package slot

import (
	"context"
	"fmt"
	"io"
	"math"
	"sync"
	"time"
)

const lolim = 1e6 // lower limit, let some space to get approximate sigma

// Function to report about progress of Monte Carlo calculation
func ProgressMC(ctx context.Context, sp *ScanPar, s Simulator, calc func(io.Writer) float64, cost float64) {
	const stepdur = 1000 * time.Millisecond
	var t0 = time.Now()
	var steps = time.Tick(stepdur)
	fmt.Printf("calculation started...\r")
	var (
		dur     time.Duration
		N, S, Q float64
		RTP     float64
		VI      float64
		ΔRTP    float64
		total   float64
	)
	var param = func() {
		dur = time.Since(t0)
		N, S, Q = s.NSQ(cost)
		RTP = calc(io.Discard)
		VI = GetZ(sp.Conf) * math.Sqrt(N*Q-S*S) / N
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

func MonteCarlo(ctx context.Context, sp *ScanPar, s Simulator, g SlotGeneric, reels Reelx) {
	var tn = CorrectThrNum(sp.TN)
	var tn64 = uint64(tn)
	var total = max(sp.Total, lolim)
	var wg sync.WaitGroup
	wg.Add(tn)
	for range tn64 {
		var gt = g.Clone()
		go func() {
			defer wg.Done()

			var wins Wins
			var N uint64
			for N = 0; N < total/tn64; N++ {
				if N%CtxGranulation == 0 {
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
						var N, S, Q = s.NSQ(gt.Cost())
						var VI = GetZ(sp.Conf) * math.Sqrt(N*Q-S*S) / N
						var t2 = VI / sp.Prec
						tp = uint64(t2 * t2)
					}
					total = max(tc, tp)
				}
				gt.SpinReels(reels)
				s.Simulate(gt, reels, &wins)
				wins.Reset()
			}
		}()
	}
	wg.Wait()
}
