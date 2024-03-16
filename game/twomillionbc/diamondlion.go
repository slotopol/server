package twomillionbc

import "math/rand/v2"

// Average gain = 175
var DiamondLion = [...]float64{
	50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150, 160, 170, 180, 190, 200, 210, 220, 230, 240, 250, 260, 270, 280, 290, 300,
}

func DiamondLionSpawn(bet float64) float64 {
	var mult = DiamondLion[rand.N(len(DiamondLion))]
	return bet * mult
}
