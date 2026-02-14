//go:build !prod || full || ct

package wildhills

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed wildhills_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Wild Hills", LNum: 10, Date: game.Date(2020, 11, 26)},        // see: https://www.slotsmate.com/software/ct-interactive/wild-hills
		{Prov: "CT Interactive", Name: "The Great Cabaret", LNum: 10, Date: game.Date(2020, 11, 26)}, // see: https://www.slotsmate.com/software/ct-interactive/the-great-cabaret
		{Prov: "CT Interactive", Name: "Magician Dreaming", LNum: 10, Date: game.Date(2020, 11, 26)}, // see: https://www.slotsmate.com/software/ct-interactive/magician-dreaming
		{Prov: "CT Interactive", Name: "Jade Heaven", LNum: 15, Date: game.Date(2020, 11, 26)},       // see: https://www.slotsmate.com/software/ct-interactive/jade-heaven
		{Prov: "CT Interactive", Name: "Navy Girl", LNum: 10, Date: game.Date(2020, 11, 26)},         // see: https://www.slotsmate.com/software/ct-interactive/navy-girl
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
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStatReg)
	game.DataRouter["ctinteractive/wildhills/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
