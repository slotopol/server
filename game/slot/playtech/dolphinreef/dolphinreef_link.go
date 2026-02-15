//go:build !prod || full || playtech

package dolphinreef

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed dolphinreef_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Playtech", Name: "Dolphin Reef", LNum: 20, Date: game.Date(2009, 2, 28)}, // see: https://www.slotsmate.com/software/playtech/dolphin-reef
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgonce |
			game.GPscat |
			game.GPwild |
			game.GPrwild,
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
	game.DataRouter["playtech/dolphinreef/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
