package config

import (
	"context"

	"github.com/spf13/pflag"
)

var FlagsSetters = []func(*pflag.FlagSet){}

var ScatIters = []func(*pflag.FlagSet, context.Context){}

var GameAliases = map[string]string{}

var GameFactory = map[string]func(string) any{}
