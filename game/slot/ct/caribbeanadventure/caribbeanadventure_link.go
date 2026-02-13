//go:build !prod || full || ct

package caribbeanadventure

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed caribbeanadventure_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Caribbean Adventure", LNum: 10, Date: game.Date(2020, 11, 25)}, // see: https://www.slotsmate.com/software/ct-interactive/caribbean-adventure
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPcpay |
			game.GPfgseq |
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
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStatReg)
	game.DataRouter["ctinteractive/caribbeanadventure/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
