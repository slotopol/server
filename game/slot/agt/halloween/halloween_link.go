//go:build !prod || full || agt

package halloween

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed halloween_reel.yaml
var reels []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Halloween"},
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPfgno,
		SX: 3,
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["agt/halloween/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, reels)
}
