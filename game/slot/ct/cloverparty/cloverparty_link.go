//go:build !prod || full || ct

package cloverparty

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed cloverparty_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Clover Party", Date: game.Date(2020, 11, 26)},    // see: https://www.slotsmate.com/software/ct-interactive/clover-party
		{Prov: "CT Interactive", Name: "20 Clovers Hot", Date: game.Date(2020, 11, 30)},  // see: https://www.slotsmate.com/software/ct-interactive/20-clovers-hot
		{Prov: "CT Interactive", Name: "20 Shining Coins", Date: game.Date(2020, 1, 31)}, // see: https://www.slotsmate.com/software/ct-interactive/20-shining-coins
		{Prov: "CT Interactive", Name: "20 Mega Slot", Date: game.Date(2020, 12, 14)},    // see: https://www.slotsmate.com/software/ct-interactive/20-mega-slot
		{Prov: "CT Interactive", Name: "20 Mega Fresh", Date: game.Date(2021, 7, 7)},     // see: https://www.slotsmate.com/software/ct-interactive/20-mega-fresh
		{Prov: "CT Interactive", Name: "20 Mega Star", Date: game.Date(2024, 12, 31)},    // see: https://www.slotsmate.com/software/ct-interactive/20-mega-star
		{Prov: "CT Interactive", Name: "20 Roosters", Date: game.Date(2019, 12, 31)},     // see: https://www.slotsmate.com/software/ct-interactive/20-roosters
		{Prov: "CT Interactive", Name: "Egg and Rooster", Date: game.Date(2020, 11, 26)}, // see: https://www.slotsmate.com/software/ct-interactive/egg-and-rooster
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgno |
			game.GPscat |
			game.GPwild,
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
	game.DataRouter["ctinteractive/cloverparty/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
