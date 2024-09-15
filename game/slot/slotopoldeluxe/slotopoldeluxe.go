package slotopoldeluxe

import (
	"math"

	slot "github.com/slotopol/server/game/slot"
	"github.com/slotopol/server/game/slot/slotopol"
	"github.com/slotopol/server/util"
)

// reels lengths [32, 32, 32, 32, 32], total reshuffles 33554432
// symbols: 76.288(lined) + 2.7777(scatter) = 79.065740%
// spin1 bonuses: count1 20736, rtp = 6.550598%
// spin3 bonuses: count3 3120, rtp = 2.956867%
// spin6 bonuses: count6 720, rtp = 1.364708%
// monopoly bonuses: count 6720, rtp = 5.739904%
// jackpots: count 32, frequency 1/1048576
// RTP = 79.066(sym) + 10.872(mje) + 5.7399(mjm) = 95.677817%
var Reels957 = slot.Reels5x{
	{1, 5, 13, 4, 13, 3, 13, 5, 9, 7, 8, 13, 10, 13, 12, 11, 13, 12, 11, 13, 13, 2, 4, 5, 2, 6, 7, 9, 8, 3, 10, 6},
	{13, 2, 12, 9, 13, 4, 5, 6, 9, 7, 13, 10, 12, 13, 11, 13, 13, 11, 12, 13, 3, 4, 5, 2, 8, 7, 10, 4, 6, 8, 3, 1},
	{2, 1, 12, 3, 4, 5, 2, 6, 7, 10, 8, 4, 5, 13, 12, 11, 13, 12, 11, 13, 12, 3, 5, 13, 9, 6, 7, 10, 13, 8, 9, 13},
	{1, 2, 12, 4, 3, 6, 4, 2, 7, 8, 3, 5, 13, 12, 11, 13, 12, 11, 13, 12, 2, 4, 5, 3, 12, 6, 10, 7, 13, 8, 9, 13},
	{1, 2, 12, 4, 3, 5, 12, 6, 7, 3, 8, 12, 2, 13, 12, 11, 13, 12, 11, 13, 12, 3, 4, 5, 2, 6, 7, 10, 13, 8, 9, 5},
}

// reels lengths [32, 32, 32, 32, 32], total reshuffles 33554432
// symbols: 77.605(lined) + 2.7777(scatter) = 80.382288%
// spin1 bonuses: count1 20736, rtp = 6.550598%
// spin3 bonuses: count3 3120, rtp = 2.956867%
// spin6 bonuses: count6 720, rtp = 1.364708%
// monopoly bonuses: count 6720, rtp = 5.739904%
// jackpots: count 32, frequency 1/1048576
// RTP = 80.382(sym) + 10.872(mje) + 5.7399(mjm) = 96.994365%
var Reels970 = slot.Reels5x{
	{1, 5, 13, 4, 13, 3, 13, 5, 9, 7, 8, 13, 10, 13, 12, 11, 13, 12, 11, 13, 13, 2, 4, 5, 2, 6, 7, 9, 8, 3, 10, 6},
	{13, 2, 12, 9, 13, 4, 5, 6, 9, 7, 13, 10, 12, 13, 11, 13, 13, 11, 12, 13, 3, 4, 5, 2, 8, 7, 10, 4, 6, 8, 3, 1},
	{2, 1, 12, 3, 4, 5, 2, 6, 7, 10, 8, 4, 5, 13, 12, 11, 13, 12, 11, 13, 12, 3, 5, 13, 9, 6, 7, 10, 13, 8, 9, 13},
	{1, 2, 12, 4, 3, 6, 4, 2, 7, 8, 10, 5, 13, 12, 11, 13, 12, 11, 13, 12, 2, 4, 5, 3, 12, 6, 10, 7, 13, 8, 9, 13},
	{1, 2, 12, 4, 3, 5, 12, 6, 7, 3, 8, 12, 9, 13, 12, 11, 13, 12, 11, 13, 12, 3, 4, 5, 2, 6, 7, 10, 13, 8, 9, 5},
}

// reels lengths [32, 32, 32, 32, 32], total reshuffles 33554432
// symbols: 79.254(lined) + 2.7777(scatter) = 82.031429%
// spin1 bonuses: count1 20736, rtp = 6.550598%
// spin3 bonuses: count3 3120, rtp = 2.956867%
// spin6 bonuses: count6 720, rtp = 1.364708%
// monopoly bonuses: count 6720, rtp = 5.739904%
// jackpots: count 32, frequency 1/1048576
// RTP = 82.031(sym) + 10.872(mje) + 5.7399(mjm) = 98.643506%
var Reels986 = slot.Reels5x{
	{1, 5, 13, 4, 13, 3, 13, 5, 9, 7, 8, 13, 10, 13, 12, 11, 13, 12, 11, 13, 13, 2, 4, 5, 2, 6, 7, 9, 8, 3, 10, 6},
	{13, 2, 12, 9, 13, 4, 5, 6, 9, 7, 13, 10, 12, 13, 11, 13, 13, 11, 12, 13, 3, 4, 5, 2, 8, 7, 10, 4, 6, 8, 3, 1},
	{2, 1, 12, 3, 4, 5, 2, 6, 7, 10, 8, 4, 5, 13, 12, 11, 13, 12, 11, 13, 12, 3, 5, 13, 9, 6, 7, 10, 13, 8, 9, 13},
	{1, 2, 12, 4, 3, 6, 4, 2, 7, 8, 10, 5, 13, 12, 11, 13, 12, 11, 13, 12, 9, 4, 5, 3, 12, 6, 10, 7, 13, 8, 9, 13},
	{1, 2, 12, 4, 3, 5, 12, 6, 7, 10, 8, 12, 9, 13, 12, 11, 13, 12, 11, 13, 12, 3, 4, 5, 2, 6, 7, 10, 13, 8, 9, 5},
}

// Original reels.
// symbols: 81.213(lined) + 2.7777(scatter) = 83.990312%
// spin1 bonuses: count1 28672, rtp = 9.057617%
// spin3 bonuses: count3 3328, rtp = 3.153992%
// spin6 bonuses: count6 768, rtp = 1.455688%
// monopoly bonuses: count 6912, rtp = 5.903901%
// jackpots: count 32, frequency 1/1048576
// RTP = 83.99(sym) + 13.667(mje) + 5.9039(mjm) = 103.561510%
var Reels104 = slot.Reels5x{
	{1, 2, 13, 4, 13, 3, 13, 5, 9, 7, 8, 13, 10, 13, 12, 11, 13, 12, 11, 13, 13, 2, 4, 5, 2, 6, 7, 9, 8, 3, 10, 6},
	{13, 2, 12, 9, 13, 2, 5, 6, 9, 7, 13, 10, 12, 13, 11, 12, 13, 11, 12, 13, 3, 4, 5, 2, 8, 7, 10, 4, 6, 8, 3, 1},
	{2, 1, 12, 3, 4, 5, 2, 6, 7, 10, 8, 4, 3, 13, 12, 11, 13, 12, 11, 13, 12, 3, 5, 13, 9, 6, 7, 10, 13, 8, 9, 13},
	{1, 2, 12, 4, 3, 6, 4, 2, 7, 8, 10, 5, 13, 12, 11, 13, 12, 11, 13, 12, 9, 4, 5, 3, 13, 6, 10, 7, 13, 8, 9, 13},
	{1, 2, 12, 4, 3, 5, 12, 6, 7, 10, 8, 12, 9, 13, 12, 11, 13, 12, 11, 13, 12, 3, 4, 5, 2, 6, 7, 10, 13, 8, 9, 5},
}

// Map with available reels.
var ReelsMap = map[float64]*slot.Reels5x{
	95.677817:  &Reels957,
	96.994365:  &Reels970,
	98.643506:  &Reels986,
	103.561510: &Reels104, // original
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
var LinePay = [13][5]float64{
	{0, 0, 0, 0, 0},           //  1 dollar
	{0, 2, 5, 25, 100},        //  2 cherry
	{0, 2, 5, 25, 100},        //  3 plum
	{0, 0, 5, 25, 100},        //  4 wmelon
	{0, 0, 5, 25, 100},        //  5 grapes
	{0, 0, 10, 100, 250},      //  6 ananas
	{0, 0, 10, 100, 250},      //  7 lemon
	{0, 0, 10, 100, 250},      //  8 drink
	{0, 2, 10, 100, 500},      //  9 palm
	{0, 2, 10, 100, 500},      // 10 yacht
	{0, 10, 200, 2000, 10000}, // 11 eldorado
	{0, 0, 0, 0, 0},           // 12 spin
	{0, 0, 0, 0, 0},           // 13 dice
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 2, 20, 1000} // 1 dollar

const (
	mje1 = 1 // Eldorado9
	mje3 = 2 // Eldorado9
	mje6 = 3 // Eldorado9
	mje9 = 4 // Eldorado9
	mjm  = 5 // Monopoly
	mjc  = 6 // Champagne
)

// Lined bonus games
var LineBonus = [13][5]int{
	{0, 0, 0, 0, 0},          //  1
	{0, 0, 0, 0, 0},          //  2
	{0, 0, 0, 0, 0},          //  3
	{0, 0, 0, 0, 0},          //  4
	{0, 0, 0, 0, 0},          //  5
	{0, 0, 0, 0, 0},          //  6
	{0, 0, 0, 0, 0},          //  7
	{0, 0, 0, 0, 0},          //  8
	{0, 0, 0, 0, 0},          //  9
	{0, 0, 0, 0, 0},          // 10
	{0, 0, 0, 0, 0},          // 11
	{0, 0, mje1, mje3, mje6}, // 12 Eldorado1, Eldorado3, Eldorado6
	{0, 0, 0, 0, mjm},        // 13 Monopoly
}

type Game struct {
	slot.Slot5x3 `yaml:",inline"`
}

func NewGame() *Game {
	return &Game{
		Slot5x3: slot.Slot5x3{
			Sel: util.MakeBitNum(5, 1),
			Bet: 1,
		},
	}
}

const (
	jid = 1 // jackpot ID
)

// Not from lined paytable.
var special = [13]bool{
	true,  //  1
	false, //  2
	false, //  3
	false, //  4
	false, //  5
	false, //  6
	false, //  7
	false, //  8
	false, //  9
	false, // 10
	false, // 11
	true,  // 12
	true,  // 13
}

const wild, scat = 11, 1

var bl = slot.BetLinesMgj

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
		var mw float64 = 1 // mult wild
		for x := 1; x <= 5; x++ {
			var sx = screen.Pos(x, line)
			if sx == wild {
				if syml == 0 {
					numw = x
				} else if special[syml-1] {
					numl = x - 1
					break
				}
				mw = 2
			} else if numw > 0 && special[sx-1] {
				numl = x - 1
				break
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
				Jack: slotopol.Jackpot[wild-1][numw-1],
			})
		} else if syml > 0 && numl > 0 && LineBonus[syml-1][numl-1] > 0 {
			*wins = append(*wins, slot.WinItem{
				Mult: 1,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyL(numl),
				BID:  LineBonus[syml-1][numl-1],
			})
		}
	}
}

// Scatters calculation.
func (g *Game) ScanScatters(screen slot.Screen, wins *slot.Wins) {
	if count := screen.ScatNum(scat); count >= 3 {
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

func (g *Game) Spawn(screen slot.Screen, wins slot.Wins) {
	for i, wi := range wins {
		switch wi.BID {
		case mje1:
			wins[i].Bon, wins[i].Pay = slotopol.EldoradoSpawn(g.Bet, 1)
		case mje3:
			wins[i].Bon, wins[i].Pay = slotopol.EldoradoSpawn(g.Bet, 3)
		case mje6:
			wins[i].Bon, wins[i].Pay = slotopol.EldoradoSpawn(g.Bet, 6)
		case mje9:
			wins[i].Bon, wins[i].Pay = slotopol.EldoradoSpawn(g.Bet, 9)
		case mjm:
			wins[i].Bon, wins[i].Pay = slotopol.MonopolySpawn(g.Bet)
		}
	}
}

func (g *Game) SetSel(sel slot.Bitset) error {
	var mask slot.Bitset = (1<<len(bl) - 1) << 1
	if sel == 0 {
		return slot.ErrNoLineset
	}
	if sel&^mask != 0 {
		return slot.ErrLinesetOut
	}
	if g.FreeSpins() > 0 {
		return slot.ErrNoFeature
	}
	g.Sel = sel
	return nil
}
