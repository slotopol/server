package slot

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"
)

// Function to report about progress of calculation by brute force
func ProgressBF(ctx context.Context, s Simulator, calc func(io.Writer) float64, total float64) {
	const stepdur = 1000 * time.Millisecond
	var t0 = time.Now()
	var steps = time.Tick(stepdur)
	fmt.Printf("calculation started...\r")
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-steps:
			var dur = time.Since(t0)
			var N, _, _ = s.NSQ(1)
			var RTP = calc(io.Discard)
			var exp = time.Duration(float64(dur) * total / N)
			fmt.Printf("processed %.1fm/%.1fm, ready %2.2f%% (%v / %v), RTP = %2.2f%%  \r",
				N/1e6, total/1e6, N/total*100,
				dur.Truncate(stepdur), exp.Truncate(stepdur),
				RTP*100)
		}
	}

	// report results
	var dur = time.Since(t0)
	var N, _, _ = s.NSQ(1)
	fmt.Printf("completed %.5g%% (%d), time spent %v                    \n",
		N/total*100, int(N), dur)
}

func BruteForcex(ctx context.Context, s Simulator, g SlotGeneric, reels Reelx) {
	var total = reels.Reshuffles()
	var tn = CorrectThrNum()
	var tn64 = uint64(tn)
	if tn%len(reels[0]) == 0 {
		panic("BruteForcex: thread number equals to 1-st reel length")
	}
	// Number of reels.
	var rn = len(reels)
	// Precompute reel lengths.
	var rlen = make([]int, rn)
	for x := range rn {
		rlen[x] = len(reels[x])
	}
	// Precompute delimiters for position calculation.
	var delim = make([]uint64, rn)
	delim[0] = 1
	for x := 1; x < rn; x++ {
		delim[x] = delim[x-1] * uint64(rlen[x-1])
	}

	var wg sync.WaitGroup
	wg.Add(tn)
	for ti := range tn64 {
		var gt = g.Clone().(SlotGeneric)
		var _, iscascade = gt.(SlotCascade)
		go func() {
			defer wg.Done()
			var wins Wins
			var rpos = make([]int, rn)
			for x := range rn {
				rpos[x] = -1
			}
			var i uint64
			var x int
			var pos int
			// Using one general loop instead of five nested loops
			// gives ~30% performance increase.
			for i = ti; i < total; i += tn64 {
				if (i/tn64)%CtxGranulation == 0 {
					select {
					case <-ctx.Done():
						return
					default:
					}
				}
				for x = range rn {
					pos = int(i/delim[x]) % rlen[x]
					if rpos[x] != pos {
						gt.SetCol(Pos(x+1), reels[x], pos)
						rpos[x] = pos
					} else {
						break
					}
				}
				s.Simulate(gt, reels, &wins)
				if iscascade && len(wins) > 0 {
					for x := range rn {
						rpos[x] = -1
					}
				}
				wins.Reset()
			}
		}()
	}
	wg.Wait()
}

func BruteForce5x3Big(ctx context.Context, s Simulator, g SlotGame, r1, rb, r5 []Sym) {
	// var total = uint64(len(r1)) * uint64(len(rb)) * uint64(len(r5))
	var tn = CorrectThrNum()
	var tn64 = uint64(tn)
	var wg sync.WaitGroup
	wg.Add(tn)
	for ti := range tn64 {
		var gt = g.Clone().(SlotGeneric)
		var cb = gt.(Bigger)
		var N uint64
		go func() {
			defer wg.Done()

			var wins Wins
			for i1 := range r1 {
				gt.SetCol(1, r1, i1)
				for _, big := range rb {
					cb.SetBig(big)
					for i5 := range r5 {
						N++
						if N%CtxGranulation == 0 {
							select {
							case <-ctx.Done():
								return
							default:
							}
						}
						if N%tn64 != ti {
							continue
						}
						gt.SetCol(5, r5, i5)
						s.Simulate(gt, nil, &wins)
						wins.Reset()
					}
				}
			}
		}()
	}
	wg.Wait()
}
