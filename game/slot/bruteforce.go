package slot

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Function to report about progress of calculation by brute force
func ReporterBF(ctx context.Context, snapch []chan Simulator2,
	sp *ScanPar, ss Simulator2, g SlotGeneric, calc CalcFunc) {
	var t0 = time.Now()
	const stepdur = 1000 * time.Millisecond
	var steps = time.Tick(stepdur)
	if sel := g.GetSel(); sel > 0 {
		fmt.Printf("selected %d lines\n", sel)
	}
	fmt.Printf("calculation started...\r")
	var total = float64(sp.Total)
	var (
		dur time.Duration
		N   float64
	)
	var latest = make([]Simulator2, len(snapch))
	var closed = make([]bool, len(snapch))
	var activeCount = len(snapch)
	for activeCount > 0 {
		select {
		case <-ctx.Done():
			activeCount = 0
		case <-steps:
		}
		for i, ch := range snapch {
			if !closed[i] {
				select {
				case snap, ok := <-ch:
					if !ok {
						closed[i] = true
						activeCount--
					} else {
						latest[i] = snap
					}
				default:
				}
			}
		}

		ss.Clear()
		for i := range latest {
			if latest[i] != nil {
				ss.Add(latest[i])
			}
		}

		dur = time.Since(t0)
		N = ss.Count()
		var RTP, _ = calc(io.Discard)

		var exp string
		if dur > 2*stepdur {
			exp = time.Duration(float64(dur-stepdur) * total / N).Truncate(time.Second).String()
		} else {
			exp = "?"
		}
		fmt.Printf("processed %.1fm/%.1fm, ready %2.2f%% (%v / %v), RTP = %2.2f%%  \r",
			N/1e6, total/1e6, N/total*100, dur.Truncate(time.Second), exp, RTP*100)
	}

	// report results
	fmt.Printf("completed %.5g%% (%d), time spent %v                    \n",
		N/total*100, int(N), dur)
	var reels = g.GetReels(sp.MRTP)
	fmt.Println()
	fmt.Printf("reels lengths %s, total reshuffles %d\n", reels.String(), reels.Reshuffles())
	calc(os.Stdout)
}

func WorkerBF(ctx context.Context, outch chan Simulator2, ti, tn uint64,
	sp *ScanPar, local Simulator2, wgame SlotGeneric) {
	var reels = wgame.GetReels(sp.MRTP)
	var total = reels.Reshuffles()
	var _, iscascade = wgame.(SlotCascade)

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
	var wins Wins
	var rpos = make([]int, rn)
	for x := range rn {
		rpos[x] = -1
	}
	var i, N uint64
	var x int
	var pos int
	// Using one general loop instead of five nested loops
	// gives ~30% performance increase.
	for i = ti; i < total; i += tn {
		N++
		if N&0xFFF == 0 {
			select {
			case <-ctx.Done():
				return
			case outch <- local.Clone():
			default:
			}
		}
		for x = range rn {
			pos = int(i/delim[x]) % rlen[x]
			if rpos[x] != pos {
				wgame.SetCol(Pos(x+1), reels[x], pos)
				rpos[x] = pos
			} else {
				break
			}
		}
		local.Simulate(wgame, &wins, sp.MRTP)
		if iscascade && len(wins) > 0 {
			for x := range rn {
				rpos[x] = -1
			}
		}
		wins.Reset()
	}
	// clear channel buffer
	select {
	case <-outch:
	default:
	}
	// push final result to channel
	outch <- local
}

func BruteForcex2(ctx context.Context,
	sp *ScanPar, ss Simulator2, mgame SlotGeneric, calc CalcFunc) {
	var reels = mgame.GetReels(sp.MRTP)
	if sp.TN%len(reels[0]) == 0 {
		panic("BruteForcex: thread number equals to 1-st reel length")
	}

	var snapch = make([]chan Simulator2, sp.TN)
	for i := range snapch {
		snapch[i] = make(chan Simulator2, 1)
	}

	var wg1, wg2 sync.WaitGroup
	func() {
		if calc != nil {
			var ctxrep, cancel = context.WithCancel(ctx)
			defer cancel()
			wg2.Go(func() {
				ReporterBF(ctxrep, snapch, sp, ss, mgame, calc)
			})
		}
		for ti, ch := range snapch {
			var local = ss.Clone()
			var wgame = mgame.Clone()
			wg1.Go(func() {
				defer close(ch)
				WorkerBF(ctx, ch, uint64(ti), uint64(len(snapch)), sp, local, wgame)
			})
		}
		wg1.Wait()
	}()
	wg2.Wait()
}

// Function to report about progress of calculation by brute force
func ProgressBF(ctx context.Context, sp *ScanPar, s Simulator, calc func(io.Writer) (float64, float64)) {
	const stepdur = 1000 * time.Millisecond
	var t0 = time.Now()
	var steps = time.Tick(stepdur)
	fmt.Printf("calculation started...\r")
	var total = float64(sp.Total)
	var (
		dur time.Duration
		N   float64
		RTP float64
	)
	var param = func() {
		dur = time.Since(t0)
		N = s.Count()
		RTP, _ = calc(io.Discard)
	}
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-steps:
			param()
			var exp = time.Duration(float64(dur) * total / N)
			fmt.Printf("processed %.1fm/%.1fm, ready %2.2f%% (%v / %v), RTP = %2.2f%%  \r",
				N/1e6, total/1e6, N/total*100,
				dur.Truncate(stepdur), exp.Truncate(stepdur),
				RTP*100)
		}
	}

	// report results
	param()
	fmt.Printf("completed %.5g%% (%d), time spent %v                    \n",
		N/total*100, int(N), dur)
}

func BruteForcex(ctx context.Context, sp *ScanPar, s Simulator, g SlotGeneric) {
	var reels = g.GetReels(sp.MRTP)
	var total = reels.Reshuffles()
	var tn64 = uint64(sp.TN)
	if sp.TN%len(reels[0]) == 0 {
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
	wg.Add(sp.TN)
	for ti := range tn64 {
		var gt = g.Clone()
		var _, iscascade = gt.(SlotCascade)
		go func() {
			defer wg.Done()
			var wins Wins
			var rpos = make([]int, rn)
			for x := range rn {
				rpos[x] = -1
			}
			var i, N uint64
			var x int
			var pos int
			// Using one general loop instead of five nested loops
			// gives ~30% performance increase.
			for i = ti; i < total; i += tn64 {
				N++
				if N&0xFFF == 0 {
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
				s.Simulate(gt, &wins, sp.MRTP)
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

func BruteForce5x3Big(ctx context.Context, sp *ScanPar, s Simulator, g SlotGeneric, r1, rb, r5 []Sym) {
	// var total = uint64(len(r1)) * uint64(len(rb)) * uint64(len(r5))
	var tn64 = uint64(sp.TN)
	var wg sync.WaitGroup
	wg.Add(sp.TN)
	for ti := range tn64 {
		var gt = g.Clone()
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
						if N&0xFFF == 0 {
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
						s.Simulate(gt, &wins, sp.MRTP)
						wins.Reset()
					}
				}
			}
		}()
	}
	wg.Wait()
}
