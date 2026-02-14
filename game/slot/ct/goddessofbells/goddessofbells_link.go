//go:build !prod || full || ct

package goddessofbells

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed goddessofbells_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Goddess of Bells", LNum: 9, Date: game.Date(2025, 2, 15)}, // see: https://www.slotsmate.com/software/ct-interactive/goddess-of-bells
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
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
	game.DataRouter["ctinteractive/goddessofbells/bon"] = &ReelsBon
	game.DataRouter["ctinteractive/goddessofbells/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
