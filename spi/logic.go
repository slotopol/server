package spi

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/slotopol/server/game"
	"github.com/slotopol/server/util"
)

// Club means independent bank into which gambles some users.
type Club struct {
	CID   uint64    `xorm:"pk autoincr" json:"cid" yaml:"cid" xml:"cid,attr"`
	CTime time.Time `xorm:"created 'ctime'" json:"ctime" yaml:"ctime" xml:"ctime"`
	Name  string    `json:"name,omitempty" yaml:"name,omitempty" xml:"name,omitempty"`
	Bank  float64   `xorm:"notnull" json:"bank" yaml:"bank" xml:"bank"` // users win/lost balance, in coins
	Fund  float64   `xorm:"notnull" json:"fund" yaml:"fund" xml:"fund"` // jackpot fund, in coins
	Lock  float64   `xorm:"notnull" json:"lock" yaml:"lock" xml:"lock"` // not changed deposit within games

	JptRate float64 `xorm:"'jptrate' notnull default 0.015" json:"jptrate" yaml:"jptrate" xml:"jptrate"`
	GainRTP float64 `xorm:"'gainrtp' notnull default 95.00" json:"gainrtp" yaml:"gainrtp" xml:"gainrtp"`

	mux sync.RWMutex
}

// User means registration of somebody. Each user can have splitted
// wallet with some coins balance in each Club. User can opens several
// games without any limitation.
type User struct {
	UID    uint64    `xorm:"pk autoincr" json:"uid" yaml:"uid" xml:"uid,attr"`
	CTime  time.Time `xorm:"created 'ctime'" json:"ctime" yaml:"ctime" xml:"ctime"`
	Email  string    `xorm:"notnull unique index" json:"email" yaml:"email" xml:"email"`
	Secret string    `xorm:"notnull" json:"secret" yaml:"secret" xml:"secret"` // auth password
	Name   string    `json:"name,omitempty" yaml:"name,omitempty" xml:"name,omitempty"`
	GAL    AL        `json:"gal,omitempty" yaml:"gal,omitempty" xml:"gal,omitempty"` // global access level
	games  util.RWMap[uint64, OpenGame]
	props  util.RWMap[uint64, *Props]
}

// State structure represents full current game state.
type State struct {
	Game         game.SlotGame `json:"game" yaml:"game" xml:"game"`
	Scrn         game.Screen   `json:"scrn" yaml:"scrn" xml:"scrn"`
	game.WinScan `yaml:",inline"`
}

// OpenGame is opened game for user with UID at club with CID.
// Each instance of game have own GID. Alias - is game type identifier.
type OpenGame struct {
	GID   uint64    `xorm:"pk autoincr" json:"gid" yaml:"gid" xml:"gid,attr"`
	CTime time.Time `xorm:"created 'ctime'" json:"ctime" yaml:"ctime" xml:"ctime"`
	CID   uint64    `xorm:"notnull" json:"cid" yaml:"cid" xml:"cid,attr"`
	UID   uint64    `xorm:"notnull" json:"uid" yaml:"uid" xml:"uid,attr"`
	Alias string    `xorm:"notnull" json:"alias" yaml:"alias" xml:"alias"`
	State `xorm:"-" yaml:",inline"`
}

func (OpenGame) TableName() string {
	return "game"
}

// Access level.
type AL uint

const (
	ALban   AL = 1 << iota // user have no access to club
	ALgame                 // can change club game settings
	ALuser                 // can change user balance and move user money to/from club deposit
	ALclub                 // can change club bank, fund, deposit
	ALadmin                // can change same access levels to other users
	ALall   = ALgame | ALuser | ALclub | ALadmin
)

// Props contains properties for user at some club.
// Any property can be zero by default, or if object does not created at DB.
type Props struct {
	CID    uint64  `xorm:"notnull index(bid)" json:"cid" yaml:"cid" xml:"cid,attr"`
	UID    uint64  `xorm:"notnull index(bid)" json:"uid" yaml:"uid" xml:"uid,attr"`
	Wallet float64 `xorm:"notnull" json:"wallet" yaml:"wallet" xml:"wallet"` // in coins
	Access AL      `xorm:"notnull" json:"access" yaml:"access" xml:"access"`
}

type Spinlog struct {
	ID     uint64    `xorm:"pk autoincr" json:"id" yaml:"id" xml:"id,attr"`
	CTime  time.Time `xorm:"created 'ctime'" json:"ctime" yaml:"ctime" xml:"ctime"`
	GID    uint64    `xorm:"notnull" json:"gid" yaml:"gid" xml:"gid,attr"`
	Game   string    `xorm:"notnull" json:"game" yaml:"game" xml:"game"`
	Screen string    `xorm:"notnull" json:"screen,omitempty" yaml:"screen,omitempty" xml:"screen,omitempty"`
	Wins   string    `json:"wins,omitempty" yaml:"wins,omitempty" xml:"wins,omitempty"`
	Gain   float64   `xorm:"notnull" json:"gain" yaml:"gain" xml:"gain"`
	Wallet float64   `xorm:"notnull" json:"wallet" yaml:"wallet" xml:"wallet"`
}

type Walletlog struct {
	ID     uint64    `xorm:"pk autoincr" json:"id" yaml:"id" xml:"id,attr"`
	CTime  time.Time `xorm:"created 'ctime'" json:"ctime" yaml:"ctime" xml:"ctime"`
	CID    uint64    `xorm:"notnull index(bid)" json:"cid" yaml:"cid" xml:"cid,attr"`
	UID    uint64    `xorm:"notnull index(bid)" json:"uid" yaml:"uid" xml:"uid,attr"`
	AdmID  uint64    `xorm:"notnull" json:"admid" yaml:"admid" xml:"admid"`
	Wallet float64   `xorm:"notnull" json:"wallet" yaml:"wallet" xml:"wallet"` // in coins
	Addend float64   `xorm:"notnull" json:"addend" yaml:"addend" xml:"addend"`
}

type Banklog struct {
	ID      uint64    `xorm:"pk autoincr" json:"id" yaml:"id" xml:"id,attr"`
	CTime   time.Time `xorm:"created 'ctime'" json:"ctime" yaml:"ctime" xml:"ctime"`
	Bank    float64   `xorm:"notnull 'bank'" json:"bank" yaml:"bank" xml:"bank"`
	Fund    float64   `xorm:"notnull 'fund'" json:"fund" yaml:"fund" xml:"fund"`
	Lock    float64   `xorm:"notnull 'lock'" json:"lock" yaml:"lock" xml:"lock"`
	BankSum float64   `xorm:"notnull 'banksum'" json:"banksum" yaml:"banksum" xml:"banksum" form:"banksum"`
	FundSum float64   `xorm:"notnull 'fundsum'" json:"fundsum" yaml:"fundsum" xml:"fundsum" form:"fundsum"`
	LockSum float64   `xorm:"notnull 'locksum'" json:"locksum" yaml:"locksum" xml:"locksum" form:"locksum"`
}

// All created clubs, by CID.
var Clubs util.RWMap[uint64, *Club]

// All registered users, by UID.
var Users util.RWMap[uint64, *User]

// All opened games, by GID.
var OpenGames util.RWMap[uint64, OpenGame]

func (user *User) Init() {
	user.games.Init(0)
	user.props.Init(0)
}

func (user *User) GetWallet(cid uint64) float64 {
	if props, ok := user.props.Get(cid); ok {
		return props.Wallet
	}
	return 0
}

func (user *User) GetAL(cid uint64) AL {
	if props, ok := user.props.Get(cid); ok {
		return props.Access
	}
	return 0
}

func (user *User) InsertProps(props *Props) {
	user.props.Set(props.CID, props)
}

// GetAdmin always returns User pointer for authorized
// requests, and access level for it.
func GetAdmin(c *gin.Context, cid uint64) (*User, AL) {
	var admin = c.MustGet(userKey).(*User)
	return admin, admin.GAL | admin.GetAL(cid)
}

func (sl *Spinlog) MarshalState(s *State) (err error) {
	var b []byte
	if b, err = json.Marshal(s.Game); err != nil {
		return
	}
	sl.Game = util.B2S(b)
	if b, err = json.Marshal(s.Scrn); err != nil {
		return
	}
	sl.Screen = util.B2S(b)
	if len(s.WinScan.Wins) > 0 {
		if b, err = json.Marshal(s.WinScan.Wins); err != nil {
			return
		}
		sl.Wins = util.B2S(b)
	}
	return
}

func init() {
	Clubs.Init(0)
	Users.Init(0)
	OpenGames.Init(0)
}
