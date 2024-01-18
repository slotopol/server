package spi

import (
	"encoding/xml"
	"net/http"

	"github.com/gin-gonic/gin"
	cfg "github.com/slotopol/server/config"
)

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

	if _, err = cfg.XormStorage.Cols("name").Update(&User{UID: arg.UID, Name: arg.Name}); err != nil {
		Ret500(c, SEC_user_rename_update, ErrNoUser)
		return
	}
	user.Name = arg.Name

	c.Status(http.StatusOK)
}
