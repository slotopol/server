package greatblue

// See: https://freeslotshub.com/playtech/great-blue/
// See: https://www.slotsmate.com/software/playtech/great-blue

import (
	"math/rand/v2"

	"github.com/slotopol/server/game/slot"
)

var ReelsMap slot.ReelsMap[*slot.Reels5x]

// Lined payment.
var LinePay = [13][5]float64{
	{0, 10, 250, 2500, 10000}, //  1 wild
	{0, 2, 25, 125, 750},      //  2 dolphin
	{0, 2, 25, 125, 750},      //  3 turtle
	{0, 0, 20, 100, 400},      //  4 fish
	{0, 0, 15, 75, 250},       //  5 seahorse
	{0, 0, 15, 75, 250},       //  6 starfish
	{0, 0, 10, 50, 150},       //  7 ace
	{0, 0, 10, 50, 150},       //  8 king
	{0, 0, 5, 25, 100},        //  9 queen
	{0, 0, 5, 25, 100},        // 10 jack
	{0, 0, 5, 25, 100},        // 11 ten
	{0, 2, 5, 25, 100},        // 12 nine
	{},                        // 13 scatter
}

// Scatters payment.
var ScatPay = [5]float64{0, 2, 5, 20, 500} // 13 scatter

// Bet lines
var BetLines = slot.BetLinesPlt5x3[:25]

const (
	shell_x5   = 1
	shell_x8   = 2
	shell_fs7  = 3
	shell_fs10 = 4
	shell_fs15 = 5
)

type Seashells struct {
	Sel1 int     `json:"sel1" yaml:"sel1" xml:"sel1"`
	Sel2 int     `json:"sel2" yaml:"sel2" xml:"sel2"`
	Mult float64 `json:"mult" yaml:"mult" xml:"mult"`
	FS   int     `json:"fs" yaml:"fs" xml:"fs"`
}

func (s *Seashells) SetupShell(shell int) {
	switch shell {
	case shell_x5:
		s.Mult += 5
	case shell_x8:
		s.Mult += 8
	case shell_fs7:
		s.FS += 7
	case shell_fs10:
		s.FS += 10
	case shell_fs15:
		s.FS += 15
	}
}

type Game struct {
	slot.Screen5x3 `yaml:",inline"`
	slot.Slotx     `yaml:",inline"`
	// multiplier on freespins
	M float64 `json:"m,omitempty" yaml:"m,omitempty" xml:"m,omitempty"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slotx: slot.Slotx{
			Sel: len(BetLines),
			Bet: 1,
		},
		M: 0,
	}
}

func (g *Game) Clone() slot.SlotGame {
	var clone = *g
	return &clone
}

const wild, scat = 1, 13

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	for li, line := range BetLines[:g.Sel] {
		var mw float64 = 1 // mult wild
		var numw, numl slot.Pos = 0, 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = g.LY(x, line)
			if sx == wild {
				if syml == 0 {
					numw = x
				}
				mw = 2
			} else if syml == 0 {
				syml = sx
			} else if sx != syml {
				numl = x - 1
				break
			}
		}

		var payw, payl float64
		if numw >= 2 {
			payw = LinePay[wild-1][numw-1]
		}
		if numl >= 2 && syml > 0 {
			payl = LinePay[syml-1][numl-1]
		}
		if payl*mw > payw {
			var mm float64 = 1 // mult mode
			if g.FSR > 0 {
				mm = g.M
			}
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * payl,
				MP:  mw * mm,
				Sym: syml,
				Num: numl,
				LI:  li + 1,
				XY:  line.HitxL(numl),
			})
		} else if payw > 0 {
			var mm float64 = 1 // mult mode
			if g.FSR > 0 {
				mm = g.M
			}
			*wins = append(*wins, slot.WinItem{
				Pay: g.Bet * payw,
				MP:  mm,
				Sym: wild,
				Num: numw,
				LI:  li + 1,
				XY:  line.HitxL(numw),
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.ScatNum(scat); count >= 2 {
		var pay, fs = ScatPay[count-1], 0
		var mm float64 = 1 // mult mode
		if g.FSR > 0 {
			mm, fs = g.M, 15
		} else if count >= 3 {
			fs = 8
		}
		*wins = append(*wins, slot.WinItem{
			Pay: g.Bet * float64(g.Sel) * pay,
			MP:  mm,
			Sym: scat,
			Num: count,
			XY:  g.ScatPos(scat),
			FS:  fs,
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.ReelSpin(reels)
}

func (g *Game) Spawn(wins slot.Wins, fund, mrtp float64) {
	if g.FSR > 0 {
		return
	}
	for i := range wins {
		if wi := &wins[i]; wi.Sym == scat {
			var idx = []int{shell_x5, shell_x8, shell_fs7, shell_fs10, shell_fs15}
			rand.Shuffle(len(idx), func(i, j int) {
				idx[i], idx[j] = idx[j], idx[i]
			})
			var bon = Seashells{
				Sel1: idx[0],
				Sel2: idx[1],
				Mult: 2,
				FS:   8,
			}
			bon.SetupShell(idx[0])
			bon.SetupShell(idx[1])
			wi.MP = 1
			wi.FS = bon.FS
			wi.Bon = bon
		}
	}
}

func (g *Game) Prepare() {
	if g.FSR == 0 {
		g.M = 0
	}
}

func (g *Game) Apply(wins slot.Wins) {
	if g.FSR != 0 {
		g.Gain += wins.Gain()
		g.FSN++
	} else {
		g.Gain = wins.Gain()
		g.FSN = 0
	}

	if g.FSR > 0 {
		g.FSR--
	}
	for _, wi := range wins {
		if wi.Sym == scat {
			if g.M > 0 {
				g.FSR += wi.FS
			} else {
				var bon = wi.Bon.(Seashells)
				g.FSR = bon.FS
				g.M = bon.Mult
			}
		}
	}
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
