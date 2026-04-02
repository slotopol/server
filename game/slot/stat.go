package slot

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"sync"

	"github.com/slotopol/server/game"

	"go.uber.org/atomic"
	"gopkg.in/yaml.v3"
)

type Simulator interface {
	// Returns number of spins, sum of pays, sum of squares of pays by spin,
	// normalized by spin cost.
	NSQ(float64) (float64, float64, float64)
	// Performs spin simulation, calculates results, update statistics.
	Simulate(SlotGame, Reelx, *Wins)
}

type Counter interface {
	Simulator
	// Returns free spins quantifier.
	FSQ() float64
	// Quantifier of free games per reshuffles.
	FGQ() float64
	// The sum of the weights of free spins series occurrences
	// used for dispersion calculation.
	ΣPL(Sym, []int) float64
}

// Returns plain math expectation (µ = S/N) and plain dispersion (D = Q/N - µ*µ).
func EvD(s Counter, cost float64) (float64, float64) {
	var N, S, Q = s.NSQ(cost)
	var µ = S / N
	var D = Q/N - µ*µ
	return µ, D
}

type ScanPar = game.ScanPar

type Uint64 struct {
	atomic.Uint64
}

func (x Uint64) String() string {
	return strconv.FormatUint(x.Load(), 10)
}

func (x Uint64) IsZero() bool {
	return x.Load() == 0
}

// MarshalXML encodes the wrapped uint64 into XML.
func (x *Uint64) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(x.Load(), start)
}

// UnmarshalXML decodes a uint64 from XML.
func (x *Uint64) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v uint64
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	x.Store(v)
	return nil
}

// MarshalYAML encodes the wrapped uint64 into YAML.
func (x Uint64) MarshalYAML() (any, error) {
	return x.Load(), nil
}

// UnmarshalYAML decodes a uint64 from YAML.
func (x *Uint64) UnmarshalYAML(value *yaml.Node) error {
	var v uint64
	if err := value.Decode(&v); err != nil {
		return err
	}
	x.Store(v)
	return nil
}

type Float64 struct {
	atomic.Float64
}

func (x Float64) String() string {
	return strconv.FormatFloat(x.Load(), 'f', -1, 64)
}

func (x Float64) IsZero() bool {
	return x.Load() == 0
}

// MarshalXML encodes the wrapped float64 into XML.
func (x *Float64) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(x.String(), start)
}

// UnmarshalXML decodes a float64 from XML.
func (x *Float64) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v float64
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	x.Store(v)
	return nil
}

// MarshalYAML encodes the wrapped float64 into YAML.
func (x Float64) MarshalYAML() (any, error) {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: x.String(),
	}, nil
}

// UnmarshalYAML decodes a float64 from YAML.
func (x *Float64) UnmarshalYAML(value *yaml.Node) error {
	var v float64
	if err := value.Decode(&v); err != nil {
		return err
	}
	x.Store(v)
	return nil
}

type Counts[T fmt.Stringer] [][]T

func (c Counts[T]) MarshalYAML() (any, error) {
	var node1 = &yaml.Node{
		Kind:  yaml.SequenceNode,
		Style: yaml.FoldedStyle,
	}
	for sym, line := range c {
		var node2 = &yaml.Node{
			Kind:  yaml.SequenceNode,
			Style: yaml.FlowStyle,
		}
		node2.LineComment = strconv.Itoa(sym + 1)

		for _, v := range line {
			node2.Content = append(node2.Content, &yaml.Node{
				Kind:  yaml.ScalarNode,
				Value: v.String(),
			})
		}

		node1.Content = append(node1.Content, node2)
	}
	return node1, nil
}

type SlotCounter struct {
	N   Uint64          // number of processed grid reshuffles, including with no wins
	C   Counts[Uint64]  // hit counter
	S   Counts[Float64] // sum of pays by symbols
	FSC Uint64          `yaml:",omitempty"`      // free spins count
	FGH Uint64          `yaml:",omitempty"`      // free games hits count
	BH  []Uint64        `yaml:",flow,omitempty"` // bonus hits count
	BS  []Float64       `yaml:",flow,omitempty"` // sum of bonuses pays
	JH  []Uint64        `yaml:",flow,omitempty"` // progressive jackpot hits count
}

func (c *SlotCounter) IsZero() bool {
	return c.N.Load() == 0 // N can not be zero for non-empty structure
}

func (c *SlotCounter) SymDim(sym Sym, pn int) {
	c.C[sym-1] = make([]Uint64, pn)
	c.S[sym-1] = make([]Float64, pn)
}

func (c *SlotCounter) CntDim(sn, pn int) {
	c.C = make([][]Uint64, sn)
	c.S = make([][]Float64, sn)
	for sym := range c.S {
		c.SymDim(Sym(sym+1), pn)
	}
}

func (c *SlotCounter) BonDim(n int) {
	c.BH = make([]Uint64, n)
	c.BS = make([]Float64, n)
}

func (c *SlotCounter) JackDim(n int) {
	c.JH = make([]Uint64, n)
}

func (c *SlotCounter) Update(wins Wins) (pay float64) {
	for _, wi := range wins {
		if wi.Sym > 0 && wi.Num > 0 {
			c.C[wi.Sym-1][wi.Num-1].Inc()
		}
		if wi.Pay != 0 {
			var p = wi.Pay * wi.MP
			if wi.Sym > 0 { // symbols pay
				c.S[wi.Sym-1][wi.Num-1].Add(p)
			} else { // bonus pay, only 2 cases can be
				c.BS[wi.BID-1].Add(p)
			}
			pay += p
		}
		if wi.FS != 0 {
			c.FSC.Add(uint64(wi.FS))
			c.FGH.Inc()
		}
		if wi.BID != 0 {
			c.BH[wi.BID-1].Inc()
		}
		if wi.JID != 0 {
			c.JH[wi.JID-1].Inc()
		}
	}
	c.N.Inc()
	return
}

func (c *SlotCounter) SymPays(sym Sym) (sum float64) {
	var pays = c.S[sym-1]
	for i := range pays {
		sum += pays[i].Load()
	}
	return
}

func (c *SlotCounter) SumPays() (sum float64) {
	for _, pays := range c.S {
		for i := range pays {
			sum += pays[i].Load()
		}
	}
	for i := range c.BS {
		sum += c.BS[i].Load()
	}
	return
}

type StatGeneric struct {
	SlotCounter `yaml:",inline"`
	Q           Float64 // sum of squares of pays by symbols
	EC          Uint64  `yaml:",omitempty"` // errors count
}

// Declare conformity with Counter interface.
var _ Counter = (*StatGeneric)(nil)

func NewStatGeneric(sn, pn int) *StatGeneric {
	var s StatGeneric
	s.CntDim(sn, pn)
	return &s
}

func (s *StatGeneric) Count() float64 {
	return float64(s.N.Load())
}

func (s *StatGeneric) RTPsym(cost float64, scat Sym) (lrtp, srtp float64) {
	for sym, pays := range s.S {
		for i := range pays {
			if Sym(sym) != scat {
				lrtp += pays[i].Load()
			} else {
				srtp += pays[i].Load()
			}
		}
	}
	var N = s.Count()
	lrtp /= N * cost
	srtp /= N * cost
	return
}

func (s *StatGeneric) NSQ(cost float64) (N float64, S float64, Q float64) {
	N = s.Count()
	S = s.SumPays() / cost
	Q = s.Q.Load() / cost / cost
	return
}

// Returns free spins quantifier.
func (s *StatGeneric) FSQ() float64 {
	return float64(s.FSC.Load()) / s.Count()
}

// Quantifier of free games per reshuffles.
func (s *StatGeneric) FGQ() float64 {
	return float64(s.FGH.Load()) / s.Count()
}

// The sum of the weights of free spins series occurrences
// used for dispersion calculation.
func (s *StatGeneric) ΣPL(scat Sym, L []int) (sum float64) {
	var N = s.Count()
	for i, Li := range L {
		var Pfgi = float64(s.C[scat-1][i].Load()) / N
		sum += Pfgi * float64(Li)
	}
	return
}

// Free Games hit rate: average number of reshuffles per free games hits.
func (s *StatGeneric) HRfg() float64 {
	return s.Count() / float64(s.FGH.Load())
}

func (s *StatGeneric) BonusHits(bid int) float64 {
	return float64(s.BH[bid-1].Load())
}

func (s *StatGeneric) JackHits(jid int) float64 {
	return float64(s.JH[jid-1].Load())
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
	Casc [FallLimit]SlotCounter
	Q    Float64 // sum of squares of pays by symbols
	EC   Uint64  `yaml:",omitempty"` // errors count
}

// Declare conformity with Counter interface.
var _ Counter = (*StatCascade)(nil)

func NewStatCascade(sn, pn int) *StatCascade {
	var s StatCascade
	s.CntDim(sn, pn)
	return &s
}

func (s *StatCascade) MarshalYAML() (any, error) {
	type stat struct {
		Casc []SlotCounter
		Q    float64
		EC   uint64 `yaml:",omitempty"`
	}
	return &stat{
		Casc: s.Casc[:s.Ncascmax()],
		Q:    s.Q.Load(),
		EC:   s.EC.Load(),
	}, nil
}

func (s *StatCascade) SymDim(sym Sym, pn int) {
	for cfn := range s.Casc {
		s.Casc[cfn].SymDim(sym, pn)
	}
}

func (s *StatCascade) CntDim(sn, pn int) {
	for cfn := range s.Casc {
		s.Casc[cfn].CntDim(sn, pn)
	}
}

func (s *StatCascade) BonDim(n int) {
	for cfn := range s.Casc {
		s.Casc[cfn].BonDim(n)
	}
}

func (s *StatCascade) JackDim(n int) {
	for cfn := range s.Casc {
		s.Casc[cfn].JackDim(n)
	}
}

func (s *StatCascade) Count() float64 {
	return float64(s.Casc[0].N.Load())
}

func (s *StatCascade) SymPays(sym Sym) (sum float64) {
	for cfn := range s.Casc {
		sum += s.Casc[cfn].SymPays(sym)
	}
	return
}

func (s *StatCascade) SumPays() (sum float64) {
	for cfn := range s.Casc {
		sum += s.Casc[cfn].SumPays()
	}
	return
}

func (s *StatCascade) SumFSC() (sum uint64) {
	for cfn := range s.Casc {
		sum += s.Casc[cfn].FSC.Load()
	}
	return
}

func (s *StatCascade) SumFGH() (sum uint64) {
	for cfn := range s.Casc {
		sum += s.Casc[cfn].FGH.Load()
	}
	return
}

func (s *StatCascade) SumBH(bid int) (sum uint64) {
	for cfn := range s.Casc {
		sum += s.Casc[cfn].BH[bid-1].Load()
	}
	return
}

func (s *StatCascade) SumJH(jid int) (sum uint64) {
	for cfn := range s.Casc {
		sum += s.Casc[cfn].JH[jid-1].Load()
	}
	return
}

func (s *StatCascade) RTPsym(cost float64, scat Sym) (lrtp, srtp float64) {
	for cfn := range s.Casc {
		var c = &s.Casc[cfn]
		for sym, pays := range c.S {
			for i := range pays {
				if Sym(sym) != scat {
					lrtp += pays[i].Load()
				} else {
					srtp += pays[i].Load()
				}
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
	S = s.SumPays() / cost
	Q = s.Q.Load() / cost / cost
	return
}

// Returns free spins quantifier.
func (s *StatCascade) FSQ() float64 {
	return float64(s.SumFSC()) / s.Count()
}

// Quantifier of free games per reshuffles.
func (s *StatCascade) FGQ() float64 {
	return float64(s.SumFGH()) / s.Count()
}

// The sum of the weights of free spins series occurrences
// used for dispersion calculation.
func (s *StatCascade) ΣPL(scat Sym, L []int) (sum float64) {
	var N = s.Count()
	for i, Li := range L {
		var c float64
		for cfn := range s.Casc {
			c += float64(s.Casc[cfn].C[scat-1][i].Load())
		}
		var Pfgi = c / N
		sum += Pfgi * float64(Li)
	}
	return
}

// Free Games hit rate: average number of reshuffles per free games hit.
func (s *StatCascade) HRfg() float64 {
	return s.Count() / float64(s.SumFGH())
}

// Cascade multiplier.
func (s *StatCascade) Mcascade() float64 {
	var pay1 = s.Casc[0].SumPays()
	var pays float64
	for cfn := range s.Casc {
		var payi = s.Casc[cfn].SumPays()
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
	for i := 1; i < len(s.Casc); i++ {
		sum += s.Casc[i].N.Load()
	}
	return float64(sum) / float64(s.Casc[1].N.Load())
}

// Inverse coefficient of fading (C1/Cn)^(1/(n-1)).
func (s *StatCascade) Kfading() float64 {
	var N1 = s.Casc[0].N.Load()
	var Nn uint64
	var cfn int
	for cfn = range s.Casc {
		var Ni = s.Casc[cfn].N.Load()
		if Ni == 0 {
			break
		}
		Nn = Ni
	}
	return math.Pow(float64(N1)/float64(Nn), 1/(float64(cfn)-1))
}

// Maximum number of cascades in avalanche.
func (s *StatCascade) Ncascmax() int {
	for cfn := range s.Casc {
		if s.Casc[cfn].N.Load() == 0 {
			return cfn
		}
	}
	return len(s.Casc)
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

type (
	CalcFunc = func(io.Writer) (float64, float64)
	CalcAlg  = func(ctx context.Context, sp *ScanPar, s Simulator, g SlotGeneric, reels Reelx)
)

const (
	CtxGranulation = 1000 // check context every N reshuffles
	FallLimit      = 20   // maximum 20 cascades for reels ~100 symbols length
)

var (
	ErrAvalanche = errors.New("too many cascading falls")
	ErrReelCount = errors.New("unexpected number of reels")
)

func ScanReels(ctx context.Context, sp *ScanPar, s Simulator, g SlotGeneric, reels Reelx,
	bruteforce, montecarlo CalcAlg, calc CalcFunc) (float64, float64) {
	if sx, sy := g.Dim(); len(reels) != int(sx) {
		panic(fmt.Errorf("%w: %d reels provided for %dx%d slot", ErrReelCount, len(reels), sx, sy))
	}
	if sel := g.GetSel(); sel > 0 {
		fmt.Printf("selected %d lines\n", sel)
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		var ctx2, cancel2 = context.WithCancel(ctx)
		defer cancel2()
		switch sp.Method {
		case game.CMbruteforce:
			go func() {
				defer wg.Done()
				ProgressBF(ctx2, sp, s, calc, float64(reels.Reshuffles()))
			}()
			bruteforce(ctx2, sp, s, g, reels)
		case game.CMmontecarlo:
			go func() {
				defer wg.Done()
				ProgressMC(ctx2, sp, s, calc, g.Cost())
			}()
			montecarlo(ctx2, sp, s, g, reels)
		}
	}()
	wg.Wait()

	fmt.Println()
	fmt.Printf("reels lengths %s, total reshuffles %d\n", reels.String(), reels.Reshuffles())
	return calc(os.Stdout)
}

func ScanReelsCommon(ctx context.Context,
	sp *ScanPar, s Simulator, g SlotGeneric, reels Reelx, calc CalcFunc) (float64, float64) {
	return ScanReels(ctx, sp, s, g, reels, BruteForcex, MonteCarlo, calc)
}

func GetStatGeneric(id string, sn, pn int) (s *StatGeneric, ok bool) {
	var v any
	if v, ok = game.DataRouter[id]; ok {
		s, ok = v.(*StatGeneric)
		return
	}
	s = NewStatGeneric(sn, pn)
	return
}

func FindStatGeneric(pattern string, ref float64, sn, pn int) (s *StatGeneric, ok bool) {
	var v any
	var id string
	if v, id = game.ClosestRoute(pattern, ref); len(id) > 0 {
		s, ok = v.(*StatGeneric)
		return
	}
	s = NewStatGeneric(sn, pn)
	return
}
