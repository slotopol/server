package slot

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"runtime"
	"sync"

	"go.uber.org/atomic"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/game"
)

type Simulator interface {
	// Returns number of spins, sum of pays, sum of squares of pays by spin,
	// normalized by spin cost.
	NSQ(float64) (float64, float64, float64)
	// Performs spin simulation, calculates results, update statistics.
	Simulate(SlotGame, Reelx, *Wins)
}

type ScanPar = game.ScanPar

// Maximum possible number of symbols at any game
const MaxSymNum = 20

type StatCounter struct {
	N   atomic.Uint64             // number of processed grid reshuffles, including with no wins
	S   [MaxSymNum]atomic.Float64 // sum of pays by symbols
	FSC atomic.Uint64             // free spins count
	FHC atomic.Uint64             // free games hits count
	BHC [8]atomic.Uint64          // bonus hits count
	JHC [4]atomic.Uint64          // jackpot hits count
}

func (c *StatCounter) Update(wins Wins) (pay float64) {
	for _, wi := range wins {
		if wi.Pay != 0 {
			var p = wi.Pay * wi.MP
			c.S[wi.Sym].Add(p)
			pay += p
		}
		if wi.FS != 0 {
			c.FSC.Add(uint64(wi.FS))
			c.FHC.Inc()
		}
		if wi.BID != 0 {
			c.BHC[wi.BID].Inc()
		}
		if wi.JID != 0 {
			c.JHC[wi.JID].Inc()
		}
	}
	c.N.Inc()
	return
}

func (c *StatCounter) SumS() (S float64) {
	for sym := range c.S {
		S += c.S[sym].Load()
	}
	return
}

type StatGeneric struct {
	StatCounter
	Q  atomic.Float64 // sum of squares of pays by symbols
	EC atomic.Uint64  // errors count
}

// Declare conformity with Stater interface.
var _ Simulator = (*StatGeneric)(nil)

func (s *StatGeneric) Errors() uint64 {
	return s.EC.Load()
}

func (s *StatGeneric) Count() float64 {
	return float64(s.N.Load())
}

func (s *StatGeneric) RTPsym(cost float64, scat Sym) (lrtp, srtp float64) {
	var sym Sym
	for sym = range MaxSymNum {
		if sym != scat {
			lrtp += s.S[sym].Load()
		} else {
			srtp += s.S[sym].Load()
		}
	}
	var N = s.Count()
	lrtp /= N * cost
	srtp /= N * cost
	return
}

func (s *StatGeneric) RTPsym2(cost float64, scat1, scat2 Sym) (lrtp, srtp float64) {
	var sym Sym
	for sym = range MaxSymNum {
		if sym != scat1 && sym != scat2 {
			lrtp += s.S[sym].Load()
		} else {
			srtp += s.S[sym].Load()
		}
	}
	var N = s.Count()
	lrtp /= N * cost
	srtp /= N * cost
	return
}

func (s *StatGeneric) NSQ(cost float64) (N float64, S float64, Q float64) {
	N = s.Count()
	S = s.SumS() / cost
	Q = s.Q.Load() / cost / cost
	return
}

func (s *StatGeneric) SymVariance(cost float64) float64 {
	var N, S, Q = s.NSQ(cost)
	// another way: Q/N - S*S/N/N
	return N*Q - S*S/N/N
}

func (s *StatGeneric) SymSigma(cost float64) float64 {
	var N, S, Q = s.NSQ(cost)
	// another way: math.Sqrt(Q/N - S*S/N/N)
	return math.Sqrt(N*Q-S*S) / N
}

// Returns (q, sq), where q = free spins quantifier, sq = 1/(1-q)
// sum of a decreasing geometric progression for retriggered free spins.
func (s *StatGeneric) FSQ() (q float64, sq float64) {
	q = float64(s.FSC.Load()) / s.Count()
	sq = 1 / (1 - q)
	return
}

// Quantifier of free games per reshuffles.
func (s *StatGeneric) FGQ() float64 {
	return float64(s.FHC.Load()) / s.Count()
}

// Free Games Frequency: average number of reshuffles per free games hit.
func (s *StatGeneric) FGF() float64 {
	return s.Count() / float64(s.FHC.Load())
}

func (s *StatGeneric) BonusHitsF(bid int) float64 {
	return float64(s.BHC[bid].Load())
}

func (s *StatGeneric) JackHitsF(jid int) float64 {
	return float64(s.JHC[jid].Load())
}

func (s *StatGeneric) Simulate(g SlotGame, reels Reelx, wins *Wins) {
	if g.Scanner(wins) != nil {
		s.EC.Inc()
		return
	}
	if pay := s.Update(*wins); pay != 0 {
		s.Q.Add(pay * pay)
	}
}

type StatCascade struct {
	Casc [FallLimit]StatCounter
	Q    atomic.Float64 // sum of squares of pays by symbols
	EC   atomic.Uint64  // errors count
}

// Declare conformity with Stater interface.
var _ Simulator = (*StatCascade)(nil)

func (s *StatCascade) Errors() uint64 {
	return s.EC.Load()
}

func (s *StatCascade) Count() float64 {
	return float64(s.Casc[0].N.Load())
}

func (s *StatCascade) SumFreeCount() uint64 {
	var sum uint64
	for i := range FallLimit {
		sum += s.Casc[i].FSC.Load()
	}
	return sum
}

func (s *StatCascade) SumFreeHits() uint64 {
	var sum uint64
	for i := range FallLimit {
		sum += s.Casc[i].FHC.Load()
	}
	return sum
}

func (s *StatCascade) SumBonusHits(bid int) uint64 {
	var sum uint64
	for i := range FallLimit {
		sum += s.Casc[i].BHC[bid].Load()
	}
	return sum
}

func (s *StatCascade) SumJackHits(jid int) uint64 {
	var sum uint64
	for i := range FallLimit {
		sum += s.Casc[i].JHC[jid].Load()
	}
	return sum
}

func (s *StatCascade) RTPsym(cost float64, scat Sym) (lrtp, srtp float64) {
	var sym Sym
	for i := range FallLimit {
		var c = &s.Casc[i]
		for sym = range MaxSymNum {
			if sym != scat {
				lrtp += c.S[sym].Load()
			} else {
				srtp += c.S[sym].Load()
			}
		}
	}
	var N = s.Count()
	lrtp /= N * cost
	srtp /= N * cost
	return
}

func (s *StatCascade) NSQ(cost float64) (N float64, S float64, Q float64) {
	N = s.Count()
	for i := range FallLimit {
		S += s.Casc[i].SumS()
	}
	S /= cost
	Q = s.Q.Load() / cost / cost
	return
}

// Returns (q, sq), where q = free spins quantifier, sq = 1/(1-q)
// sum of a decreasing geometric progression for retriggered free spins.
func (s *StatCascade) FSQ() (q float64, sq float64) {
	q = float64(s.SumFreeCount()) / s.Count()
	sq = 1 / (1 - q)
	return
}

// Quantifier of free games per reshuffles.
func (s *StatCascade) FGQ() float64 {
	return float64(s.SumFreeHits()) / s.Count()
}

// Free Games Frequency: average number of reshuffles per free games hit.
func (s *StatCascade) FGF() float64 {
	return s.Count() / float64(s.SumFreeHits())
}

// Cascade multiplier.
func (s *StatCascade) Mcascade() float64 {
	var pay1 = s.Casc[0].SumS()
	var pays float64
	for i := range FallLimit {
		var payi = s.Casc[i].SumS()
		pays += payi
		if payi == 0 {
			break
		}
	}
	return pays / pay1
}

// Average Cascade Length.
func (s *StatCascade) ACL() float64 {
	var sum uint64
	for i := 1; i < FallLimit; i++ {
		sum += s.Casc[i].N.Load()
	}
	return float64(sum) / float64(s.Casc[1].N.Load())
}

// Inverse coefficient of fading (C1/Cn)^(1/(n-1)).
func (s *StatCascade) Kfading() float64 {
	var N1 = s.Casc[0].N.Load()
	var Nn uint64
	var i int
	for i = range FallLimit {
		var Ni = s.Casc[i].N.Load()
		if Ni == 0 {
			break
		}
		Nn = Ni
	}
	return math.Pow(float64(N1)/float64(Nn), 1/(float64(i)-1))
}

// Maximum number of cascades in avalanche.
func (s *StatCascade) Ncascmax() int {
	for i := range FallLimit {
		if s.Casc[i].N.Load() == 0 {
			return i
		}
	}
	return FallLimit
}

func (s *StatCascade) Simulate(g SlotGame, reels Reelx, wins *Wins) {
	var sc = g.(Cascade)
	var err error
	var cfn int
	var pay float64
	for {
		sc.UntoFall()
		if cfn++; cfn > FallLimit {
			err = ErrAvalanche
			break
		}
		var wp = len(*wins)
		if err = g.Scanner(wins); err != nil {
			break
		}
		pay += s.Casc[cfn-1].Update((*wins)[wp:])
		sc.Strike((*wins)[wp:])
		if len(*wins) == wp {
			break
		}
		sc.PushFall(reels)
	}
	if pay > 0 {
		s.Q.Add(pay * pay)
	}
	if err != nil {
		s.EC.Inc()
		return
	}
}

type CalcAlg = func(ctx context.Context, sp *ScanPar, s Simulator, g SlotGeneric, reels Reelx)

const (
	CtxGranulation = 1000 // check context every N reshuffles
	FallLimit      = 20   // maximum 20 cascades for reels ~100 symbols length
)

var (
	ErrAvalanche = errors.New("too many cascading falls")
	ErrReelCount = errors.New("unexpected number of reels")
)

func CorrectThrNum(tn int) int {
	if tn < 1 {
		return runtime.GOMAXPROCS(0)
	}
	return tn
}

func ScanReels(ctx context.Context, sp *ScanPar, s Simulator, g SlotGeneric, reels Reelx,
	bruteforce, montecarlo CalcAlg,
	calc func(io.Writer) float64) float64 {
	if sx, sy := g.Dim(); len(reels) != int(sx) {
		panic(fmt.Errorf("%w: %d reels provided for %dx%d slot", ErrReelCount, len(reels), sx, sy))
	}
	fmt.Printf("selected %d lines\n", g.GetSel())
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		var ctx2, cancel2 = context.WithCancel(ctx)
		defer cancel2()
		if cfg.MCCount > 0 || sp.Prec > 0 {
			go func() {
				defer wg.Done()
				ProgressMC(ctx2, sp, s, calc, g.Cost())
			}()
			montecarlo(ctx2, sp, s, g, reels)
		} else {
			go func() {
				defer wg.Done()
				ProgressBF(ctx2, sp, s, calc, float64(reels.Reshuffles()))
			}()
			bruteforce(ctx2, sp, s, g, reels)
		}
	}()
	wg.Wait()

	fmt.Printf("reels lengths %s, total reshuffles %d\n", reels.String(), reels.Reshuffles())
	return calc(os.Stdout)
}

func ScanReelsCommon(ctx context.Context, sp *ScanPar, s Simulator, g SlotGeneric, reels Reelx,
	calc func(io.Writer) float64) float64 {
	return ScanReels(ctx, sp, s, g, reels, BruteForcex, MonteCarlo, calc)
}
