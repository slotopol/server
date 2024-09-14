package champagne

import (
	"math"

	slot "github.com/slotopol/server/game/slot"
	"github.com/slotopol/server/util"
)

// Original reels.
// *bonus reels calculations*
// reels lengths [32, 32, 32, 32, 32], total reshuffles 33554432
// symbols: 135.5(lined) + 0.7242(scatter) = 136.224294%
// free games 7171740, q = 0.21373, sq = 1/(1-q) = 1.271835
// free games frequency: 1/140.36
// champagne bonuses: count 11025, rtp = 6.362796%
// jackpots: count 32, frequency 1/1048576
// RTP = sq*(rtp(sym)+rtp(mjc)) = 1.2718*(136.22+6.3628) = 181.347256%
// *regular reels calculations*
// reels lengths [32, 32, 32, 32, 32], total reshuffles 33554432
// symbols: 69.974(lined) + 0.7242(scatter) = 70.698214%
// free games 3585870, q = 0.106867
// free games frequency: 1/140.36
// champagne bonuses: count 11025, rtp = 6.362796%
// jackpots: count 32, frequency 1/1048576
// RTP = rtp(sym) + rtp(mjc) + q*rtp(fg) = 70.698 + 6.3628 + 0.10687*181.35 = 96.441093%
var Reels964 = slot.Reels5x{
	{12, 1, 5, 2, 12, 11, 2, 11, 12, 3, 2, 8, 12, 3, 4, 6, 12, 2, 5, 10, 3, 9, 7, 8, 4, 3, 7, 9, 2, 3, 4, 6},
	{2, 5, 10, 12, 9, 6, 3, 4, 12, 2, 6, 8, 3, 12, 11, 2, 11, 12, 5, 7, 4, 6, 3, 4, 12, 2, 5, 8, 2, 7, 1, 9},
	{12, 5, 10, 12, 9, 6, 3, 4, 12, 2, 6, 8, 3, 12, 11, 2, 11, 12, 5, 7, 4, 6, 3, 4, 12, 2, 5, 8, 12, 7, 1, 9},
	{12, 8, 2, 12, 6, 5, 2, 4, 12, 2, 1, 3, 2, 9, 7, 12, 11, 11, 11, 11, 12, 5, 2, 12, 8, 6, 2, 3, 10, 12, 2, 4},
	{12, 11, 7, 12, 6, 4, 12, 3, 2, 12, 3, 7, 12, 3, 5, 1, 12, 3, 8, 9, 12, 4, 3, 2, 12, 5, 3, 10, 2, 12, 3, 6},
}

// Map with available reels.
var ReelsMap = map[float64]*slot.Reels5x{
	96.441093: &Reels964, // original
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
var LinePay = [12][5]float64{
	{0, 0, 0, 0, 0},           //  1 dollar
	{0, 3, 5, 20, 100},        //  2 cherry
	{0, 3, 5, 20, 100},        //  3 plum
	{0, 0, 5, 20, 100},        //  4 wmelon
	{0, 0, 5, 20, 100},        //  5 grapes
	{0, 0, 5, 20, 100},        //  6 ananas
	{0, 0, 5, 20, 100},        //  7 lemon
	{0, 0, 5, 20, 100},        //  8 drink
	{0, 5, 10, 20, 1000},      //  9 palm
	{0, 7, 10, 20, 1000},      // 10 yacht
	{0, 10, 100, 2000, 10000}, // 11 eldorado
	{0, 0, 0, 0, 0},           // 12 fizz
}

// Scatters payment.
var ScatPay = [5]float64{0, 0, 0, 0, 1000} // 1 dollar

// Scatter freespins table
var ScatFreespinReg = [5]int{0, 0, 15, 15, 15} // 1 dollar

// Scatter freespins table
var ScatFreespinBon = [5]int{0, 0, 30, 30, 30} // 1 dollar

const (
	mje1 = 1 // Eldorado9
	mje3 = 2 // Eldorado9
	mje6 = 3 // Eldorado9
	mje9 = 4 // Eldorado9
	mjm  = 5 // Monopoly
	mjc  = 6 // Champagne
)

// Lined bonus games
var LineBonus = [12][5]int{
	{0, 0, 0, 0, 0},   //  1
	{0, 0, 0, 0, 0},   //  2
	{0, 0, 0, 0, 0},   //  3
	{0, 0, 0, 0, 0},   //  4
	{0, 0, 0, 0, 0},   //  5
	{0, 0, 0, 0, 0},   //  6
	{0, 0, 0, 0, 0},   //  7
	{0, 0, 0, 0, 0},   //  8
	{0, 0, 0, 0, 0},   //  9
	{0, 0, 0, 0, 0},   // 10
	{0, 0, 0, 0, 0},   // 11
	{0, 0, 0, 0, mjc}, // 12 Champagne
}

const (
	jid = 1 // jackpot ID
)

// Jackpot win combinations.
var Jackpot = [12][5]int{
	{0, 0, 0, 0, 0},   // //  1 dollar
	{0, 0, 0, 0, 0},   // //  2 cherry
	{0, 0, 0, 0, 0},   // //  3 plum
	{0, 0, 0, 0, 0},   // //  4 wmelon
	{0, 0, 0, 0, 0},   // //  5 grapes
	{0, 0, 0, 0, 0},   // //  6 ananas
	{0, 0, 0, 0, 0},   // //  7 lemon
	{0, 0, 0, 0, 0},   // //  8 drink
	{0, 0, 0, 0, 0},   // //  9 palm
	{0, 0, 0, 0, 0},   // // 10 yacht
	{0, 0, 0, 0, jid}, // // 11 eldorado
	{0, 0, 0, 0, 0},   // // 12 fizz
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

// Not from lined paytable.
var Special = [12]bool{
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
}

const wild, scat = 11, 1

var bl = slot.BetLinesMgj

func (g *Game) Scanner(screen slot.Screen, wins *slot.Wins) {
	g.ScanLined(screen, wins)
	g.ScanScatters(screen, wins)
}

// Lined symbols calculation.
func (g *Game) ScanLined(screen slot.Screen, wins *slot.Wins) {
	var mm float64 = 1 // mult mode
	if g.FS > 0 {
		mm = 2
	}

	for li := range g.SBL.Bits() {
		var line = bl.Line(li)

		var numw, numl = 0, 5
		var syml slot.Sym
		for x := 1; x <= 5; x++ {
			var sx = screen.Pos(x, line)
			if sx == wild {
				if syml == 0 {
					numw = x
				} else if Special[syml-1] {
					numl = x - 1
					break
				}
			} else if numw > 0 && Special[sx-1] {
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
		if payl > payw {
			*wins = append(*wins, slot.WinItem{
				Pay:  g.Bet * payl,
				Mult: mm,
				Sym:  syml,
				Num:  numl,
				Line: li,
				XY:   line.CopyL(numl),
			})
		} else if payw > 0 {
			if syml > 0 {
				*wins = append(*wins, slot.WinItem{
					Pay:  g.Bet * payw,
					Mult: mm,
					Sym:  wild,
					Num:  numw,
					Line: li,
					XY:   line.CopyL(numw),
				})
			} else {
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
		} else if syml > 0 && numl > 0 && LineBonus[syml-1][numl-1] > 0 {
			*wins = append(*wins, slot.WinItem{
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
		var fs int
		if g.FS > 0 {
			fs = ScatFreespinBon[count-1]
		} else {
			fs = ScatFreespinReg[count-1]
		}
		var pay = ScatPay[count-1]
		*wins = append(*wins, slot.WinItem{
			Pay:  g.Bet * float64(g.SBL.Num()) * pay,
			Mult: 1,
			Sym:  scat,
			Num:  count,
			XY:   screen.ScatPos(scat),
			Free: fs,
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
		case mjc:
			wins[i].Bon, wins[i].Pay = ChampagneSpawn(g.Bet)
		}
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
