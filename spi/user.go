package spi

import (
	"encoding/xml"
	"net/http"

	"github.com/gin-gonic/gin"
	cfg "github.com/slotopol/server/config"
	"xorm.io/xorm"
)

// Changes 'Name' of given user.
func SpiUserRename(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid,attr" form:"uid"`
		Name    string   `json:"name" yaml:"name" xml:"name,attr" form:"name"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_user_rename_nobind, err)
		return
	}
	if arg.UID == 0 {
		Ret400(c, SEC_user_rename_nouid, ErrNoUID)
		return
	}
	if arg.Name == "" {
		Ret400(c, SEC_user_rename_noname, ErrNoData)
		return
	}

	var user *User
	if user, ok = Users.Get(arg.UID); !ok {
		Ret404(c, SEC_user_rename_nouser, ErrNoUser)
		return
	}

	var admin, al = GetAdmin(c, 0)
	if admin != user && al&ALadmin == 0 {
		Ret403(c, SEC_prop_rename_noaccess, ErrNoAccess)
		return
	}

	if _, err = cfg.XormStorage.Cols("name").Update(&User{UID: arg.UID, Name: arg.Name}); err != nil {
		Ret500(c, SEC_user_rename_update, err)
		return
	}
	user.Name = arg.Name

	c.Status(http.StatusOK)
}

func SpiUserSecret(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName   xml.Name `json:"-" yaml:"-" xml:"arg"`
		UID       uint64   `json:"uid" yaml:"uid" xml:"uid,attr" form:"uid"`
		OldSecret string   `json:"oldsecret" yaml:"oldsecret" xml:"oldsecret" form:"oldsecret"`
		NewSecret string   `json:"newsecret" yaml:"newsecret" xml:"newsecret" form:"newsecret"`
	}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_user_secret_nobind, err)
		return
	}
	if arg.UID == 0 {
		Ret400(c, SEC_user_secret_nouid, ErrNoUID)
		return
	}
	if len(arg.NewSecret) < 6 {
		Ret400(c, SEC_user_secret_smallsec, ErrSmallKey)
		return
	}

	var user *User
	if user, ok = Users.Get(arg.UID); !ok {
		Ret404(c, SEC_user_secret_nouser, ErrNoUser)
		return
	}

	var admin, al = GetAdmin(c, 0)
	if admin != user && al&ALadmin == 0 {
		Ret403(c, SEC_prop_secret_noaccess, ErrNoAccess)
		return
	}

	if arg.OldSecret != user.Secret && al&ALadmin == 0 {
		Ret403(c, SEC_prop_secret_nosecret, ErrNoSecret)
		return
	}

	if _, err = cfg.XormStorage.Cols("secret").Update(&User{UID: arg.UID, Secret: arg.NewSecret}); err != nil {
		Ret500(c, SEC_user_secret_update, err)
		return
	}
	user.Secret = arg.NewSecret

	c.Status(http.StatusOK)
}

// Deletes registration, drops user and all linked records from database,
// and moves all remained coins at wallets to rooms deposits.
func SpiUserDelete(c *gin.Context) {
	var err error
	var ok bool
	var arg struct {
		XMLName xml.Name `json:"-" yaml:"-" xml:"arg"`
		UID     uint64   `json:"uid" yaml:"uid" xml:"uid,attr" form:"uid"`
		Secret  string   `json:"secret" yaml:"secret" xml:"secret" form:"secret"`
	}
	var ret struct {
		XMLName xml.Name       `json:"-" yaml:"-" xml:"ret"`
		Wallets map[uint64]int `json:"wallets" yaml:"wallets" xml:"wallets"`
	}
	ret.Wallets = map[uint64]int{}

	if err = c.ShouldBind(&arg); err != nil {
		Ret400(c, SEC_user_delete_nobind, err)
		return
	}
	if arg.UID == 0 {
		Ret400(c, SEC_user_delete_nouid, ErrNoUID)
		return
	}

	var user *User
	if user, ok = Users.Get(arg.UID); !ok {
		Ret404(c, SEC_user_delete_nouser, ErrNoUser)
		return
	}

	var admin, al = GetAdmin(c, 0)
	if admin != user && al&ALadmin == 0 {
		Ret403(c, SEC_prop_delete_noaccess, ErrNoAccess)
		return
	}

	if arg.Secret != user.Secret && al&ALadmin == 0 {
		Ret403(c, SEC_prop_delete_nosecret, ErrNoSecret)
		return
	}

	// write gain and total bet as transaction
	if _, err = cfg.XormStorage.Transaction(func(session *xorm.Session) (_ interface{}, err error) {
		defer func() {
			if err != nil {
				session.Rollback()
			}
		}()

		if _, err = session.ID(arg.UID).Delete(user); err != nil {
			Ret500(c, SEC_prop_delete_sqluser, err)
			return
		}

		const sql1 = `UPDATE room SET lock=lock+? WHERE rid=?`
		if user.props.Range(func(rid uint64, props *Props) bool {
			if props.Wallet != 0 {
				if _, err = session.Exec(sql1, props.Wallet, rid); err != nil {
					Ret500(c, SEC_game_delete_sqllock, err)
					return false
				}
			}
			return true
		}); err != nil {
			return
		}

		if _, err = session.Where("uid=?", arg.UID).Delete(&Props{}); err != nil {
			Ret500(c, SEC_prop_delete_sqlprops, err)
			return
		}

		if _, err = session.Where("uid=?", arg.UID).Delete(&OpenGame{}); err != nil {
			Ret500(c, SEC_prop_delete_sqlgames, err)
			return
		}

		return
	}); err != nil {
		return
	}

	user.props.Range(func(rid uint64, props *Props) bool {
		ret.Wallets[rid] = props.Wallet
		if room, ok := Rooms.Get(rid); ok && props.Wallet != 0 {
			room.Lock += float64(props.Wallet)
			ret.Wallets[rid] = props.Wallet
		}
		return true
	})
	user.games.Range(func(gid uint64, og OpenGame) bool {
		OpenGames.Delete(gid)
		return true
	})
	Users.Delete(arg.UID)

	RetOk(c, ret)
}
