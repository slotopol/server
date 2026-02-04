package slot

import "math"

// Quantile (for Volatility Index)
func GetZ(confidence float64) float64 {
	// probability for one side (p = 1 - alpha/2)
	p := (1.0 + confidence) / 2.0
	// formula using the inverse error function is: Z = sqrt(2) * ErfInv(2p - 1)
	return math.Sqrt(2) * math.Erfinv(2*p-1)
}

// Volatility Index Class with 3 gradations
func VIclass3(vi float64) int {
	switch {
	case vi < 10:
		return 1
	case vi < 25:
		return 2
	default:
		return 3
	}
}

// Volatility Index Class with 6 gradations
func VIclass6(vi float64) int {
	switch {
	case vi < 7:
		return 1
	case vi < 12:
		return 2
	case vi < 18:
		return 3
	case vi < 25:
		return 4
	case vi < 45:
		return 5
	default:
		return 6
	}
}

var VIname3 = map[int]string{
	1: "Low",
	2: "Medium",
	3: "High",
}

var VIname6 = map[int]string{
	1: "Low",
	2: "Medium-Low",
	3: "Medium",
	4: "Medium-High",
	5: "High",
	6: "Very High",
}

// Index of Convergence. The number of spins after which,
// with a given confidence, the player will not be in profit on
// this particular slot machine.
func CI(confidence, rtp, sigma float64) float64 {
	var nh = GetZ(confidence) * sigma / (1 - rtp)
	return nh * nh
}

// Bankroll management formula to protect user against ruin for N spins.
func BankrollUser(confidence, rtp, sigma, N float64) float64 {
	return GetZ(confidence)*sigma*math.Sqrt(N) + N*(1-rtp)
}

type GameGroup struct {
	K     int     // number of players
	RTP   float64 // reels RTP
	Sigma float64 // volatility (sigma)
}

// Minimum bankroll requirement to ensure long-term
// payment capabilities of the bank.
func Bankroll(confidence float64, groups []GameGroup) float64 {
	var totalVariance float64
	var totalHouseEdge float64
	for _, g := range groups {
		totalVariance += float64(g.K) * g.Sigma * g.Sigma
		totalHouseEdge += float64(g.K) * (1 - g.RTP)
	}
	return (totalVariance / (2 * totalHouseEdge)) * math.Log(1/(1-confidence))
}
