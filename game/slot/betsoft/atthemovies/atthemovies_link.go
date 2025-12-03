//go:build !prod || full || betsoft

package atthemovies

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed atthemovies_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "BetSoft", Name: "At the Movies", Date: game.Year(2012)}, // see: https://www.slotsmate.com/software/betsoft/at-the-movies
		{Prov: "BetSoft", Name: "Sushi Bar", Date: game.Year(2014)},     // see: https://www.slotsmate.com/software/betsoft/sushi-bar
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPlsel |
			game.GPfgseq |
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
	game.DataRouter["betsoft/atthemovies/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
