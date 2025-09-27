//go:build !prod || full || agt

package merrychristmas

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed merrychristmas_reel.yaml
var reels []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Merry Christmas"},
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPretrig |
			game.GPscat |
			game.GPwild,
		SX: 3,
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStatReg)
	game.DataRouter["agt/merrychristmas/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, reels)
}
