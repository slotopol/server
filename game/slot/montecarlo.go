package slot

import (
	"context"
	"fmt"
	"io"
	"math"
	"os"
	"sync"
	"time"

	"github.com/slotopol/server/game"
	"github.com/slotopol/server/prng"
)

const lolim = 1e6 // lower limit, let some space to get approximate sigma

// Function to report about progress of Monte Carlo calculation
func ReporterMC(ctx context.Context, snapch []chan Simulator2,
	sp *ScanPar, ss Simulator2, g SlotGeneric, calc CalcFunc) {
	const stepdur = 1000 * time.Millisecond
	var t0 = time.Now()
	var steps = time.Tick(stepdur)
	if sel := g.GetSel(); sel > 0 {
		fmt.Printf("selected %d lines\n", sel)
	}
	fmt.Printf("calculation started...\r")
	var (
		dur   time.Duration
		N     float64
		ΔRTP  float64
		total float64
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
		var RTP, D float64
		if RTP, D = calc(io.Discard); math.IsNaN(D) {
			_, D = EvD(ss, g.Cost())
		}
		var VI = game.GetZ(sp.Conf) * math.Sqrt(D)
		N = ss.Count()
		ΔRTP = VI / math.Sqrt(N)
		var tc, tp float64
		tc = max(float64(sp.Total), lolim)
		if sp.Prec > 0 {
			var t2 = VI / sp.Prec
			tp = t2 * t2
		}
		total = max(tc, tp)

		var exp string
		if dur > 2*stepdur {
			exp = time.Duration(float64(dur-stepdur) * total / N).Truncate(time.Second).String()
		} else {
			exp = "?"
		}
		fmt.Printf("processed %.1fm/%.1fm, ready %2.2f%% (%v / %v), RTP = %2.2f%%, Δ[%.4g%%] = %.4g%%  \r",
			N/1e6, total/1e6, N/total*100,
			dur.Truncate(time.Second), exp,
			RTP*100, sp.Conf*100, ΔRTP*100)
	}

	// report results
	fmt.Printf("completed %.5g%% (%d), time spent %v, precision Δ[%.4g%%] = %.4g%%                \n",
		N/total*100, int(N), dur, sp.Conf*100, ΔRTP*100)
	var reels = g.GetReels(sp.MRTP)
	fmt.Println()
	fmt.Printf("reels lengths %s, total reshuffles %d\n", reels.String(), reels.Reshuffles())
	calc(os.Stdout)
}

func WorkerMC(ctx context.Context, outch chan Simulator2, rng *prng.Rand, tn uint64,
	sp *ScanPar, local Simulator2, wgame SlotGeneric) {
	var reels = wgame.GetReels(sp.MRTP)
	var total = max(sp.Total, lolim)

	var wins Wins
	var i, N uint64
	for i = 0; i < total; i += tn {
		N++
		if N&0xFFF == 0 {
			// check on break
			select {
			case <-ctx.Done():
				return
			case outch <- local.Clone():
			default:
			}
			// recalculate total iterations number
			var tc, tp uint64
			tc = max(sp.Total, lolim)
			if sp.Prec > 0 {
				var _, D = EvD(local, wgame.Cost())
				var VI = game.GetZ(sp.Conf) * math.Sqrt(D)
				var t2 = VI / sp.Prec
				tp = uint64(t2 * t2)
			}
			total = max(tc, tp)
		}
		wgame.SpinReels(rng, reels)
		local.Simulate(wgame, &wins, sp.MRTP)
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

func MonteCarlo2(ctx context.Context, sp *ScanPar, ss Simulator2, mgame SlotGeneric, calc CalcFunc) {
	var snapch = make([]chan Simulator2, sp.TN)
	for i := range snapch {
		snapch[i] = make(chan Simulator2, 1)
	}

	var msrc = prng.NewXoshiro256PP(123456789) // master prng source

	var wg1, wg2 sync.WaitGroup
	func() {
		if calc != nil {
			var ctxrep, cancel = context.WithCancel(ctx)
			defer cancel()
			wg2.Go(func() {
				ReporterMC(ctxrep, snapch, sp, ss, mgame, calc)
			})
		}
		for _, ch := range snapch {
			var local = ss.Clone()
			var wgame = mgame.Clone() // worker game
			var wsrc = msrc.Clone()   // worker prng source
			msrc.Jump()               // shift the master forward by 2^128 steps
			var rng = prng.New(wsrc)

			wg1.Go(func() {
				defer close(ch)
				WorkerMC(ctx, ch, rng, uint64(len(snapch)), sp, local, wgame)
			})
		}
		wg1.Wait()
	}()
	wg2.Wait()
}

// Function to report about progress of Monte Carlo calculation
func ProgressMC(ctx context.Context, sp *ScanPar, s Simulator, calc CalcFunc, cost float64) {
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
			_, D = EvD(s, cost)
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

func MonteCarlo(ctx context.Context, sp *ScanPar, s Simulator, mgame SlotGeneric) {
	var reels = mgame.GetReels(sp.MRTP)
	var tn64 = uint64(sp.TN)
	var total = max(sp.Total, lolim)
	var wg sync.WaitGroup

	var msrc = prng.NewXoshiro256PP(123456789) // master prng source

	for range tn64 {
		var wgame = mgame.Clone() // worker game
		var wsrc = msrc.Clone()   // worker prng source
		msrc.Jump()               // shift the master forward by 2^128 steps
		var rng = prng.New(wsrc)

		wg.Go(func() {
			var wins Wins
			var N uint64
			for N = 0; N < total; N += tn64 {
				if N&0xFFF == 0 {
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
						var _, D = EvD(s, wgame.Cost())
						var VI = game.GetZ(sp.Conf) * math.Sqrt(D)
						var t2 = VI / sp.Prec
						tp = uint64(t2 * t2)
					}
					total = max(tc, tp)
				}
				wgame.SpinReels(rng, reels)
				s.Simulate(wgame, &wins, sp.MRTP)
				wins.Reset()
			}
		})
	}
	wg.Wait()
}
