package slot

import (
	"context"
	"fmt"
	"io"
	"math"
	"sync"
	"time"

	cfg "github.com/slotopol/server/config"
)

// Function to report about progress of Monte Carlo calculation
func ProgressMC(ctx context.Context, s Simulator, calc func(io.Writer) float64, cost float64) {
	var confidence = 0.95
	var totalmin float64 = 1e6
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
		VI = GetZ(confidence) * math.Sqrt(N*Q-S*S) / N
		ΔRTP = VI / math.Sqrt(N)
		if cfg.MCCount > 0 {
			total = float64(cfg.MCCount) * 1e6
		} else {
			var t2 = VI / cfg.MCPrec * 100
			total = max(t2*t2, totalmin)
		}
	}
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-steps:
			param()
			var exp = time.Duration(float64(dur) * total / N)
			fmt.Printf("processed %.1fm/%.1fm, ready %2.2f%% (%v / %v), RTP = %2.2f%%, Δ[%.2g%%] = %2.2f%%  \r",
				N/1e6, total/1e6, N/total*100,
				dur.Truncate(stepdur), exp.Truncate(stepdur),
				RTP*100, confidence*100, ΔRTP*100)
		}
	}

	// report results
	param()
	fmt.Printf("completed %.5g%% (%d), time spent %v, precision Δ[%.2g%%] = %.4g%%                \n",
		N/total*100, int(N), dur, confidence*100, ΔRTP*100)
}

func MonteCarlo(ctx context.Context, s Simulator, g SlotGeneric, reels Reelx) {
	var confidence = 0.95
	var totalmin uint64 = 1e6 // let some space to get approximate sigma
	var tn = CorrectThrNum()
	var tn64 = uint64(tn)
	var total uint64
	if cfg.MCCount > 0 {
		total = cfg.MCCount * 1e6
	} else {
		total = totalmin
	}
	var wg sync.WaitGroup
	wg.Add(tn)
	for range tn64 {
		var gt = g.Clone().(SlotGeneric)
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
					if cfg.MCPrec > 0 {
						var N, S, Q = s.NSQ(gt.Cost())
						var VI = GetZ(confidence) * math.Sqrt(N*Q-S*S) / N
						var t2 = VI / cfg.MCPrec * 100
						total = max(uint64(t2*t2), totalmin)
					}
				}
				gt.SpinReels(reels)
				s.Simulate(gt, reels, &wins)
				wins.Reset()
			}
		}()
	}
	wg.Wait()
}
