package core

import (
	"encoding/json"
	"errors"
	"iter"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/atomic"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/game"
	"github.com/slotopol/server/util"
)

const (
	// Pointer to User object stored at gin context
	// after successful authorization.
	userKey = "user"
)

var (
	ErrNotOpened = errors.New("game with given ID is not opened")
	ErrNoAlias   = errors.New("no game alias")
)

type ClubData struct {
	CID   uint64    `xorm:"pk autoincr" json:"cid" yaml:"cid" xml:"cid,attr"`                                             // club ID
	CTime time.Time `xorm:"created 'ctime' notnull default CURRENT_TIMESTAMP" json:"ctime" yaml:"ctime" xml:"ctime,attr"` // creation time
	UTime time.Time `xorm:"updated 'utime' notnull default CURRENT_TIMESTAMP" json:"utime" yaml:"utime" xml:"utime,attr"` // update time
	Name  string    `xorm:"notnull" json:"name,omitempty" yaml:"name,omitempty" xml:"name,omitempty"`
	Bank  float64   `xorm:"notnull default 0" json:"bank" yaml:"bank" xml:"bank"`          // users win/lost balance, in coins
	Fund  float64   `xorm:"notnull default 0" json:"fund" yaml:"fund" xml:"fund"`          // jackpot fund, in coins
	Lock  float64   `xorm:"notnull default 0" json:"lock" yaml:"lock" xml:"lock"`          // not changed deposit within games
	Rate  float64   `xorm:"'rate' notnull default 2.5" json:"rate" yaml:"rate" xml:"rate"` // jackpot rate for games with progressive jackpot
	MRTP  float64   `xorm:"'mrtp' notnull default 0" json:"mrtp" yaml:"mrtp" xml:"mrtp"`   // master RTP
}

func (ClubData) TableName() string {
	return "club"
}

// Club means independent bank into which gambles some users.
type Club struct {
	data ClubData
	mux  sync.RWMutex
}

// User flag.
type UF uint

const (
	UFactivated UF = 1 << iota // account is activated
	UFsigncode                 // sign-in required code
)

// User means registration of somebody. Each user can have splitted
// wallet with some coins balance in each Club. User can opens several
// games without any limitation.
type User struct {
	UID    uint64    `xorm:"pk autoincr" json:"uid" yaml:"uid" xml:"uid,attr"`                                             // user ID
	CTime  time.Time `xorm:"created 'ctime' notnull default CURRENT_TIMESTAMP" json:"ctime" yaml:"ctime" xml:"ctime,attr"` // creation time
	UTime  time.Time `xorm:"updated 'utime' notnull default CURRENT_TIMESTAMP" json:"utime" yaml:"utime" xml:"utime,attr"` // update time
	Email  string    `xorm:"notnull unique index" json:"email" yaml:"email" xml:"email"`                                   // unique user email
	Secret string    `xorm:"notnull" json:"secret" yaml:"secret" xml:"secret"`                                             // auth password
	Name   string    `xorm:"notnull" json:"name,omitempty" yaml:"name,omitempty" xml:"name,omitempty"`                     // user name
	Code   uint32    `xorm:"notnull default 0" json:"code,omitempty" yaml:"code,omitempty" xml:"code,omitempty"`           // verification code
	Status UF        `xorm:"notnull default 0" json:"status,omitempty" yaml:"status,omitempty" xml:"status,omitempty"`     // account status
	GAL    AL        `xorm:"notnull default 0" json:"gal,omitempty" yaml:"gal,omitempty" xml:"gal,omitempty"`              // global access level
	props  util.RWMap[uint64, *Props]
}

// Story is opened game for user with UID at club with CID.
// Each instance of game have own GID. Alias - is game type identifier.
type Story struct {
	GID   uint64    `xorm:"pk" json:"gid" yaml:"gid" xml:"gid,attr"`                                                      // game ID
	CID   uint64    `xorm:"notnull" json:"cid" yaml:"cid" xml:"cid,attr"`                                                 // club ID
	UID   uint64    `xorm:"notnull" json:"uid" yaml:"uid" xml:"uid,attr"`                                                 // user ID
	Alias string    `xorm:"notnull" json:"alias" yaml:"alias" xml:"alias"`                                                // game type identifier
	CTime time.Time `xorm:"created 'ctime' notnull default CURRENT_TIMESTAMP" json:"ctime" yaml:"ctime" xml:"ctime,attr"` // creation time
}

var StoryCounter atomic.Uint64 // last GID

// Scene represents game with all the connected environment.
type Scene struct {
	Story `yaml:",inline"`
	SLID  uint64      `json:"slid" yaml:"slid" xml:"slid,attr"` // last spin log ID
	Game  game.Gamble `json:"game" yaml:"game" xml:"game"`
}

// Access level.
type AL uint

const (
	ALmember AL = 1 << iota // user have access to club
	ALdealer                // can change club game settings and users gameplay
	ALbooker                // can change user properties and move user money to/from club deposit
	ALmaster                // can change club bank, fund, deposit
	ALadmin                 // can change same access levels to other users
	ALall    = ALmember | ALdealer | ALbooker | ALmaster | ALadmin
)

// Props contains properties for user at some club.
// Any property can be zero by default, or if object does not created at DB.
type Props struct {
	CID    uint64    `xorm:"notnull index(bid)" json:"cid" yaml:"cid" xml:"cid,attr"`                                      // club ID
	UID    uint64    `xorm:"notnull index(bid)" json:"uid" yaml:"uid" xml:"uid,attr"`                                      // user ID
	CTime  time.Time `xorm:"created 'ctime' notnull default CURRENT_TIMESTAMP" json:"ctime" yaml:"ctime" xml:"ctime,attr"` // creation time
	UTime  time.Time `xorm:"updated 'utime' notnull default CURRENT_TIMESTAMP" json:"utime" yaml:"utime" xml:"utime,attr"` // update time
	Wallet float64   `xorm:"notnull default 0" json:"wallet" yaml:"wallet" xml:"wallet"`                                   // in coins
	Access AL        `xorm:"notnull default 0" json:"access" yaml:"access" xml:"access"`                                   // access level
	MRTP   float64   `xorm:"notnull default 0" json:"mrtp" yaml:"mrtp" xml:"mrtp"`                                         // personal master RTP
}

// Properties master for new registered user.
var PropMaster []Props

type Spinlog struct {
	SLID   uint64    `xorm:"pk" json:"slid" yaml:"slid" xml:"slid,attr"`                                                   // spin log ID
	GID    uint64    `xorm:"notnull" json:"gid" yaml:"gid" xml:"gid,attr"`                                                 // game ID
	CTime  time.Time `xorm:"created 'ctime' notnull default CURRENT_TIMESTAMP" json:"ctime" yaml:"ctime" xml:"ctime,attr"` // creation time
	MRTP   float64   `xorm:"notnull" json:"mrtp" yaml:"mrtp" xml:"mrtp"`                                                   // master RTP
	Game   string    `xorm:"notnull" json:"game" yaml:"game" xml:"game"`                                                   // game data
	Wins   string    `xorm:"text" json:"wins,omitempty" yaml:"wins,omitempty" xml:"wins,omitempty"`                        // list of wins marshaled to JSON
	Gain   float64   `xorm:"notnull" json:"gain" yaml:"gain" xml:"gain"`                                                   // total gain at last spin
	Wallet float64   `xorm:"notnull" json:"wallet" yaml:"wallet" xml:"wallet"`
}

var SpinCounter atomic.Uint64 // last spin log ID

type Multlog struct {
	MLID   uint64    `xorm:"pk" json:"mlid" yaml:"mlid" xml:"mlid,attr"`
	GID    uint64    `xorm:"notnull" json:"gid" yaml:"gid" xml:"gid,attr"` // game ID
	CTime  time.Time `xorm:"created 'ctime' notnull default CURRENT_TIMESTAMP" json:"ctime" yaml:"ctime" xml:"ctime,attr"`
	MRTP   float64   `xorm:"notnull" json:"mrtp" yaml:"mrtp" xml:"mrtp,attr"`  // master RTP
	Mult   float64   `xorm:"notnull" json:"mult" yaml:"mult" xml:"mult"`       // multiplier
	Risk   float64   `xorm:"notnull" json:"risk" yaml:"risk" xml:"risk"`       // the amount that is being gambled out
	Win    bool      `xorm:"notnull" json:"win" yaml:"win" xml:"win"`          // double up is win
	Gain   float64   `xorm:"notnull" json:"gain" yaml:"gain" xml:"gain"`       // total gain after double up
	Wallet float64   `xorm:"notnull" json:"wallet" yaml:"wallet" xml:"wallet"` // wallet after double up
}

var MultCounter atomic.Uint64 // last multiplier log ID

type Walletlog struct {
	WID    uint64    `xorm:"pk autoincr" json:"wid" yaml:"wid" xml:"wid,attr"`
	CID    uint64    `xorm:"notnull index(bid)" json:"cid" yaml:"cid" xml:"cid,attr"`                                      // club ID
	UID    uint64    `xorm:"notnull index(bid)" json:"uid" yaml:"uid" xml:"uid,attr"`                                      // user ID
	AID    uint64    `xorm:"notnull" json:"aid" yaml:"aid" xml:"aid,attr"`                                                 // admin ID
	CTime  time.Time `xorm:"created 'ctime' notnull default CURRENT_TIMESTAMP" json:"ctime" yaml:"ctime" xml:"ctime,attr"` // creation time
	Wallet float64   `xorm:"notnull" json:"wallet" yaml:"wallet" xml:"wallet"`                                             // new value in coins
	Sum    float64   `xorm:"notnull" json:"sum" yaml:"sum" xml:"sum"`
}

type Banklog struct {
	BID     uint64    `xorm:"pk autoincr" json:"bid" yaml:"bid" xml:"bid,attr"`
	CTime   time.Time `xorm:"created 'ctime' notnull default CURRENT_TIMESTAMP" json:"ctime" yaml:"ctime" xml:"ctime,attr"`
	Bank    float64   `xorm:"notnull 'bank'" json:"bank" yaml:"bank" xml:"bank"`
	Fund    float64   `xorm:"notnull 'fund'" json:"fund" yaml:"fund" xml:"fund"`
	Lock    float64   `xorm:"notnull 'lock'" json:"lock" yaml:"lock" xml:"lock"`
	BankSum float64   `xorm:"notnull 'banksum'" json:"banksum" yaml:"banksum" xml:"banksum" form:"banksum"`
	FundSum float64   `xorm:"notnull 'fundsum'" json:"fundsum" yaml:"fundsum" xml:"fundsum" form:"fundsum"`
	LockSum float64   `xorm:"notnull 'locksum'" json:"locksum" yaml:"locksum" xml:"locksum" form:"locksum"`
}

var SessCounter atomic.Uint64 // last session ID

type Session struct {
	SID    uint64    `xorm:"pk" json:"sid" yaml:"sid" xml:"sid,attr"`
	CTime  time.Time `xorm:"created 'ctime' notnull default CURRENT_TIMESTAMP" json:"ctime" yaml:"ctime" xml:"ctime,attr"`
	UTime  time.Time `xorm:"updated 'utime' notnull default CURRENT_TIMESTAMP" json:"utime" yaml:"utime" xml:"utime,attr"`
	Access string    `xorm:"notnull unique index" json:"access" yaml:"access" xml:"access"`
	Refrsh string    `xorm:"notnull unique index" json:"refrsh" yaml:"refrsh" xml:"refrsh"`
	UAgent string    `xorm:"notnull" json:"uagent" yaml:"uagent" xml:"uagent"`
}

// All created clubs, by CID.
var Clubs util.RWMap[uint64, *Club]

// All registered users, by UID.
var Users util.RWMap[uint64, *User]

// Scenes cache, by GID.
var Scenes util.RWMap[uint64, *Scene]

func MakeClub(cd ClubData) *Club {
	return &Club{data: cd}
}

func (club *Club) Get() ClubData {
	club.mux.RLock()
	defer club.mux.RUnlock()
	return club.data
}

func (club *Club) CID() uint64 { // read only
	return club.data.CID
}

func (club *Club) Name() string {
	club.mux.RLock()
	defer club.mux.RUnlock()
	return club.data.Name
}

func (club *Club) SetName(name string) {
	club.mux.Lock()
	defer club.mux.Unlock()
	club.data.Name = name
}

func (club *Club) Bank() float64 {
	club.mux.RLock()
	defer club.mux.RUnlock()
	return club.data.Bank
}

func (club *Club) Fund() float64 {
	club.mux.RLock()
	defer club.mux.RUnlock()
	return club.data.Fund
}

func (club *Club) Deposit() float64 {
	club.mux.RLock()
	defer club.mux.RUnlock()
	return club.data.Lock
}

func (club *Club) GetCash() (bank, fund, deposit float64) {
	club.mux.RLock()
	defer club.mux.RUnlock()
	return club.data.Bank, club.data.Fund, club.data.Lock
}

func (club *Club) AddCash(bank, fund, deposit float64) {
	club.mux.Lock()
	defer club.mux.Unlock()
	club.data.Bank += bank
	club.data.Fund += fund
	club.data.Lock += deposit
}

func (club *Club) AddBank(bank float64) {
	club.mux.Lock()
	defer club.mux.Unlock()
	club.data.Bank += bank
}

func (club *Club) AddFund(fund float64) {
	club.mux.Lock()
	defer club.mux.Unlock()
	club.data.Fund += fund
}

func (club *Club) AddDeposit(deposit float64) {
	club.mux.Lock()
	defer club.mux.Unlock()
	club.data.Lock += deposit
}

func (club *Club) Rate() float64 {
	club.mux.RLock()
	defer club.mux.RUnlock()
	return club.data.Rate
}

func (club *Club) MRTP() float64 {
	club.mux.RLock()
	defer club.mux.RUnlock()
	return club.data.MRTP
}

func (user *User) Init() {
	user.props.Init(0)
}

func (user *User) GetProps(cid uint64) (*Props, bool) {
	return user.props.Get(cid)
}

func (user *User) Items() iter.Seq2[uint64, *Props] {
	return user.props.Items()
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

func (user *User) GetRTP(cid uint64) float64 {
	if props, ok := user.props.Get(cid); ok {
		return props.MRTP
	}
	return 0
}

func (user *User) InsertProps(props *Props) {
	user.props.Set(props.CID, props)
}

// GetAdmin returns User pointer for authorized requests,
// and access level for it. Or nil pointer for unauthorized requests.
// It called after Auth(false) middleware.
func GetAdmin(c *gin.Context, cid uint64) (*User, AL) {
	if value, exists := c.Get(userKey); exists {
		var admin = value.(*User)
		return admin, admin.GAL | admin.GetAL(cid)
	}
	return nil, 0
}

// MustAdmin always returns User pointer for authorized
// requests, and access level for it.
// It called after Auth(true) middleware.
func MustAdmin(c *gin.Context, cid uint64) (*User, AL) {
	var admin = c.MustGet(userKey).(*User)
	return admin, admin.GAL | admin.GetAL(cid)
}

func GetRTP(user *User, club *Club) float64 {
	if user != nil {
		if props, ok := user.props.Get(club.CID()); ok && props.MRTP != 0 {
			return props.MRTP
		}
	}
	if mrtp := club.MRTP(); mrtp != 0 {
		return mrtp
	}
	return cfg.DefMRTP // default master RTP if no others found
}

func GetScene(gid uint64) (scene *Scene, err error) {
	var ok bool
	if scene, ok = Scenes.Get(gid); ok {
		return
	}

	var tmp Scene
	if ok, _ = XormStorage.ID(gid).Get(&tmp.Story); !ok {
		err = ErrNotOpened
		return
	}
	var maker func() game.Gamble
	if maker, ok = game.GameFactory[tmp.Alias]; !ok {
		err = ErrNoAlias
		return
	}
	tmp.Game = maker()

	scene = &tmp
	Scenes.Set(gid, scene)

	if !Cfg.UseSpinLog {
		return
	}

	var rec Spinlog
	if ok, _ = XormSpinlog.Where("gid = ?", gid).Desc("ctime").Get(&rec); !ok {
		return
	}
	scene.SLID = rec.SLID
	err = json.Unmarshal(util.S2B(rec.Game), scene.Game)
	return
}

func init() {
	Clubs.Init(0)
	Users.Init(0)
	Scenes.Init(0)
}
