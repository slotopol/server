package neonbananas

import "github.com/slotopol/server/game/slot"

// Win multipliers for Lucky Slot
// (this reel helps to increase frequency of Lucky Slot bonuses)
var LuckySlotMult = [...]float64{
	10, 10, 10, 10, 10, 25, 25, 25, 50, 50, 50, 100, 200, 300, 1000,
}

// Expected value of LuckySlotMult
const Els = 125

func LuckySlotSpawn(rng *slot.Rand, bet float64, spins int) (any, float64) {
	var cash float64
	var res = make([]float64, spins)
	for i := range res {
		res[i] = LuckySlotMult[rng.N(len(LuckySlotMult))]
		cash += res[i]
	}
	return res, bet * cash
}
