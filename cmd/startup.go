package cmd

import (
	"fmt"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/spi"
	"github.com/slotopol/server/util"

	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

const (
	slotclubfile = "slot-club.sqlite"
	slotspinfile = "slot-spin.sqlite"
)

var (
	Cfg = cfg.Cfg
)

func InitStorage() (err error) {
	if cfg.XormStorage, err = xorm.NewEngine(Cfg.XormDriverName, util.JoinPath(cfg.SqlPath, slotclubfile)); err != nil {
		return
	}
	cfg.XormStorage.SetMapper(names.GonicMapper{})

	var session = cfg.XormStorage.NewSession()
	defer session.Close()

	if err = session.Sync(
		&spi.Club{}, &spi.User{}, &spi.Props{},
		&spi.OpenGame{}, spi.Walletlog{}, spi.Banklog{},
	); err != nil {
		return
	}

	var ok bool
	if ok, err = session.IsTableEmpty(&spi.Club{}); err != nil {
		return
	}
	if ok {
		if _, err = session.Insert(&spi.Club{
			CID:     1,
			Bank:    10000,
			Fund:    1000000,
			Lock:    0,
			Name:    "virtual",
			JptRate: 0.015,
			GainRTP: 95.00,
		}); err != nil {
			return
		}
		if _, err = session.Insert(&[]spi.User{
			{
				UID:    1,
				Email:  "admin@example.org",
				Secret: "pGjkSD",
				Name:   "admin",
				GAL:    spi.ALall,
			},
			{
				UID:    2,
				Email:  "dealer@example.org",
				Secret: "jpTyD4",
				Name:   "dealer",
				GAL:    spi.ALgame,
			},
			{
				UID:    3,
				Email:  "player@example.org",
				Secret: "Et7oAm",
				Name:   "player",
				GAL:    0,
			},
		}); err != nil {
			return
		}
		if _, err = session.Insert(&[]spi.Props{
			{
				UID:    2,
				CID:    1,
				Wallet: 12000,
				Access: spi.ALuser | spi.ALclub,
			},
			{
				UID:    3,
				CID:    1,
				Wallet: 1000,
				Access: 0,
			},
		}); err != nil {
			return
		}
	}

	const limit = 256

	var offset = 0
	for {
		var chunk []*spi.Club
		if err = session.Limit(limit, offset).Find(&chunk); err != nil {
			return
		}
		offset += limit
		for _, club := range chunk {
			spi.Clubs.Set(club.CID, club)
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
			if !spi.Clubs.Has(props.CID) {
				return fmt.Errorf("found props without club linkage, UID=%d, CID=%d, value=%g", props.UID, props.CID, props.Wallet)
			}
			var user, ok = spi.Users.Get(props.UID)
			if !ok {
				return fmt.Errorf("found props without user linkage, UID=%d, CID=%d, value=%g", props.UID, props.CID, props.Wallet)
			}
			user.InsertProps(props)
		}
		if limit > len(chunk) {
			break
		}
	}
	return
}

func InitSpinlog() (err error) {
	if cfg.XormSpinlog, err = xorm.NewEngine(Cfg.XormDriverName, util.JoinPath(cfg.SqlPath, slotspinfile)); err != nil {
		return
	}
	cfg.XormSpinlog.SetMapper(names.GonicMapper{})

	var session = cfg.XormSpinlog.NewSession()
	defer session.Close()

	if err = session.Sync(&spi.Spinlog{}); err != nil {
		return
	}
	return
}

func Init() (err error) {
	if err = InitStorage(); err != nil {
		return fmt.Errorf("can not init XORM records storage: %w", err)
	}
	if err = InitSpinlog(); err != nil {
		return fmt.Errorf("can not init XORM spins log storage: %w", err)
	}
	return
}
