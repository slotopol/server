//go:build !prod || full || novomatic

package alwayshot

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed alwayshot_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Always Hot", LNum: 5, Date: game.Date(2008, 11, 17)},
		{Prov: "Novomatic", Name: "Always Hot Deluxe", LNum: 5, Date: game.Date(2008, 11, 18)}, // see: https://www.slotsmate.com/software/novomatic/always-hot-deluxe
		{Prov: "Novomatic", Name: "Always American", LNum: 5, Date: game.Date(2016, 4, 15)},    // see: https://www.slotsmate.com/software/novomatic/always-american
		{Prov: "AGT", Name: "Tropic Hot", LNum: 5},                                             // see: https://agtsoftware.com/games/agt/tropichot
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlsel |
			game.GPfgno,
		SX: 3,
		SY: 3,
		SN: sn,
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["novomatic/alwayshot/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
