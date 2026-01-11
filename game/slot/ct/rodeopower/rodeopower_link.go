//go:build !prod || full || ct

package rodeopower

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed rodeopower_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Rodeo Power", Date: game.Date(2024, 10, 31)}, // see: https://www.slotsmate.com/software/ct-interactive/rodeo-power
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgseq |
			game.GPscat |
			game.GPwild |
			game.GPwmult,
		SX: 5,
		SY: 3,
		SN: sn,
		WN: 243,
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame() }, CalcStatReg)
	game.DataRouter["ctinteractive/rodeopower/bon"] = &ReelsBon
	game.DataRouter["ctinteractive/rodeopower/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
