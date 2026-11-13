package core

import (
	"xorm.io/xorm"

	"github.com/slotopol/server/util"
)

const (
	sqlbank1 = `UPDATE club SET bank=bank+?, utime=CURRENT_TIMESTAMP WHERE cid=?`
	sqlbank2 = `UPDATE props SET wallet=wallet+?, utime=CURRENT_TIMESTAMP WHERE uid=? AND cid=?`
	sqlbank3 = `UPDATE props SET access=?, utime=CURRENT_TIMESTAMP WHERE uid=? AND cid=?`
	sqlbank4 = `UPDATE props SET mrtp=?, utime=CURRENT_TIMESTAMP WHERE uid=? AND cid=?`
)

type BatBank struct {
	cid    uint64
	updlim int
	inslim int

	yield  float64
	uidsum map[uint64]float64
	uidal  map[uint64]AL
	uidrtp map[uint64]float64
	log    []Walletlog
}

var _ util.Batcher = (*BatBank)(nil)

func (b *BatBank) Init(cid uint64, updlim int, inslim int) *BatBank {
	b.cid = cid
	b.updlim = updlim
	b.inslim = inslim

	b.yield = 0
	b.uidsum = make(map[uint64]float64, updlim)
	b.uidal = make(map[uint64]AL, updlim)
	b.uidrtp = make(map[uint64]float64, updlim)
	b.log = make([]Walletlog, 0, inslim)
	return b
}

func (b *BatBank) IsEmpty() bool {
	return b.yield == 0 &&
		len(b.uidsum) == 0 &&
		len(b.uidal) == 0 &&
		len(b.uidrtp) == 0 &&
		len(b.log) == 0
}

func (b *BatBank) IsFull() bool {
	return len(b.log) >= b.inslim ||
		len(b.uidsum)+len(b.uidal)+len(b.uidrtp) >= b.updlim
}

func (b *BatBank) Clear() {
	b.yield = 0
	clear(b.uidsum)
	clear(b.uidal)
	clear(b.uidrtp)
	b.log = b.log[:0]
}

func (b *BatBank) Flush(session *xorm.Session) (err error) {
	if b.yield != 0 {
		if _, err = session.Exec(sqlbank1, b.yield, b.cid); err != nil {
			return
		}
	}
	for uid, sum := range b.uidsum {
		if _, err = session.Exec(sqlbank2, sum, uid, b.cid); err != nil {
			return
		}
	}
	for uid, access := range b.uidal {
		if _, err = session.Exec(sqlbank3, access, uid, b.cid); err != nil {
			return
		}
	}
	for uid, mrtp := range b.uidrtp {
		if _, err = session.Exec(sqlbank4, mrtp, uid, b.cid); err != nil {
			return
		}
	}
	if len(b.log) > 0 {
		if _, err = session.InsertMulti(&b.log); err != nil {
			return
		}
	}
	return
}

type PackBank struct {
	util.BatPack[*BatBank]
}

func (p *PackBank) Init(eng *xorm.Engine, cid uint64, updlim int, inslim int) *PackBank {
	p.BatPack.Init(eng, (&BatBank{}).Init(cid, updlim, inslim))
	return p
}

func (p *PackBank) AddYield(uid uint64, sum float64) (flush bool, err error) {
	return p.Upsert(func(b *BatBank) {
		b.yield += sum
		b.uidsum[uid] -= sum
	})
}

func (p *PackBank) AddCredit(uid, aid uint64, wallet, sum float64) (flush bool, err error) {
	return p.Upsert(func(b *BatBank) {
		b.uidsum[uid] += sum
		b.log = append(b.log, Walletlog{
			CID:    b.cid,
			UID:    uid,
			AID:    aid,
			Wallet: wallet,
			Sum:    sum,
		})
	})
}

func (p *PackBank) SetAL(uid uint64, access AL) (flush bool, err error) {
	return p.Upsert(func(b *BatBank) {
		b.uidal[uid] = access
	})
}

func (p *PackBank) SetRTP(uid uint64, mrtp float64) (flush bool, err error) {
	return p.Upsert(func(b *BatBank) {
		b.uidrtp[uid] = mrtp
	})
}
