//go:build !prod || full || novomatic

package justjewels

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed justjewels_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Just Jewels", LNum: 9, Date: game.Year(2005)},               // see: https://www.slotsmate.com/software/novomatic/just-jewels
		{Prov: "Novomatic", Name: "Just Jewels Deluxe", LNum: 10, Date: game.Date(2008, 7, 1)}, // see: https://www.slotsmate.com/software/novomatic/just-jewels-deluxe
		{Prov: "Novomatic", Name: "Just Fruits", LNum: 9, Date: game.Year(2001)},               // see: https://www.slotsmate.com/software/novomatic/just-fruits
		{Prov: "Novomatic", Name: "Royal Jewels", LNum: 10},                                    // see: https://casino.ru/garden-of-riches-novomatic/
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPapay |
			game.GPlsel |
			game.GPfgno |
			game.GPscat,
		SX: 5,
		SY: 3,
		SN: sn,
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["novomatic/justjewels/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
