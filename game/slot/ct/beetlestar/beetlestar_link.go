//go:build !prod || full || ct

package beetlestar

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed beetlestar_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Beetle Star", LNum: 10, Date: game.Date(2018, 12, 31)}, // see: https://www.livebet2.com/casino/slots/ct-interactive/beetle-star
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgseq |
			game.GPfgmult |
			game.GPscat |
			game.GPwild |
			game.GPwmult,
		SX: 5,
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStatReg)
	game.DataRouter["ctinteractive/beetlestar/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
