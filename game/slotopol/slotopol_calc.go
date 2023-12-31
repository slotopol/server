package slotopol

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/slotopol/server/game"
)

var Emje float64 // Eldorado game 1 spin calculated expectation

// Eldorado expectation.
func ExpEldorado() float64 {
	var sum = 0
	for _, v := range Eldorado {
		sum += v
	}
	var E = float64(sum) / float64(len(Eldorado))
	fmt.Printf("eldorado 1 spin: count = %d, E = %g\n", len(Eldorado), E)
	return E
}

var Emjm float64 // Monopoly game calculated expectation

func ExpMonopoly() float64 {
	var dices = [7]int{1, 1, 1, 1, 1, 1, 1}
	var sumi, sumii, sumj int
	var count, zcount int = 6 * 6 * 6 * 6 * 6 * 6 * 6, 0
	var pos, win int
	for i1 := 1; i1 <= 6; i1++ {
		for i2 := 1; i2 <= 6; i2++ {
			for i3 := 1; i3 <= 6; i3++ {
				for i4 := 1; i4 <= 6; i4++ {
					for i5 := 1; i5 <= 6; i5++ {
						for i6 := 1; i6 <= 6; i6++ {
							for i7 := 1; i7 <= 6; i7++ {
								pos, sumj = 1, 0
								for j := 0; j < 7; j++ {
									pos = (pos+dices[j]-1)%20 + 1
									if Monopoly[pos-1].Jump > 0 {
										win = Monopoly[Monopoly[pos-1].Jump-1].Mult
									} else {
										win = Monopoly[pos-1].Mult
									}
									if Monopoly[pos-1].Dice {
										win *= dices[j]
									}
									if Monopoly[pos-1].Jump > 0 {
										pos = Monopoly[pos-1].Jump
									}
									sumj += win
								}
								if sumj == 0 {
									sumj += 5000
									zcount++
								}
								sumi += sumj
								sumii += sumj * sumj

								dices[6] = dices[6]%6 + 1
							}
							dices[5] = dices[5]%6 + 1
						}
						dices[4] = dices[4]%6 + 1
					}
					dices[3] = dices[3]%6 + 1
				}
				dices[2] = dices[2]%6 + 1
			}
			dices[1] = dices[1]%6 + 1
		}
		dices[0] = dices[0]%6 + 1
	}
	var E = float64(sumi) / float64(count)
	var v = float64(sumii)/float64(count) - E*E
	var sigma = math.Sqrt(v)
	fmt.Printf("monopoly: count = %d, sum = %d, zerocount = %d, p(zero) = 1/%d, E = %g\n", count, sumi, zcount, int(float64(count)/float64(zcount)), E)
	fmt.Printf("monopoly: variance = %.6g, sigma = %.6g, limits = %.6g ... %.6g\n", v, sigma, E-sigma, E+sigma)
	return E
}

func CalcStat(ctx context.Context, rn string) float64 {
	fmt.Printf("*bonus games calculations*\n")
	Emje = ExpEldorado()
	Emjm = ExpMonopoly()
	fmt.Printf("*reels calculations*\n")
	var reels *game.Reels5x
	if rn != "" {
		var ok bool
		if reels, ok = ReelsMap[rn]; !ok {
			return 0
		}
	} else {
		reels = &Reels100
	}
	var g = NewGame(reels)
	var sbl = float64(len(g.SBL))
	var s game.Stat

	var dur = func() time.Duration {
		var t0 = time.Now()
		var ctx2, cancel2 = context.WithCancel(ctx)
		defer cancel2()
		go s.Progress(ctx2, time.NewTicker(2*time.Second), sbl, float64(g.Reels.Reshuffles()))
		s.BruteForce5x(ctx2, g, g.Reels)
		return time.Since(t0)
	}()

	var n = float64(s.Reshuffles)
	var lrtp, srtp = float64(s.LinePay) / n / sbl * 100, float64(s.ScatPay) / n * 100
	var rtpsym = lrtp + srtp
	var qmje9 = float64(s.BonusCount[mje9]) / n / sbl
	var rtpmje9 = Emje * 9 * qmje9 * 100
	var qmjm = float64(s.BonusCount[mjm]) / n / sbl
	var rtpmjm = Emjm * qmjm * 100
	var rtp = rtpsym + rtpmje9 + rtpmjm
	fmt.Printf("completed %.5g%%, selected %d lines, time spent %v\n", float64(s.Reshuffles)/float64(g.Reels.Reshuffles())*100, len(g.SBL), dur)
	fmt.Printf("reels lengths [%d, %d, %d, %d, %d], total reshuffles %d\n",
		len(g.Reels.Reel(1)), len(g.Reels.Reel(2)), len(g.Reels.Reel(3)), len(g.Reels.Reel(4)), len(g.Reels.Reel(5)), g.Reels.Reshuffles())
	fmt.Printf("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%\n", lrtp, srtp, rtpsym)
	fmt.Printf("spin9 bonuses: count %d, rtp = %.6f%%\n", s.BonusCount[mje9], rtpmje9)
	fmt.Printf("monopoly bonuses: count %d, rtp = %.6f%%\n", s.BonusCount[mjm], rtpmjm)
	fmt.Printf("jackpots: count %d, frequency 1/%d\n", s.JackCount[jid], int(n/float64(s.JackCount[jid])))
	fmt.Printf("RTP = %.5g(sym) + %.5g(mje9) + %.5g(mjm) = %.6f%%\n", rtpsym, rtpmje9, rtpmjm, rtp)
	return rtp
}
