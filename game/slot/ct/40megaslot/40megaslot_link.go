//go:build !prod || full || ct

package fortymegaslot

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed 40megaslot_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "40 Mega Slot", LNum: 40, Date: game.Date(2020, 12, 19)},    // see: https://www.slotsmate.com/software/ct-interactive/40-mega-slot
		{Prov: "CT Interactive", Name: "40 Roosters", LNum: 40, Date: game.Date(2020, 1, 3)},       // see: https://www.slotsmate.com/software/ct-interactive/40-roosters
		{Prov: "CT Interactive", Name: "40 Shining Coins", LNum: 40, Date: game.Date(2019, 11, 1)}, // see: https://www.livebet2.com/casino/slots/ct-interactive/40-shining-coins
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgno |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 4,
		SN: sn,
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["ctinteractive/40megaslot/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
