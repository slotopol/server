//go:build !prod || full || ct

package hypercuber

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed hypercuber_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Hyper Cuber", Date: game.Date(2022, 10, 10)}, // see: https://www.livebet2.com/casino/slots/ct-interactive/hyper-cuber
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPcpay |
			game.GPcasc |
			game.GPfgseq |
			game.GPfgreel |
			game.GPfgmult |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 3,
		SN: sn,
		LN: 0,
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame() }, CalcStatReg)
	game.DataRouter["ctinteractive/hypercuber/bon"] = &ReelsBon
	game.DataRouter["ctinteractive/hypercuber/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
