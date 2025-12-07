//go:build !prod || full || novomatic

package sizzlinghot

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed sizzlinghot_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Sizzling Hot", LNum: 5, Date: game.Date(2003, 3, 6)},          // see: https://www.slotsmate.com/software/novomatic/novomatic-sizzling-hot
		{Prov: "Novomatic", Name: "Sizzling Hot Deluxe", LNum: 5, Date: game.Date(2007, 11, 13)}, // see: https://www.slotsmate.com/software/novomatic/sizzling-hot-deluxe
		{Prov: "Novomatic", Name: "Age of Heroes", LNum: 5, Date: game.Date(2016, 4, 15)},        // see: https://www.slotsmate.com/software/novomatic/age-of-heroes
		{Prov: "Novomatic", Name: "Hot Cubes", LNum: 5, Date: game.Date(2007, 6, 15)},            // see: https://www.slotsmate.com/software/novomatic/hot-cubes
		{Prov: "Novomatic", Name: "Diamond 7", LNum: 5, Date: game.Date(2012, 11, 15)},           // see: https://www.slotsmate.com/software/novomatic/diamond-7
		{Prov: "Novomatic", Name: "Fruits 'n Royals", LNum: 5, Date: game.Date(2009, 12, 29)},    // see: https://www.slotsmate.com/software/novomatic/fruits-n-royals
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgno |
			game.GPscat,
		SX: 5,
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["novomatic/sizzlinghot/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
