//go:build !prod || full || netent

package trolls

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed trolls_reel.yaml
var reels []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "NetEnt", Name: "Trolls", Date: game.Year(2009)},            // see: https://casino.ru/trolls-netent/
		{Prov: "NetEnt", Name: "Excalibur", Date: game.Date(2013, 11, 11)}, // see: https://www.slotsmate.com/software/netent/excalibur
		{Prov: "NetEnt", Name: "Pandora's Box", Date: game.Year(2009)},
		{Prov: "NetEnt", Name: "Wild Witches", Date: game.Year(2010)},
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPretrig |
			game.GPfgmult |
			game.GPscat |
			game.GPwild |
			game.GPwmult,
		SX: 5,
		SY: 3,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStat)
	game.DataRouter["netent/trolls/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, reels)
}
