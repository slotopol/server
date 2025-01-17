package api

import (
	"encoding/xml"
	"sync/atomic"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/game"
	"github.com/slotopol/server/game/keno"
	"github.com/slotopol/server/game/slot"
	"github.com/slotopol/server/util"

	"github.com/gin-gonic/gin"
)

var (
	SpinBuf util.SqlBuf[Spinlog]
	MultBuf util.SqlBuf[Multlog]
	BankBat = map[uint64]*SqlBank{}
	JoinBuf = SqlStory{}
)

// Joins to game and creates new instance of game.
func ApiGameJoin(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr" form:"cid" binding:"required"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid,attr" form:"uid" binding:"required"`
		Alias   string   `json:"alias" yaml:"alias" xml:"alias" form:"alias" binding:"required"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr"`
		Game    any      `json:"game" yaml:"game" xml:"game"`
		Wallet  float64  `json:"wallet" yaml:"wallet" xml:"wallet"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_game_join_nobind, err)
		return
	}

	var club *Club
	if club, ok = Clubs.Get(arg.CID); !ok {
		Ret404(c, AEC_game_join_noclub, ErrNoClub)
		return
	}

	var user *User
	if user, ok = Users.Get(arg.UID); !ok {
		Ret404(c, AEC_game_join_nouser, ErrNoUser)
		return
	}

	var admin, al = MustAdmin(c, arg.CID)
	if (al&ALmem == 0) || (admin != user && al&ALgame == 0) {
		Ret403(c, AEC_game_join_noaccess, ErrNoAccess)
		return
	}

	var alias = util.ToID(arg.Alias)
	var maker, has = game.GameFactory[alias]
	if !has {
		Ret400(c, AEC_game_join_noalias, ErrNoAliase)
		return
	}

	var anygame = maker()
	var gid = atomic.AddUint64(&StoryCounter, 1)
	var scene = &Scene{
		Story: Story{
			GID:   gid,
			Alias: alias,
			CID:   arg.CID,
			UID:   arg.UID,
		},
		Game: anygame,
	}

	// insert new story entry
	if Cfg.ClubInsertBuffer > 1 {
		go JoinBuf.Join(cfg.XormStorage, &scene.Story)
	} else if err = JoinBuf.Join(cfg.XormStorage, &scene.Story); err != nil {
		Ret500(c, AEC_game_join_sql, err)
		return
	}

	Scenes.Set(scene.GID, scene)

	ret.GID = gid
	ret.Game = anygame
	// make game screen object
	var rtp = GetRTP(user, club)
	switch game := anygame.(type) {
	case slot.SlotGame:
		game.Spin(rtp)
	case keno.KenoGame:
		game.Spin(rtp)
	}
	ret.Wallet = user.GetWallet(arg.CID)

	RetOk(c, ret)
}

// Removes instance of opened game.
func ApiGamePart(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" form:"gid" binding:"required"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_game_part_nobind, err)
		return
	}

	var scene *Scene
	if scene, err = GetScene(arg.GID); err != nil {
		Ret404(c, AEC_game_part_noscene, err)
		return
	}

	var user *User
	if user, ok = Users.Get(scene.UID); !ok {
		Ret500(c, AEC_game_part_nouser, ErrNoUser)
		return
	}

	var admin, al = MustAdmin(c, scene.CID)
	if admin != user && al&ALgame == 0 {
		Ret403(c, AEC_game_part_noaccess, ErrNoAccess)
		return
	}

	Scenes.Delete(arg.GID)

	Ret204(c)
}

// Returns full info of game scene with given GID, and balance on wallet.
func ApiGameInfo(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" form:"gid" binding:"required"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr"`
		Alias   string   `json:"alias" yaml:"alias" xml:"alias"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid,attr"`
		SID     uint64   `json:"sid" yaml:"sid" xml:"sid,attr"`
		Game    any      `json:"game" yaml:"game" xml:"game"`
		Wallet  float64  `json:"wallet" yaml:"wallet" xml:"wallet"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_game_info_nobind, err)
		return
	}

	var scene *Scene
	if scene, err = GetScene(arg.GID); err != nil {
		Ret404(c, AEC_game_info_noscene, err)
		return
	}

	var user *User
	if user, ok = Users.Get(scene.UID); !ok {
		Ret500(c, AEC_game_info_nouser, ErrNoUser)
		return
	}

	var admin, al = MustAdmin(c, scene.CID)
	if admin != user && al&ALgame == 0 {
		Ret403(c, AEC_game_info_noaccess, ErrNoAccess)
		return
	}

	var props *Props
	if props, ok = user.props.Get(scene.CID); !ok {
		Ret500(c, AEC_game_info_noprops, ErrNoProps)
		return
	}

	ret.GID = arg.GID
	ret.Alias = scene.Alias
	ret.CID = scene.CID
	ret.UID = scene.UID
	ret.SID = scene.SID
	ret.Game = scene.Game
	ret.Wallet = props.Wallet

	RetOk(c, ret)
}

// Returns master RTP for given GID.
func ApiGameRtpGet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" form:"gid" binding:"required"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		MRTP    float64  `json:"mrtp" yaml:"mrtp" xml:"mrtp"`
		RTP     float64  `json:"rtp" yaml:"rtp" xml:"rtp"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_game_rtpget_nobind, err)
		return
	}

	var scene *Scene
	if scene, err = GetScene(arg.GID); err != nil {
		Ret404(c, AEC_game_rtpget_noscene, err)
		return
	}

	var gi *game.GameInfo
	if gi = game.GetInfo(scene.Alias); gi == nil {
		Ret500(c, AEC_game_rtpget_noinfo, ErrNoAliase)
		return
	}

	var club *Club
	if club, ok = Clubs.Get(scene.CID); !ok {
		Ret500(c, AEC_game_rtpget_noclub, ErrNoClub)
		return
	}

	var user *User
	if user, ok = Users.Get(scene.UID); !ok {
		Ret500(c, AEC_game_rtpget_nouser, ErrNoUser)
		return
	}

	var admin, al = MustAdmin(c, scene.CID)
	if admin != user && al&ALgame == 0 {
		Ret403(c, AEC_game_rtpget_noaccess, ErrNoAccess)
		return
	}

	ret.MRTP = GetRTP(user, club)
	ret.RTP = gi.FindRTP(ret.MRTP)

	RetOk(c, ret)
}
