//go:build !prod || full || novomatic

package config

import (
	"context"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/game/sizzlinghot"
	"github.com/spf13/pflag"
)

func init() {
	cfg.FlagsSetters = append(cfg.FlagsSetters, func(flags *pflag.FlagSet) {
		flags.Bool("sizzlinghot", false, "'Sizzling Hot' Novomatic 5x3 slots")
	})
	cfg.ScatIters = append(cfg.ScatIters, func(flags *pflag.FlagSet, ctx context.Context) {
		if is, _ := flags.GetBool("sizzlinghot"); is {
			var rn, _ = flags.GetString("reels")
			sizzlinghot.CalcStat(ctx, rn)
		}
	})

	for _, alias := range []string{
		"sizzlinghot",
		"sizzlinghotdeluxe",
	} {
		cfg.GameAliases[alias] = "sizzlinghot"
	}

	cfg.GameFactory["sizzlinghot"] = func(rd string) any {
		if _, ok := sizzlinghot.ReelsMap[rd]; ok {
			return sizzlinghot.NewGame(rd)
		}
		return nil
	}
}
