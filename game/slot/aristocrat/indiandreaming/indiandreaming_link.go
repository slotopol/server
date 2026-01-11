//go:build !prod || full || aristocrat

package indiandreaming

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed indiandreaming_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Aristocrat", Name: "Indian Cash Catcher", Date: game.Date(2014, 3, 15)},
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
		SN: sn,
		WN: 243,
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame() }, CalcStatReg)
	game.DataRouter["aristocrat/indiancashcatcher/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
