//go:build !prod || full || agt

package iceqween

import (
	"github.com/slotopol/server/game"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{ID: "agt/icequeen", Prov: "AGT", Name: "Ice Queen"},            // see: https://demo.agtsoftware.com/games/agt/iceqween
		{ID: "agt/stalker", Prov: "AGT", Name: "STALKER"},               // see: https://demo.agtsoftware.com/games/agt/stalker
		{ID: "agt/bigfive", Prov: "AGT", Name: "Big Five"},              // see: https://demo.agtsoftware.com/games/agt/bigfive
		{ID: "agt/arabiannights", Prov: "AGT", Name: "Arabian Nights"},  // see: https://demo.agtsoftware.com/games/agt/arabiannights
		{ID: "agt/anonymous", Prov: "AGT", Name: "Anonymous"},           // see: https://demo.agtsoftware.com/games/agt/anonymous
		{ID: "agt/grandtheft", Prov: "AGT", Name: "Grand Theft"},        // see: https://demo.agtsoftware.com/games/agt/bankofny
		{ID: "agt/firefighters", Prov: "AGT", Name: "Firefighters"},     // see: https://demo.agtsoftware.com/games/agt/firefighters
		{ID: "agt/timemachineii", Prov: "AGT", Name: "Time Machine II"}, // see: https://demo.agtsoftware.com/games/agt/timemachine2
		{ID: "agt/bitcoin", Prov: "AGT", Name: "Bitcoin"},               // see: https://demo.agtsoftware.com/games/agt/bitcoin
		{ID: "agt/piratesgold", Prov: "AGT", Name: "Pirates Gold"},      // see: https://demo.agtsoftware.com/games/agt/piratesgold
		{ID: "agt/theleprechaun", Prov: "AGT", Name: "The Leprechaun"},  // see: https://demo.agtsoftware.com/games/agt/leprechaun
	},
	GP: game.GPsel |
		game.GPretrig |
		game.GPfgmult |
		game.GPfgreel |
		game.GPscat |
		game.GPwild,
	SX:  5,
	SY:  3,
	SN:  len(LinePay),
	LN:  len(BetLines),
	BN:  0,
	RTP: game.MakeRtpList(ReelsMap),
}

func init() {
	game.GameList = append(game.GameList, &Info)
	for _, ga := range Info.Aliases {
		game.ScanFactory[ga.ID] = CalcStatReg
		game.GameFactory[ga.ID] = func() any {
			return NewGame()
		}
	}
}
