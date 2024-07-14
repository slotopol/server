package fortuneteller

const (
	type1 = 1
	type2 = 2
	type3 = 3
	type4 = 4
)

func CardsWin(c1, c2, c3 int) float64 {
	if c1 == type1 && c2 == type1 && c3 == type1 {
		return 1000
	}
	if c1 == type2 && c2 == type2 && c3 == type2 {
		return 200
	}
	if c1 == type3 && c2 == type3 && c3 == type3 {
		return 50
	}
	if c1 == type4 && c2 == type4 && c3 == type4 {
		return 20
	}
	if (c1 == type1 && c2 == type1) || (c2 == type1 && c3 == type1) || (c1 == type1 && c3 == type1) {
		return 100
	}
	if (c1 == type2 && c2 == type2) || (c2 == type2 && c3 == type2) || (c1 == type2 && c3 == type2) {
		return 20
	}
	if (c1 == type3 && c2 == type3) || (c2 == type3 && c3 == type3) || (c1 == type3 && c3 == type3) {
		return 10
	}
	if (c1 == type4 && c2 == type4) || (c2 == type4 && c3 == type4) || (c1 == type4 && c3 == type4) {
		return 5
	}
	return 5
}
