//go:build !prod || full || megajack

package config

import (
	"context"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/game/champagne"
	"github.com/spf13/pflag"
)

func init() {
	cfg.FlagsSetters = append(cfg.FlagsSetters, func(flags *pflag.FlagSet) {
		flags.Bool("champagne", false, "'Champagne' Megajack 5x3 slots")
	})
	cfg.ScatIters = append(cfg.ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
		if is, _ := flags.GetBool("champagne"); is {
			var rn, _ = flags.GetString("reels")
			champagne.CalcStatReg(ctx, rn)
		}
	})

	for _, alias := range []string{
		"champagne",
	} {
		cfg.GameAliases[alias] = "champagne"
	}

	cfg.GameFactory["champagne"] = func(rd string) any {
		if _, ok := champagne.ReelsMap[rd]; ok {
			return champagne.NewGame(rd)
		}
		return nil
	}
}
