package cmd

import (
	"fmt"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/spi"

	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

const (
	slotfile = "slot.sqlite"
)

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

	if err = session.Sync(&spi.Room{}, &spi.User{}, &spi.Props{}, &spi.Spinlog{}); err != nil {
		return
	}

	var ok bool
	if ok, err = session.IsTableEmpty(&spi.Room{}); err != nil {
		return
	}
	if ok {
		if _, err = session.Insert(&spi.Room{
			RID:  1,
			Bank: 10000,
			Fund: 1000000,
			Lock: 0,
			Name: "virtual",
		}); err != nil {
			return
		}
		if _, err = session.Insert(&spi.User{
			UID:    1,
			Email:  "example@example.org",
			Secret: "pGjkSD",
			Name:   "admin",
		}); err != nil {
			return
		}
		if _, err = session.Insert(&spi.Props{
			UID:    1,
			RID:    0,
			Wallet: 0,
			Access: spi.ALall,
		}); err != nil {
			return
		}
		if _, err = session.Insert(&spi.Props{
			UID:    1,
			RID:    1,
			Wallet: 1000,
			Access: spi.ALall,
		}); err != nil {
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
		var chunk []*spi.Props
		if err = session.Limit(limit, offset).Find(&chunk); err != nil {
			return
		}
		offset += limit
		for _, props := range chunk {
			if props.RID != 0 && !spi.Rooms.Has(props.RID) {
				return fmt.Errorf("found props without room linkage, UID=%d, RID=%d, value=%d", props.UID, props.RID, props.Wallet)
			}
			var user, ok = spi.Users.Get(props.UID)
			if !ok {
				return fmt.Errorf("found props without user linkage, UID=%d, RID=%d, value=%d", props.UID, props.RID, props.Wallet)
			}
			user.InsertProps(props)
		}
		if limit > len(chunk) {
			break
		}
	}
	return
}

func Init() (err error) {
	if err = InitStorage(); err != nil {
		return fmt.Errorf("can not init XORM records storage: %w", err)
	}
	return
}
