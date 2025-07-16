package fairyqueen

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/slotopol/server/game/slot"
)

func BruteForce5x3es2(ctx context.Context, s slot.Stater, g *Game, reels *slot.Reels5x, es slot.Sym) {
	var tn = slot.CorrectThrNum()
	var tn64 = uint64(tn)
	var r3 = reels.Reel(3)
	var r4 = reels.Reel(4)
	var r5 = reels.Reel(5)
	var wg sync.WaitGroup
	wg.Add(tn)
	for ti := range tn64 {
		var sg = g.Clone().(slot.ClassicSlot) // classic slot game
		var reshuf uint64
		go func() {
			defer wg.Done()

			var wins slot.Wins
			var x slot.Pos
			for x = 1; x <= 2; x++ {
				if es != scat {
					g.SetSym(x, 1, es)
					g.SetSym(x, 2, es)
					g.SetSym(x, 3, es)
				} else {
					g.SetSym(x, 1, 0)
					g.SetSym(x, 2, scat)
					g.SetSym(x, 3, 0)
				}
			}

			for i3 := range r3 {
				sg.SetCol(3, r3, i3)
				for i4 := range r4 {
					sg.SetCol(4, r4, i4)
					for i5 := range r5 {
						reshuf++
						if reshuf%slot.CtxGranulation == 0 {
							select {
							case <-ctx.Done():
								return
							default:
							}
						}
						if reshuf%tn64 != ti {
							continue
						}
						sg.SetCol(5, r5, i5)
						sg.Scanner(&wins)
						s.Update(wins, 1)
						wins.Reset()
					}
				}
			}
		}()
	}
	wg.Wait()
}

func BruteForce5x3es3(ctx context.Context, s slot.Stater, g *Game, reels *slot.Reels5x, es slot.Sym) {
	var tn = slot.CorrectThrNum()
	var tn64 = uint64(tn)
	var r4 = reels.Reel(4)
	var r5 = reels.Reel(5)
	var wg sync.WaitGroup
	wg.Add(tn)
	for ti := range tn64 {
		var sg = g.Clone().(slot.ClassicSlot) // classic slot game
		var reshuf uint64
		go func() {
			defer wg.Done()

			var wins slot.Wins
			var x slot.Pos
			for x = 1; x <= 3; x++ {
				if es != scat {
					g.SetSym(x, 1, es)
					g.SetSym(x, 2, es)
					g.SetSym(x, 3, es)
				} else {
					g.SetSym(x, 1, 0)
					g.SetSym(x, 2, scat)
					g.SetSym(x, 3, 0)
				}
			}

			for i4 := range r4 {
				sg.SetCol(4, r4, i4)
				for i5 := range r5 {
					reshuf++
					if reshuf%slot.CtxGranulation == 0 {
						select {
						case <-ctx.Done():
							return
						default:
						}
					}
					if reshuf%tn64 != ti {
						continue
					}
					sg.SetCol(5, r5, i5)
					sg.Scanner(&wins)
					s.Update(wins, 1)
					wins.Reset()
				}
			}
		}()
	}
	wg.Wait()
}

func CalcStatBon(ctx context.Context, es slot.Sym) (float64, float64) {
	var reels = ReelsBon
	var g = NewGame()
	g.Sel = 1
	g.FSR = 10 // set free spins mode
	g.ES = es
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var reshuf = s.Count()
		var cost, _ = g.Cost()
		var lrtp, srtp = s.LineRTP(cost), s.ScatRTP(cost)
		var rtpsym = lrtp + srtp
		var q = s.FreeCount() / reshuf
		var sq = 1 / (1 - q)
		if q > 0 {
			fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
			fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCountU(), q, sq)
			fmt.Fprintf(w, "free games frequency: 1/%.5g\n", reshuf/s.FreeHits())
		}
		fmt.Fprintf(w, "RTP[%d] = %.6f%%\n", es, rtpsym)
		return rtpsym
	}

	func() {
		var ctx2, cancel2 = context.WithCancel(ctx)
		defer cancel2()
		if ReelNumBon[g.ES-1] == 2 {
			s.SetPlan(uint64(len(reels.Reel(3))) * uint64(len(reels.Reel(4))) * uint64(len(reels.Reel(5))))
			BruteForce5x3es2(ctx2, &s, g, reels, g.ES)
		} else {
			s.SetPlan(uint64(len(reels.Reel(4))) * uint64(len(reels.Reel(5))))
			BruteForce5x3es3(ctx2, &s, g, reels, g.ES)
		}
	}()
	return calc(os.Stdout), s.FreeCount() / s.Count()
}

func CalcStatReg(ctx context.Context, mrtp float64) float64 {
	fmt.Printf("*bonus reels calculations*\n")
	var rtpe = map[slot.Sym]float64{}
	var qe = map[slot.Sym]float64{}
	var es slot.Sym
	for es = 2; es <= scat; es++ {
		fmt.Printf("*calculations for expanding symbol [%d]*\n", es)
		rtpe[es], qe[es] = CalcStatBon(ctx, es)
		if ctx.Err() != nil {
			return 0
		}
	}
	var rtpsym, qfs float64
	for _, es := range ExpSymReel {
		rtpsym += rtpe[es]
		qfs += qe[es]
	}
	rtpsym /= float64(len(ExpSymReel))
	qfs /= float64(len(ExpSymReel))
	var sqfs = 1 / (1 - qfs)
	var rtpfs = sqfs * rtpsym
	fmt.Printf("free spins: q = %.5g, sq = 1/(1-q) = %.6f\n", qfs, sqfs)
	fmt.Printf("free games frequency: 1/%.5g\n", 10/qfs)
	fmt.Printf("RTPfs = sq*rtp(sym) = %.5g*%.5g = %.6f%%\n", sqfs, rtpsym, rtpfs)

	fmt.Printf("*regular reels calculations*\n")
	var reels, _ = slot.FindClosest(ReelsMap, mrtp)
	var g = NewGame()
	g.Sel = 1
	var s slot.Stat

	var calc = func(w io.Writer) float64 {
		var reshuf = s.Count()
		var cost, _ = g.Cost()
		var lrtp, srtp = s.LineRTP(cost), s.ScatRTP(cost)
		var rtpsym = lrtp + srtp
		var q = s.FreeCount() / reshuf
		var sq = 1 / (1 - q)
		var rtp = rtpsym + q*rtpfs
		fmt.Fprintf(w, "symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FreeCountU(), q, sq)
		fmt.Fprintf(w, "free games frequency: 1/%.5g\n", reshuf/s.FreeHits())
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%\n", rtpsym, q, rtpfs, rtp)
		return rtp
	}

	return slot.ScanReels5x(ctx, &s, g, reels, calc)
}
