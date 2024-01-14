package spi

import (
	"sync"

	"github.com/slotopol/server/game"
)

type Room struct {
	RID  uint64  `xorm:"pk autoincr" json:"rid" yaml:"rid" xml:"rid,attr"`
	Bank float64 `xorm:"notnull" json:"bank" yaml:"bank" xml:"bank"` // users win/lost balance, in coins
	Fund float64 `xorm:"notnull" json:"fund" yaml:"fund" xml:"fund"` // jackpot fund, in coins
	Name string  `json:"name,omitempty" yaml:"name,omitempty" xml:"name,omitempty"`
	mux  sync.RWMutex
}

type User struct {
	UID    uint64 `xorm:"pk autoincr" json:"uid" yaml:"uid" xml:"uid,attr"`
	Email  string `xorm:"notnull unique index" json:"email" yaml:"email" xml:"email"`
	Secret string `xorm:"notnull" json:"secret" yaml:"secret" xml:"secret"` // auth password
	Name   string `json:"name,omitempty" yaml:"name,omitempty" xml:"name,omitempty"`

	games   RWMap[uint64, OpenGame]
	balance RWMap[uint64, *Balance]
}

type Balance struct {
	UID   uint64 `xorm:"notnull index(bid)" json:"uid" yaml:"uid" xml:"uid,attr"`
	RID   uint64 `xorm:"notnull index(bid)" json:"rid" yaml:"rid" xml:"rid,attr"`
	Value int    `xorm:"notnull" json:"value" yaml:"value" xml:"value"` // in coins
}

type OpenGame struct {
	GID   uint64 `xorm:"pk autoincr" json:"gid" yaml:"gid" xml:"gid,attr"`
	UID   uint64 `xorm:"notnull" json:"uid" yaml:"uid" xml:"uid,attr"`
	RID   uint64 `xorm:"notnull" json:"rid" yaml:"rid" xml:"rid,attr"`
	Alias string `xorm:"notnull" json:"alias" yaml:"alias" xml:"alias,attr"`
	game  game.SlotGame
}

var Rooms RWMap[uint64, *Room]

var Users RWMap[uint64, *User]

var GIDcounter uint64
var OpenGames RWMap[uint64, OpenGame]

func (user *User) Init() {
	user.games.Init(0)
	user.balance.Init(0)
}

func (user *User) InsertBalance(bal *Balance) {
	user.balance.Set(bal.RID, bal)
}

func init() {
	Rooms.Init(0)
	Users.Init(0)
	OpenGames.Init(0)
}
