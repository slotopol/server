package core

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

func SafeTransaction(engine *xorm.Engine, f func(*xorm.Session) error) (err error) {
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
	engine *xorm.Engine
	mux    sync.Mutex

	cid     uint64
	banksum float64
	usersum map[uint64]float64
	useral  map[uint64]AL
	userrtp map[uint64]float64
	log     []Walletlog
	usercap int
	logsize int
	lft     time.Time // last flush time
}

func (sb *SqlBank) clear() {
	sb.banksum = 0
	clear(sb.usersum)
	clear(sb.useral)
	clear(sb.userrtp)
	sb.log = sb.log[:0]
	sb.lft = time.Now()
}

func (sb *SqlBank) isEmpty() bool {
	return sb.banksum == 0 &&
		len(sb.usersum) == 0 &&
		len(sb.useral) == 0 &&
		len(sb.userrtp) == 0 &&
		len(sb.log) == 0
}

func (sb *SqlBank) isFull() bool {
	return len(sb.usersum) >= sb.usercap ||
		len(sb.useral) >= sb.usercap ||
		len(sb.userrtp) >= sb.usercap ||
		len(sb.log) >= sb.logsize
}

func (sb *SqlBank) safeTransaction() (err error) {
	if sb.isEmpty() {
		return
	}
	if err = SafeTransaction(sb.engine, sb.transaction); err != nil {
		return
	}
	sb.clear()
	return
}

func (sb *SqlBank) transaction(session *xorm.Session) (err error) {
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

func (sb *SqlBank) Init(engine *xorm.Engine, cid uint64, capacity, logsize int) {
	sb.mux.Lock()
	defer sb.mux.Unlock()

	sb.engine = engine
	sb.cid = cid
	sb.banksum = 0
	sb.usersum = make(map[uint64]float64, capacity)
	sb.useral = make(map[uint64]AL, capacity)
	sb.userrtp = make(map[uint64]float64, capacity)
	sb.log = make([]Walletlog, 0, logsize)
	sb.usercap = capacity
	sb.logsize = logsize
}

func (sb *SqlBank) Last() time.Time {
	sb.mux.Lock()
	defer sb.mux.Unlock()
	return sb.lft
}

func (sb *SqlBank) Put(uid uint64, sum float64) (err error) {
	sb.mux.Lock()
	defer sb.mux.Unlock()

	sb.banksum += sum
	sb.usersum[uid] -= sum

	if sb.isFull() {
		return sb.safeTransaction()
	}
	return
}

func (sb *SqlBank) Add(uid, aid uint64, wallet, sum float64) (err error) {
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

	if sb.isFull() {
		return sb.safeTransaction()
	}
	return
}

func (sb *SqlBank) Access(uid uint64, access AL) (err error) {
	sb.mux.Lock()
	defer sb.mux.Unlock()

	sb.useral[uid] = access

	if sb.isFull() {
		return sb.safeTransaction()
	}
	return
}

func (sb *SqlBank) MRTP(uid uint64, mrtp float64) (err error) {
	sb.mux.Lock()
	defer sb.mux.Unlock()

	sb.userrtp[uid] = mrtp

	if sb.isFull() {
		return sb.safeTransaction()
	}
	return
}

func (sb *SqlBank) Flush(d time.Duration) (err error) {
	sb.mux.Lock()
	defer sb.mux.Unlock()

	if d == 0 || time.Since(sb.lft) >= d {
		return sb.safeTransaction()
	}
	return
}

type SqlStory struct {
	engine *xorm.Engine
	mux    sync.Mutex

	log     []*Story
	logsize int
	lft     time.Time // last flush time
}

func (ss *SqlStory) clear() {
	ss.log = ss.log[:0]
	ss.lft = time.Now()
}

func (ss *SqlStory) isEmpty() bool {
	return len(ss.log) == 0
}

func (ss *SqlStory) isFull() bool {
	return len(ss.log) >= ss.logsize
}

func (ss *SqlStory) safeTransaction() (err error) {
	if ss.isEmpty() {
		return
	}
	if err = SafeTransaction(ss.engine, ss.transaction); err != nil {
		return
	}
	ss.clear()
	return
}

func (ss *SqlStory) transaction(session *xorm.Session) (err error) {
	if len(ss.log) > 0 {
		if _, err = session.InsertMulti(&ss.log); err != nil {
			return
		}
	}
	return
}

func (ss *SqlStory) Init(engine *xorm.Engine, logsize int) {
	ss.mux.Lock()
	defer ss.mux.Unlock()
	ss.engine = engine
	ss.log = make([]*Story, 0, logsize)
	ss.logsize = logsize
}

func (ss *SqlStory) Last() time.Time {
	ss.mux.Lock()
	defer ss.mux.Unlock()
	return ss.lft
}

func (ss *SqlStory) Join(s *Story) (err error) {
	ss.mux.Lock()
	defer ss.mux.Unlock()

	ss.log = append(ss.log, s)

	if ss.isFull() {
		return ss.safeTransaction()
	}
	return
}

func (ss *SqlStory) Flush(d time.Duration) (err error) {
	ss.mux.Lock()
	defer ss.mux.Unlock()

	if d == 0 || time.Since(ss.lft) >= d {
		return ss.safeTransaction()
	}
	return
}
