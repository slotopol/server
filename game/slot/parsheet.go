package slot

import (
	"fmt"
	"io"
	"math"
	"sort"
)

// Maximum RTP to get convergence point
const RTPconv = 0.995 // 99.5%

func print_vi(w io.Writer, sp *ScanPar, sigma float64) {
	var vi = GetZ(sp.Conf) * sigma
	fmt.Fprintf(w, "sigma = %.6g, VI[%.4g%%] = %.6g (%s)\n", sigma, sp.Conf*100, vi, VIname5[VIclass5(sigma)])
}

func print_ci(w io.Writer, sp *ScanPar, rtp, sigma float64) {
	if rtp > RTPconv {
		return
	}
	var ci = CI(sp.Conf, rtp, sigma)
	var BRci = BankrollPlayer(sp.Conf, rtp, sigma, ci)
	fmt.Fprintf(w, "CI[%.4g%%] = %d, bankroll[CI] = %.6g\n", sp.Conf*100, int(ci+0.5), BRci)
}

func print_ranges(w io.Writer, sp *ScanPar, rtp, sigma float64) {
	fmt.Fprintln(w)
	fmt.Fprintf(w, "RTP ranges for spins number with confidence %.4g%%:\n", sp.Conf*100)
	var N = []int{1e3, 1e4, 1e5, 1e6, 1e7}
	var vi = GetZ(sp.Conf) * sigma
	var ci = CI(sp.Conf, rtp, sigma)
	if ci < 1e7 {
		N = append(N, int(ci+0.5))
		sort.Ints(N)
	}
	for _, n := range N {
		var Δ = vi / math.Sqrt(float64(n))
		fmt.Fprintf(w, "%8d: %.2f%% ... %.2f%%\n", n, (rtp-Δ)*100, (rtp+Δ)*100)
	}
}

func print_contribution_generic(w io.Writer, sp *ScanPar, s *StatGeneric, rtp float64) {
	fmt.Fprintln(w)
	fmt.Fprintf(w, "symbols contribution to payouts:\n")
	fmt.Fprintf(w, "sym   rate   rtp\n")
	var sum = s.SumPays()
	var sym Sym
	for sym = 1; sym < Sym(len(s.S)); sym++ {
		var c = s.SymPays(sym) / sum
		fmt.Fprintf(w, "%2d: %5.2f%% %5.2f%%\n", sym, c*100, rtp*c*100)
	}
}

func print_contribution_cascade(w io.Writer, sp *ScanPar, s *StatCascade, rtp float64) {
	fmt.Fprintln(w)
	fmt.Fprintf(w, "symbols contribution to payouts:\n")
	fmt.Fprintf(w, "sym   rate   rtp\n")
	var sum = s.SumPays()
	var sym Sym
	for sym = 1; sym < Sym(len(s.Casc[0].S)); sym++ {
		var c = s.SymPays(sym) / sum
		fmt.Fprintf(w, "%2d: %5.2f%% %5.2f%%\n", sym, c*100, rtp*c*100)
	}
}

func print_contribution_falls(w io.Writer, sp *ScanPar, s *StatCascade, rtp float64) {
	fmt.Fprintln(w)
	fmt.Fprintf(w, "cascades contribution to payouts:\n")
	fmt.Fprintf(w, "fall  rate   rtp\n")
	var sum = s.SumPays()
	for cfn := range FallLimit {
		var c = s.Casc[cfn].SumPays() / sum
		fmt.Fprintf(w, "%2d: %5.2f%% %5.2f%%\n", cfn+1, c*100, rtp*c*100)
		if c == 0 {
			break
		}
	}
}

// Parsheet for simple generic slot (without free games and bonuses).
func Parsheet_generic_simple(w io.Writer, sp *ScanPar, s *StatGeneric, cost float64) (float64, float64) {
	var N, S, Q = s.NSQ(cost)
	var µ = S / N
	var sigma = math.Sqrt(Q/N - µ*µ)
	fmt.Fprintf(w, "RTP = %.8g%%\n", µ*100)
	print_vi(w, sp, sigma)
	print_ci(w, sp, µ, sigma)
	print_ranges(w, sp, µ, sigma)
	print_contribution_generic(w, sp, s, µ)
	return µ, sigma
}

// Parsheet for generic slot with freegames
// with `m` multiplier on freegames (m=1 if no multiplier).
// Each hit of freegames series has `L` freespins.
func Parsheet_generic_freegames(w io.Writer, sp *ScanPar, s *StatGeneric, cost, m float64, L int) (float64, float64) {
	var N, S, Q = s.NSQ(cost)
	var µ = S / N
	var Dsym = Q/N - µ*µ
	var Pfg = s.FGQ()
	var q, sq = s.FSQ()
	var rtpfs = m * sq * µ
	var rtp = µ + q*rtpfs
	var Eser, Dser = float64(L) * sq, float64(L) * q * sq * sq * sq // Galton-Watson process
	var sigma = math.Sqrt(Dsym + m*m*Pfg*(Eser*Dsym+µ*µ*Dser))      // Wald's equation
	fmt.Fprintf(w, "symbols: µ = %.8g%%, sigma(sym) = %.6g\n", µ*100, math.Sqrt(Dsym))
	fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FSC.Load(), q, sq)
	fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", s.FGF())
	fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.8g%%\n", µ*100, q, rtpfs*100, rtp*100)
	print_vi(w, sp, sigma)
	print_ci(w, sp, rtp, sigma)
	print_ranges(w, sp, rtp, sigma)
	print_contribution_generic(w, sp, s, rtp)
	return rtp, sigma
}

func Parsheet_generic_fgseries(w io.Writer, sp *ScanPar, s *StatGeneric, cost, m float64, L []int, scat Sym) (float64, float64) {
	var N, S, Q = s.NSQ(cost)
	var µ = S / N
	var Dsym = Q/N - µ*µ
	var Pfg = make([]float64, len(L))
	for i := range L {
		Pfg[i] = float64(s.C[scat][i+1].Load()) / N
	}
	var ΣPL float64
	for i, Li := range L {
		var Pfgi = float64(s.C[scat][i+1].Load()) / N
		ΣPL += Pfgi * float64(Li)
	}
	var q, sq = s.FSQ()
	var rtpfs = m * sq * µ
	var rtp = µ + q*rtpfs
	var sigma = math.Sqrt(Dsym + m*m*ΣPL*(Dsym*sq+µ*µ*q*sq*sq*sq))
	fmt.Fprintf(w, "symbols: µ = %.8g%%, sigma(sym) = %.6g\n", µ*100, math.Sqrt(Dsym))
	fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FSC.Load(), q, sq)
	fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", s.FGF())
	fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.8g%%\n", µ*100, q, rtpfs*100, rtp*100)
	print_vi(w, sp, sigma)
	print_ci(w, sp, rtp, sigma)
	print_ranges(w, sp, rtp, sigma)
	print_contribution_generic(w, sp, s, µ)
	return rtp, sigma
}

// Parsheet for simple cascade slot (without free games and bonuses).
func Parsheet_cascade_simple(w io.Writer, sp *ScanPar, s *StatCascade, cost float64) (float64, float64) {
	var N, S, Q = s.NSQ(cost)
	var N2 = float64(s.Casc[1].N.Load())
	var N3 = float64(s.Casc[2].N.Load())
	var N4 = float64(s.Casc[3].N.Load())
	var N5 = float64(s.Casc[4].N.Load())
	var µ = S / N
	var sigma = math.Sqrt(Q/N - µ*µ)
	fmt.Fprintf(w, "fall[2] = %.10g, Ec2 = Kf2 = 1/%.5g\n", N2, N/N2)
	fmt.Fprintf(w, "fall[3] = %.10g, Ec3 = 1/%.5g, Kf3 = 1/%.5g\n", N3, N/N3, N2/N3)
	fmt.Fprintf(w, "fall[4] = %.10g, Ec4 = 1/%.5g, Kf4 = 1/%.5g\n", N4, N/N4, N3/N4)
	fmt.Fprintf(w, "fall[5] = %.10g, Ec5 = 1/%.5g, Kf5 = 1/%.5g\n", N5, N/N5, N4/N5)
	fmt.Fprintf(w, "Mcascade = %.5g, ACL = %.5g, Kfading = 1/%.5g, Ncascmax = %d\n", s.Mcascade(), s.ACL(), s.Kfading(), s.Ncascmax())
	fmt.Fprintf(w, "RTP = %.8g%%\n", µ*100)
	print_vi(w, sp, sigma)
	print_ci(w, sp, µ, sigma)
	print_ranges(w, sp, µ, sigma)
	print_contribution_cascade(w, sp, s, µ)
	print_contribution_falls(w, sp, s, µ)
	return µ, sigma
}
