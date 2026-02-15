//go:build !prod || full || igt

package wolfrun

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed wolfrun_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "IGT", Name: "Wolf Run", LNum: 40, Date: game.Year(2013)}, // see: https://www.slotsmate.com/software/igt/wolf-run
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
			game.GPfgreel |
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
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStatReg)
	game.DataRouter["igt/wolfrun/bon"] = &ReelsBon
	game.DataRouter["igt/wolfrun/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
