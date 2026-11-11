package slot

import (
	"fmt"
	"io"
	"math"
	"sort"

	"github.com/slotopol/server/game"
	"github.com/slotopol/server/util"
	"go.yaml.in/yaml/v3"
)

func Print_vi(w io.Writer, sp *ScanPar, D float64) {
	if !sp.IsVI() {
		return
	}
	var sigma = math.Sqrt(D)
	var vi = game.GetZ(sp.Conf) * sigma
	fmt.Fprintf(w, "sigma = %.6g, VI[%.4g%%] = %.6g (%s)\n", sigma, sp.Conf*100, vi, game.VIname5[game.VIclass5(sigma)])
}

func Print_ci(w io.Writer, sp *ScanPar, µ, D float64) {
	if !sp.IsCI() {
		return
	}
	if µ > game.RTPconv {
		return
	}
	var sigma = math.Sqrt(D)
	var ci = game.CI(sp.Conf, µ, sigma)
	var BRci = game.BankrollPlayer(sp.Conf, µ, sigma, ci)
	fmt.Fprintf(w, "CI[%.4g%%] = %d, bankroll[CI] = %.6g\n", sp.Conf*100, int(ci+0.5), BRci)
}

func Print_spread(w io.Writer, sp *ScanPar, µ, D float64) {
	if !sp.IsSpread() {
		return
	}
	fmt.Fprintln(w)
	fmt.Fprintf(w, "RTP spread for spins number with confidence %.4g%%:\n", sp.Conf*100)
	var N = []int{1e3, 1e4, 1e5, 1e6, 1e7}
	var sigma = math.Sqrt(D)
	var vi = game.GetZ(sp.Conf) * sigma
	var ci = game.CI(sp.Conf, µ, sigma)
	if ci < 1e7 {
		N = append(N, int(ci+0.5))
		sort.Ints(N)
	}
	for _, n := range N {
		var Δ = vi / math.Sqrt(float64(n))
		fmt.Fprintf(w, "%8d: %.2f%% ... %.2f%%\n", n, (µ-Δ)*100, (µ+Δ)*100)
	}
}

func p5f(p float64) string {
	if p != 0 {
		return fmt.Sprintf("%5.2f", p)
	} else {
		return "    0"
	}
}

func Print_symbols_generic(w io.Writer, sp *ScanPar, s *StatGeneric, µ float64) {
	if !sp.IsSym() {
		return
	}
	fmt.Fprintln(w)
	fmt.Fprintf(w, "symbols contribution to payouts:\n")
	fmt.Fprintf(w, "sym rate%%  rtp%% |")
	for x := range s.S[0] {
		fmt.Fprintf(w, " %5d", x+1)
	}
	fmt.Fprintf(w, "\n")
	var sum = s.SumPays()
	for sym, pays := range s.S {
		var cs = s.SymPays(Sym(sym+1)) / sum
		fmt.Fprintf(w, "%2d: %s %s |", sym+1, p5f(cs*100), p5f(µ*cs*100))
		for x := range pays {
			var cx = pays[x].Load() / sum
			fmt.Fprintf(w, " %s", p5f(cx*100))
		}
		fmt.Fprintf(w, "\n")
	}
}

func Print_symbols_cascade(w io.Writer, sp *ScanPar, s *StatCascade, µ float64) {
	if !sp.IsSym() {
		return
	}
	fmt.Fprintln(w)
	fmt.Fprintf(w, "symbols contribution to payouts:\n")
	fmt.Fprintf(w, "sym rate%%  rtp%% |")
	for x := range s.Casc[0].S[0] {
		fmt.Fprintf(w, " %5d", x+1)
	}
	fmt.Fprintf(w, "\n")
	var sum = s.SumPays()
	for sym, pays := range s.Casc[0].S {
		var cs = s.SymPays(Sym(sym+1)) / sum
		fmt.Fprintf(w, "%2d: %s %s |", sym+1, p5f(cs*100), p5f(µ*cs*100))
		for x := range pays {
			var cx float64
			for cfn := range s.Casc {
				cx += s.Casc[cfn].S[sym][x].Load()
			}
			cx /= sum
			fmt.Fprintf(w, " %s", p5f(cx*100))
		}
		fmt.Fprintf(w, "\n")
	}
}

func Print_contribution_falls(w io.Writer, sp *ScanPar, s *StatCascade, µ float64) {
	if !sp.IsCasc() {
		return
	}
	fmt.Fprintln(w)
	fmt.Fprintf(w, "cascades contribution to payouts:\n")
	fmt.Fprintf(w, "cfn rate%%  rtp%%\n")
	var sum = s.SumPays()
	for cfn := range s.Casc {
		var c = s.Casc[cfn].SumPays() / sum
		fmt.Fprintf(w, "%2d: %s %s\n", cfn, p5f(c*100), p5f(µ*c*100))
		if c == 0 {
			break
		}
	}
}

func Print_raw(w io.Writer, sp *ScanPar, s Simulator) {
	if !sp.IsRaw() {
		return
	}
	fmt.Fprintln(w)
	fmt.Fprintf(w, "simulator raw data:\n")
	var b, err = yaml.Marshal(s)
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}
	fmt.Fprintln(w, util.B2S(b))
}

func Print_cascmetrics(w io.Writer, sp *ScanPar, s *StatCascade) {
	if !sp.IsCM() {
		return
	}
	fmt.Fprintln(w)
	fmt.Fprintf(w, "cascade metrics:\n")
	var N = s.Count()
	var Nc = float64(s.Casc[1].N.Load()) // current number
	fmt.Fprintf(w, "N[1] = %.10g, Ec1 = Kf1 = 1/%.5g\n", Nc, N/Nc)
	var Np float64 // previous number
	for cfn := 2; cfn < len(s.Casc); cfn++ {
		Np = Nc
		Nc = float64(s.Casc[cfn].N.Load())
		if Nc == 0 {
			break
		}
		if N/Nc < 1e5 {
			fmt.Fprintf(w, "N[%d] = %.10g, Ec%d = 1/%.5g, Kf%d = 1/%.5g\n", cfn, Nc, cfn, N/Nc, cfn, Np/Nc)
		} else {
			fmt.Fprintf(w, "N[%d] = %.10g, Ec%d = 1/%d, Kf%d = 1/%.5g\n", cfn, Nc, cfn, int(N/Nc), cfn, Np/Nc)
		}
	}
	fmt.Fprintf(w, "Mcascade = %.5g, ACL = %.5g, Kfading = 1/%.5g, Ncascmax = %d\n", s.Mcascade(), s.ACL(), s.Kfading(), s.Ncascmax())
}

// Print_all assembles all print functions of parsheets.
func Print_all(w io.Writer, sp *ScanPar, s Simulator, µ, D float64) {
	Print_vi(w, sp, D)
	Print_ci(w, sp, µ, D)
	Print_spread(w, sp, µ, D)
	switch stat := s.(type) {
	case *StatGeneric:
		Print_symbols_generic(w, sp, stat, µ)
	case *StatCascade:
		Print_cascmetrics(w, sp, stat)
		Print_symbols_cascade(w, sp, stat, µ)
		Print_contribution_falls(w, sp, stat, µ)
	}
	Print_raw(w, sp, s)
}

// Parsheet for simple slot model (without free games and bonuses).
func Parsheet_simple(w io.Writer, sp *ScanPar, s Counter, cost float64) (float64, float64) {
	var µ, D = s.EvD(cost)
	if sp.IsMain() {
		fmt.Fprintf(w, "RTP = %.8g%%\n", µ*100)
	}
	Print_all(w, sp, s, µ, D)
	return µ, D
}

// Parsheet for model with identical rules and symbol payout structures
// between base paid spins and free bonus spins, recursive re-triggers
// during free spins, and a global multiplier `m` applied during bonus rounds.
func Parsheet_solereelrecur(w io.Writer, sp *ScanPar, s Simulator, cost float64, m float64) (float64, float64) {
	var EVbase, Dbase = s.EvD(cost)
	var HRfg, EK, EK2 = s.Efs()
	var VarK = EK2 - EK*EK
	var SEK = 1 / (1 - EK)
	var VarNfg = VarK * SEK * SEK * SEK
	var q = EK / (1 - EK)
	var EV = EVbase * (1 + m*q)
	var D = Dbase*(1+q*m*m) + m*m*EVbase*EVbase*VarNfg
	if sp.IsMain() {
		fmt.Fprintf(w, "base: µ = %.8g%%, sigma(base) = %.6g\n", EVbase*100, math.Sqrt(Dbase))
		fmt.Fprintf(w, "free: HRfg = 1/%.5g, E[K] = %.5g, q = %.5g\n", 1/HRfg, EK, q)
		if m > 1 {
			fmt.Fprintf(w, "RTP = %.5g*(1 + %.5g*%.5g) = %.8g%%\n", EVbase*100, q, m, EV*100)
		} else {
			fmt.Fprintf(w, "RTP = %.5g*(1 + %.5g) = %.8g%%\n", EVbase*100, q, EV*100)
		}
	}
	Print_all(w, sp, s, EV, D)
	return EV, D
}

// Parsheet to calculate bonus reel in dual-reel set model with recursive re-triggers.
func Parsheet_affixrecur(w io.Writer, sp *ScanPar, s Simulator, cost float64) (float64, float64) {
	const m = 1
	var EVbase, Dbase = s.EvD(cost)
	var HRfg, EK, EK2 = s.Efs()
	var VarK = EK2 - EK*EK
	var SEK = 1 / (1 - EK)
	var VarNfg = VarK * SEK * SEK * SEK
	var q = EK / (1 - EK)
	var EV = EVbase * (1 + m*q)
	var D = Dbase*(1+q*m*m) + m*m*EVbase*EVbase*VarNfg
	if sp.IsFG() {
		fmt.Fprintf(w, "base: µ = %.8g%%, sigma(base) = %.6g\n", EVbase*100, math.Sqrt(Dbase))
		fmt.Fprintf(w, "free: HRfg = 1/%.5g, E[K] = %.5g, q = %.5g\n", 1/HRfg, EK, q)
		if m > 1 {
			fmt.Fprintf(w, "RTP = %.5g*(1 + %.5g*%.5g) = %.8g%%\n", EVbase*100, q, m, EV*100)
		} else {
			fmt.Fprintf(w, "RTP = %.5g*(1 + %.5g) = %.8g%%\n", EVbase*100, q, EV*100)
		}
	}
	Print_vi(w, sp, D)
	Print_ci(w, sp, EV, D)
	return EV, D
}

// Parsheet for model with different reel lengths, independent counters
// for each set, recursive re-triggers within bonus spins (which evaluate
// on the bonus reel set), and a multiplier `m` applied to bonus spins.
func Parsheet_dualreelrecur(w io.Writer, sp *ScanPar, sr Simulator, sb Simulator, cost float64) (float64, float64) {
	const m = 1
	// regular reels parameters
	var EVr, Dr = sr.EvD(cost)
	var HRfgr, EKr, EK2r = sr.Efs()
	var VarKr = EK2r - EKr*EKr
	// bonus reels parameters
	var EVb, Db = sb.EvD(cost)
	var HRfgb, EKb, EK2b = sb.Efs()
	var VarKb = EK2b - EKb*EKb
	// aggregate parameters
	var a = 1 / (1 - EKb)
	var q = EKr * a
	var VarNfg = VarKr*a*a + VarKb*EKr*a*a*a
	var EV = EVr + q*m*EVb
	var D = Dr + m*m*(EVb*EVb*VarNfg+q*Db)
	if sp.IsFG() {
		fmt.Fprintf(w, "*bonus reels*\n")
		fmt.Fprintf(w, "µ = %.8g%%, sigma(base) = %.6g\n", EVb*100, math.Sqrt(Db))
		if EKb > 0 {
			fmt.Fprintf(w, "free: HRfg = 1/%.5g, E[K] = %.5g\n", 1/HRfgb, EKb)
		}
	}
	if sp.IsMain() {
		fmt.Fprintf(w, "*regular reels*\n")
		fmt.Fprintf(w, "µ = %.8g%%, sigma(base) = %.6g\n", EVr*100, math.Sqrt(Dr))
		fmt.Fprintf(w, "free: HRfg = 1/%.5g, E[K] = %.5g, q = %.5g\n", 1/HRfgr, EKr, q)
		fmt.Fprintf(w, "RTP = %.5g(r) + %.5g*%.5g(b) = %.8g%%\n", EVr*100, q, EVb*100, EV*100)
	}
	Print_all(w, sp, sr, EV, D)
	return EV, D
}

// Parsheet to calculate bonus reel in dual-reel set model.
func Parsheet_affixonce(w io.Writer, sp *ScanPar, s Simulator, cost float64, Lr, Lb float64) (float64, float64) {
	const m = 1
	var EVbase, Dbase = s.EvD(cost)
	var HRfg, EK, EK2 = s.Efs()
	var pb = HRfg                    // P
	var Pre = 1 - math.Pow(1-pb, Lr) // P(A)=1−(1−P)^N
	var EV = m * EVbase * (1 + Pre*Lb/Lr)
	var Etotal = Lr + Lb*Pre
	var Vlen = EVbase * EVbase * Lb * Lb * Pre * (1 - Pre)
	var D = m * m * (Etotal*Dbase + Vlen)
	if sp.IsFG() {
		fmt.Fprintf(w, "symbols: µ = %.8g%%, sigma(sym) = %.6g\n", EVbase*100, math.Sqrt(Dbase))
		fmt.Fprintf(w, "free: HRfg = 1/%.5g, E[K] = %.5g, E[K2] = %.5g\n", 1/HRfg, EK, EK2)
		fmt.Fprintf(w, "RTP = %.5g*(1+%.5g) = %.8g%%\n", EVbase*100, Pre*Lb/Lr, EV*100)
	}
	Print_vi(w, sp, D)
	Print_ci(w, sp, EV, D)
	return EVbase, Dbase
}

// Parsheet for model utilizing two independent reel sets under the constraint
// that Free Spins can re-trigger at most once during a bonus feature.
func Parsheet_dualreelonce(w io.Writer, sp *ScanPar, sr Simulator, sb Simulator, cost float64, Lr, Lb float64) (float64, float64) {
	const m = 1
	// regular reels parameters
	var EVr, Dr = sr.EvD(cost)
	var HRfg, EKr, EK2r = sr.Efs()
	var VarKr = EK2r - EKr*EKr
	// bonus reels parameters
	var EVb, Db = sb.EvD(cost)
	var pb, _, _ = sb.Efs()
	// aggregate parameters
	var Pre = 1 - math.Pow(1-pb, Lr)
	var q = EKr + HRfg*Lb*Pre
	var VarLex = HRfg*Lb*Lb*Pre*(1-Pre) + Pre*Pre*Lb*Lb*HRfg*(1-HRfg)
	var CovKrLex = Lr * Lb * Pre * HRfg * (1 - HRfg)
	var VarNfg = VarKr + 2*CovKrLex + VarLex
	var EV = EVr + q*m*EVb
	var D = Dr + m*m*(EVb*EVb*VarNfg+q*Db)
	if sp.IsMain() {
		fmt.Fprintf(w, "µ = %.8g%%, sigma(base) = %.6g\n", EVr*100, math.Sqrt(Dr))
		fmt.Fprintf(w, "free: HRfg = 1/%.5g, E[K] = %.5g, q = %.5g\n", 1/HRfg, EKr, q)
		fmt.Fprintf(w, "Pb = %.6g, Pre = %.6g, E[Lex] = %.5g\n", pb, Pre, Lb*Pre)
		fmt.Fprintf(w, "RTP = %.5g(r) + %.5g*%.5g(b) = %.8g%%\n", EVr*100, q, EVb*100, EV*100)
	}
	Print_all(w, sp, sr, EV, D)
	return EV, D
}
