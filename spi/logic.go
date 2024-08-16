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
	CID   uint64    `xorm:"pk autoincr" json:"cid" yaml:"cid" xml:"cid,attr"`      // club ID
	CTime time.Time `xorm:"created 'ctime'" json:"ctime" yaml:"ctime" xml:"ctime"` // creation time
	UTime time.Time `xorm:"updated 'utime'" json:"utime" yaml:"utime" xml:"utime"` // update time
	Name  string    `json:"name,omitempty" yaml:"name,omitempty" xml:"name,omitempty"`
	Bank  float64   `xorm:"notnull" json:"bank" yaml:"bank" xml:"bank"` // users win/lost balance, in coins
	Fund  float64   `xorm:"notnull" json:"fund" yaml:"fund" xml:"fund"` // jackpot fund, in coins
	Lock  float64   `xorm:"notnull" json:"lock" yaml:"lock" xml:"lock"` // not changed deposit within games

	JptRate float64 `xorm:"'jptrate' notnull default 0.015" json:"jptrate" yaml:"jptrate" xml:"jptrate"`
	GainRTP float64 `xorm:"'gainrtp' notnull default 95" json:"gainrtp" yaml:"gainrtp" xml:"gainrtp"`

	mux sync.RWMutex
}

// User means registration of somebody. Each user can have splitted
// wallet with some coins balance in each Club. User can opens several
// games without any limitation.
type User struct {
	UID    uint64    `xorm:"pk autoincr" json:"uid" yaml:"uid" xml:"uid,attr"`      // user ID
	CTime  time.Time `xorm:"created 'ctime'" json:"ctime" yaml:"ctime" xml:"ctime"` // creation time
	UTime  time.Time `xorm:"updated 'utime'" json:"utime" yaml:"utime" xml:"utime"` // update time
	Email  string    `xorm:"notnull unique index" json:"email" yaml:"email" xml:"email"`
	Secret string    `xorm:"notnull" json:"secret" yaml:"secret" xml:"secret"` // auth password
	Name   string    `json:"name,omitempty" yaml:"name,omitempty" xml:"name,omitempty"`
	GAL    AL        `json:"gal,omitempty" yaml:"gal,omitempty" xml:"gal,omitempty"` // global access level
	games  util.RWMap[uint64, *Scene]
	props  util.RWMap[uint64, *Props]
}

// Story is opened game for user with UID at club with CID.
// Each instance of game have own GID. Alias - is game type identifier.
type Story struct {
	GID   uint64    `xorm:"pk autoincr" json:"gid" yaml:"gid" xml:"gid,attr"`      // game ID
	CTime time.Time `xorm:"created 'ctime'" json:"ctime" yaml:"ctime" xml:"ctime"` // creation time
	UTime time.Time `xorm:"updated 'utime'" json:"utime" yaml:"utime" xml:"utime"` // update time
	Alias string    `xorm:"notnull" json:"alias" yaml:"alias" xml:"alias"`         // game type identifier
	CID   uint64    `xorm:"notnull" json:"cid" yaml:"cid" xml:"cid,attr"`          // club ID
	UID   uint64    `xorm:"notnull" json:"uid" yaml:"uid" xml:"uid,attr"`          // user ID
	Flow  bool      `xorm:"notnull" json:"flow" yaml:"flow" xml:"flow,attr"`       // game is not closed
}

// Scene represents game with all the connected environment.
type Scene struct {
	Story `yaml:",inline"`
	SID   uint64        `json:"sid" yaml:"sid" xml:"sid,attr"` // last spin ID
	Game  game.SlotGame `json:"game" yaml:"game" xml:"game"`
	Scrn  game.Screen   `json:"scrn" yaml:"scrn" xml:"scrn"`
	Wins  game.Wins     `json:"wins,omitempty" yaml:"wins,omitempty" xml:"wins,omitempty"`
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
	CID    uint64    `xorm:"notnull index(bid)" json:"cid" yaml:"cid" xml:"cid,attr"`    // club ID
	UID    uint64    `xorm:"notnull index(bid)" json:"uid" yaml:"uid" xml:"uid,attr"`    // user ID
	CTime  time.Time `xorm:"created 'ctime'" json:"ctime" yaml:"ctime" xml:"ctime"`      // creation time
	UTime  time.Time `xorm:"updated 'utime'" json:"utime" yaml:"utime" xml:"utime"`      // update time
	Wallet float64   `xorm:"notnull default 0" json:"wallet" yaml:"wallet" xml:"wallet"` // in coins
	Access AL        `xorm:"notnull default 0" json:"access" yaml:"access" xml:"access"` // access level
}

type Spinlog struct {
	SID    uint64    `xorm:"pk" json:"sid" yaml:"sid" xml:"sid,attr"`                                        // spin ID
	CTime  time.Time `xorm:"created 'ctime'" json:"ctime" yaml:"ctime" xml:"ctime"`                          // creation time
	GID    uint64    `xorm:"notnull" json:"gid" yaml:"gid" xml:"gid,attr"`                                   // game ID
	Game   string    `xorm:"notnull" json:"game" yaml:"game" xml:"game"`                                     // game data
	Screen string    `xorm:"notnull" json:"screen,omitempty" yaml:"screen,omitempty" xml:"screen,omitempty"` // game screen marshaled to JSON
	Wins   string    `xorm:"text" json:"wins,omitempty" yaml:"wins,omitempty" xml:"wins,omitempty"`          // list of wins marshaled to JSON
	Gain   float64   `xorm:"notnull" json:"gain" yaml:"gain" xml:"gain"`                                     // total gain at last spin
	Wallet float64   `xorm:"notnull" json:"wallet" yaml:"wallet" xml:"wallet"`
}

var SpinCounter uint64 // last spin log ID

type Multlog struct {
	ID     uint64    `xorm:"pk" json:"id" yaml:"id" xml:"id,attr"`
	CTime  time.Time `xorm:"created 'ctime'" json:"ctime" yaml:"ctime" xml:"ctime"`
	GID    uint64    `xorm:"notnull" json:"gid" yaml:"gid" xml:"gid,attr"` // game ID
	Mult   int       `xorm:"notnull" json:"mult" yaml:"mult" xml:"mult"`   // multiplier
	Risk   float64   `xorm:"notnull" json:"risk" yaml:"risk" xml:"risk"`
	Gain   float64   `xorm:"notnull" json:"gain" yaml:"gain" xml:"gain"`
	Wallet float64   `xorm:"notnull" json:"wallet" yaml:"wallet" xml:"wallet"`
}

var MultCounter uint64 // last multiplier log ID

type Walletlog struct {
	ID     uint64    `xorm:"pk autoincr" json:"id" yaml:"id" xml:"id,attr"`
	CTime  time.Time `xorm:"created 'ctime'" json:"ctime" yaml:"ctime" xml:"ctime"`   // creation time
	CID    uint64    `xorm:"notnull index(bid)" json:"cid" yaml:"cid" xml:"cid,attr"` // club ID
	UID    uint64    `xorm:"notnull index(bid)" json:"uid" yaml:"uid" xml:"uid,attr"` // user ID
	AID    uint64    `xorm:"notnull" json:"aid" yaml:"aid" xml:"aid"`                 // admin ID
	Wallet float64   `xorm:"notnull" json:"wallet" yaml:"wallet" xml:"wallet"`        // new value in coins
	Sum    float64   `xorm:"notnull" json:"sum" yaml:"sum" xml:"sum"`
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

// All scenes, by GID.
var Scenes util.RWMap[uint64, *Scene]

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

func (sl *Spinlog) MarshalState(scene *Scene) (err error) {
	var b []byte
	if b, err = json.Marshal(scene.Game); err != nil {
		return
	}
	sl.Game = util.B2S(b)
	if b, err = json.Marshal(scene.Scrn); err != nil {
		return
	}
	sl.Screen = util.B2S(b)
	if len(scene.Wins) > 0 {
		if b, err = json.Marshal(scene.Wins); err != nil {
			return
		}
		sl.Wins = util.B2S(b)
	}
	return
}

func init() {
	Clubs.Init(0)
	Users.Init(0)
	Scenes.Init(0)
}
