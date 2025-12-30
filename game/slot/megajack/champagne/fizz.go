package champagne

import (
	"math/rand/v2"
)

// len = 5, count = 10, avr bottle gain = 68, M = 168
var Bottles = [...]float64{
	10, 10, 20, 150, 150,
}

type WinBottle struct {
	Mult float64 `json:"mult" yaml:"mult" xml:"mult,attr"` // bet multiplier
	Pay  float64 `json:"pay" yaml:"pay" xml:"pay,attr"`    // pay by this cell
}

func ChampagneSpawn(bet float64) (any, float64) {
	var res [5]WinBottle
	var cash float64

	var p = Bottles
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
