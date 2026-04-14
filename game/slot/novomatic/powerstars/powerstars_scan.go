package powerstars

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/slotopol/server/game/slot"
)

type futureResult struct {
	expectedFuture     float64
	secondMomentFuture float64
}

// Calculates the probabilities of transitioning to new states containing wilds.
func calculateTransitionProbs(currentMask int, p2, p3, p4 float64, E, G map[int]float64, muCurrent float64) futureResult {
	res := futureResult{}

	// Checking the free reels (2, 3, 4)
	// For each free reel, the probability of a wild appearing is pi
	// For each occupied slot: 1 (it is already recorded there).
	for nextS := 1; nextS < 8; nextS++ {
		if (nextS|currentMask) == nextS && nextS != currentMask {
			prob := 1.0
			if (currentMask & 1) == 0 {
				if (nextS & 1) > 0 {
					prob *= p2
				} else {
					prob *= (1 - p2)
				}
			}
			if (currentMask & 2) == 0 {
				if (nextS & 2) > 0 {
					prob *= p3
				} else {
					prob *= (1 - p3)
				}
			}
			if (currentMask & 4) == 0 {
				if (nextS & 4) > 0 {
					prob *= p4
				} else {
					prob *= (1 - p4)
				}
			}

			res.expectedFuture += prob * E[nextS]
			res.secondMomentFuture += prob * G[nextS]
		}
	}
	return res
}

func calculateInitialSpin(mu0, d0, p2, p3, p4 float64, E, G map[int]float64) (float64, float64) {
	exp := mu0
	secMom := d0 + mu0*mu0

	for s := 1; s < 8; s++ {
		prob := 1.0
		if (s & 1) > 0 {
			prob *= p2
		} else {
			prob *= (1 - p2)
		}
		if (s & 2) > 0 {
			prob *= p3
		} else {
			prob *= (1 - p3)
		}
		if (s & 4) > 0 {
			prob *= p4
		} else {
			prob *= (1 - p4)
		}

		exp += prob * E[s]
		secMom += prob*G[s] + 2*mu0*prob*E[s]
	}
	return exp, secMom
}

func wcsym(wc bool) byte {
	if wc {
		return '*'
	}
	return '-'
}

// custom parsheet
func CalcStat(ctx context.Context, sp *slot.ScanPar) (rtp, D float64) {
	var wc, _ = ChanceMap.FindClosest(sp.MRTP) // wild chance
	var p2, p3, p4 = wc, wc, wc

	var mu, d [8]float64
	var sr *slot.StatGeneric
	for mask := range mu {
		var wc2, wc3, wc4 = mask&4 > 0, mask&2 > 0, mask&1 > 0

		fmt.Printf("\n(%d/8) calculations of star combinations [%c%c%c]\n", mask+1, wcsym(wc2), wcsym(wc3), wcsym(wc4))
		var s = slot.NewStatGeneric(sn, 5)
		var g = NewGame(sp.Sel)
		if wc2 {
			g.PRW[1] = 1
		}
		if wc3 {
			g.PRW[2] = 1
		}
		if wc4 {
			g.PRW[3] = 1
		}
		var calc = func(w io.Writer) (float64, float64) {
			var µ, D = slot.EvD(s, g.Cost())
			if sp.IsFG() {
				fmt.Fprintf(w, "RTP[%c%c%c] = %.8g%%\n", wcsym(wc2), wcsym(wc3), wcsym(wc4), µ*100)
				slot.Print_all(w, sp, s, µ, D)
			}
			return µ, D
		}
		mu[mask], d[mask] = slot.ScanReelsCommon(ctx, sp, s, g, Reels, calc)
		if mask == 0 {
			sr = s
			// break
		}
	}
	// µ, D values sample:
	// mu = [8]float64{0.6906290681685502, 3.799393045573862, 5.650865445506066, 32.61763641928931, 3.799393045573862, 15.98352672732838, 33.33277680385118, 93.66391184573003}
	// d = [8]float64{3.653645637859768, 31.40956337151818, 48.65385656738518, 501.7239875904347, 32.092689819232646, 319.2589878905468, 462.4906731837034, 4202.498615000493}

	// Calculation results for respin chains
	// E[i] — Expected value of the entire chain of winnings, starting FROM a respin in state i
	// G[i] — Second moment (E[W^2]) of this chain
	var E = make(map[int]float64)
	var G = make(map[int]float64)
	// Iterate in reverse order: from the 3 wilds to the 1.
	E[7] = 0
	G[7] = 0

	// Iterate through all states containing at least one wild symbol (respins sequences).
	for s := 6; s >= 1; s-- {
		currentMu := mu[s]
		currentSecondMoment := d[s] + mu[s]*mu[s]
		var pNew = calculateTransitionProbs(s, p2, p3, p4, E, G, currentMu)
		E[s] = currentMu + pNew.expectedFuture
		G[s] = currentSecondMoment + pNew.secondMomentFuture + 2*currentMu*pNew.expectedFuture
	}

	// Final calculation for the base game (state 0)
	// Here, we calculate the probability of transitioning from state 0 to any state s > 0.
	var totalE, totalG = calculateInitialSpin(mu[0], d[0], p2, p3, p4, E, G)
	rtp = totalE
	D = totalG - totalE*totalE
	if sp.IsMain() {
		fmt.Printf("wild chance: 1/%.5g\n", 1/wc)
		fmt.Printf("RTP = %.8g%%\n", rtp*100)
	}
	var w = os.Stdout
	slot.Print_all(w, sp, sr, rtp, D)
	return rtp, D
}
