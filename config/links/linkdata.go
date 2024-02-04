package links

import (
	"context"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/pflag"
)

var FlagsSetters = []func(*pflag.FlagSet){}

var ScatIters = []func(*pflag.FlagSet, context.Context){}

var GameAliases = map[string]string{}

var GameFactory = map[string]func(string) any{}
