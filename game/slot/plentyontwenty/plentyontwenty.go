package plentyontwenty

import (
	slot "github.com/slotopol/server/game/slot"
)

// Lined payment.
var LinePay = [8][5]float64{
	{0, 0, 40, 400, 1000}, // seven
	{0, 0, 20, 80, 400},   // bell
	{0, 0, 20, 40, 200},   // melon
	{0, 0, 20, 40, 200},   // plum
	{0, 0, 10, 20, 100},   // orange
	{0, 0, 10, 20, 100},   // lemon
	{0, 0, 10, 20, 100},   // cherry
	{0, 0, 0, 0, 0},       // star
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 5, 20, 500} // star

const (
	jid = 1 // jackpot ID
)

// Jackpot win combinations.
var Jackpot = [8][5]int{
	{0, 0, 0, 0, 0}, // seven
	{0, 0, 0, 0, 0}, // bell
	{0, 0, 0, 0, 0}, // melon
	{0, 0, 0, 0, 0}, // plum
	{0, 0, 0, 0, 0}, // orange
	{0, 0, 0, 0, 0}, // lemon
	{0, 0, 0, 0, 0}, // cherry
	{0, 0, 0, 0, 0}, // star
}

type Game struct {
	slot.Slot5x3 `yaml:",inline"`
}

func NewGame() *Game {
	return &Game{
		Slot5x3: slot.Slot5x3{
			Sel: slot.MakeBitNum(20, 1),
			Bet: 1,
		},
	}
}

const wild, scat = 1, 8

var bl = slot.BetLinesNvm20

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	g.ScanLined(screen, wins)
	g.ScanScatters(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen slot.Screen, wins *slot.Wins) {
	for li := range g.Sel.Bits() {
		var line = bl.Line(li)

		var numw, numl = 0, 5
		var syml slot.Sym
		for x := 1; x <= 5; x++ {
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
				Jack: Jackpot[wild-1][numw-1],
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen slot.Screen, wins *slot.Wins) {
	if count := screen.ScatNum(scat); count >= 2 {
		var pay = ScatPay[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel.Num()) * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
		})
	}
}

func (g *Game) Spin(screen slot.Screen, mrtp float64) {
	var _, reels = FindReels(mrtp)
	screen.Spin(reels)
}

func (g *Game) SetSel(sel slot.Bitset) error {
	if sel.IsZero() {
		return slot.ErrNoLineset
	}
	if bs := sel; !bs.AndNot(slot.MakeBitNum(len(bl), 1)).IsZero() {
		return slot.ErrLinesetOut
	}
	if g.FreeSpins() > 0 {
		return slot.ErrNoFeature
	}
	g.Sel = sel
	return nil
}
