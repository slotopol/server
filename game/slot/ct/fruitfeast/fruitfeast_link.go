//go:build !prod || full || ct

package fruitfeast

import (
	_ "embed"

	"github.com/slotopol/server/game"
)

//go:embed fruitfeast_data.yaml
var data []byte

var Info = game.AlgInfo{
	Aliases: []game.GameAlias{
		{Prov: "CT Interactive", Name: "Fruit Feast", LNum: 40, Date: game.Date(2020, 9, 1)},           // see: https://www.slotsmate.com/software/ct-interactive/fruit-feast
		{Prov: "CT Interactive", Name: "Golden Acorn", LNum: 40, Date: game.Date(2020, 11, 24)},        // see: https://www.slotsmate.com/software/ct-interactive/golden-acorn
		{Prov: "CT Interactive", Name: "Wet and Juicy", LNum: 40, Date: game.Date(2020, 11, 25)},       // see: https://www.slotsmate.com/software/ct-interactive/wet-and-juicy
		{Prov: "CT Interactive", Name: "Mountain Heroes", LNum: 40, Date: game.Date(2021, 10, 31)},     // see: https://www.slotsmate.com/software/ct-interactive/mountain-heroes
		{Prov: "CT Interactive", Name: "40 Brilliants", LNum: 40, Date: game.Date(2020, 6, 30)},        // see: https://www.slotsmate.com/software/ct-interactive/40-brilliants
		{Prov: "CT Interactive", Name: "40 Fruitata Wins", LNum: 40, Date: game.Date(2023, 9, 14)},     // see: https://www.slotsmate.com/software/ct-interactive/40-fruitata-wins
		{Prov: "CT Interactive", Name: "40 Diamond Treasures", LNum: 40, Date: game.Date(2023, 5, 12)}, // see: https://www.livebet2.com/casino/slots/ct-interactive/40-diamond-treasures
		{Prov: "CT Interactive", Name: "40 Shining jewels", LNum: 40, Date: game.Date(2020, 12, 28)},   // see: https://www.livebet2.com/casino/slots/ct-interactive/40-shining-jewels
		{Prov: "CT Interactive", Name: "40 Dice Treasures", LNum: 40, Date: game.Date(2020, 11, 26)},   // see: https://www.livebet2.com/casino/slots/ct-interactive/40-dice-treasures
		{Prov: "CT Interactive", Name: "40 Hell's Cherries", LNum: 40, Date: game.Date(2025, 11, 2)},   // see: https://www.livebet2.com/casino/slots/ct-interactive/40-hell-s-cherries
		{Prov: "CT Interactive", Name: "Easter Bonanza", LNum: 40, Date: game.Date(2023, 4, 3)},        // see: https://www.livebet2.com/casino/slots/ct-interactive/easter-bonanza
	},
	AlgDescr: game.AlgDescr{
		GT: game.GTslot,
		GP: game.GPlpay |
			game.GPfgno |
			game.GPscat |
			game.GPwild,
		SX: 5,
		SY: 4,
		SN: len(LinePay),
		LN: len(BetLines),
		BN: 0,
	},
	Update: func(ai *game.AlgInfo) { ai.RTP = game.MakeRtpList(ReelsMap) },
}

func init() {
	Info.SetupFactory(func(sel int) game.Gamble { return NewGame(sel) }, CalcStat)
	game.DataRouter["ctinteractive/fruitfeast/reel"] = &ReelsMap
	game.LoadMap = append(game.LoadMap, data)
}
