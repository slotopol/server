package api

import (
	"encoding/xml"

	"github.com/gin-gonic/gin"
	cfg "github.com/slotopol/server/config"
)

// Returns all properties for pointed user at pointed club.
func ApiPropsGet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr" form:"cid" binding:"required"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid,attr" form:"uid" binding:"required"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		Wallet  float64  `json:"wallet" yaml:"wallet" xml:"wallet"`
		Access  AL       `json:"access" yaml:"access" xml:"access"`
		MRTP    float64  `json:"mrtp" yaml:"mrtp" xml:"mrtp"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_prop_get_nobind, err)
		return
	}

	var club *Club
	if club, ok = Clubs.Get(arg.CID); !ok {
		Ret404(c, AEC_prop_get_noclub, ErrNoClub)
		return
	}
	_ = club

	var user *User
	if user, ok = Users.Get(arg.UID); !ok {
		Ret404(c, AEC_prop_get_nouser, ErrNoUser)
		return
	}

	var admin, al = MustAdmin(c, arg.CID)
	if admin != user && al&ALbooker == 0 {
		Ret403(c, AEC_prop_get_noaccess, ErrNoAccess)
		return
	}

	if props, ok := user.props.Get(arg.CID); ok {
		ret.Wallet = props.Wallet
		ret.Access = props.Access
		ret.MRTP = props.MRTP
	}

	RetOk(c, ret)
}

// Returns balance at wallet for pointed user at pointed club.
func ApiPropsWalletGet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr" form:"cid" binding:"required"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid,attr" form:"uid" binding:"required"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		Wallet  float64  `json:"wallet" yaml:"wallet" xml:"wallet"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_prop_walletget_nobind, err)
		return
	}

	var club *Club
	if club, ok = Clubs.Get(arg.CID); !ok {
		Ret404(c, AEC_prop_walletget_noclub, ErrNoClub)
		return
	}
	_ = club

	var user *User
	if user, ok = Users.Get(arg.UID); !ok {
		Ret404(c, AEC_prop_walletget_nouser, ErrNoUser)
		return
	}

	var admin, al = MustAdmin(c, arg.CID)
	if admin != user && al&ALbooker == 0 {
		Ret403(c, AEC_prop_walletget_noaccess, ErrNoAccess)
		return
	}

	ret.Wallet = user.GetWallet(arg.CID)

	RetOk(c, ret)
}

// Adds some coins to user wallet. Sum can be < 0 to remove some coins.
func ApiPropsWalletAdd(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr" form:"cid" binding:"required"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid,attr" form:"uid" binding:"required"`
		Sum     float64  `json:"sum" yaml:"sum" xml:"sum" binding:"required"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		Wallet  float64  `json:"wallet" yaml:"wallet" xml:"wallet"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_prop_walletadd_nobind, err)
		return
	}
	if arg.Sum > cfg.Cfg.AdjunctLimit || arg.Sum < -cfg.Cfg.AdjunctLimit {
		Ret400(c, AEC_prop_walletadd_limit, ErrTooBig)
		return
	}

	var club *Club
	if club, ok = Clubs.Get(arg.CID); !ok {
		Ret404(c, AEC_prop_walletadd_noclub, ErrNoClub)
		return
	}
	_ = club

	var user *User
	if user, ok = Users.Get(arg.UID); !ok {
		Ret404(c, AEC_prop_walletadd_nouser, ErrNoUser)
		return
	}

	var admin, al = MustAdmin(c, arg.CID)
	if al&ALbooker == 0 {
		Ret403(c, AEC_prop_walletadd_noaccess, ErrNoAccess)
		return
	}

	var props *Props
	if props, ok = user.props.Get(arg.CID); !ok {
		Ret500(c, AEC_prop_walletadd_noprops, ErrNoProps)
		return
	}
	if props.Wallet+arg.Sum < 0 {
		Ret403(c, AEC_prop_walletadd_nomoney, ErrNoMoney)
		return
	}

	// update wallet as transaction
	if Cfg.ClubInsertBuffer > 1 {
		go BankBat[arg.CID].Add(cfg.XormStorage, arg.UID, admin.UID, props.Wallet+arg.Sum, arg.Sum)
	} else if err = BankBat[arg.CID].Add(cfg.XormStorage, arg.UID, admin.UID, props.Wallet+arg.Sum, arg.Sum); err != nil {
		Ret500(c, AEC_prop_walletadd_sql, err)
		return
	}

	// make changes to memory data
	props.Wallet += arg.Sum

	ret.Wallet = props.Wallet

	RetOk(c, ret)
}

// Returns personal access level for pointed user at pointed club.
func ApiPropsAlGet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr" form:"cid"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid,attr" form:"uid" binding:"required"`
		All     bool     `json:"all" yaml:"all" xml:"all,attr" form:"all"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		Access  AL       `json:"access" yaml:"access" xml:"access"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_prop_alget_nobind, err)
		return
	}

	var club *Club
	if club, ok = Clubs.Get(arg.CID); !ok && !arg.All {
		Ret404(c, AEC_prop_alget_noclub, ErrNoClub)
		return
	}
	_ = club

	var user *User
	if user, ok = Users.Get(arg.UID); !ok {
		Ret404(c, AEC_prop_alget_nouser, ErrNoUser)
		return
	}

	var admin, al = MustAdmin(c, arg.CID)
	if admin != user && al&(ALbooker+ALadmin) == 0 {
		Ret403(c, AEC_prop_alget_noaccess, ErrNoAccess)
		return
	}

	ret.Access = user.GetAL(arg.CID)
	if arg.All {
		ret.Access |= user.GAL
	}

	RetOk(c, ret)
}

// Set personal access level for given user at given club.
func ApiPropsAlSet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr" form:"cid" binding:"required"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid,attr" form:"uid" binding:"required"`
		Access  AL       `json:"access" yaml:"access" xml:"access"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_prop_alset_nobind, err)
		return
	}

	var club *Club
	if club, ok = Clubs.Get(arg.CID); !ok {
		Ret404(c, AEC_prop_alset_noclub, ErrNoClub)
		return
	}
	_ = club

	var user *User
	if user, ok = Users.Get(arg.UID); !ok {
		Ret404(c, AEC_prop_alset_nouser, ErrNoUser)
		return
	}

	var admin, al = MustAdmin(c, arg.CID)
	if al&ALadmin == 0 {
		Ret403(c, AEC_prop_alset_noaccess, ErrNoAccess)
		return
	}
	_ = admin

	var props *Props
	if props, ok = user.props.Get(arg.CID); !ok {
		Ret500(c, AEC_prop_alset_noprops, ErrNoProps)
		return
	}
	if al&arg.Access != arg.Access {
		Ret403(c, AEC_prop_alset_nolevel, ErrNoLevel)
		return
	}

	// update access level as transaction
	if Cfg.ClubInsertBuffer > 1 {
		go BankBat[arg.CID].Access(cfg.XormStorage, arg.UID, arg.Access)
	} else if err = BankBat[arg.CID].Access(cfg.XormStorage, arg.UID, arg.Access); err != nil {
		Ret500(c, AEC_prop_rtpset_sql, err)
		return
	}

	// make changes to memory data
	props.Access = arg.Access

	Ret204(c)
}

// Returns master RTP for pointed user at pointed club.
// This RTP if it set have more priority then club RTP.
func ApiPropsRtpGet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr" form:"cid" binding:"required"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid,attr" form:"uid"`
		All     bool     `json:"all" yaml:"all" xml:"all,attr" form:"all"`
	}
	var ret struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"ret"`
		MRTP    float64  `json:"mrtp" yaml:"mrtp" xml:"mrtp"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_prop_rtpget_nobind, err)
		return
	}

	var club *Club
	if club, ok = Clubs.Get(arg.CID); !ok {
		Ret404(c, AEC_prop_rtpget_noclub, ErrNoClub)
		return
	}
	_ = club

	var user *User
	if user, ok = Users.Get(arg.UID); !ok && !arg.All {
		Ret404(c, AEC_prop_rtpget_nouser, ErrNoUser)
		return
	}

	var admin, al = MustAdmin(c, arg.CID)
	if admin != user && al&ALbooker == 0 {
		Ret403(c, AEC_prop_rtpget_noaccess, ErrNoAccess)
		return
	}

	if arg.All {
		ret.MRTP = GetRTP(user, club)
	} else {
		ret.MRTP = user.GetRTP(arg.CID)
	}

	RetOk(c, ret)
}

// Set personal master RTP for given user at given club.
func ApiPropsRtpSet(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		CID     uint64   `json:"cid" yaml:"cid" xml:"cid,attr" form:"cid"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid,attr" form:"uid"`
		MRTP    float64  `json:"mrtp" yaml:"mrtp" xml:"mrtp"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, AEC_prop_rtpset_nobind, err)
		return
	}

	var club *Club
	if club, ok = Clubs.Get(arg.CID); !ok {
		Ret404(c, AEC_prop_rtpset_noclub, ErrNoClub)
		return
	}
	_ = club

	var user *User
	if user, ok = Users.Get(arg.UID); !ok {
		Ret404(c, AEC_prop_rtpset_nouser, ErrNoUser)
		return
	}

	var admin, al = MustAdmin(c, arg.CID)
	if al&ALbooker == 0 {
		Ret403(c, AEC_prop_rtpset_noaccess, ErrNoAccess)
		return
	}
	_ = admin

	var props *Props
	if props, ok = user.props.Get(arg.CID); !ok {
		Ret500(c, AEC_prop_rtpset_noprops, ErrNoProps)
		return
	}

	// update master RTP as transaction
	if Cfg.ClubInsertBuffer > 1 {
		go BankBat[arg.CID].MRTP(cfg.XormStorage, arg.UID, arg.MRTP)
	} else if err = BankBat[arg.CID].MRTP(cfg.XormStorage, arg.UID, arg.MRTP); err != nil {
		Ret500(c, AEC_prop_rtpset_sql, err)
		return
	}

	// make changes to memory data
	props.MRTP = arg.MRTP

	Ret204(c)
}
