package windows

import (
	"sync"
	"time"
)

type RollingWindowOption func(rollingWindow *RollingWindow)

type RollingWindow struct {
	l             sync.RWMutex
	startTime     time.Time
	size          int
	timec         TimecIface
	pos           int
	lastTime      time.Time
	interval      time.Duration
	statics       WindowStatic
	datas         []*Item
	ignoreCurrent bool // 忽略当前组的数据
}

func NewRollingWindow(size int, interval time.Duration, opts ...RollingWindowOption) *RollingWindow {
	timec := &DefaultTimec{}
	w := &RollingWindow{
		size:      size,
		lastTime:  timec.Now(),
		startTime: timec.Now(),
		timec:     timec,
		interval:  interval,
		datas:     make([]*Item, size),
	}
	for i := 0; i < size; i++ {
		w.datas[i] = new(Item)
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

func (rw *RollingWindow) Add(v float64) {
	rw.l.Lock()
	defer rw.l.Unlock()
	rw.updatePos()
	rw.add(v)
}

func (rw *RollingWindow) GetLastVisit() time.Time {
	return rw.lastTime
}

func (rw *RollingWindow) GetStaticItemNum() int {
	if rw.ignoreCurrent {
		return rw.size -1
	}
	return rw.size
}

func (rw *RollingWindow) offset() int {
	offset := int(rw.timec.Now().Sub(rw.lastTime) / rw.interval)
	if 0 <= offset && offset < rw.size {
		return offset
	}
	return rw.size
}

func (rw *RollingWindow) Cal(fn func(b *Item)) {
	rw.l.RLock()
	defer rw.l.RUnlock()

	var count int
	span := rw.offset()
	// ignore current bucket, because of partial data
	if span == 0 && rw.ignoreCurrent {
		count = rw.size - 1
	} else {
		count = rw.size - span
	}
	if count > 0 {
		start := (rw.pos + span + 1) % rw.size
		for i := 0; i < count; i++ {
			fn(rw.datas[(start+i)%rw.size])
		}
	}
}

func (rw *RollingWindow) Static() (float64, int) {
	if rw.statics != nil {
		// 更新加全局锁
		rw.l.Lock()
		rw.updatePos()
		rw.l.Unlock()
		// 取数据加读写锁
		rw.l.RLock()
		defer rw.l.RUnlock()
		return rw.statics.Static(rw.datas[rw.pos])
	}
	return 0, 0
}

func (rw *RollingWindow) StaticValue() float64 {
	if rw.statics != nil {
		// 更新加全局锁
		rw.l.Lock()
		rw.updatePos()
		rw.l.Unlock()
		// 取数据加读写锁
		rw.l.RLock()
		defer rw.l.RUnlock()
		return rw.statics.Value(rw.datas[rw.pos])
	}
	return 0
}

func (rw *RollingWindow) StaticTotal() int {
	if rw.statics != nil {
		// 更新加全局锁
		rw.l.Lock()
		rw.updatePos()
		rw.l.Unlock()
		// 取数据加读写锁
		rw.l.RLock()
		defer rw.l.RUnlock()
		return rw.statics.Total(rw.datas[rw.pos])
	}
	return 0
}

func (rw *RollingWindow) updatePos() {
	offset := rw.offset()
	if offset > 0 {
		pos := rw.pos
		// reset expired buckets
		start := pos + 1
		steps := start + offset
		var remainder int
		if steps > rw.size {
			remainder = steps - rw.size
			steps = rw.size
		}
		for i := start; i < steps; i++ {
			rw.resetData(i)
			pos = i
		}
		for i := 0; i < remainder; i++ {
			rw.resetData(i)
			pos = i
		}
		rw.pos = pos
		rw.lastTime = rw.timec.Now()
	}
}

func (rw *RollingWindow) resetData(i int) {
	if rw.statics != nil {
		rw.statics.Reset(rw.datas[i], i)
	}
	rw.datas[i].Reset()
}

func (rw *RollingWindow) add(v float64) {
	if rw.statics != nil {
		rw.statics.Add(v, rw.pos)
	}
	rw.datas[rw.pos].Add(v)
}

func (rw *RollingWindow) UpdateOpts(opts ...RollingWindowOption) {
	for _, opt := range opts {
		opt(rw)
	}
}

type Item struct {
	Val   float64
	Total int
}

func (it *Item) Add(v float64) {
	it.Val += v
	it.Total++
}

func (it *Item) Reset() {
	it.Val = 0
	it.Total = 0
}

func IgnoreCurrentBucket(b bool) RollingWindowOption {
	return func(w *RollingWindow) {
		w.ignoreCurrent = b
		if w.statics != nil {
			w.statics.SetIgnoreCurrent(b)
		}
	}
}

func WithStatic(staticType StaticType) RollingWindowOption {
	return func(w *RollingWindow) {
		w.statics = NewStatic(staticType, w.size)
		w.statics.SetIgnoreCurrent(w.ignoreCurrent)
	}
}

// suggest only use Init
func WithTimeC(timec TimecIface) RollingWindowOption {
	return func(w *RollingWindow) {
		w.timec = timec
		w.startTime = timec.Now()
		w.lastTime = timec.Now()
	}
}
