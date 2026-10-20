package api

import (
	"encoding/xml"
	"fmt"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/slotopol/server/core"
	"github.com/slotopol/server/game"
	"github.com/slotopol/server/util"
)

// Returns full list of all available algorithms.
func ApiGameAlgs(c *gin.Context) {
	RetOk(c, game.AlgList)
}

// List of available games selected by filters.
func ApiGameList(c *gin.Context) {
	var err error
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		Include string   `json:"include" yaml:"include" xml:"include" form:"inc"`
		Exclude string   `json:"exclude" yaml:"exclude" xml:"exclude" form:"exc"`
		Sort    bool     `json:"sort" yaml:"sort" xml:"sort" form:"sort"`
	}
	var ret struct {
		XMLName xml.Name         `json:"-" yaml:"-" xml:"ret"`
		List    []*game.GameInfo `json:"list" yaml:"list" xml:"list>gi"`
		AlgNum  int              `json:"algnum" yaml:"algnum" xml:"algnum"`
		PrvNum  int              `json:"prvnum" yaml:"prvnum" xml:"prvnum"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_game_list_nobind, err)
		return
	}
	if len(arg.Include) == 0 {
		arg.Include = "all"
	}
	var include = strings.Split(arg.Include, " ")
	var exclude = strings.Split(arg.Exclude, " ")

	var finclist, fexclist [][]game.Filter
	var f game.Filter
	var flist []game.Filter
	for _, inc := range include {
		var keys = strings.Split(inc, "&")
		flist = nil
		for _, key := range keys {
			if f = game.GetFilter(key); f == nil {
				Ret400(c, AEC_game_list_inc, fmt.Errorf("filter with name '%s' does not recognized", key))
				return
			}
			flist = append(flist, f)
		}
		finclist = append(finclist, flist)
	}
	for _, exc := range exclude {
		var keys = strings.Split(exc, "&")
		flist = nil
		for _, key := range keys {
			if f = game.GetFilter(key); f == nil {
				Ret400(c, AEC_game_list_exc, fmt.Errorf("filter with name '%s' does not recognized", key))
				return
			}
		}
		fexclist = append(fexclist, flist)
	}

	var alg = map[*game.AlgDescr]int{}
	var prov = map[string]int{}
	var gamelist = make([]*game.GameInfo, 0, 256)
	for _, gi := range game.InfoMap {
		if game.Passes(gi, finclist, fexclist) {
			alg[gi.AlgDescr]++
			prov[util.ToID(gi.Prov)]++
			gamelist = append(gamelist, gi)
		}
	}

	sort.Slice(gamelist, func(i, j int) bool {
		var gii, gij = gamelist[i], gamelist[j]
		if arg.Sort {
			if gii.Prov == gij.Prov {
				return gii.Name < gij.Name
			}
			return gii.Prov < gij.Prov
		} else {
			if gii.Name == gij.Name {
				return gii.Prov < gij.Prov
			}
			return gii.Name < gij.Name
		}
	})

	ret.List = gamelist
	ret.AlgNum = len(alg)
	ret.PrvNum = len(prov)

	RetOk(c, ret)
}

// Creates new instance of game.
func ApiGameNew(c *gin.Context) {
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

	var club *core.Club
	if club, ok = core.Clubs.Get(arg.CID); !ok {
		Ret404(c, AEC_game_new_noclub, ErrNoClub)
		return
	}
	_ = club

	var user *core.User
	if user, ok = core.Users.Get(arg.UID); !ok {
		Ret404(c, AEC_game_new_nouser, ErrNoUser)
		return
	}

	var admin, al = core.MustAdmin(c, arg.CID)
	if (al&core.ALmember == 0) || (admin != user && al&core.ALdealer == 0) {
		Ret403(c, AEC_game_new_noaccess, ErrNoAccess)
		return
	}

	var scene *core.Scene
	var alias = util.ToID(arg.Alias)
	var maker, has = game.GameFactory[alias]
	if !has {
		Ret400(c, AEC_game_new_noalias, ErrNoAlias)
		return
	}

	var anygame = maker()
	var gid = core.StoryCounter.Inc()
	scene = &core.Scene{
		Story: core.Story{
			GID:   gid,
			CID:   arg.CID,
			UID:   arg.UID,
			Alias: alias,
		},
		Game: anygame,
	}

	// insert new story entry
	if Cfg.ClubInsertBuffer > 1 {
		core.XormWG.Go(func() (err error) {
			return core.JoinBuf.Join(&scene.Story)
		})
	} else if err = core.JoinBuf.Join(&scene.Story); err != nil {
		Ret500(c, AEC_game_new_sql, err)
		return
	}

	core.Scenes.Set(scene.GID, scene)

	ret.GID = scene.GID
	ret.Game = scene.Game
	ret.Wallet = user.GetWallet(arg.CID)

	RetOk(c, ret)
}

// Joins to game and creates new instance of game.
func ApiGameJoin(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr" form:"cid" binding:"required"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid,attr" form:"uid" binding:"required"`
		GID     uint64   `json:"gid" yaml:"gid" xml:"gid,attr" form:"gid" binding:"required"`
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

	var user *core.User
	if user, ok = core.Users.Get(arg.UID); !ok {
		Ret404(c, AEC_game_join_nouser, ErrNoUser)
		return
	}

	var admin, al = core.MustAdmin(c, arg.CID)
	if (al&core.ALmember == 0) || (admin != user && al&core.ALdealer == 0) {
		Ret403(c, AEC_game_join_noaccess, ErrNoAccess)
		return
	}

	var scene *core.Scene
	if scene, err = core.GetScene(arg.GID); err != nil {
		Ret404(c, AEC_game_join_noscene, err)
		return
	}

	ret.GID = scene.GID
	ret.Game = scene.Game
	ret.Wallet = user.GetWallet(arg.CID)

	RetOk(c, ret)
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
		SLID    uint64   `json:"slid" yaml:"slid" xml:"slid,attr"`
		Game    any      `json:"game" yaml:"game" xml:"game"`
		Wallet  float64  `json:"wallet" yaml:"wallet" xml:"wallet"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_game_info_nobind, err)
		return
	}

	var scene *core.Scene
	if scene, err = core.GetScene(arg.GID); err != nil {
		Ret404(c, AEC_game_info_noscene, err)
		return
	}

	var user *core.User
	if user, ok = core.Users.Get(scene.UID); !ok {
		Ret500(c, AEC_game_info_nouser, ErrNoUser)
		return
	}

	var admin, al = core.MustAdmin(c, scene.CID)
	if admin != user && al&core.ALdealer == 0 {
		Ret403(c, AEC_game_info_noaccess, ErrNoAccess)
		return
	}

	var props *core.Props
	if props, ok = user.GetProps(scene.CID); !ok {
		Ret500(c, AEC_game_info_noprops, ErrNoProps)
		return
	}

	ret.GID = arg.GID
	ret.Alias = scene.Alias
	ret.CID = scene.CID
	ret.UID = scene.UID
	ret.SLID = scene.SLID
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

	var scene *core.Scene
	if scene, err = core.GetScene(arg.GID); err != nil {
		Ret404(c, AEC_game_rtpget_noscene, err)
		return
	}

	var gi *game.GameInfo
	if gi, ok = game.InfoMap[scene.Alias]; !ok {
		Ret500(c, AEC_game_rtpget_noinfo, ErrBadAlias)
		return
	}

	var club *core.Club
	if club, ok = core.Clubs.Get(scene.CID); !ok {
		Ret500(c, AEC_game_rtpget_noclub, ErrNoClub)
		return
	}

	var user *core.User
	if user, ok = core.Users.Get(scene.UID); !ok {
		Ret500(c, AEC_game_rtpget_nouser, ErrNoUser)
		return
	}

	var admin, al = core.MustAdmin(c, scene.CID)
	if admin != user && al&core.ALdealer == 0 {
		Ret403(c, AEC_game_rtpget_noaccess, ErrNoAccess)
		return
	}

	ret.MRTP = core.GetRTP(user, club)
	ret.RTP = gi.FindClosest(ret.MRTP)

	RetOk(c, ret)
}
