//go:build !prod || full || ct

package halloweenhot

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed halloweenhot_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Halloween Hot", LNum: 5, Date: game.Date(2021, 10, 15)}, // see: https://www.slotsmate.com/software/ct-interactive/halloween-hot
		{Prov: "CT Interactive", Name: "Fire King", LNum: 5, Date: game.Date(2021, 9, 1)},       // see: https://www.livebet.com/casino/slots/ct-interactive/fire-king
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfill |
			game.GPfgno |
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
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["ctinteractive/halloweenhot/rmap"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
