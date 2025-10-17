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
		{Prov: "CT Interactive", Name: "Wild Hills", Date: game.Date(2020, 11, 26)},        // see: https://www.slotsmate.com/software/ct-interactive/wild-hills
		{Prov: "CT Interactive", Name: "The Great Cabaret", Date: game.Date(2020, 11, 26)}, // see: https://www.slotsmate.com/software/ct-interactive/the-great-cabaret
		{Prov: "CT Interactive", Name: "Magician Dreaming", Date: game.Date(2020, 11, 26)}, // see: https://www.slotsmate.com/software/ct-interactive/magician-dreaming
		{Prov: "CT Interactive", Name: "Forest Nymph", Date: game.Date(2020, 11, 26)},      // see: https://www.slotsmate.com/software/ct-interactive/forest-nymph
		{Prov: "CT Interactive", Name: "Jade Heaven", Date: game.Date(2020, 11, 26)},       // see: https://www.slotsmate.com/software/ct-interactive/jade-heaven
		{Prov: "CT Interactive", Name: "Navy Girl", Date: game.Date(2020, 11, 26)},         // see: https://www.slotsmate.com/software/ct-interactive/navy-girl
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
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
	Info.SetupFactory(func() game.Gamble { return NewGame() }, CalcStatReg)
	game.DataRouter["ctinteractive/wildhills/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
