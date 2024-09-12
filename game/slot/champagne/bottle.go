package champagne

import (
	"math/rand/v2"
	"slices"
)

// len = 36, count = 630, avr bottle gain = 90.555556, M = 193.65079365079
var Bottles = [36]float64{
	10, 10, 10, 10, 10, 10, // 6
	20, 20, 20, 20, 20, 20, // 6
	30, 30, 30, 30, 30, 30, // 6
	50, 50, 50, 50, 50, 50, // 6
	100, 100, 100, 100, // 4
	150, 150, 150, 150, // 4
	300, 300, // 2
	500, 500, // 2
}

type WinBottle struct {
	Mult float64 `json:"mult" yaml:"mult" xml:"mult,attr"` // bet multiplier
	Pay  float64 `json:"pay" yaml:"pay" xml:"pay,attr"`    // pay by this cell
}

func ChampagneSpawn(bet float64) (any, float64) {
	var res [5]WinBottle
	var cash float64

	var p = slices.Clone(Bottles[:])
	rand.Shuffle(len(p), func(i, j int) {
		p[i], p[j] = p[j], p[i]
	})
	for i := range res {
		res[i].Mult = p[i]
		res[i].Pay = bet * p[i]
	}
	cash = p[0] + p[1]
	if p[0] == p[1] {
		cash *= 2
	}
	return res, bet * cash
}
