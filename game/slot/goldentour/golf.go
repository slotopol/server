package goldentour

import "math/rand/v2"

// Average gain = 16
var Golf = [...]float64{
	5, 6, 7, 8, 9, 10, 11, 12, 14, 15, 16, 18, 20, 25, 30, 50, // E = 16
	// 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 20, 25, 40, 50, 100, // E = 20
}

func GolfSpawn(totalbet float64) (pay float64) {
	var mult = Golf[rand.N(len(Golf))]
	return totalbet * mult
}
