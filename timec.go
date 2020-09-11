package windows

import "time"

type TimecIface interface {
	Now() time.Time
	Next(d time.Duration)
}

type DefaultTimec struct {
}

func (*DefaultTimec) Now() time.Time {
	return time.Now()
}

func (*DefaultTimec) Next(d time.Duration) {
	time.Sleep(d)
}

type forTestTimec struct {
	now time.Time
}

func newTestTimec() *forTestTimec {
	return &forTestTimec{now: time.Now()}
}

func (tc *forTestTimec) Now() time.Time {
	return tc.now
}

func (tc *forTestTimec) Next(d time.Duration) {
	tc.now = tc.now.Add(d)
}
