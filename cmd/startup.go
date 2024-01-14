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

	if err = session.Sync(&spi.Room{}, &spi.User{}, &spi.Balance{}); err != nil {
		return
	}

	var ok bool
	if ok, err = session.IsTableEmpty(&spi.Room{}); err != nil {
		return
	}
	if ok {
		var room = spi.Room{
			RID:  1,
			Bank: 10000,
			Fund: 1000000,
			Name: "virtual",
		}
		if _, err = session.Insert(&room); err != nil {
			return
		}
		var user = spi.User{
			UID:    1,
			Email:  "example@example.org",
			Secret: "pGjkSD",
			Name:   "admin",
		}
		if _, err = session.Insert(&user); err != nil {
			return
		}
		var balance = spi.Balance{
			UID:   1,
			RID:   1,
			Value: 1000,
		}
		if _, err = session.Insert(&balance); err != nil {
			return
		}
	}

	const limit = 256

	var offset = 0
	for {
		var chunk []*spi.Room
		if err = session.Limit(limit, offset).Find(&chunk); err != nil {
			return
		}
		offset += limit
		for _, room := range chunk {
			spi.Rooms.Set(room.RID, room)
		}
		if limit > len(chunk) {
			break
		}
	}

	offset = 0
	for {
		var chunk []*spi.User
		if err = session.Limit(limit, offset).Find(&chunk); err != nil {
			return
		}
		offset += limit
		for _, user := range chunk {
			user.Init()
			spi.Users.Set(user.UID, user)
		}
		if limit > len(chunk) {
			break
		}
	}

	offset = 0
	for {
		var chunk []*spi.Balance
		if err = session.Limit(limit, offset).Find(&chunk); err != nil {
			return
		}
		offset += limit
		for _, bal := range chunk {
			if !spi.Rooms.Has(bal.RID) {
				return fmt.Errorf("found balance without room linkage, UID=%d, RID=%d, value=%d", bal.UID, bal.RID, bal.Value)
			}
			var user, ok = spi.Users.Get(bal.UID)
			if !ok {
				return fmt.Errorf("found balance without user linkage, UID=%d, RID=%d, value=%d", bal.UID, bal.RID, bal.Value)
			}
			user.InsertBalance(bal)
		}
		if limit > len(chunk) {
			break
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
