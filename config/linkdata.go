package config

import (
	"context"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/pflag"
	"xorm.io/xorm"
)

var FlagsSetters = []func(*pflag.FlagSet){}

var ScatIters = []func(*pflag.FlagSet, context.Context){}

var GameAliases = map[string]string{}

var GameFactory = map[string]func(string) any{}

var (
	XormStorage *xorm.Engine
	XormSpinlog *xorm.Engine
)
