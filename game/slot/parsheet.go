package slot

import (
	"fmt"
	"io"
	"math"
)

// Parsheet for simple generic slot (without free games and bonuses).
func Parsheet_generic_simple(w io.Writer, sp *ScanPar, s *StatGeneric, cost float64) (float64, float64) {
	var N, S, Q = s.NSQ(cost)
	var µ = S / N
	var sigma = math.Sqrt(Q/N - µ*µ)
	var vi = GetZ(sp.Conf) * sigma
	var ci = CI(sp.Conf, µ, sigma)
	var BRci = BankrollPlayer(sp.Conf, µ, sigma, ci)
	fmt.Fprintf(w, "RTP = %.8g%%\n", µ*100)
	fmt.Fprintf(w, "sigma = %.6g, VI[%.4g%%] = %.6g (%s)\n", sigma, sp.Conf*100, vi, VIname6[VIclass6(vi)])
	fmt.Fprintf(w, "CI[%.4g%%] = %d, bankroll[CI] = %.6g\n", sp.Conf*100, int(ci), BRci)
	return µ, sigma
}

// Parsheet for generic slot with freegames.
// Each hit of freegames series has `L` freespins.
func Parsheet_generic_freegames(w io.Writer, sp *ScanPar, s *StatGeneric, cost, L float64) (float64, float64) {
	var N, S, Q = s.NSQ(cost)
	var µ = S / N
	var Dsym = Q/N - µ*µ
	var Pfg = s.FGQ()
	var q, sq = s.FSQ()
	var rtpfs = sq * µ
	var rtp = µ + q*rtpfs
	var Eser, Dser = L * sq, L * sq * sq * sq              // Galton-Watson process
	var sigma = math.Sqrt(Dsym + Pfg*(Eser*Dsym+µ*µ*Dser)) // Wald's equation
	var vi = GetZ(sp.Conf) * sigma
	var ci = CI(sp.Conf, rtp, sigma)
	var BRci = BankrollPlayer(sp.Conf, rtp, sigma, ci)
	fmt.Fprintf(w, "symbols: µ = %.8g%%, sigma(sym) = %.6g\n", µ*100, math.Sqrt(Dsym))
	fmt.Fprintf(w, "free spins %d, q = %.5g, sq = 1/(1-q) = %.6f\n", s.FSC.Load(), q, sq)
	fmt.Fprintf(w, "free games hit rate: 1/%.5g\n", s.FGF())
	fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.8g%%\n", µ*100, q, rtpfs*100, rtp*100)
	fmt.Fprintf(w, "sigma = %.6g, VI[%.4g%%] = %.6g (%s)\n", sigma, sp.Conf*100, vi, VIname6[VIclass6(vi)])
	fmt.Fprintf(w, "CI[%.4g%%] = %d, bankroll[CI] = %.6g\n", sp.Conf*100, int(ci), BRci)
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
	var vi = GetZ(sp.Conf) * sigma
	var ci = CI(sp.Conf, µ, sigma)
	var BRci = BankrollPlayer(sp.Conf, µ, sigma, ci)
	fmt.Fprintf(w, "fall[2] = %.10g, Ec2 = Kf2 = 1/%.5g\n", N2, N/N2)
	fmt.Fprintf(w, "fall[3] = %.10g, Ec3 = 1/%.5g, Kf3 = 1/%.5g\n", N3, N/N3, N2/N3)
	fmt.Fprintf(w, "fall[4] = %.10g, Ec4 = 1/%.5g, Kf4 = 1/%.5g\n", N4, N/N4, N3/N4)
	fmt.Fprintf(w, "fall[5] = %.10g, Ec5 = 1/%.5g, Kf5 = 1/%.5g\n", N5, N/N5, N4/N5)
	fmt.Fprintf(w, "Mcascade = %.5g, ACL = %.5g, Kfading = 1/%.5g, Ncascmax = %d\n", s.Mcascade(), s.ACL(), s.Kfading(), s.Ncascmax())
	fmt.Fprintf(w, "RTP = %.8g%%\n", µ*100)
	fmt.Fprintf(w, "sigma = %.6g, VI[%.4g%%] = %.6g (%s)\n", sigma, sp.Conf*100, vi, VIname6[VIclass6(vi)])
	fmt.Fprintf(w, "CI[%.4g%%] = %d, bankroll[CI] = %.6g\n", sp.Conf*100, int(ci), BRci)
	return µ, sigma
}
