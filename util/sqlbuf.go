package util

import (
	"sync"
	"time"

	"xorm.io/xorm"
)

type SqlBuf[T any] struct {
	engine *xorm.Engine
	mux    sync.Mutex

	buf []T
	lft time.Time // last flush time
}

func (sb *SqlBuf[T]) clear() {
	sb.buf = sb.buf[:0]
	sb.lft = time.Now()
}

func (sb *SqlBuf[T]) isEmpty() bool {
	return len(sb.buf) == 0
}

func (sb *SqlBuf[T]) isFull() bool {
	return len(sb.buf) >= cap(sb.buf)
}

func (sb *SqlBuf[T]) Init(engine *xorm.Engine, capacity int) {
	sb.mux.Lock()
	defer sb.mux.Unlock()
	sb.engine = engine
	sb.buf = make([]T, 0, capacity)
}

func (sb *SqlBuf[T]) Len() int {
	sb.mux.Lock()
	defer sb.mux.Unlock()
	return len(sb.buf)
}

func (sb *SqlBuf[T]) Last() time.Time {
	sb.mux.Lock()
	defer sb.mux.Unlock()
	return sb.lft
}

func (sb *SqlBuf[T]) Put(val T) (err error) {
	sb.mux.Lock()
	defer sb.mux.Unlock()

	sb.buf = append(sb.buf, val)

	if sb.isFull() {
		_, err = sb.engine.Insert(&sb.buf)
		sb.clear()
	}
	return
}

func (sb *SqlBuf[T]) Flush(d time.Duration) (err error) {
	sb.mux.Lock()
	defer sb.mux.Unlock()

	if sb.isEmpty() {
		return
	}
	if d == 0 || time.Since(sb.lft) >= d {
		_, err = sb.engine.Insert(&sb.buf)
		sb.clear()
	}
	return
}
