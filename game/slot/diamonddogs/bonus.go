package diamonddogs

import "math/rand/v2"

var Pays = [...]float64{
	10, 20, 20, 40, 40, 80, 120, 120, 450,
}

// BonusSpawn returns array with shuffled pays, where first element
// always non-zero and others anywhere can be zero (3 zero total).
// On average it returns 3 elements before zero.
// E = 3*M = 3*100 = 300
func BonusSpawn(bet float64) (any, float64) {
	var dogs [12]float64
	copy(dogs[:], Pays[:])
	var first = rand.N(9)
	dogs[0], dogs[first] = dogs[first], dogs[0]
	rand.Shuffle(11, func(i, j int) {
		dogs[i+1], dogs[j+1] = dogs[j+1], dogs[i+1]
	})
	var pay float64
	for _, v := range dogs {
		pay += v
		if v == 0 {
			break
		}
	}
	return dogs, pay
}
