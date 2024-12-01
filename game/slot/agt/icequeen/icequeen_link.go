//go:build !prod || full || agt

package iceqween

import (
	"github.com/slotopol/server/game"
	"github.com/slotopol/server/util"
)

var Info = game.GameInfo{
	Aliases: []game.GameAlias{
		{Prov: "AGT", Name: "Ice Queen"},       // see: https://demo.agtsoftware.com/games/agt/iceqween
		{Prov: "AGT", Name: "STALKER"},         // see: https://demo.agtsoftware.com/games/agt/stalker
		{Prov: "AGT", Name: "Big Five"},        // see: https://demo.agtsoftware.com/games/agt/bigfive
		{Prov: "AGT", Name: "Arabian Nights"},  // see: https://demo.agtsoftware.com/games/agt/arabiannights
		{Prov: "AGT", Name: "Anonymous"},       // see: https://demo.agtsoftware.com/games/agt/anonymous
		{Prov: "AGT", Name: "Grand Theft"},     // see: https://demo.agtsoftware.com/games/agt/bankofny
		{Prov: "AGT", Name: "Firefighters"},    // see: https://demo.agtsoftware.com/games/agt/firefighters
		{Prov: "AGT", Name: "Time Machine II"}, // see: https://demo.agtsoftware.com/games/agt/timemachine2
		{Prov: "AGT", Name: "Bitcoin"},         // see: https://demo.agtsoftware.com/games/agt/bitcoin
		{Prov: "AGT", Name: "Pirates Gold"},    // see: https://demo.agtsoftware.com/games/agt/piratesgold
		{Prov: "AGT", Name: "The Leprechaun"},  // see: https://demo.agtsoftware.com/games/agt/leprechaun
	},
	GP: game.GPsel |
		game.GPretrig |
		game.GPfgmult |
		game.GPfgreel |
		game.GPscat |
		game.GPwild |
		game.GPwmult,
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
		var aid = util.ToID(ga.Prov + "/" + ga.Name)
		game.ScanFactory[aid] = CalcStatReg
		game.GameFactory[aid] = func() any { return NewGame() }
	}
}
