package twomillionbc

import "math/rand/v2"

// Average gain = 11
var Acorn = [...]int{
	4, 5, 6, 7, 8, 10, 10, 10, 12, 15, 20, 25,
}

func AcornSpawn(acornbet int) (pay int) {
	var mult = Acorn[rand.N(len(Acorn))]
	return acornbet * mult
}
