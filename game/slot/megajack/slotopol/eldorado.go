package slotopol

import "github.com/slotopol/server/game/slot"

// Average gain = 106
var Eldorado = []float64{
	1000, // 1
	10,   // 2
	25,   // 3
	100,  // 4
	10,   // 5
	50,   // 6
	300,  // 7
	10,   // 8
	50,   // 9
	200,  // 10
	10,   // 11
	25,   // 12
	100,  // 13
	25,   // 14
	10,   // 15
	100,  // 16
	10,   // 17
	50,   // 18
	10,   // 19
	25,   // 20
}

type WinElSeg struct {
	Pos  int     `json:"pos" yaml:"pos" xml:"pos,attr"`    // segment number, starts from 1
	Pay  float64 `json:"pay" yaml:"pay" xml:"pay,attr"`    // pay by this segment
	Mult float64 `json:"mult" yaml:"mult" xml:"mult,attr"` // bet multiplier
}

func EldoradoSpawn(rng *slot.Rand, bet float64, spins int) (any, float64) {
	var res = make([]WinElSeg, spins)
	var pay float64
	for i := range res {
		var n = rng.N(len(Eldorado))
		res[i].Pos = n + 1
		res[i].Pay = bet * Eldorado[n]
		res[i].Mult = Eldorado[n]
		pay += res[i].Pay
	}
	return res, pay
}
