//go:build !prod || full || ct

package mightyrex

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed mightyrex_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Mighty Rex", LNum: 25, Date: game.Date(2014, 1, 15)},           // see: https://www.slotsmate.com/software/ct-interactive/mighty-rex
		{Prov: "CT Interactive", Name: "Bavarian Forest", LNum: 25, Date: game.Year(2010)},             // see: https://www.slotsmate.com/software/ct-interactive/bavarian-forest
		{Prov: "CT Interactive", Name: "Mountain Song Quechua", LNum: 20, Date: game.Date(2016, 6, 1)}, // see: https://www.livebet.com/casino/slots/ct-interactive/mountain-song-quechua
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgtwic |
			game.GPfgreel |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 3,
		SN: sn,
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["ctinteractive/mightyrex/bon"] = &ReelsBon
	game.DataRouter["ctinteractive/mightyrex/rmap"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
