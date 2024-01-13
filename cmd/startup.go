package cmd

import (
	"fmt"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/spi"

	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

const slotfile = "slot.sqlite"

var (
	Cfg = cfg.Cfg
)

func InitStorage() (err error) {
	if cfg.XormStorage, err = xorm.NewEngine(Cfg.XormDriverName, cfg.JoinPath(cfg.CfgPath, slotfile)); err != nil {
		return
	}
	cfg.XormStorage.SetMapper(names.GonicMapper{})

	var session = cfg.XormStorage.NewSession()
	defer session.Close()

	if err = session.Sync(&spi.Room{}, &spi.User{}); err != nil {
		return
	}

	var ok bool
	if ok, err = session.IsTableEmpty(&spi.Room{}); err != nil {
		return
	}
	if ok {
		var room = spi.Room{
			RID:  1,
			Bank: 0,
			Fund: 0,
		}
		if _, err = session.Insert(&room); err != nil {
			return
		}
		var user = spi.User{
			UID:     1,
			RID:     1,
			Balance: 1000,
			Email:   "example@example.org",
			Secret:  "admin",
			Name:    "admin",
		}
		if _, err = session.Insert(&user); err != nil {
			return
		}
	}
	return
}

func Init() (err error) {
	if err = InitStorage(); err != nil {
		return fmt.Errorf("can not init XORM storage: %w", err)
	}
	return
}
