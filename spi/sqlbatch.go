package spi

import (
	"sync"
	"time"

	"xorm.io/xorm"
)

const (
	sqlbank1 = `UPDATE club SET bank=bank+? WHERE cid=?`
	sqlbank2 = `UPDATE props SET wallet=wallet+? WHERE uid=? AND cid=?`
)

func SafeTransaction(engine *xorm.Engine, f func(*Session) error) (err error) {
	var session = engine.NewSession()
	defer session.Close()
	defer func() {
		if err != nil {
			session.Rollback()
		} else {
			err = session.Commit()
		}
	}()

	if err = session.Begin(); err != nil {
		return
	}
	if err = f(session); err != nil {
		return
	}
	return
}

type SqlBank struct {
	cid     uint64
	banksum float64
	usersum map[uint64]float64
	usercap int
	lft     time.Time // last flush time
	mux     sync.Mutex
}

func (sb *SqlBank) Init(cid uint64, capacity int) {
	sb.mux.Lock()
	defer sb.mux.Unlock()
	sb.cid = cid
	sb.banksum = 0
	sb.usersum = make(map[uint64]float64, capacity)
	sb.usercap = capacity
}

func (sb *SqlBank) Put(engine *xorm.Engine, uid uint64, sum float64) (err error) {
	sb.mux.Lock()
	defer sb.mux.Unlock()
	sb.banksum += sum
	sb.usersum[uid] -= sum
	if len(sb.usersum) >= sb.usercap {
		if err = SafeTransaction(engine, func(session *Session) (err error) {
			if _, err = session.Exec(sqlbank1, sb.banksum, sb.cid); err != nil {
				return
			}
			for uid, sum := range sb.usersum {
				if _, err = session.Exec(sqlbank2, sum, uid, sb.cid); err != nil {
					return
				}
			}
			return
		}); err != nil {
			return
		}
		sb.banksum = 0
		clear(sb.usersum)
		sb.lft = time.Now()
	}
	return
}

func (sb *SqlBank) Flush(engine *xorm.Engine, d time.Duration) (err error) {
	sb.mux.Lock()
	defer sb.mux.Unlock()
	if len(sb.usersum) == 0 {
		return
	}
	if d == 0 || time.Since(sb.lft) >= d {
		if err = SafeTransaction(engine, func(session *Session) (err error) {
			if _, err = session.Exec(sqlbank1, sb.banksum, sb.cid); err != nil {
				return
			}
			for uid, sum := range sb.usersum {
				if _, err = session.Exec(sqlbank2, sum, uid, sb.cid); err != nil {
					return
				}
			}
			return
		}); err != nil {
			return
		}
		sb.banksum = 0
		clear(sb.usersum)
		sb.lft = time.Now()
	}
	return
}
