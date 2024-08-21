package links

import (
	"context"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/pflag"
)

type GameAlias struct {
	ID   string `json:"id" yaml:"id" xml:"id"`
	Name string `json:"name" yaml:"name" xml:"name"`
}

type GameInfo struct {
	Aliases  []GameAlias `json:"aliases" yaml:"aliases" xml:"aliases"`
	Provider string      `json:"provider" yaml:"provider" xml:"provider"`
	ScrnX    int         `json:"scrnx" yaml:"scrnx" xml:"scrnx"`
	ScrnY    int         `json:"scrny" yaml:"scrny" xml:"scrny"`
	RtpList  []float64   `json:"rtplist" yaml:"rtplist" xml:"rtplist"`
}

var GameList = []GameInfo{}

var FlagsSetters = []func(*pflag.FlagSet){}

var ScanIters = []func(*pflag.FlagSet, context.Context){}

var GameFactory = map[string]func(float64) any{}

func Atof(s string) (f float64) {
	f, _ = strconv.ParseFloat(s, 64)
	return
}
