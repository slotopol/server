//go:build !prod || full || ct

package burningbrilliants

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed burningbrilliants_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "100 Burning Brilliants", LNum: 100, Date: game.Date(2026, 4, 2)}, // see: https://www.slotsmate.com/software/ct-interactive/100-burning-brilliants
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPcasc |
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
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["ctinteractive/100burningbrilliants/rmap"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
