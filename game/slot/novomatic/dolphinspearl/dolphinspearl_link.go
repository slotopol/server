//go:build !prod || full || novomatic

package dolphinspearl

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed dolphinspearl_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Dolphins Pearl", LNum: 9, Date: game.Date(2001, 4, 25)},         // see: https://www.slotsmate.com/software/novomatic/dolphins-pearl
		{Prov: "Novomatic", Name: "Dolphins Pearl Deluxe", LNum: 10, Date: game.Date(2009, 4, 28)}, // see: https://www.slotsmate.com/software/novomatic/dolphins-pearl-deluxe
		{Prov: "Novomatic", Name: "Attila", LNum: 9, Date: game.Year(2009)},                        // see: https://casino.ru/attila-novomatic/
		{Prov: "Novomatic", Name: "Banana Splash", LNum: 9, Date: game.Year(2009)},                 // see: https://casino.ru/banana-splash-novomatic/
		{Prov: "Novomatic", Name: "Dynasty Of Ming", LNum: 9, Date: game.Date(2008, 2, 20)},
		{Prov: "Novomatic", Name: "Gryphons Gold", LNum: 9, Date: game.Year(2009)},               // see: https://www.slotsmate.com/software/novomatic/gryphons-gold
		{Prov: "Novomatic", Name: "Gryphons Gold Deluxe", LNum: 10, Date: game.Date(2017, 4, 1)}, // see: https://www.slotsmate.com/software/novomatic/gryphons-gold-deluxe
		{Prov: "Novomatic", Name: "Joker Dolphin", LNum: 9, Date: game.Year(2007)},
		{Prov: "Novomatic", Name: "King Of Cards", LNum: 9, Date: game.Date(2012, 2, 2)},
		{Prov: "Novomatic", Name: "Lucky Lady's Charm", LNum: 9, Date: game.Year(2001)},
		{Prov: "Novomatic", Name: "Lucky Lady's Charm Deluxe", LNum: 10, Date: game.Year(2008)}, // see: https://www.slotsmate.com/software/novomatic/lucky-ladys-charm-deluxe
		{Prov: "Novomatic", Name: "Pharaoh's Gold II", LNum: 9, Date: game.Year(2008)},          // see: https://casino.ru/pharaohs-gold-2-novomatic/
		{Prov: "Novomatic", Name: "Pharaoh's Gold III", LNum: 9, Date: game.Year(2011)},
		{Prov: "Novomatic", Name: "Polar Fox", LNum: 9, Date: game.Year(2008)}, // see: https://casino.ru/polar-fox-novomatic/
		{Prov: "Novomatic", Name: "Ramses II", LNum: 9, Date: game.Date(2011, 4, 15)},
		{Prov: "Novomatic", Name: "Royal Treasures", LNum: 9, Date: game.Year(2012)},              // see: https://www.slotsmate.com/software/novomatic/novomatic-royal-treasures
		{Prov: "Novomatic", Name: "Secret Forest", LNum: 9, Date: game.Year(2013)},                // see: https://www.slotsmate.com/software/novomatic/secret-forest
		{Prov: "Novomatic", Name: "The Money Game", LNum: 9, Date: game.Year(2009)},               // see: https://www.slotsmate.com/software/novomatic/the-money-game
		{Prov: "Novomatic", Name: "The Money Game Deluxe", LNum: 10, Date: game.Date(2018, 7, 1)}, // see: https://www.slotsmate.com/software/novomatic/the-money-game-deluxe
		{Prov: "Novomatic", Name: "Unicorn Magic", LNum: 9, Date: game.Year(2006)},                // see: https://casino.ru/unicorn-magic-novomatic/
		{Prov: "Novomatic", Name: "Cold Spell", LNum: 9, Date: game.Date(2018, 10, 18)},           // see: https://www.slotsmate.com/software/novomatic/cold-spell
		{Prov: "Novomatic", Name: "Mermaid's Pearl", LNum: 9, Date: game.Date(2014, 8, 15)},       // see: https://www.slotsmate.com/software/novomatic/mermaids-pearl
		{Prov: "Aristocrat", Name: "Dolphin Treasure", LNum: 10, Date: game.Date(1996, 12, 1)},    // See: https://freeslotshub.com/aristocrat/dolphin-treasure/
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
			game.GPfgreel |
			game.GPfgmult |
			game.GPscat |
			game.GPwild |
			game.GPwmult,
		SX: 5,
		SY: 3,
		SN: sn,
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStatReg)
	game.DataRouter["novomatic/dolphinspearl/bon"] = &ReelsBon
	game.DataRouter["novomatic/dolphinspearl/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
