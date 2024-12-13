package valkyrie

// See: https://demo.agtsoftware.com/games/agt/valkyrie

import (
	_ "embed"
	"math/rand/v2"

	"github.com/slotopol/server/game/slot"
)

//go:embed valkyrie_bon.yaml
var rbon []byte

var BonusReel = slot.ReadBon[[]slot.Sym](rbon)

//go:embed valkyrie_reel.yaml
var reels []byte

var ReelsMap = slot.ReadReelsMap[*slot.Reels5x](reels)

// Lined payment.
var LinePay = [13][5]float64{
	{0, 2, 50, 500, 1000}, //  1 wild
	{},                    //  2 scatter
	{0, 2, 25, 250, 500},  //  3 warrior
	{0, 2, 25, 100, 200},  //  4 helmet
	{0, 0, 20, 100, 200},  //  5 shield
	{0, 0, 15, 50, 100},   //  6 axe
	{0, 0, 15, 50, 100},   //  7 mug
	{0, 0, 10, 25, 50},    //  8 ace
	{0, 0, 10, 25, 50},    //  9 king
	{0, 0, 5, 10, 25},     // 10 queen
	{0, 0, 5, 10, 25},     // 11 jack
	{0, 0, 5, 10, 25},     // 12 ten
	{0, 2, 5, 10, 25},     // 13 nine
}

// Bet lines
var BetLines = slot.BetLinesAgt5x3[:30]

type Game struct {
	slot.Slot5x3 `yaml:",inline"`
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slot5x3: slot.Slot5x3{
			Sel: len(BetLines),
			Bet: 1,
		},
	}
}

const wild, scat = 1, 2

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	g.ScanLined(screen, wins)
	g.ScanScatters(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen slot.Screen, wins *slot.Wins) {
	for li := 1; li <= g.Sel; li++ {
		var line = BetLines[li-1]

		var numw, numl slot.Pos = 0, 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = screen.Pos(x, line)
			if sx == wild {
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
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen slot.Screen, wins *slot.Wins) {
	if count := screen.ScatNum(scat); count >= 3 {
		*wins = append(*wins, slot.WinItem{
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
			Free: 15,
		})
	}
}

func (g *Game) Spin(screen slot.Screen, mrtp float64) {
	var reels, _ = slot.FindReels(ReelsMap, mrtp)
	if g.FSR == 0 {
		screen.Spin(reels)
	} else {
		var reel []slot.Sym
		var hit int
		// set 1 reel
		reel = reels.Reel(1)
		hit = rand.N(len(reel))
		screen.SetCol(1, reel, hit)
		// set center
		var big = BonusReel[rand.N(len(BonusReel))]
		var x slot.Pos
		for x = 2; x <= 4; x++ {
			screen.Set(x, 1, big)
			screen.Set(x, 2, big)
			screen.Set(x, 3, big)
		}
		// set 5 reel
		reel = reels.Reel(5)
		hit = rand.N(len(reel))
		screen.SetCol(5, reel, hit)
	}
}

func (g *Game) SetSel(sel int) error {
	return g.SetSelNum(sel, len(BetLines))
}
