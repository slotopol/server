//go:build !prod || full || netent

package diamonddogs

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed diamonddogs_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "NetEnt", LNum: 25, Name: "Diamond Dogs", Date: game.Year(2013)},
		{Prov: "NetEnt", LNum: 25, Name: "Voodoo Vibes", Date: game.Year(2009)},
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
			game.GPfgreel |
			game.GPfgmult |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 3,
		SN: sn,
		LN: len(BetLines),
		BN: 1,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStatReg)
	game.DataRouter["netent/diamonddogs/bon"] = &ReelsBon
	game.DataRouter["netent/diamonddogs/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
