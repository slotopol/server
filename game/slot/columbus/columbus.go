package columbus

import (
	"math"

	slot "github.com/slotopol/server/game/slot"
	"github.com/slotopol/server/util"
)

// reels lengths [32, 28, 32, 28, 32], total reshuffles 25690112
// symbols: 60.236(lined) + 0(scatter) = 60.236366%
// free spins 2540160, q = 0.098877, sq = 1/(1-q) = 1.109726
// free games frequency: 1/101.14
// RTP = 60.236(sym) + 0.098877*251.43(fg) = 85.096797%
var ReelsReg85 = slot.Reels5x{
	{1, 5, 7, 9, 8, 2, 10, 6, 4, 10, 6, 3, 4, 8, 9, 10, 7, 8, 9, 7, 8, 9, 3, 8, 9, 7, 5, 9, 7, 2, 5, 9},
	{1, 7, 6, 3, 7, 4, 6, 9, 3, 8, 5, 6, 8, 9, 6, 8, 5, 2, 8, 5, 8, 6, 9, 8, 4, 5, 8, 6},
	{1, 6, 8, 4, 7, 5, 9, 6, 7, 9, 5, 6, 8, 7, 5, 3, 9, 7, 8, 7, 6, 9, 7, 5, 8, 4, 10, 7, 9, 2, 10, 9},
	{1, 6, 9, 8, 6, 7, 9, 4, 6, 8, 9, 6, 8, 7, 3, 6, 9, 5, 7, 6, 8, 9, 2, 9, 7, 6, 5, 8},
	{10, 6, 5, 9, 6, 10, 5, 6, 4, 1, 8, 6, 1, 8, 9, 3, 7, 2, 4, 3, 9, 4, 2, 9, 3, 4, 5, 3, 9, 4, 7, 5},
}

// reels lengths [32, 28, 32, 28, 32], total reshuffles 25690112
// symbols: 63.094(lined) + 0(scatter) = 63.094159%
// free spins 2540160, q = 0.098877, sq = 1/(1-q) = 1.109726
// free games frequency: 1/101.14
// RTP = 63.094(sym) + 0.098877*251.43(fg) = 87.954589%
var ReelsReg88 = slot.Reels5x{
	{1, 5, 7, 9, 8, 2, 10, 6, 4, 10, 6, 3, 4, 8, 9, 10, 7, 8, 9, 7, 8, 9, 3, 8, 9, 7, 5, 9, 7, 2, 5, 9},
	{1, 7, 6, 3, 7, 4, 6, 9, 3, 8, 5, 6, 8, 9, 6, 8, 6, 2, 8, 5, 8, 6, 9, 8, 4, 5, 8, 6},
	{1, 6, 8, 4, 7, 5, 9, 6, 7, 9, 5, 6, 8, 7, 5, 3, 9, 7, 8, 7, 5, 9, 7, 5, 8, 4, 10, 7, 9, 2, 10, 9},
	{1, 6, 9, 8, 6, 7, 9, 4, 6, 8, 1, 6, 8, 7, 3, 6, 9, 5, 7, 6, 8, 9, 2, 9, 7, 6, 5, 8},
	{10, 6, 5, 9, 6, 10, 5, 6, 4, 1, 8, 6, 1, 8, 9, 3, 7, 2, 4, 3, 9, 4, 2, 9, 3, 4, 5, 3, 9, 4, 7, 5},
}

// reels lengths [32, 28, 32, 28, 32], total reshuffles 25690112
// symbols: 65.101(lined) + 0(scatter) = 65.101001%
// free spins 2540160, q = 0.098877, sq = 1/(1-q) = 1.109726
// free games frequency: 1/101.14
// RTP = 65.101(sym) + 0.098877*251.43(fg) = 89.961431%
var ReelsReg90 = slot.Reels5x{
	{1, 5, 7, 9, 8, 2, 10, 6, 4, 10, 6, 3, 4, 8, 9, 10, 7, 8, 9, 7, 8, 9, 3, 8, 9, 7, 5, 9, 7, 2, 5, 9},
	{1, 7, 6, 3, 7, 4, 6, 9, 3, 8, 5, 6, 8, 9, 6, 8, 6, 2, 8, 5, 8, 6, 9, 8, 4, 5, 8, 6},
	{1, 6, 8, 4, 7, 5, 9, 6, 7, 9, 5, 6, 8, 7, 5, 3, 9, 7, 8, 7, 5, 9, 7, 5, 8, 4, 10, 7, 9, 2, 10, 9},
	{1, 6, 9, 5, 6, 7, 9, 4, 6, 8, 1, 6, 8, 7, 3, 6, 1, 5, 7, 6, 8, 9, 2, 9, 7, 6, 5, 8},
	{10, 6, 5, 9, 6, 10, 5, 6, 4, 9, 8, 6, 1, 8, 9, 3, 7, 2, 4, 3, 9, 4, 2, 9, 3, 4, 8, 3, 9, 4, 7, 5},
}

// reels lengths [32, 28, 32, 28, 32], total reshuffles 25690112
// symbols: 67.18(lined) + 0(scatter) = 67.180478%
// free spins 2540160, q = 0.098877, sq = 1/(1-q) = 1.109726
// free games frequency: 1/101.14
// RTP = 67.18(sym) + 0.098877*251.43(fg) = 92.040908%
var ReelsReg92 = slot.Reels5x{
	{1, 5, 7, 9, 8, 2, 10, 6, 4, 10, 6, 3, 4, 8, 9, 10, 7, 8, 9, 7, 8, 9, 3, 8, 9, 7, 5, 9, 7, 2, 5, 9},
	{1, 7, 6, 3, 7, 4, 6, 9, 3, 8, 5, 6, 8, 9, 6, 8, 5, 2, 8, 5, 8, 6, 9, 8, 4, 5, 8, 6},
	{1, 6, 8, 4, 7, 5, 9, 6, 7, 9, 5, 6, 8, 7, 5, 3, 9, 6, 8, 7, 6, 9, 7, 5, 8, 4, 10, 7, 9, 2, 10, 9},
	{1, 6, 9, 8, 6, 7, 9, 4, 6, 8, 9, 5, 8, 7, 3, 6, 9, 5, 7, 1, 8, 9, 2, 9, 7, 6, 5, 8},
	{10, 6, 1, 9, 6, 10, 5, 6, 4, 1, 8, 6, 1, 8, 9, 3, 7, 2, 6, 3, 9, 4, 2, 9, 3, 4, 5, 3, 9, 4, 7, 5},
}

// reels lengths [32, 28, 32, 28, 32], total reshuffles 25690112
// symbols: 69.142(lined) + 0(scatter) = 69.142322%
// free spins 2540160, q = 0.098877, sq = 1/(1-q) = 1.109726
// free games frequency: 1/101.14
// RTP = 69.142(sym) + 0.098877*251.43(fg) = 94.002752%
var ReelsReg94 = slot.Reels5x{
	{1, 5, 7, 9, 8, 2, 10, 6, 4, 10, 6, 3, 4, 8, 9, 10, 7, 8, 9, 7, 8, 9, 3, 8, 9, 7, 5, 9, 7, 2, 5, 9},
	{1, 7, 6, 3, 7, 4, 6, 9, 3, 8, 5, 6, 8, 9, 6, 8, 5, 2, 8, 5, 8, 6, 9, 8, 4, 5, 8, 6},
	{1, 6, 8, 4, 7, 5, 9, 6, 7, 9, 5, 6, 8, 7, 5, 3, 9, 6, 8, 7, 6, 9, 7, 5, 8, 4, 10, 7, 9, 2, 10, 9},
	{1, 6, 9, 8, 1, 7, 9, 4, 6, 8, 9, 3, 5, 2, 6, 3, 9, 5, 7, 1, 8, 9, 2, 9, 7, 4, 5, 8},
	{10, 6, 8, 9, 6, 10, 5, 6, 4, 1, 8, 7, 1, 8, 9, 3, 7, 2, 6, 8, 9, 4, 7, 9, 3, 8, 5, 3, 9, 4, 7, 5},
}

// reels lengths [32, 28, 32, 28, 32], total reshuffles 25690112
// symbols: 70.162(lined) + 0(scatter) = 70.162364%
// free spins 2540160, q = 0.098877, sq = 1/(1-q) = 1.109726
// free games frequency: 1/101.14
// RTP = 70.162(sym) + 0.098877*251.43(fg) = 95.022795%
var ReelsReg95 = slot.Reels5x{
	{1, 5, 7, 9, 8, 2, 10, 6, 4, 10, 6, 3, 4, 8, 9, 10, 7, 8, 9, 7, 8, 9, 3, 8, 9, 7, 5, 9, 7, 2, 5, 9},
	{1, 7, 6, 3, 7, 4, 6, 9, 3, 8, 5, 6, 2, 9, 6, 8, 5, 2, 8, 5, 8, 6, 9, 8, 4, 5, 8, 6},
	{1, 6, 8, 4, 7, 5, 9, 6, 7, 9, 5, 6, 8, 7, 5, 3, 9, 6, 8, 7, 6, 9, 7, 5, 8, 4, 10, 7, 9, 2, 10, 9},
	{1, 6, 4, 8, 1, 7, 9, 4, 6, 8, 9, 3, 5, 2, 6, 3, 9, 5, 7, 1, 8, 9, 2, 5, 7, 4, 5, 8},
	{10, 6, 8, 9, 6, 10, 5, 6, 4, 1, 8, 7, 1, 8, 9, 3, 7, 2, 6, 8, 9, 4, 7, 9, 3, 7, 5, 3, 9, 4, 7, 5},
}

// reels lengths [32, 28, 32, 28, 32], total reshuffles 25690112
// symbols: 71.226(lined) + 0(scatter) = 71.225614%
// free spins 2540160, q = 0.098877, sq = 1/(1-q) = 1.109726
// free games frequency: 1/101.14
// RTP = 71.226(sym) + 0.098877*251.43(fg) = 96.086044%
var ReelsReg96 = slot.Reels5x{
	{1, 5, 7, 9, 8, 2, 10, 6, 4, 10, 6, 3, 4, 8, 9, 10, 7, 8, 9, 7, 8, 9, 3, 8, 9, 7, 5, 9, 7, 2, 5, 9},
	{1, 7, 6, 3, 7, 4, 6, 9, 3, 8, 5, 6, 2, 9, 6, 8, 5, 2, 8, 5, 8, 6, 9, 8, 4, 5, 8, 6},
	{1, 6, 8, 4, 7, 5, 9, 6, 7, 9, 5, 6, 8, 7, 5, 3, 9, 6, 8, 7, 6, 9, 7, 5, 8, 4, 10, 7, 9, 2, 10, 9},
	{1, 6, 9, 8, 1, 7, 9, 4, 6, 8, 9, 3, 5, 2, 6, 3, 9, 5, 7, 1, 8, 9, 2, 9, 7, 4, 5, 8},
	{10, 6, 8, 9, 6, 10, 5, 6, 4, 1, 8, 7, 1, 8, 9, 3, 7, 2, 6, 8, 9, 4, 7, 9, 3, 8, 5, 3, 9, 4, 7, 5},
}

// reels lengths [32, 28, 32, 28, 32], total reshuffles 25690112
// symbols: 72.169(lined) + 0(scatter) = 72.168623%
// free spins 2540160, q = 0.098877, sq = 1/(1-q) = 1.109726
// free games frequency: 1/101.14
// RTP = 72.169(sym) + 0.098877*251.43(fg) = 97.029053%
var ReelsReg97 = slot.Reels5x{
	{1, 5, 7, 9, 8, 2, 10, 6, 4, 10, 6, 3, 4, 8, 9, 10, 7, 8, 9, 7, 8, 9, 3, 8, 9, 7, 5, 9, 7, 2, 5, 9},
	{1, 7, 6, 3, 7, 4, 6, 9, 3, 8, 5, 6, 2, 9, 6, 8, 5, 2, 8, 5, 8, 6, 9, 8, 4, 5, 8, 6},
	{1, 6, 8, 4, 7, 5, 9, 6, 7, 9, 5, 6, 8, 7, 5, 3, 9, 6, 8, 7, 6, 9, 7, 5, 8, 4, 10, 7, 9, 2, 10, 9},
	{1, 6, 9, 8, 1, 7, 9, 8, 6, 8, 9, 3, 5, 2, 6, 3, 9, 5, 7, 1, 8, 9, 2, 9, 7, 4, 5, 8},
	{10, 6, 8, 9, 6, 10, 5, 6, 4, 1, 8, 7, 1, 8, 9, 3, 7, 2, 6, 8, 9, 4, 7, 9, 8, 3, 5, 8, 9, 4, 7, 5},
}

// reels lengths [32, 28, 32, 28, 32], total reshuffles 25690112
// symbols: 68.866(lined) + 0(scatter) = 68.865562%
// free spins 7620480, q = 0.29663, sq = 1/(1-q) = 1.421729
// free games frequency: 1/33.712
// RTP = 68.866(sym) + 0.29663*251.43(fg) = 143.446853%
var ReelsReg143 = slot.Reels5x{
	{1, 5, 7, 10, 8, 2, 10, 6, 4, 10, 6, 3, 4, 8, 9, 10, 7, 8, 9, 7, 8, 9, 3, 8, 9, 7, 5, 9, 7, 2, 5, 9},
	{1, 7, 6, 3, 7, 4, 6, 9, 3, 8, 5, 6, 2, 9, 6, 8, 5, 2, 8, 5, 8, 6, 9, 8, 4, 5, 8, 6},
	{1, 6, 8, 4, 7, 5, 9, 6, 10, 9, 5, 6, 8, 7, 5, 3, 9, 6, 8, 7, 6, 9, 7, 5, 8, 4, 10, 7, 9, 2, 10, 9},
	{1, 6, 9, 8, 1, 7, 9, 8, 6, 8, 9, 3, 5, 2, 6, 3, 9, 5, 7, 1, 8, 9, 2, 9, 7, 4, 5, 8},
	{10, 6, 8, 9, 6, 10, 5, 6, 4, 1, 8, 7, 1, 8, 9, 3, 7, 2, 6, 8, 9, 4, 7, 9, 8, 3, 5, 10, 9, 4, 7, 5},
}

// reels lengths [25, 25, 24, 25, 24], total reshuffles 9000000
// symbols: 194.86(lined) + 0(scatter) = 194.856667%
// free spins 2025000, q = 0.225, sq = 1/(1-q) = 1.290323
// free games frequency: 1/44.444
// RTP = sq*rtp(sym) = 1.2903*194.86 = 251.427957%
var ReelsBon = slot.Reels5x{
	{5, 1, 9, 8, 2, 10, 6, 5, 10, 8, 3, 4, 8, 5, 10, 7, 8, 5, 9, 8, 5, 3, 6, 5, 8},
	{7, 9, 3, 7, 4, 6, 9, 7, 8, 9, 7, 8, 9, 7, 8, 5, 2, 8, 5, 4, 9, 6, 7, 9, 1},
	{6, 8, 4, 7, 5, 9, 6, 2, 9, 5, 6, 8, 10, 6, 3, 9, 6, 8, 7, 5, 10, 6, 8, 1},
	{6, 9, 1, 6, 7, 9, 4, 2, 8, 9, 3, 8, 7, 3, 4, 9, 5, 4, 6, 3, 5, 8, 3, 9, 1},
	{7, 4, 9, 6, 10, 9, 6, 4, 9, 8, 7, 1, 8, 9, 3, 7, 2, 9, 3, 8, 10, 4, 5, 2},
}

// Map with available reels.
var ReelsMap = map[float64]*slot.Reels5x{
	85.096797:  &ReelsReg85,
	87.954589:  &ReelsReg88,
	89.961431:  &ReelsReg90,
	92.040908:  &ReelsReg92,
	94.002752:  &ReelsReg94,
	95.022795:  &ReelsReg95,
	96.086044:  &ReelsReg96,
	97.029053:  &ReelsReg97,
	143.446853: &ReelsReg143,
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
var LinePay = [10][5]float64{
	{0, 10, 100, 1000, 5000}, //  1 columbus
	{0, 5, 50, 200, 1000},    //  2 spain
	{0, 5, 25, 100, 500},     //  3 necklace
	{0, 5, 15, 75, 250},      //  4 sextant
	{0, 0, 10, 40, 150},      //  5 ace
	{0, 0, 10, 40, 150},      //  6 king
	{0, 0, 10, 40, 150},      //  7 queen
	{0, 0, 5, 20, 100},       //  8 jack
	{0, 0, 5, 20, 100},       //  9 ten
	{0, 0, 0, 0, 0},          // 10 ship
}

const (
	jid = 1 // jackpot ID
)

// Jackpot win combinations.
var Jackpot = [10][5]int{
	{0, 0, 0, 0, 0}, //  1 columbus
	{0, 0, 0, 0, 0}, //  2 spain
	{0, 0, 0, 0, 0}, //  3 necklace
	{0, 0, 0, 0, 0}, //  4 sextant
	{0, 0, 0, 0, 0}, //  5 ace
	{0, 0, 0, 0, 0}, //  6 king
	{0, 0, 0, 0, 0}, //  7 queen
	{0, 0, 0, 0, 0}, //  8 jack
	{0, 0, 0, 0, 0}, //  9 ten
	{0, 0, 0, 0, 0}, // 10 ship
}

type Game struct {
	slot.Slot5x3 `yaml:",inline"`
	// free spin number
	FS int `json:"fs,omitempty" yaml:"fs,omitempty" xml:"fs,omitempty"`
}

func NewGame() *Game {
	return &Game{
		Slot5x3: slot.Slot5x3{
			SBL: util.MakeBitNum(5, 1),
			Bet: 1,
		},
		FS: 0,
	}
}

const wild, scat = 1, 10

var bl = slot.BetLinesNvm10

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	g.ScanLined(screen, wins)
	g.ScanScatters(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen slot.Screen, wins *slot.Wins) {
	var wbon slot.Sym
	if g.FS > 0 {
		wbon = scat
	}

	for li := range g.SBL.Bits() {
		var line = bl.Line(li)

		var numw, numl = 0, 5
		var syml slot.Sym
		for x := 1; x <= 5; x++ {
			var sx = screen.Pos(x, line)
			if sx == wild || sx == wbon {
				if syml == 0 {
					numw = x
				}
			} else if syml == 0 && sx != scat {
				syml = sx
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		var payw, payl float64
		if numw > 0 {
			payw = LinePay[wild-1][numw-1]
		}
		if numl > 0 && syml > 0 {
			payl = LinePay[syml-1][numl-1]
		}
		if payl > payw {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payl,
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyL(numl),
			})
		} else if payw > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payw,
				Mult: 1,
				Sym:  wild,
				Num:  numw,
				Line: li,
				XY:   line.CopyL(numw),
				Jack: Jackpot[wild-1][numw-1],
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen slot.Screen, wins *slot.Wins) {
	if count := screen.ScatNumOdd(scat); count >= 3 {
		*wins = append(*wins, slot.WinItem{
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPosOdd(scat),
			Free: 10,
		})
	}
}

func (g *Game) Spin(screen slot.Screen, mrtp float64) {
	if g.FS == 0 {
		var _, reels = FindReels(mrtp)
		screen.Spin(reels)
	} else {
		screen.Spin(&ReelsBon)
	}
}

func (g *Game) Apply(screen slot.Screen, wins slot.Wins) {
	if g.FS > 0 {
		g.Gain += wins.Gain()
	} else {
		g.Gain = wins.Gain()
	}

	if g.FS > 0 {
		g.FS--
	}
	for _, wi := range wins {
		if wi.Free > 0 {
			g.FS += wi.Free
		}
	}
}

func (g *Game) FreeSpins() int {
	return g.FS
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
