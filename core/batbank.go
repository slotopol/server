package core

import (
	"slices"
	"strconv"
	"strings"

	"xorm.io/xorm"

	"github.com/slotopol/server/util"
)

const (
	sqlbank1 = `UPDATE club SET bank=bank+?, utime=CURRENT_TIMESTAMP WHERE cid=?`
	sqlbank2 = `UPDATE props SET access=?, utime=CURRENT_TIMESTAMP WHERE uid=? AND cid=?`
	sqlbank3 = `UPDATE props SET mrtp=?, utime=CURRENT_TIMESTAMP WHERE uid=? AND cid=?`
)

// Formsts SQL-query for batch wallet updates like
//   UPDATE props
//   SET
//       wallet = wallet + CASE uid
//           WHEN 101 THEN 50
//           WHEN 102 THEN -20
//           WHEN 103 THEN 150
//       END,
//       utime = CURRENT_TIMESTAMP
//   WHERE cid = 2
//     AND uid IN (101, 102, 103);
func formatWalletQuery(cid uint64, uids []uint64, updates map[uint64]float64) string {
	var caseBuf strings.Builder
	var inBuf strings.Builder
	caseBuf.Grow(len(uids) * 35)
	inBuf.Grow(len(uids) * 12)

	var numBuf = make([]byte, 0, 32)
	for i, uid := range uids {
		var sum = updates[uid]
		// Format: WHEN <uid> THEN <delta>
		caseBuf.WriteString(" WHEN ")
		numBuf = strconv.AppendUint(numBuf[:0], uid, 10)
		caseBuf.Write(numBuf)
		caseBuf.WriteString(" THEN ")
		numBuf = strconv.AppendFloat(numBuf[:0], sum, 'f', -1, 64) // 'f' is safer for SQL
		caseBuf.Write(numBuf)
		// Format list for: IN (<uid>,...)
		if i > 0 {
			inBuf.WriteByte(',')
		}
		numBuf = strconv.AppendUint(numBuf[:0], uid, 10)
		inBuf.Write(numBuf)
	}

	var query strings.Builder
	query.Grow(caseBuf.Len() + inBuf.Len() + 120)
	query.WriteString("UPDATE props SET wallet = wallet + CASE uid")
	query.WriteString(caseBuf.String())
	query.WriteString(" END, utime = CURRENT_TIMESTAMP WHERE cid = ")
	numBuf = strconv.AppendUint(numBuf[:0], cid, 10)
	query.Write(numBuf)
	query.WriteString(" AND uid IN (")
	query.WriteString(inBuf.String())
	query.WriteByte(')')

	return query.String()
}

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
	if len(b.uidsum) > 0 {
		var uids = make([]uint64, 0, len(b.uidsum))
		for uid, sum := range b.uidsum {
			if sum == 0 {
				continue
			}
			uids = append(uids, uid)
		}
		if len(uids) > 0 {
			slices.Sort(uids)
			var query = formatWalletQuery(b.cid, uids, b.uidsum)
			if _, err = session.Exec(query); err != nil {
				return
			}
		}
	}
	for uid, access := range b.uidal {
		if _, err = session.Exec(sqlbank2, access, uid, b.cid); err != nil {
			return
		}
	}
	for uid, mrtp := range b.uidrtp {
		if _, err = session.Exec(sqlbank3, mrtp, uid, b.cid); err != nil {
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
