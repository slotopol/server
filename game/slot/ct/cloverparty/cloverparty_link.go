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
		{Prov: "CT Interactive", Name: "Clover Party", LNum: 20, Date: game.Date(2020, 11, 26)},     // see: https://www.slotsmate.com/software/ct-interactive/clover-party
		{Prov: "CT Interactive", Name: "20 Clovers Hot", LNum: 20, Date: game.Date(2020, 11, 30)},   // see: https://www.slotsmate.com/software/ct-interactive/20-clovers-hot
		{Prov: "CT Interactive", Name: "20 Shining Coins", LNum: 20, Date: game.Date(2020, 1, 31)},  // see: https://www.slotsmate.com/software/ct-interactive/20-shining-coins
		{Prov: "CT Interactive", Name: "20 Mega Slot", LNum: 20, Date: game.Date(2020, 12, 14)},     // see: https://www.slotsmate.com/software/ct-interactive/20-mega-slot
		{Prov: "CT Interactive", Name: "20 Mega Fresh", LNum: 20, Date: game.Date(2021, 7, 7)},      // see: https://www.slotsmate.com/software/ct-interactive/20-mega-fresh
		{Prov: "CT Interactive", Name: "20 Mega Star", LNum: 20, Date: game.Date(2024, 12, 31)},     // see: https://www.slotsmate.com/software/ct-interactive/20-mega-star
		{Prov: "CT Interactive", Name: "20 Roosters", LNum: 20, Date: game.Date(2019, 12, 31)},      // see: https://www.slotsmate.com/software/ct-interactive/20-roosters
		{Prov: "CT Interactive", Name: "Egg and Rooster", LNum: 20, Date: game.Date(2020, 11, 26)},  // see: https://www.slotsmate.com/software/ct-interactive/egg-and-rooster
		{Prov: "CT Interactive", Name: "20 Dice Party", LNum: 20, Date: game.Date(2023, 5, 12)},     // see: https://www.livebet2.com/casino/slots/ct-interactive/20-dice-party
		{Prov: "CT Interactive", Name: "20 Star Party", LNum: 20, Date: game.Date(2020, 12, 22)},    // see: https://www.livebet2.com/casino/slots/ct-interactive/20-star-party
		{Prov: "CT Interactive", Name: "20 Fruitata Wins", LNum: 20, Date: game.Date(2023, 11, 15)}, // see: https://www.livebet2.com/casino/slots/ct-interactive/20-fruitata-wins
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgno |
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
	game.DataRouter["ctinteractive/cloverparty/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
