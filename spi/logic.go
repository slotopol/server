package spi

import (
	"sync"
	"time"

	"github.com/slotopol/server/game"
)

type Room struct {
	RID  uint64  `xorm:"pk autoincr" json:"rid" yaml:"rid" xml:"rid,attr"`
	Bank float64 `xorm:"notnull" json:"bank" yaml:"bank" xml:"bank"` // users win/lost balance, in coins
	Fund float64 `xorm:"notnull" json:"fund" yaml:"fund" xml:"fund"` // jackpot fund, in coins
	Lock float64 `xorm:"notnull" json:"lock" yaml:"lock" xml:"lock"` // not changed deposit within games
	Name string  `json:"name,omitempty" yaml:"name,omitempty" xml:"name,omitempty"`
	mux  sync.RWMutex
}

type User struct {
	UID    uint64 `xorm:"pk autoincr" json:"uid" yaml:"uid" xml:"uid,attr"`
	Email  string `xorm:"notnull unique index" json:"email" yaml:"email" xml:"email"`
	Secret string `xorm:"notnull" json:"secret" yaml:"secret" xml:"secret"` // auth password
	Name   string `json:"name,omitempty" yaml:"name,omitempty" xml:"name,omitempty"`
	games  RWMap[uint64, OpenGame]
	props  RWMap[uint64, *Props]
}

type OpenGame struct {
	GID   uint64 `xorm:"pk autoincr" json:"gid" yaml:"gid" xml:"gid,attr"`
	UID   uint64 `xorm:"notnull" json:"uid" yaml:"uid" xml:"uid,attr"`
	RID   uint64 `xorm:"notnull" json:"rid" yaml:"rid" xml:"rid,attr"`
	Alias string `xorm:"notnull" json:"alias" yaml:"alias" xml:"alias"`
	game  game.SlotGame
}

func (OpenGame) TableName() string {
	return "game"
}

type AL uint

const (
	ALgame  AL = 1 << iota // can change room game settings
	ALuser                 // can change user balance and move user money to/from room deposit
	ALbank                 // can change room bank, fund, deposit
	ALadmin                // can change same access levels to other users
	ALall   AL = 0xffff
)

type Props struct {
	UID    uint64 `xorm:"notnull index(bid)" json:"uid" yaml:"uid" xml:"uid,attr"`
	RID    uint64 `xorm:"notnull index(bid)" json:"rid" yaml:"rid" xml:"rid,attr"`
	Wallet int    `xorm:"notnull" json:"wallet" yaml:"wallet" xml:"wallet"` // in coins
	Access AL     `xorm:"notnull" json:"access" yaml:"access" xml:"access"`
}

type Spinlog struct {
	SID    uint64    `xorm:"pk autoincr" json:"sid" yaml:"sid" xml:"sid,attr"`
	GID    uint64    `xorm:"notnull" json:"gid" yaml:"gid" xml:"gid,attr"`
	Game   string    `xorm:"notnull" json:"game" yaml:"game" xml:"game"`
	Screen string    `xorm:"notnull" json:"screen" yaml:"screen" xml:"screen"`
	Wins   string    `json:"wins" yaml:"wins" xml:"wins"`
	Gain   int       `xorm:"notnull" json:"gain" yaml:"gain" xml:"gain"`
	Wallet int       `xorm:"notnull" json:"wallet" yaml:"wallet" xml:"wallet"`
	CTime  time.Time `xorm:"created" json:"ctime" yaml:"ctime" xml:"ctime"`
}

var Rooms RWMap[uint64, *Room]

var Users RWMap[uint64, *User]

var OpenGames RWMap[uint64, OpenGame]

func (user *User) Init() {
	user.games.Init(0)
	user.props.Init(0)
}

func (user *User) InsertProps(props *Props) {
	user.props.Set(props.RID, props)
}

func init() {
	Rooms.Init(0)
	Users.Init(0)
	OpenGames.Init(0)
}
