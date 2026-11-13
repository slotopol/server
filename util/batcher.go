package util

import (
	"sync"
	"time"

	"xorm.io/xorm"
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

type Batcher interface {
	IsEmpty() bool
	IsFull() bool
	Clear()
	Flush(*xorm.Session) error
}

type BatPack[T Batcher] struct {
	mux sync.Mutex
	lft time.Time // last flush time
	eng *xorm.Engine
	bat T
}

func (bp *BatPack[T]) Init(eng *xorm.Engine, bat T) *BatPack[T] {
	bp.mux.Lock()
	defer bp.mux.Unlock()
	bp.bat = bat
	bp.eng = eng
	bp.lft = time.Now()
	return bp
}

func (bp *BatPack[T]) flush() (err error) {
	if err = SafeTransaction(bp.eng, bp.bat.Flush); err != nil {
		return
	}
	bp.bat.Clear()
	bp.lft = time.Now()
	return
}

func (bp *BatPack[T]) TryFlush(d time.Duration) (err error) {
	bp.mux.Lock()
	defer bp.mux.Unlock()
	if (d <= 0 || time.Since(bp.lft) >= d) && !bp.bat.IsEmpty() {
		err = bp.flush()
	}
	return
}

func (bp *BatPack[T]) Upsert(uf func(T)) (flush bool, err error) {
	bp.mux.Lock()
	defer bp.mux.Unlock()
	uf(bp.bat)
	if bp.bat.IsFull() {
		flush, err = true, bp.flush()
	}
	return
}

type BatLog[T any] struct {
	limit int
	log   []T
}

var _ Batcher = (*BatLog[struct{}])(nil)

func (b *BatLog[T]) Init(limit int) *BatLog[T] {
	b.limit = limit
	b.log = make([]T, 0, limit)
	return b
}

func (b *BatLog[T]) IsEmpty() bool {
	return len(b.log) == 0
}

func (b *BatLog[T]) IsFull() bool {
	return len(b.log) >= b.limit
}

func (b *BatLog[T]) Clear() {
	b.log = b.log[:0]
}

func (b *BatLog[T]) Flush(session *xorm.Session) (err error) {
	_, err = session.InsertMulti(&b.log)
	return
}

func (b *BatLog[T]) Put(v T) {
	b.log = append(b.log, v)
}

type PackLog[T any] struct {
	BatPack[*BatLog[T]]
}

func (p *PackLog[T]) Init(eng *xorm.Engine, limit int) {
	p.BatPack.Init(eng, (&BatLog[T]{}).Init(limit))
}

func (p *PackLog[T]) Put(v T) (flush bool, err error) {
	return p.Upsert(func(b *BatLog[T]) {
		b.Put(v)
	})
}
