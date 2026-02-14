//go:build !prod || full || ct

package treasurekingdom

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed treasurekingdom_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Treasure Kingdom", LNum: 9, Date: game.Date(2020, 11, 26)}, // see: https://www.slotsmate.com/software/ct-interactive/treasure-kingdom
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgseq |
			game.GPfgreel |
			game.GPfgmult |
			game.GPscat |
			game.GPwild |
			game.GPwmult,
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
	game.DataRouter["ctinteractive/treasurekingdom/bon"] = &ReelsBon
	game.DataRouter["ctinteractive/treasurekingdom/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
