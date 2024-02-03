package twomillionbc

import "math/rand"

// Average gain = 11 (on 3 acorns)
var Acorn = [...]int{
	1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 3, 3, 3, 3, 4, 4, 4, 5, 6, 7, 8, 10, 12,
}

func AcornSpawn(acornbet int) (pay int) {
	var mult = Acorn[rand.Intn(len(Acorn))]
	return acornbet * mult
}
