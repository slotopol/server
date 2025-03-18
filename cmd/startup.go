package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/slotopol/server/api"
	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/util"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v3"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

var (
	Cfg = cfg.Cfg // shortcut
)

var (
	ErrNoClubName = errors.New("name of 'club' database does not provided at data source name")
	ErrNoSpinName = errors.New("name of 'spin' database does not provided at data source name")
)

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
	switch Cfg.DriverName {
	case "sqlite3":
		var fpath string
		if Cfg.ClubSourceName != ":memory:" {
			fpath = util.JoinPath(cfg.SqlPath, Cfg.ClubSourceName)
		} else {
			fpath = Cfg.ClubSourceName
		}
		if cfg.XormStorage, err = xorm.NewEngine(Cfg.DriverName, fpath); err != nil {
			return
		}
		if Cfg.ClubSourceName != ":memory:" {
			log.Println("club db: sqlite")
		} else {
			log.Println("club db: memory")
		}

	case "mysql", "postgres":
		if cfg.XormStorage, err = xorm.NewEngine(Cfg.DriverName, Cfg.ClubSourceName); err != nil {
			return
		}
		log.Printf("club db: %s\n", Cfg.DriverName)
	}
	cfg.XormStorage.SetMapper(names.GonicMapper{})

	var session = cfg.XormStorage.NewSession()
	defer session.Close()

	if err = session.Sync(
		&api.ClubData{}, &api.User{}, &api.Props{},
		&api.Story{}, api.Walletlog{}, api.Banklog{},
	); err != nil {
		return
	}

	var ok bool
	if ok, err = session.IsTableEmpty(&api.ClubData{}); err != nil {
		return
	}
	if ok {
		var body []byte
		if body, err = os.ReadFile(util.JoinFilePath(cfg.CfgPath, "slot-clubinit.sql")); err != nil {
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
	if body, err = os.ReadFile(util.JoinFilePath(cfg.CfgPath, "slot-newuser.yaml")); err != nil {
		log.Printf("can not open YAML-file with properties initialization for new user: %s", err.Error())
		err = nil // remove error
	} else if err = yaml.Unmarshal(body, &api.PropMaster); err != nil {
		log.Printf("can not unmarshal 'slot-newuser.yaml': %s", err.Error())
		err = nil // remove error
	}

	const limit = 256

	var offset = 0
	for {
		var chunk []api.ClubData
		if err = session.Limit(limit, offset).Find(&chunk); err != nil {
			return
		}
		offset += limit
		for _, cd := range chunk {
			api.Clubs.Set(cd.CID, api.MakeClub(cd))
			var bat = &api.SqlBank{}
			bat.Init(cd.CID, Cfg.ClubUpdateBuffer, Cfg.ClubInsertBuffer)
			api.BankBat[cd.CID] = bat
		}
		if limit > len(chunk) {
			break
		}
	}
	log.Printf("loaded %d clubs\n", api.Clubs.Len())

	offset = 0
	for {
		var chunk []*api.User
		if err = session.Limit(limit, offset).Find(&chunk); err != nil {
			return
		}
		offset += limit
		for _, user := range chunk {
			user.Init()
			api.Users.Set(user.UID, user)
		}
		if limit > len(chunk) {
			break
		}
	}
	log.Printf("loaded %d users\n", api.Users.Len())

	offset = 0
	for {
		var chunk []*api.Props
		if err = session.Limit(limit, offset).Find(&chunk); err != nil {
			return
		}
		offset += limit
		for _, props := range chunk {
			if !api.Clubs.Has(props.CID) {
				return fmt.Errorf("found props without club linkage, UID=%d, CID=%d, value=%g", props.UID, props.CID, props.Wallet)
			}
			var user, ok = api.Users.Get(props.UID)
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
	if i64, err = session.Count(&api.Story{}); err != nil {
		return
	}
	api.StoryCounter = uint64(i64)

	api.JoinBuf.Init(Cfg.ClubInsertBuffer)
	return
}

func InitSpinlog() (err error) {
	switch Cfg.DriverName {
	case "sqlite3":
		var fpath string
		if Cfg.SpinSourceName != ":memory:" {
			fpath = util.JoinPath(cfg.SqlPath, Cfg.SpinSourceName)
		} else {
			fpath = Cfg.SpinSourceName
		}
		if cfg.XormSpinlog, err = xorm.NewEngine(Cfg.DriverName, fpath); err != nil {
			return
		}
		if Cfg.ClubSourceName != ":memory:" {
			log.Println("spin db: sqlite")
		} else {
			log.Println("spin db: memory")
		}

	case "mysql", "postgres":
		if cfg.XormSpinlog, err = xorm.NewEngine(Cfg.DriverName, Cfg.SpinSourceName); err != nil {
			return
		}
		log.Printf("spin db: %s\n", Cfg.DriverName)
	}
	cfg.XormSpinlog.SetMapper(names.GonicMapper{})

	var session = cfg.XormSpinlog.NewSession()
	defer session.Close()

	if err = session.Sync(&api.Spinlog{}, &api.Multlog{}); err != nil {
		return
	}
	var i64 int64
	if i64, err = session.Count(&api.Spinlog{}); err != nil {
		return
	}
	api.SpinCounter = uint64(i64)
	if i64, err = session.Count(&api.Multlog{}); err != nil {
		return
	}
	api.MultCounter = uint64(i64)

	api.SpinBuf.Init(Cfg.SpinInsertBuffer)
	api.MultBuf.Init(Cfg.SpinInsertBuffer)
	return
}

func SqlLoop(exitctx context.Context) {
	var fd = Cfg.SqlFlushTick
	var flush = time.Tick(fd)
	var passers = time.Tick(time.Hour * 8)
	for {
		select {
		case <-flush:
			for cid, bat := range api.BankBat {
				if err := bat.Flush(cfg.XormStorage, fd); err != nil {
					log.Printf("can not update bank for cid=%d: %s", cid, err.Error())
				}
			}
			if err := api.JoinBuf.Flush(cfg.XormStorage, fd); err != nil {
				log.Printf("can not write to story log: %s", err.Error())
			}
			if Cfg.UseSpinLog {
				if err := api.SpinBuf.Flush(cfg.XormSpinlog, fd); err != nil {
					log.Printf("can not write to spin log: %s", err.Error())
				}
				if err := api.MultBuf.Flush(cfg.XormSpinlog, fd); err != nil {
					log.Printf("can not write to mult log: %s", err.Error())
				}
			}
		case <-passers:
			cfg.XormStorage.Where("ctime<? AND status=0", time.Now().Add(-time.Hour*3*24).Format(time.DateTime)).Delete(&api.User{})
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
	if Cfg.SpinSourceName == "" {
		Cfg.UseSpinLog = false
	}
	if Cfg.UseSpinLog {
		if err = InitSpinlog(); err != nil {
			err = fmt.Errorf("can not init XORM spins log storage: %w", err)
			return
		}
	}
	return
}

func Done() (err error) {
	var errs []error
	for _, bat := range api.BankBat {
		errs = append(errs, bat.Flush(cfg.XormStorage, 0))
	}
	errs = append(errs, api.JoinBuf.Flush(cfg.XormStorage, 0))
	errs = append(errs, cfg.XormStorage.Close())

	if Cfg.UseSpinLog {
		errs = append(errs, api.SpinBuf.Flush(cfg.XormSpinlog, 0))
		errs = append(errs, api.MultBuf.Flush(cfg.XormSpinlog, 0))
		errs = append(errs, cfg.XormSpinlog.Close())
	}
	return errors.Join(errs...)
}
