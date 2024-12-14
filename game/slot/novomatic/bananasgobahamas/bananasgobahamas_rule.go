package bananasgobahamas

import (
	_ "embed"

	"github.com/slotopol/server/game/slot"
)

//go:embed bananasgobahamas_bon.yaml
var rbon []byte

var ReelsBon = slot.ReadBon[*slot.Reels5x](rbon)

//go:embed bananasgobahamas_reel.yaml
var reels []byte

var ReelsMap = slot.ReadReelsMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [13][5]float64{
	{0, 10, 250, 2500, 9000}, //  1 banana
	{0, 2, 30, 120, 800},     //  2 strawberry
	{0, 2, 30, 120, 800},     //  3 water
	{0, 0, 20, 100, 400},     //  4 pineapple
	{0, 0, 20, 70, 250},      //  5 mango
	{0, 0, 20, 70, 250},      //  6 coconu
	{0, 0, 10, 50, 120},      //  7 ace
	{0, 0, 10, 50, 120},      //  8 king
	{0, 0, 4, 30, 100},       //  9 queen
	{0, 0, 4, 30, 100},       // 10 jack
	{0, 0, 4, 30, 100},       // 11 ten
	{0, 2, 4, 30, 100},       // 12 nine
	{},                       // 13 suitcase
}

// Scatters payment.
var ScatPay = [5]float64{0, 2, 4, 20, 500} // 13 suitcase

// Scatter freespins table
var ScatFreespin = [5]int{0, 0, 45, 45, 45} // 13 suitcase

const (
	jid = 1 // jackpot ID
)

// Jackpot win combinations.
var Jackpot = [13][5]int{
	{0, 0, 0, 0, 0}, //  1 banana
	{0, 0, 0, 0, 0}, //  2 strawberry
	{0, 0, 0, 0, 0}, //  3 water
	{0, 0, 0, 0, 0}, //  4 pineapple
	{0, 0, 0, 0, 0}, //  5 mango
	{0, 0, 0, 0, 0}, //  6 coconu
	{0, 0, 0, 0, 0}, //  7 ace
	{0, 0, 0, 0, 0}, //  8 king
	{0, 0, 0, 0, 0}, //  9 queen
	{0, 0, 0, 0, 0}, // 10 jack
	{0, 0, 0, 0, 0}, // 11 ten
	{0, 0, 0, 0, 0}, // 12 nine
	{0, 0, 0, 0, 0}, // 13 suitcase
}

// Bet lines
var BetLines = slot.BetLinesNvm10

type Game struct {
	slot.Slot5x3 `yaml:",inline"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slot5x3: slot.Slot5x3{
			Sel: 5,
			Bet: 1,
		},
	}
}

const wild, scat = 1, 13

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	g.ScanLined(screen, wins)
	g.ScanScatters(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen slot.Screen, wins *slot.Wins) {
	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]

		var mw float64 = 1 // mult wild
		var numw, numl slot.Pos = 0, 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = screen.Pos(x, line)
			if sx == wild {
				if syml == 0 {
					numw = x
				}
				mw = 2
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
		if payl*mw > payw {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payl,
				Mult: mw,
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
	if count := screen.ScatNum(scat); count >= 2 {
		var pay, fs = ScatPay[count-1], ScatFreespin[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
			Free: fs,
		})
	}
}

func (g *Game) Spin(screen slot.Screen, mrtp float64) {
	if g.FSR == 0 {
		var reels, _ = slot.FindReels(ReelsMap, mrtp)
		screen.Spin(reels)
	} else {
		screen.Spin(ReelsBon)
	}
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
