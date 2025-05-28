package api

import (
	"sync"
	"time"

	"xorm.io/xorm"
)

const (
	sqlbank1 = `UPDATE club SET bank=bank+?, utime=CURRENT_TIMESTAMP WHERE cid=?`
	sqlbank2 = `UPDATE props SET wallet=wallet+?, utime=CURRENT_TIMESTAMP WHERE uid=? AND cid=?`
	sqlbank3 = `UPDATE props SET access=?, utime=CURRENT_TIMESTAMP WHERE uid=? AND cid=?`
	sqlbank4 = `UPDATE props SET mrtp=?, utime=CURRENT_TIMESTAMP WHERE uid=? AND cid=?`
)

func SafeTransaction(engine *xorm.Engine, f func(*Session) error) (err error) {
	var session = engine.NewSession()
	defer session.Close()

	if err = session.Begin(); err != nil {
		return
	}
	if err = f(session); err != nil {
		session.Rollback()
		return
	}
	err = session.Commit()
	return
}

type SqlBank struct {
	cid     uint64
	banksum float64
	usersum map[uint64]float64
	useral  map[uint64]AL
	userrtp map[uint64]float64
	log     []Walletlog
	usercap int
	logsize int
	lft     time.Time // last flush time
	mux     sync.Mutex
}

func (sb *SqlBank) Init(cid uint64, capacity, logsize int) {
	sb.mux.Lock()
	defer sb.mux.Unlock()
	sb.cid = cid
	sb.banksum = 0
	sb.usersum = make(map[uint64]float64, capacity)
	sb.useral = make(map[uint64]AL, capacity)
	sb.userrtp = make(map[uint64]float64, capacity)
	sb.log = make([]Walletlog, 0, logsize)
	sb.usercap = capacity
	sb.logsize = logsize
}

func (sb *SqlBank) clear() {
	sb.banksum = 0
	clear(sb.usersum)
	clear(sb.useral)
	clear(sb.userrtp)
	sb.log = sb.log[:0]
	sb.lft = time.Now()
}

func (sb *SqlBank) IsEmpty() bool {
	return sb.banksum == 0 &&
		len(sb.usersum) == 0 && len(sb.useral) == 0 && len(sb.userrtp) == 0 &&
		len(sb.log) == 0
}

func (sb *SqlBank) transaction(session *Session) (err error) {
	if sb.banksum != 0 {
		if _, err = session.Exec(sqlbank1, sb.banksum, sb.cid); err != nil {
			return
		}
	}
	for uid, sum := range sb.usersum {
		if _, err = session.Exec(sqlbank2, sum, uid, sb.cid); err != nil {
			return
		}
	}
	for uid, access := range sb.useral {
		if _, err = session.Exec(sqlbank3, access, uid, sb.cid); err != nil {
			return
		}
	}
	for uid, mrtp := range sb.userrtp {
		if _, err = session.Exec(sqlbank4, mrtp, uid, sb.cid); err != nil {
			return
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

func (sb *SqlBank) Add(engine *xorm.Engine, uid, aid uint64, wallet, sum float64) (err error) {
	sb.mux.Lock()
	defer sb.mux.Unlock()
	sb.usersum[uid] += sum
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

func (sb *SqlBank) Access(engine *xorm.Engine, uid uint64, access AL) (err error) {
	sb.mux.Lock()
	defer sb.mux.Unlock()
	sb.useral[uid] = access
	if len(sb.useral) >= sb.usercap {
		if err = SafeTransaction(engine, sb.transaction); err != nil {
			return
		}
		sb.clear()
	}
	return
}

func (sb *SqlBank) MRTP(engine *xorm.Engine, uid uint64, mrtp float64) (err error) {
	sb.mux.Lock()
	defer sb.mux.Unlock()
	sb.userrtp[uid] = mrtp
	if len(sb.userrtp) >= sb.usercap {
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
	if sb.IsEmpty() {
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

type SqlStory struct {
	log     []*Story
	logsize int
	lft     time.Time // last flush time
	mux     sync.Mutex
}

func (ss *SqlStory) Init(logsize int) {
	ss.mux.Lock()
	defer ss.mux.Unlock()
	ss.log = make([]*Story, 0, logsize)
	ss.logsize = logsize
}

func (ss *SqlStory) clear() {
	ss.log = ss.log[:0]
	ss.lft = time.Now()
}

func (ss *SqlStory) IsEmpty() bool {
	return len(ss.log) == 0
}

func (ss *SqlStory) transaction(session *Session) (err error) {
	if len(ss.log) > 0 {
		if _, err = session.InsertMulti(&ss.log); err != nil {
			return
		}
	}
	return
}

func (ss *SqlStory) Join(engine *xorm.Engine, s *Story) (err error) {
	ss.mux.Lock()
	defer ss.mux.Unlock()
	ss.log = append(ss.log, s)
	if len(ss.log) >= ss.logsize {
		if err = SafeTransaction(engine, ss.transaction); err != nil {
			return
		}
		ss.clear()
	}
	return
}

func (ss *SqlStory) Flush(engine *xorm.Engine, d time.Duration) (err error) {
	ss.mux.Lock()
	defer ss.mux.Unlock()
	if ss.IsEmpty() {
		return
	}
	if d == 0 || time.Since(ss.lft) >= d {
		if err = SafeTransaction(engine, ss.transaction); err != nil {
			return
		}
		ss.clear()
	}
	return
}
