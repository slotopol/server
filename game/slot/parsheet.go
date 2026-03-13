package slot

import (
	"fmt"
	"io"
	"math"
	"sort"

	"github.com/slotopol/server/util"
	"gopkg.in/yaml.v3"
)

const ( // print flags for slots
	PF_main = 1 << iota // RTP, sigma and other main information

	PF_jack   // info about progressive jackpots
	PF_fg     // info for bonus reels
	PF_vi     // volatility index
	PF_ci     // convergence index
	PF_spread // RTP spread
	PF_sym    // symbols contribution to payouts
	PF_raw    // simulator raw data
)

// Maximum RTP to get convergence point
const RTPconv = 0.995 // 99.5%

func Print_vi(w io.Writer, sp *ScanPar, sigma float64) {
	if sp.PF&PF_vi == 0 {
		return
	}
	var vi = GetZ(sp.Conf) * sigma
	fmt.Fprintf(w, "sigma = %.6g, VI[%.4g%%] = %.6g (%s)\n", sigma, sp.Conf*100, vi, VIname5[VIclass5(sigma)])
}

func Print_ci(w io.Writer, sp *ScanPar, rtp, sigma float64) {
	if sp.PF&PF_ci == 0 {
		return
	}
	if rtp > RTPconv {
		return
	}
	var ci = CI(sp.Conf, rtp, sigma)
	var BRci = BankrollPlayer(sp.Conf, rtp, sigma, ci)
	fmt.Fprintf(w, "CI[%.4g%%] = %d, bankroll[CI] = %.6g\n", sp.Conf*100, int(ci+0.5), BRci)
}

func Print_ranges(w io.Writer, sp *ScanPar, rtp, sigma float64) {
	if sp.PF&PF_spread == 0 {
		return
	}
	fmt.Fprintln(w)
	fmt.Fprintf(w, "RTP spread for spins number with confidence %.4g%%:\n", sp.Conf*100)
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

func Print_contribution_generic(w io.Writer, sp *ScanPar, s *StatGeneric, rtp float64) {
	if sp.PF&PF_sym == 0 {
		return
	}
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

func Print_contribution_cascade(w io.Writer, sp *ScanPar, s *StatCascade, rtp float64) {
	if sp.PF&PF_sym == 0 {
		return
	}
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

func Print_contribution_falls(w io.Writer, sp *ScanPar, s *StatCascade, rtp float64) {
	if sp.PF&PF_sym == 0 {
		return
	}
	fmt.Fprintln(w)
	fmt.Fprintf(w, "cascades contribution to payouts:\n")
	fmt.Fprintf(w, "fall  rate   rtp\n")
	var sum = s.SumPays()
	for cfn := range s.Casc {
		var c = s.Casc[cfn].SumPays() / sum
		fmt.Fprintf(w, "%2d: %5.2f%% %5.2f%%\n", cfn+1, c*100, rtp*c*100)
		if c == 0 {
			break
		}
	}
}

func Print_raw(w io.Writer, sp *ScanPar, s Simulator) {
	if sp.PF&PF_raw == 0 {
		return
	}
	fmt.Fprintln(w)
	fmt.Fprintf(w, "simulator raw data:\n")
	var b, err = yaml.Marshal(s)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, util.B2S(b))
}

func Print_all(w io.Writer, sp *ScanPar, s Simulator, rtp, sigma float64) {
	Print_vi(w, sp, sigma)
	Print_ci(w, sp, rtp, sigma)
	Print_ranges(w, sp, rtp, sigma)
	switch stat := s.(type) {
	case *StatGeneric:
		Print_contribution_generic(w, sp, stat, rtp)
	case *StatCascade:
		Print_contribution_cascade(w, sp, stat, rtp)
		Print_contribution_falls(w, sp, stat, rtp)
	}
	Print_raw(w, sp, s)
}

// Parsheet for simple generic slot (without free games and bonuses).
func Parsheet_generic_simple(w io.Writer, sp *ScanPar, s *StatGeneric, cost float64) (float64, float64) {
	var N, S, Q = s.NSQ(cost)
	var µ = S / N
	var sigma = math.Sqrt(Q/N - µ*µ)
	if sp.PF&PF_main != 0 {
		fmt.Fprintf(w, "RTP = %.8g%%\n", µ*100)
	}
	Print_all(w, sp, s, µ, sigma)
	return µ, sigma
}

// Parsheet for generic slot with retriggerable freegames
// with `m` multiplier on freegames (m=1 if no multiplier).
// Each hit of freegames series has `L` freespins.
func Parsheet_generic_fgretrig(w io.Writer, sp *ScanPar, s *StatGeneric, cost, m, L float64) (float64, float64) {
	var N, S, Q = s.NSQ(cost)
	var µ = S / N
	var Dsym = Q/N - µ*µ
	var q, sq = s.FSQ()
	var Pfg = s.FGQ()
	var rtpfs = m * sq * µ
	var rtp = µ + q*rtpfs
	var Eser, Dser = L * sq, L * q * sq * sq * sq              // Galton-Watson process
	var sigma = math.Sqrt(Dsym + m*m*Pfg*(Eser*Dsym+µ*µ*Dser)) // Wald's equation
	if sp.PF&PF_main != 0 {
		fmt.Fprintf(w, "symbols: µ = %.8g%%, sigma(sym) = %.6g\n", µ*100, math.Sqrt(Dsym))
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FSC.Load(), q, sq)
		fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", s.FGF())
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.8g%%\n", µ*100, q, rtpfs*100, rtp*100)
	}
	Print_all(w, sp, s, rtp, sigma)
	return rtp, sigma
}

func Parsheet_generic_fgretrig_series(w io.Writer, sp *ScanPar, s *StatGeneric, cost, m float64, L []int, scat Sym) (float64, float64) {
	var N, S, Q = s.NSQ(cost)
	var µ = S / N
	var Dsym = Q/N - µ*µ
	var ΣPL float64
	for i, Li := range L {
		var Pfgi = float64(s.C[scat][i+1].Load()) / N
		ΣPL += Pfgi * float64(Li)
	}
	var q, sq = s.FSQ()
	var rtpfs = m * sq * µ
	var rtp = µ + q*rtpfs
	var sigma = math.Sqrt(Dsym + m*m*ΣPL*(sq*Dsym+µ*µ*q*sq*sq*sq))
	if sp.PF&PF_main != 0 {
		fmt.Fprintf(w, "symbols: µ = %.8g%%, sigma(sym) = %.6g\n", µ*100, math.Sqrt(Dsym))
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FSC.Load(), q, sq)
		fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", s.FGF())
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.8g%%\n", µ*100, q, rtpfs*100, rtp*100)
	}
	Print_all(w, sp, s, rtp, sigma)
	return rtp, sigma
}

// Parsheet for generic slot with splitted statistics for regular
// games `sr` and statistics for NON-retriggerable bonus games `sb`.
// with `m` multiplier on freegames (m=1 if no multiplier).
// Each hit of freegames series has `L` freespins.
func Parsheet_generic_fgonce_split(w io.Writer, sp *ScanPar, sr, sb *StatGeneric, cost, m, L float64) (float64, float64) {
	// bonus reels parameters
	var Nb, Sb, Qb = sb.NSQ(cost)
	var µb = Sb / Nb
	var Dsymb = Qb/Nb - µb*µb
	// regular reels parameters
	var Nr, Sr, Qr = sr.NSQ(cost)
	var µr = Sr / Nr
	var Dsymr = Qr/Nr - µr*µr
	var qr, sqr = sr.FSQ()
	var Pfg = sr.FGQ()
	// calculation
	var rtpfs = m * µb
	var rtp = µr + qr*rtpfs
	var sigma = math.Sqrt(Dsymr + m*m*Pfg*(L*Dsymb+L*L*µb*µb))
	if sp.PF&PF_fg != 0 {
		fmt.Fprintf(w, "*bonus reels*\n")
		fmt.Fprintf(w, "RTP(fg) = %.8g%%\n", rtpfs*100)
	}
	if sp.PF&PF_main != 0 {
		fmt.Fprintf(w, "*regular reels*\n")
		fmt.Fprintf(w, "symbols: µ = %.8g%%, sigma(sym) = %.6g\n", µr*100, math.Sqrt(Dsymr))
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", sr.FSC.Load(), qr, sqr)
		fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", sr.FGF())
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.8g%%\n", µr*100, qr, rtpfs*100, rtp*100)
	}
	Print_all(w, sp, sr, rtp, sigma)
	return rtp, sigma
}

// Parsheet for generic slot with splitted statistics for regular
// games `sr` and statistics for retriggerable bonus games `sb`.
// with `m` multiplier on freegames (m=1 if no multiplier).
// Each hit of freegames series has `L` freespins.
func Parsheet_generic_fgretrig_split(w io.Writer, sp *ScanPar, sr, sb *StatGeneric, cost, m, L float64) (float64, float64) {
	// bonus reels parameters
	var Nb, Sb, Qb = sb.NSQ(cost)
	var µb = Sb / Nb
	var Dsymb = Qb/Nb - µb*µb
	var qb, sqb = sb.FSQ()
	// regular reels parameters
	var Nr, Sr, Qr = sr.NSQ(cost)
	var µr = Sr / Nr
	var Dsymr = Qr/Nr - µr*µr
	var qr, sqr = sr.FSQ()
	var Pfg = sr.FGQ()
	// calculation
	var rtpfs = m * sqb * µb
	var rtp = µr + qr*rtpfs
	var Eser, Dser = L * sqb, L * qb * sqb * sqb * sqb             // Galton-Watson process
	var sigma = math.Sqrt(Dsymr + m*m*Pfg*(Eser*Dsymb+µb*µb*Dser)) // Wald's equation
	if sp.PF&PF_fg != 0 {
		fmt.Fprintf(w, "*bonus reels*\n")
		fmt.Fprintf(w, "symbols: µ = %.8g%%, sigma(sym) = %.6g\n", µb*100, math.Sqrt(Dsymb))
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", sb.FSC.Load(), qb, sqb)
		fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", sb.FGF())
		fmt.Fprintf(w, "rtp(fg) = m*sq*rtp(sym) = %g*%.5g*%.5g = %.6f%%\n", m, sqb, µb*100, rtpfs*100)
	}
	if sp.PF&PF_main != 0 {
		fmt.Fprintf(w, "*regular reels*\n")
		fmt.Fprintf(w, "symbols: µ = %.8g%%, sigma(sym) = %.6g\n", µr*100, math.Sqrt(Dsymr))
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", sr.FSC.Load(), qr, sqr)
		fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", sr.FGF())
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.8g%%\n", µr*100, qr, rtpfs*100, rtp*100)
	}
	Print_all(w, sp, sr, rtp, sigma)
	return rtp, sigma
}

func Parsheet_generic_fgretrig_split_series(w io.Writer, sp *ScanPar, sr, sb *StatGeneric, cost, m float64, L []int, scat Sym) (float64, float64) {
	// bonus reels parameters
	var Nb, Sb, Qb = sb.NSQ(cost)
	var µb = Sb / Nb
	var Dsymb = Qb/Nb - µb*µb
	var qb, sqb = sb.FSQ()
	// regular reels parameters
	var Nr, Sr, Qr = sr.NSQ(cost)
	var µr = Sr / Nr
	var Dsymr = Qr/Nr - µr*µr
	var qr, sqr = sr.FSQ()
	var ΣPL float64
	for i, Li := range L {
		var Pfgi = float64(sr.C[scat][i+1].Load()) / Nr
		ΣPL += Pfgi * float64(Li)
	}
	// calculation
	var rtpfs = m * sqb * µb
	var rtp = µr + qr*rtpfs
	var sigma = math.Sqrt(Dsymr + m*m*ΣPL*(sqb*Dsymb+µb*µb*qb*sqb*sqb*sqb))
	if sp.PF&PF_fg != 0 {
		fmt.Fprintf(w, "*bonus reels*\n")
		fmt.Fprintf(w, "symbols: µ = %.8g%%, sigma(sym) = %.6g\n", µb*100, math.Sqrt(Dsymb))
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", sb.FSC.Load(), qb, sqb)
		fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", sb.FGF())
		fmt.Fprintf(w, "rtp(fg) = m*sq*rtp(sym) = %g*%.5g*%.5g = %.6f%%\n", m, sqb, µb*100, rtpfs*100)
	}
	if sp.PF&PF_main != 0 {
		fmt.Fprintf(w, "*regular reels*\n")
		fmt.Fprintf(w, "symbols: µ = %.8g%%, sigma(sym) = %.6g\n", µr*100, math.Sqrt(Dsymr))
		fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", sr.FSC.Load(), qr, sqr)
		fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", sr.FGF())
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.8g%%\n", µr*100, qr, rtpfs*100, rtp*100)
	}
	Print_all(w, sp, sr, rtp, sigma)
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
	if sp.PF&PF_main != 0 {
		fmt.Fprintf(w, "fall[2] = %.10g, Ec2 = Kf2 = 1/%.5g\n", N2, N/N2)
		fmt.Fprintf(w, "fall[3] = %.10g, Ec3 = 1/%.5g, Kf3 = 1/%.5g\n", N3, N/N3, N2/N3)
		fmt.Fprintf(w, "fall[4] = %.10g, Ec4 = 1/%.5g, Kf4 = 1/%.5g\n", N4, N/N4, N3/N4)
		fmt.Fprintf(w, "fall[5] = %.10g, Ec5 = 1/%.5g, Kf5 = 1/%.5g\n", N5, N/N5, N4/N5)
		fmt.Fprintf(w, "Mcascade = %.5g, ACL = %.5g, Kfading = 1/%.5g, Ncascmax = %d\n", s.Mcascade(), s.ACL(), s.Kfading(), s.Ncascmax())
		fmt.Fprintf(w, "RTP = %.8g%%\n", µ*100)
	}
	Print_all(w, sp, s, µ, sigma)
	return µ, sigma
}
