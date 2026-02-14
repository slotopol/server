//go:build !prod || full || ct

package devilsfruits

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed devilsfruits_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Devil's Fruits", LNum: 5, Date: game.Date(2020, 11, 25)}, // see: https://www.slotsmate.com/software/ct-interactive/devils-fruits
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPmix |
			game.GPfgno |
			game.GPwild |
			game.GPwmult,
		SX: 3,
		SY: 3,
		SN: sn,
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["ctinteractive/devilsfruits/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
