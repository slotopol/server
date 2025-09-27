package neonbananas

// See: https://www.slotsmate.com/software/ct-interactive/neon-bananas

import (
	"github.com/slotopol/server/game/slot"
)

var ReelsMap slot.ReelsMap[*slot.Reels5x]

// Lined payment.
var LinePay = [13][5]float64{
	{0, 10, 200, 2000, 10000}, //  1 seven
	{},                        //  2 dollar
	{},                        //  3 banana
	{},                        //  4 lucky slot
	{0, 2, 10, 100, 500},      //  5 grapes
	{0, 2, 10, 100, 500},      //  6 melon
	{0, 0, 10, 100, 250},      //  7 apple
	{0, 0, 10, 100, 250},      //  8 pear
	{0, 0, 10, 100, 250},      //  9 peach
	{0, 0, 5, 20, 100},        // 10 orange
	{0, 0, 5, 20, 100},        // 11 lemon
	{0, 2, 5, 20, 100},        // 12 plum
	{0, 2, 5, 20, 100},        // 13 cherry
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 2, 20, 500} // 2 dollar

// Bet lines
var BetLines = slot.BetLinesMgj[:20]

type Game struct {
	slot.Screen5x3 `yaml:",inline"`
	slot.Slotx     `yaml:",inline"`
	BP             float64 `json:"bp" yaml:"bp" xml:"bp"` // bananas pay
	BC             int     `json:"bc" yaml:"bc" xml:"bc"` // bananas counter
	BN             int     `json:"bn" yaml:"bn" xml:"bn"` // bananas bonus number
}

// Declare conformity with SlotGame interface.
var _ slot.SlotGame = (*Game)(nil)

func NewGame() *Game {
	return &Game{
		Slotx: slot.Slotx{
			Sel: len(BetLines),
			Bet: 1,
		},
	}
}

func (g *Game) Clone() slot.SlotGame {
	var clone = *g
	return &clone
}

const (
	bbid = 1 // bananas bonus id
	lsb1 = 2 // lucky slot bonus id for 3 symbols
	lsb3 = 3 // lucky slot bonus id for 4 symbols
	lsb6 = 4 // lucky slot bonus id for 5 symbols
)
const wild, scat1, scat2, lssym = 1, 2, 3, 4

func (g *Game) Scanner(wins *slot.Wins) error {
	g.ScanLined(wins)
	g.ScanScatters(wins)
	return nil
}

// Lined symbols calculation.
func (g *Game) ScanLined(wins *slot.Wins) {
	for li, line := range BetLines[:g.Sel] {
		var numw, numl slot.Pos = 0, 5
		var syml slot.Sym
		var x slot.Pos
		for x = 1; x <= 5; x++ {
			var sx = g.LY(x, line)
			if sx == wild && syml != lssym {
				if syml == 0 {
					numw = x
				}
			} else if syml == 0 && (numw == 0 || sx != lssym) {
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
		if syml == lssym && numl >= 3 {
			*wins = append(*wins, slot.WinItem{
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li + 1,
				XY:   line.CopyL(numl),
				BID:  lsb1 + int(numl) - 3, // lucky slot bonus id
			})
		} else if payl > payw {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payl,
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li + 1,
				XY:   line.CopyL(numl),
			})
		} else if payw > 0 {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payw,
				Mult: 1,
				Sym:  wild,
				Num:  numw,
				Line: li + 1,
				XY:   line.CopyL(numw),
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(wins *slot.Wins) {
	if count := g.ScatNum(scat1); count >= 3 {
		var pay = ScatPay[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.Sel) * pay,
			Mult: 1,
			Sym:  scat1,
			Num:  count,
			XY:   g.ScatPos(scat1),
		})
	}
	if count := g.ScatNum(scat2); count >= 5 {
		*wins = append(*wins, slot.WinItem{
			Mult: 1,
			Sym:  scat2,
			Num:  count,
			BID:  bbid,
		})
	}
}

func (g *Game) Spin(mrtp float64) {
	var reels, _ = ReelsMap.FindClosest(mrtp)
	g.ReelSpin(reels)
}

func (g *Game) Spawn(wins slot.Wins, fund, mrtp float64) {
	for i, wi := range wins {
		switch wi.BID {
		case bbid:
			if bc := g.BC + 1; (g.BN == 0 && bc == 25) || (g.BN > 0 && bc == 15) {
				wins[i].Pay = g.BP + g.Bet*float64(g.Sel)
			}
		case lsb1:
			wins[i].Bon, wins[i].Pay = LuckySlotSpawn(g.Bet, 1)
		case lsb3:
			wins[i].Bon, wins[i].Pay = LuckySlotSpawn(g.Bet, 3)
		case lsb6:
			wins[i].Bon, wins[i].Pay = LuckySlotSpawn(g.Bet, 6)
		}
	}
}

func (g *Game) Apply(wins slot.Wins) {
	g.Slotx.Apply(wins)

	for _, wi := range wins {
		if wi.BID == bbid {
			if wi.Pay == 0 {
				g.BP += g.Bet * float64(g.Sel)
				g.BC++
			} else {
				g.BC = 0
				g.BN++
			}
		}
	}
}

func (g *Game) SetSel(sel int) error {
	return slot.ErrNoFeature
}
