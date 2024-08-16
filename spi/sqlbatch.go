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
	userins map[uint64]bool
	usercap int
	logsize int
	log     []Walletlog
	lft     time.Time // last flush time
	mux     sync.Mutex
}

func (sb *SqlBank) Init(cid uint64, capacity, logsize int) {
	sb.mux.Lock()
	defer sb.mux.Unlock()
	sb.cid = cid
	sb.banksum = 0
	sb.usersum = make(map[uint64]float64, capacity)
	sb.userins = make(map[uint64]bool, capacity)
	sb.usercap = capacity
	sb.logsize = logsize
	sb.log = make([]Walletlog, 0, logsize)
}

func (sb *SqlBank) clear() {
	sb.banksum = 0
	clear(sb.usersum)
	clear(sb.userins)
	sb.log = sb.log[:0]
	sb.lft = time.Now()
}

func (sb *SqlBank) transaction(session *Session) (err error) {
	if sb.banksum != 0 {
		if _, err = session.Exec(sqlbank1, sb.banksum, sb.cid); err != nil {
			return
		}
	}
	for uid, sum := range sb.usersum {
		if sb.userins[uid] {
			var props = &Props{
				CID:    sb.cid,
				UID:    uid,
				Wallet: sum,
			}
			if _, err = session.InsertOne(props); err != nil {
				return
			}
		} else {
			if _, err = session.Exec(sqlbank2, sum, uid, sb.cid); err != nil {
				return
			}
		}
	}
	if len(sb.log) > 0 {
		if _, err = session.InsertMulti(&sb.log); err != nil {
			return
		}
	}
	return
}

func (sb *SqlBank) Put(engine *xorm.Engine, uid uint64, sum float64) (err error) {
	sb.mux.Lock()
	defer sb.mux.Unlock()
	sb.banksum += sum
	sb.usersum[uid] -= sum
	if len(sb.usersum) >= sb.usercap {
		if err = SafeTransaction(engine, sb.transaction); err != nil {
			return
		}
		sb.clear()
	}
	return
}

func (sb *SqlBank) Add(engine *xorm.Engine, uid, aid uint64, wallet, sum float64, ins bool) (err error) {
	sb.mux.Lock()
	defer sb.mux.Unlock()
	sb.usersum[uid] += sum
	if ins {
		sb.userins[uid] = ins
	}
	sb.log = append(sb.log, Walletlog{
		CID:    sb.cid,
		UID:    uid,
		AID:    aid,
		Wallet: wallet,
		Sum:    sum,
	})
	if len(sb.usersum) >= sb.usercap || len(sb.log) >= sb.logsize {
		if err = SafeTransaction(engine, sb.transaction); err != nil {
			return
		}
		sb.clear()
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
		if err = SafeTransaction(engine, sb.transaction); err != nil {
			return
		}
		sb.clear()
	}
	return
}
