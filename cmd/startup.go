package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/spi"
	"github.com/slotopol/server/util"
	"gopkg.in/yaml.v3"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

var (
	Cfg = cfg.Cfg
)

var (
	ErrNoClubName = errors.New("name of 'club' database does not provided at data source name")
	ErrNoSpinName = errors.New("name of 'spin' database does not provided at data source name")
)

const sqlnewdb = "CREATE DATABASE IF NOT EXISTS `%s`"

func Startup() (exitctx context.Context) {
	//var cancel context.CancelFunc
	exitctx, cancel := context.WithCancel(context.Background())
	go func() {
		// Make exit signal on function exit.
		defer cancel()

		var sigint = make(chan os.Signal, 1)
		var sigterm = make(chan os.Signal, 1)
		// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) or SIGTERM (Ctrl+/)
		// SIGKILL, SIGQUIT will not be caught.
		signal.Notify(sigint, syscall.SIGINT)
		signal.Notify(sigterm, syscall.SIGTERM)
		// Block until we receive our signal.
		select {
		case <-exitctx.Done():
			if errors.Is(exitctx.Err(), context.DeadlineExceeded) {
				log.Println("shutting down by timeout")
			} else if errors.Is(exitctx.Err(), context.Canceled) {
				log.Println("shutting down by cancel")
			} else {
				log.Printf("shutting down by %s\n", exitctx.Err().Error())
			}
		case <-sigint:
			log.Println("shutting down by break")
		case <-sigterm:
			log.Println("shutting down by process termination")
		}
		signal.Stop(sigint)
		signal.Stop(sigterm)
	}()
	return
}

func InitStorage() (err error) {
	if Cfg.DriverName != "sqlite3" {
		var c1 = strings.Split(Cfg.ClubSourceName, "@/")
		if len(c1) < 2 {
			return ErrNoClubName
		}
		var c2 = strings.Split(c1[1], "?")
		var engine *xorm.Engine
		if engine, err = xorm.NewEngine(Cfg.DriverName, c1[0]+"@/"); err != nil {
			return
		}
		defer engine.Close()
		if _, err = engine.Exec(fmt.Sprintf(sqlnewdb, c2[0])); err != nil {
			return
		}
	}

	if Cfg.DriverName == "sqlite3" {
		cfg.XormStorage, err = xorm.NewEngine(Cfg.DriverName, util.JoinPath(cfg.SqlPath, Cfg.ClubSourceName))
	} else {
		cfg.XormStorage, err = xorm.NewEngine(Cfg.DriverName, Cfg.ClubSourceName)
	}
	if err != nil {
		return
	}
	cfg.XormStorage.SetMapper(names.GonicMapper{})

	var session = cfg.XormStorage.NewSession()
	defer session.Close()

	if err = session.Sync(
		&spi.Club{}, &spi.User{}, &spi.Props{},
		&spi.Story{}, spi.Walletlog{}, spi.Banklog{},
	); err != nil {
		return
	}

	var ok bool
	if ok, err = session.IsTableEmpty(&spi.Club{}); err != nil {
		return
	}
	if ok {
		var body []byte
		if body, err = os.ReadFile(util.JoinFilePath(cfg.CfgPath, "slot-club-init.sql")); err != nil {
			log.Printf("can not open SQL-file with initial settings: %s", err.Error())
			err = nil // remove error
		}
		var list = bytes.Split(body, []byte{';'})
		for _, cmd := range list {
			if cmd = bytes.TrimSpace(cmd); len(cmd) > 0 {
				if _, err = session.Exec(util.B2S(cmd)); err != nil {
					return
				}
			}
		}
	}

	// Read properies master for new registered user
	var body []byte
	if body, err = os.ReadFile(util.JoinFilePath(cfg.CfgPath, "slot-new-user.yaml")); err != nil {
		log.Printf("can not open YAML-file with properties initialization for new user: %s", err.Error())
		err = nil // remove error
	} else if err = yaml.Unmarshal(body, &spi.PropMaster); err != nil {
		log.Printf("can not unmarshal 'slot-new-user.yaml': %s", err.Error())
		err = nil // remove error
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
			var bat = &spi.SqlBank{}
			bat.Init(club.CID, Cfg.ClubUpdateBuffer, Cfg.ClubInsertBuffer)
			spi.BankBat[club.CID] = bat
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

	var i64 int64
	if i64, err = session.Count(&spi.Story{}); err != nil {
		return
	}
	spi.StoryCounter = uint64(i64)

	spi.JoinBuf.Init(Cfg.ClubUpdateBuffer, Cfg.ClubInsertBuffer)
	return
}

func InitSpinlog() (err error) {
	if Cfg.DriverName != "sqlite3" {
		var c1 = strings.Split(Cfg.SpinSourceName, "@/")
		if len(c1) < 2 {
			return ErrNoSpinName
		}
		var c2 = strings.Split(c1[1], "?")
		var engine *xorm.Engine
		if engine, err = xorm.NewEngine(Cfg.DriverName, c1[0]+"@/"); err != nil {
			return
		}
		defer engine.Close()
		if _, err = engine.Exec(fmt.Sprintf(sqlnewdb, c2[0])); err != nil {
			return
		}
	}

	if Cfg.DriverName == "sqlite3" {
		cfg.XormSpinlog, err = xorm.NewEngine(Cfg.DriverName, util.JoinPath(cfg.SqlPath, Cfg.SpinSourceName))
	} else {
		cfg.XormSpinlog, err = xorm.NewEngine(Cfg.DriverName, Cfg.SpinSourceName)
	}
	if err != nil {
		return
	}
	cfg.XormSpinlog.SetMapper(names.GonicMapper{})

	var session = cfg.XormSpinlog.NewSession()
	defer session.Close()

	if err = session.Sync(&spi.Spinlog{}, &spi.Multlog{}); err != nil {
		return
	}
	var i64 int64
	if i64, err = session.Count(&spi.Spinlog{}); err != nil {
		return
	}
	spi.SpinCounter = uint64(i64)
	if i64, err = session.Count(&spi.Multlog{}); err != nil {
		return
	}
	spi.MultCounter = uint64(i64)

	spi.SpinBuf.Init(Cfg.SpinInsertBuffer)
	spi.MultBuf.Init(Cfg.SpinInsertBuffer)
	return
}

func SqlLoop(exitctx context.Context) {
	var fd = Cfg.SqlFlushTick
	var flush = time.Tick(fd)
	var passers = time.Tick(time.Hour * 8)
	for {
		select {
		case <-flush:
			for cid, bat := range spi.BankBat {
				if err := bat.Flush(cfg.XormStorage, fd); err != nil {
					log.Printf("can not update bank for cid=%d: %s", cid, err.Error())
				}
			}
			if err := spi.JoinBuf.Flush(cfg.XormStorage, fd); err != nil {
				log.Printf("can not write to story log: %s", err.Error())
			}
			if err := spi.SpinBuf.Flush(cfg.XormSpinlog, fd); err != nil {
				log.Printf("can not write to spin log: %s", err.Error())
			}
			if err := spi.MultBuf.Flush(cfg.XormSpinlog, fd); err != nil {
				log.Printf("can not write to mult log: %s", err.Error())
			}
		case <-passers:
			cfg.XormStorage.Where("ctime<? AND status=0", time.Now().Add(-time.Hour*3*24).Format(time.DateTime)).Delete(&spi.User{})
		case <-exitctx.Done():
			return
		}
	}
}

func Init() (err error) {
	if err = InitStorage(); err != nil {
		err = fmt.Errorf("can not init XORM records storage: %w", err)
		return
	}
	if err = InitSpinlog(); err != nil {
		err = fmt.Errorf("can not init XORM spins log storage: %w", err)
		return
	}
	return
}

func Done() (err error) {
	var errs []error
	for _, bat := range spi.BankBat {
		errs = append(errs, bat.Flush(cfg.XormStorage, 0))
	}
	errs = append(errs, spi.JoinBuf.Flush(cfg.XormStorage, 0))

	errs = append(errs, spi.SpinBuf.Flush(cfg.XormSpinlog, 0))
	errs = append(errs, spi.MultBuf.Flush(cfg.XormSpinlog, 0))

	errs = append(errs, cfg.XormStorage.Close())
	errs = append(errs, cfg.XormSpinlog.Close())
	return errors.Join(errs...)
}
