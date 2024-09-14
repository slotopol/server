package jewels

import (
	"math"

	slot "github.com/slotopol/server/game/slot"
	"github.com/slotopol/server/util"
)

// reels lengths [27, 27, 27, 27, 27], total reshuffles 14348907
// RTP = 88.89513326694501%
var Reels89 = slot.Reels5x{
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 6, 6, 4, 4, 4, 4, 7, 7, 5, 5, 5, 5, 6, 6, 6, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 6, 6, 4, 4, 4, 4, 7, 7, 5, 5, 5, 5, 6, 6, 6, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 6, 6, 4, 4, 4, 4, 7, 7, 5, 5, 5, 5, 6, 6, 6, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 6, 6, 4, 4, 4, 4, 7, 7, 5, 5, 5, 5, 6, 6, 6, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 6, 6, 4, 4, 4, 4, 7, 7, 5, 5, 5, 5, 6, 6, 6, 7, 7, 7},
}

// reels lengths [27, 27, 25, 27, 27], total reshuffles 13286025
// RTP = 89.84989867172462%
var Reels90 = slot.Reels5x{
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 6, 6, 7, 7, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 6, 6, 7, 7, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 6, 6, 7, 7, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 6, 6, 7, 7, 7, 7, 7},
}

// reels lengths [27, 27, 21, 27, 27], total reshuffles 11160261
// RTP = 91.00109755497654%
var Reels91 = slot.Reels5x{
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 6, 6, 4, 4, 4, 4, 7, 7, 5, 5, 5, 5, 6, 6, 6, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 6, 6, 4, 4, 4, 4, 7, 7, 5, 5, 5, 5, 6, 6, 6, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 5, 5, 5, 6, 6, 6, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 6, 6, 4, 4, 4, 4, 7, 7, 5, 5, 5, 5, 6, 6, 6, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 6, 6, 4, 4, 4, 4, 7, 7, 5, 5, 5, 5, 6, 6, 6, 7, 7, 7},
}

// reels lengths [25, 25, 27, 25, 25], total reshuffles 10546875
// RTP = 92.98564740740741%
var Reels93 = slot.Reels5x{
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 6, 6, 7, 7, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7},
}

// reels lengths [27, 21, 27, 21, 27], total reshuffles 8680203
// RTP = 95.01621102640111%
var Reels95 = slot.Reels5x{
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 6, 6, 4, 4, 4, 4, 7, 7, 5, 5, 5, 5, 6, 6, 6, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 5, 5, 5, 6, 6, 6, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 6, 6, 4, 4, 4, 4, 7, 7, 5, 5, 5, 5, 6, 6, 6, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 5, 5, 5, 6, 6, 6, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 6, 6, 4, 4, 4, 4, 7, 7, 5, 5, 5, 5, 6, 6, 6, 7, 7, 7},
}

// reels lengths [25, 25, 25, 25, 25], total reshuffles 9765625
// RTP = 96.0811008%
var Reels96 = slot.Reels5x{
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7},
}

// reels lengths [21, 27, 21, 27, 21], total reshuffles 6751269
// RTP = 97.61364863405679%
var Reels98 = slot.Reels5x{
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 5, 5, 5, 6, 6, 6, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 6, 6, 4, 4, 4, 4, 7, 7, 5, 5, 5, 5, 6, 6, 6, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 5, 5, 5, 6, 6, 6, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 6, 6, 4, 4, 4, 4, 7, 7, 5, 5, 5, 5, 6, 6, 6, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 5, 5, 5, 6, 6, 6, 7, 7, 7},
}

// reels lengths [23, 23, 23, 23, 23], total reshuffles 6436343
// RTP = 99.78647812896236%
var Reels100 = slot.Reels5x{
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7},
}

// reels lengths [21, 21, 21, 21, 21], total reshuffles 4084101
// RTP = 117.80805616707323%
var Reels118 = slot.Reels5x{
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 5, 5, 5, 6, 6, 6, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 5, 5, 5, 6, 6, 6, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 5, 5, 5, 6, 6, 6, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 5, 5, 5, 6, 6, 6, 7, 7, 7},
	{1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 4, 5, 5, 5, 6, 6, 6, 7, 7, 7},
}

// Map with available reels.
var ReelsMap = map[float64]*slot.Reels5x{
	88.89513326694501:  &Reels89,
	89.84989867172462:  &Reels90,
	91.00109755497654:  &Reels91,
	92.98564740740741:  &Reels93,
	95.01621102640111:  &Reels95,
	96.0811008:         &Reels96,
	97.61364863405679:  &Reels98,
	99.78647812896236:  &Reels100,
	117.80805616707323: &Reels118,
}

func FindReels(mrtp float64) (rtp float64, reels slot.Reels) {
	for p, r := range ReelsMap {
		if math.Abs(mrtp-p) < math.Abs(mrtp-rtp) {
			rtp, reels = p, r
		}
	}
	return
}

// Lined payment.
var LinePay = [7][5]float64{
	{0, 0, 20, 200, 2000}, // 1 crown
	{0, 0, 15, 100, 500},  // 2 gold
	{0, 0, 15, 100, 500},  // 3 money
	{0, 0, 10, 50, 200},   // 4 ruby
	{0, 0, 10, 50, 200},   // 5 sapphire
	{0, 0, 5, 25, 100},    // 6 emerald
	{0, 0, 5, 25, 100},    // 7 amethyst
}

const (
	jid = 1 // jackpot ID
)

// Jackpot win combinations.
var Jackpot = [7][5]int{
	{0, 0, 0, 0, 0}, //  1 crown
	{0, 0, 0, 0, 0}, //  2 gold
	{0, 0, 0, 0, 0}, //  3 money
	{0, 0, 0, 0, 0}, //  4 ruby
	{0, 0, 0, 0, 0}, //  5 sapphire
	{0, 0, 0, 0, 0}, //  6 emerald
	{0, 0, 0, 0, 0}, //  7 amethyst
}

type Game struct {
	slot.Slot5x3 `yaml:",inline"`
}

func NewGame() *Game {
	return &Game{
		Slot5x3: slot.Slot5x3{
			SBL: util.MakeBitNum(5, 1),
			Bet: 1,
		},
	}
}

var bl = slot.BetLinesNvm10

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	g.ScanLined(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen slot.Screen, wins *slot.Wins) {
	for li := range g.SBL.Bits() {
		var line = bl.Line(li)

		var syml = screen.Pos(3, line)
		var xy = slot.NewLine5x()
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
			*wins = append(*wins, slot.WinItem{
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

func (g *Game) Spin(screen slot.Screen, mrtp float64) {
	var _, reels = FindReels(mrtp)
	screen.Spin(reels)
}

func (g *Game) SetLines(sbl slot.Bitset) error {
	var mask slot.Bitset = (1<<len(bl) - 1) << 1
	if sbl == 0 {
		return slot.ErrNoLineset
	}
	if sbl&^mask != 0 {
		return slot.ErrLinesetOut
	}
	if g.FreeSpins() > 0 {
		return slot.ErrNoFeature
	}
	g.SBL = sbl
	return nil
}
