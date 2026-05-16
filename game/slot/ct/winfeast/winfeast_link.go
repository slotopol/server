//go:build !prod || full || ct

package winfeast

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed winfeast_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Win Feast", LNum: 20, Date: game.Date(2023, 2, 6)},      // see: https://www.livebet.com/casino/slots/ct-interactive/win-feast
		{Prov: "CT Interactive", Name: "Treasure Chase", LNum: 20, Date: game.Date(2023, 1, 3)}, // see: https://www.livebet.com/casino/slots/ct-interactive/treasure-chase
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPcasc |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 3,
		SN: sn,
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["ctinteractive/winfeast/rmap"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
