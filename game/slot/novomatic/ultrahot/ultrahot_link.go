//go:build !prod || full || novomatic

package ultrahot

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed ultrahot_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "Novomatic", Name: "Ultra Hot", Date: game.Year(2002)},                // see: https://www.slotsmate.com/software/novomatic/ultra-hot
		{Prov: "Novomatic", Name: "Ultra Hot Deluxe", Date: game.Date(2008, 11, 18)}, // see: https://www.slotsmate.com/software/novomatic/ultrahot-deluxe
		{Prov: "Novomatic", Name: "Ultra Gems", Date: game.Date(2018, 1, 15)},        // see: https://www.slotsmate.com/software/novomatic/ultra-gems
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlsel |
			game.GPfill |
			game.GPfgno,
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
	game.DataRouter["novomatic/ultrahot/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
