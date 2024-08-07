package justjewels

import "github.com/slotopol/server/game"

// reels lengths [39, 39, 39, 39, 39], total reshuffles 90224199
// RTP = 114.75(lined) + 8.0152(scatter) = 122.764204%
var Reels123 = game.Reels5x{
	{1, 1, 1, 6, 6, 6, 2, 2, 2, 2, 5, 5, 5, 3, 3, 3, 3, 7, 7, 7, 8, 4, 4, 4, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7, 8, 4, 4, 4},
	{1, 1, 1, 6, 6, 6, 2, 2, 2, 2, 5, 5, 5, 3, 3, 3, 3, 7, 7, 7, 8, 4, 4, 4, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7, 8, 4, 4, 4},
	{1, 1, 1, 6, 6, 6, 2, 2, 2, 2, 5, 5, 5, 3, 3, 3, 3, 7, 7, 7, 8, 4, 4, 4, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7, 8, 4, 4, 4},
	{1, 1, 1, 6, 6, 6, 2, 2, 2, 2, 5, 5, 5, 3, 3, 3, 3, 7, 7, 7, 8, 4, 4, 4, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7, 8, 4, 4, 4},
	{1, 1, 1, 6, 6, 6, 2, 2, 2, 2, 5, 5, 5, 3, 3, 3, 3, 7, 7, 7, 8, 4, 4, 4, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7, 8, 4, 4, 4},
}

// Map with available reels.
var ReelsMap = map[string]*game.Reels5x{
	"123": &Reels123, // minimum possible percentage
}

// Lined payment.
var LinePay = [8][5]float64{
	{0, 0, 50, 500, 5000}, // 1 crown
	{0, 0, 30, 150, 500},  // 2 gold
	{0, 0, 30, 150, 500},  // 3 money
	{0, 0, 15, 50, 200},   // 4 ruby
	{0, 0, 15, 50, 200},   // 5 sapphire
	{0, 0, 10, 25, 150},   // 6 emerald
	{0, 0, 10, 25, 150},   // 7 amethyst
	{0, 0, 0, 0, 0},       // 8 euro
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 2, 10, 50} // 8 euro

const (
	jid = 1 // jackpot ID
)

// Jackpot win combinations.
var Jackpot = [8][5]int{
	{0, 0, 0, 0, 0}, //  1 crown
	{0, 0, 0, 0, 0}, //  2 gold
	{0, 0, 0, 0, 0}, //  3 money
	{0, 0, 0, 0, 0}, //  4 ruby
	{0, 0, 0, 0, 0}, //  5 sapphire
	{0, 0, 0, 0, 0}, //  6 emerald
	{0, 0, 0, 0, 0}, //  7 amethyst
	{0, 0, 0, 0, 0}, //  8 euro
}

type Game struct {
	game.Slot5x3 `yaml:",inline"`
}

func NewGame(rd string) *Game {
	return &Game{
		Slot5x3: game.Slot5x3{
			RD:  rd,
			SBL: game.MakeBitNum(5),
			Bet: 1,
		},
	}
}

const scat = 8

var bl = game.BetLinesNvm10

func (g *Game) Scanner(screen game.Screen, wins *game.Wins) {
	g.ScanLined(screen, wins)
	g.ScanScatters(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen game.Screen, wins *game.Wins) {
	for li := g.SBL.Next(0); li != 0; li = g.SBL.Next(li) {
		var line = bl.Line(li)

		var syml = screen.Pos(3, line)
		var xy = game.NewLine5x()
		var numl = 1
		xy.Set(3, line.At(3))
		if screen.Pos(2, line) == syml {
			xy.Set(2, line.At(2))
			numl++
			if screen.Pos(1, line) == syml {
				xy.Set(1, line.At(1))
				numl++
			}
		}
		if screen.Pos(4, line) == syml {
			xy.Set(4, line.At(4))
			numl++
			if screen.Pos(5, line) == syml {
				xy.Set(5, line.At(5))
				numl++
			}
		}

		if numl >= 3 {
			*wins = append(*wins, game.WinItem{
				Pay:  g.Bet * LinePay[syml-1][numl-1],
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   xy,
			})
		} else {
			xy.Free()
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen game.Screen, wins *game.Wins) {
	if count := screen.ScatNum(scat); count >= 3 {
		var pay = ScatPay[count-1]
		*wins = append(*wins, game.WinItem{
			Pay:  g.Bet * float64(g.SBL.Num()) * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
		})
	}
}

func (g *Game) Spin(screen game.Screen) {
	screen.Spin(ReelsMap[g.RD])
}

func (g *Game) SetLines(sbl game.Bitset) error {
	var mask game.Bitset = (1<<len(bl) - 1) << 1
	if sbl == 0 {
		return game.ErrNoLineset
	}
	if sbl&^mask != 0 {
		return game.ErrLinesetOut
	}
	if g.FreeSpins() > 0 {
		return game.ErrNoFeature
	}
	g.SBL = sbl
	return nil
}

func (g *Game) SetReels(rd string) error {
	if _, ok := ReelsMap[rd]; !ok {
		return game.ErrNoReels
	}
	g.RD = rd
	return nil
}
