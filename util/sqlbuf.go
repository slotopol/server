package util

import (
	"sync"
	"time"

	"xorm.io/xorm"
)

type SqlBuf[T any] struct {
	buf []T
	lft time.Time // last flush time
	mux sync.Mutex
}

func (sb *SqlBuf[T]) Init(capacity int) {
	sb.mux.Lock()
	sb.buf = make([]T, 0, capacity)
	sb.mux.Unlock()
}

func (sb *SqlBuf[T]) Len() int {
	return len(sb.buf)
}

func (sb *SqlBuf[T]) Last() time.Time {
	return sb.lft
}

func (sb *SqlBuf[T]) Flush(engine *xorm.Engine, d time.Duration) (err error) {
	sb.mux.Lock()
	defer sb.mux.Unlock()
	if len(sb.buf) == 0 {
		return
	}
	if d == 0 || time.Since(sb.lft) >= d {
		_, err = engine.Insert(&sb.buf)
		sb.buf = sb.buf[:0]
		sb.lft = time.Now()
	}
	return
}

func (sb *SqlBuf[T]) Push(val T, engine *xorm.Engine) (err error) {
	sb.mux.Lock()
	defer sb.mux.Unlock()
	sb.buf = append(sb.buf, val)
	if len(sb.buf) == cap(sb.buf) {
		_, err = engine.Insert(&sb.buf)
		sb.buf = sb.buf[:0]
		sb.lft = time.Now()
	}
	return
}
