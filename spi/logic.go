package spi

import "github.com/slotopol/server/game"

type Room struct {
	RID  uint64  `xorm:"pk autoincr" json:"rid" yaml:"rid" xml:"rid,attr"`
	Bank float64 `json:"bank" yaml:"bank" xml:"bank"` // users win/lost balance, in coins
	Fund float64 `json:"fund" yaml:"fund" xml:"fund"` // jackpot fund, in coins
}

type User struct {
	UID     uint64 `xorm:"pk autoincr" json:"uid" yaml:"uid" xml:"uid,attr"`
	RID     uint64 `json:"rid" yaml:"rid" xml:"rid,attr"`
	Balance int    `json:"balance" yaml:"balance" xml:"balance"` // in coins
	Email   string `xorm:"notnull unique index" json:"email" yaml:"email" xml:"email"`
	Secret  string `json:"secret" yaml:"secret" xml:"secret"` // auth password
	Name    string `json:"name,omitempty" yaml:"name,omitempty" xml:"name,omitempty"`
	games   map[uint64]OpenGame
}

type OpenGame struct {
	GID   uint64 `xorm:"pk autoincr" json:"gid" yaml:"gid" xml:"gid,attr"`
	UID   uint64 `json:"uid" yaml:"uid" xml:"uid,attr"`
	Alias string `json:"alias" yaml:"alias" xml:"alias,attr"`
	game  game.SlotGame
}

var UIDcounter uint64 = 1
var Users = map[uint64]*User{
	1: &User{
		UID:     1,
		Balance: 100,
		games:   map[uint64]OpenGame{},
	},
}

var GIDcounter uint64
var OpenGames = map[uint64]OpenGame{}
