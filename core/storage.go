package core

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/util"
	"go.yaml.in/yaml/v3"
	"golang.org/x/sync/errgroup"
	"xorm.io/xorm"
	"xorm.io/xorm/names"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var (
	Cfg = cfg.Cfg // shortcut
)

var (
	XormStorage *xorm.Engine
	XormSpinlog *xorm.Engine
	XormWG      errgroup.Group
)

var (
	SpinBuf util.SqlBuf[Spinlog]
	MultBuf util.SqlBuf[Multlog]
	BankBat = map[uint64]*SqlBank{}
	JoinBuf = SqlStory{}
)

func InitStorage() (err error) {
	switch Cfg.DriverName {
	case "sqlite3":
		var fpath string
		var ismem = strings.Contains(Cfg.ClubSourceName, ":memory:")
		if ismem {
			fpath = Cfg.ClubSourceName
		} else {
			fpath = util.JoinPath(cfg.SqlPath, Cfg.ClubSourceName)
		}
		if XormStorage, err = xorm.NewEngine(Cfg.DriverName, fpath); err != nil {
			return
		}
		if ismem {
			XormStorage.SetMaxOpenConns(1)
			XormStorage.SetMaxIdleConns(1)
			XormStorage.SetConnMaxLifetime(0)
			cfg.Infof("club db: memory")
		} else {
			cfg.Infof("club db: sqlite")
		}

	case "mysql", "postgres":
		if XormStorage, err = xorm.NewEngine(Cfg.DriverName, Cfg.ClubSourceName); err != nil {
			return
		}
		cfg.Infof("club db: %s", Cfg.DriverName)
	}
	XormStorage.SetMapper(names.GonicMapper{})

	var session = XormStorage.NewSession()
	defer session.Close()

	if err = session.Sync(
		&ClubData{}, &User{}, &Props{},
		&Story{}, Walletlog{}, Banklog{},
		&Session{},
	); err != nil {
		return
	}

	var ok bool
	if ok, err = session.IsTableEmpty(&ClubData{}); err != nil {
		return
	}
	if ok {
		var body []byte
		if body, err = os.ReadFile(util.JoinFilePath(cfg.CfgPath, "slot-clubinit.sql")); err != nil {
			cfg.Errorf("can not open SQL-file with initial settings: %s", err.Error())
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
		cfg.Errorf("can not open YAML-file with properties initialization for new user: %s", err.Error())
		err = nil // remove error
	} else if err = yaml.Unmarshal(body, &PropMaster); err != nil {
		cfg.Errorf("can not unmarshal 'slot-newuser.yaml': %s", err.Error())
		err = nil // remove error
	}

	const limit = 256

	var offset = 0
	for {
		var chunk []ClubData
		if err = session.Limit(limit, offset).Find(&chunk); err != nil {
			return
		}
		offset += limit
		for _, cd := range chunk {
			Clubs.Set(cd.CID, MakeClub(cd))
			var bat = &SqlBank{}
			bat.Init(XormStorage, cd.CID, Cfg.ClubUpdateBuffer, Cfg.ClubInsertBuffer)
			BankBat[cd.CID] = bat
		}
		if limit > len(chunk) {
			break
		}
	}
	cfg.Infof("loaded %d clubs", Clubs.Len())

	offset = 0
	for {
		var chunk []*User
		if err = session.Limit(limit, offset).Find(&chunk); err != nil {
			return
		}
		offset += limit
		for _, user := range chunk {
			user.Init()
			Users.Set(user.UID, user)
		}
		if limit > len(chunk) {
			break
		}
	}
	cfg.Infof("loaded %d users", Users.Len())

	offset = 0
	for {
		var chunk []*Props
		if err = session.Limit(limit, offset).Find(&chunk); err != nil {
			return
		}
		offset += limit
		for _, props := range chunk {
			if !Clubs.Has(props.CID) {
				return fmt.Errorf("found props without club linkage, UID=%d, CID=%d, value=%g", props.UID, props.CID, props.Wallet)
			}
			var user, ok = Users.Get(props.UID)
			if !ok {
				return fmt.Errorf("found props without user linkage, UID=%d, CID=%d, value=%g", props.UID, props.CID, props.Wallet)
			}
			user.InsertProps(props)
		}
		if limit > len(chunk) {
			break
		}
	}

	// Read last IDs from database to continue incrementing counters.
	var has bool
	var lastID int64
	if has, err = session.Table("story").Desc("gid").Cols("gid").Get(&lastID); err != nil {
		return
	}
	if has {
		StoryCounter.Store(uint64(lastID))
	}
	if has, err = session.Table("session").Desc("sid").Cols("sid").Get(&lastID); err != nil {
		return
	}
	if has {
		SessCounter.Store(uint64(lastID))
	}

	JoinBuf.Init(XormStorage, Cfg.ClubInsertBuffer)
	return
}

func InitSpinlog() (err error) {
	switch Cfg.DriverName {
	case "sqlite3":
		var fpath string
		var ismem = strings.Contains(Cfg.SpinSourceName, ":memory:")
		if ismem {
			fpath = Cfg.SpinSourceName
		} else {
			fpath = util.JoinPath(cfg.SqlPath, Cfg.SpinSourceName)
		}
		if XormSpinlog, err = xorm.NewEngine(Cfg.DriverName, fpath); err != nil {
			return
		}
		if ismem {
			XormSpinlog.SetMaxOpenConns(1)
			XormSpinlog.SetMaxIdleConns(1)
			XormSpinlog.SetConnMaxLifetime(0)
			cfg.Infof("spin db: memory")
		} else {
			cfg.Infof("spin db: sqlite")
		}

	case "mysql", "postgres":
		if XormSpinlog, err = xorm.NewEngine(Cfg.DriverName, Cfg.SpinSourceName); err != nil {
			return
		}
		cfg.Infof("spin db: %s", Cfg.DriverName)
	}
	XormSpinlog.SetMapper(names.GonicMapper{})

	var session = XormSpinlog.NewSession()
	defer session.Close()

	if err = session.Sync(&Spinlog{}, &Multlog{}); err != nil {
		return
	}

	// Read last IDs from database to continue incrementing counters.
	var has bool
	var lastID int64
	if has, err = session.Table("spinlog").Desc("slid").Cols("slid").Get(&lastID); err != nil {
		return
	}
	if has {
		SpinCounter.Store(uint64(lastID))
	}
	var i64 int64
	if has, err = session.Table("multlog").Desc("mlid").Cols("mlid").Get(&i64); err != nil {
		return
	}
	if has {
		MultCounter.Store(uint64(i64))
	}

	SpinBuf.Init(XormSpinlog, Cfg.SpinInsertBuffer)
	MultBuf.Init(XormSpinlog, Cfg.SpinInsertBuffer)
	return
}

func FlushAll(fd time.Duration) (err error) {
	var errs []error
	for cid, bat := range BankBat {
		if err := bat.Flush(fd); err != nil {
			errs = append(errs, fmt.Errorf("can not update bank for cid=%d: %w", cid, err))
		}
	}
	if err := JoinBuf.Flush(fd); err != nil {
		errs = append(errs, fmt.Errorf("can not write to story log: %w", err))
	}
	if Cfg.UseSpinLog {
		if err := SpinBuf.Flush(fd); err != nil {
			errs = append(errs, fmt.Errorf("can not write to spin log: %w", err))
		}
		if err := MultBuf.Flush(fd); err != nil {
			errs = append(errs, fmt.Errorf("can not write to mult log: %w", err))
		}
	}
	return errors.Join(errs...)
}

func FlushLoop(exitctx context.Context) (err error) {
	cfg.Debugf("flash loop started, tick %s", Cfg.SqlFlushTick.String())
	defer cfg.Debugf("flash loop done")

	var fd = Cfg.SqlFlushTick
	var flush = time.Tick(fd)
	var passers = time.Tick(time.Hour * 8)
	for {
		select {
		case <-flush:
			if err = FlushAll(fd); err != nil {
				cfg.Fatalf(err.Error())
			}
		case <-passers:
			XormStorage.Where("ctime<? AND status=0", time.Now().Add(-time.Hour*3*24).Format(time.DateTime)).Delete(&User{})
		case <-exitctx.Done():
			return
		}
	}
}

func InitXORM() (err error) {
	if err = InitStorage(); err != nil {
		err = fmt.Errorf("can not init xorm records storage: %w", err)
		return
	}
	if Cfg.SpinSourceName == "" {
		Cfg.UseSpinLog = false
	}
	if Cfg.UseSpinLog {
		if err = InitSpinlog(); err != nil {
			err = fmt.Errorf("can not init xorm spins log storage: %w", err)
			return
		}
	}
	cfg.Debugf("xorm initialized")
	return
}

func DoneXORM() error {
	var errs []error

	if err := XormWG.Wait(); err != nil {
		errs = append(errs, err)
	}

	if err := FlushAll(0); err != nil {
		errs = append(errs, err)
	}

	if err := XormStorage.Close(); err != nil {
		errs = append(errs, fmt.Errorf("cannot close storage: %w", err))
	}
	if Cfg.UseSpinLog {
		if err := XormSpinlog.Close(); err != nil {
			errs = append(errs, fmt.Errorf("cannot close spinlog: %w", err))
		}
	}

	var err error
	if err = errors.Join(errs...); err != nil {
		err = fmt.Errorf("errors occurred on xorm finalize: %w", err)
	} else {
		cfg.Debugf("xorm done")
	}
	return err
}
