package windows

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

const duration = time.Millisecond * 50

var timecRWList = []TimecIface{newTestTimec(), &DefaultTimec{}}

func elapse(timec TimecIface) {
	timec.Next(duration)
}

func getRandomString() string {
	return fmt.Sprintf("%d", rand.Int63())
}

func TestRollingWindowAdd(t *testing.T) {
	const size = 3
	for _, timec := range timecRWList {
		r := NewRollingWindow(size, duration, WithTimeC(timec))
		listBuckets := func() []float64 {
			var buckets []float64
			r.Cal(func(b *Item) {
				buckets = append(buckets, b.Val)
			})
			return buckets
		}
		assert.Equal(t, []float64{0, 0, 0}, listBuckets())
		r.Add(1)
		assert.Equal(t, []float64{0, 0, 1}, listBuckets())
		elapse(timec)
		r.Add(2)
		r.Add(3)
		assert.Equal(t, []float64{0, 1, 5}, listBuckets())
		elapse(timec)
		r.Add(4)
		r.Add(5)
		r.Add(6)
		assert.Equal(t, []float64{1, 5, 15}, listBuckets())
		elapse(timec)
		r.Add(7)
		assert.Equal(t, []float64{5, 15, 7}, listBuckets())
	}
}

func TestRollingWindowReset(t *testing.T) {
	const size = 3

	for _, timec := range timecRWList {
		r := NewRollingWindow(size, duration, IgnoreCurrentBucket(true), WithTimeC(timec))
		listBuckets := func() []float64 {
			var buckets []float64
			r.Cal(func(b *Item) {
				buckets = append(buckets, b.Val)
			})
			return buckets
		}
		r.Add(1)
		elapse(timec)
		assert.Equal(t, []float64{0, 1}, listBuckets())
		elapse(timec)
		assert.Equal(t, []float64{1}, listBuckets())
		elapse(timec)
		assert.Nil(t, listBuckets())

		// cross window
		r.Add(1)
		for i := 0; i <= size; i++ {
			elapse(timec)
		}
		assert.Nil(t, listBuckets())
	}
}

func TestRollingWindowReduce(t *testing.T) {
	const size = 4
	for _, timec := range timecRWList {
		tests := []struct {
			win    *RollingWindow
			expect float64
		}{
			{
				win:    NewRollingWindow(size, duration, WithTimeC(timec)),
				expect: 10,
			},
			{
				win:    NewRollingWindow(size, duration, IgnoreCurrentBucket(true), WithTimeC(timec)),
				expect: 4,
			},
		}

		for _, test := range tests {
			t.Run(getRandomString(), func(t *testing.T) {
				r := test.win
				for x := 0; x < size; x = x + 1 {
					for i := 0; i <= x; i++ {
						r.Add(float64(i))
					}
					if x < size-1 {
						elapse(timec)
					}
				}
				var result float64
				r.Cal(func(b *Item) {
					result += b.Val
				})
				assert.Equal(t, test.expect, result)
			})
		}
	}
}

func TestRollingWindowDataRace(t *testing.T) {
	const size = 3

	for _, timec := range timecRWList {
		r := NewRollingWindow(size, duration, WithTimeC(timec))
		var stop = make(chan bool)
		go func() {
			for {
				select {
				case <-stop:
					return
				default:
					r.Add(float64(rand.Int63()))
					time.Sleep(duration / 2)
				}
			}
		}()
		go func() {
			for {
				select {
				case <-stop:
					return
				default:
					r.Cal(func(b *Item) {})
				}
			}
		}()
		time.Sleep(duration * 5)
		close(stop)
	}
}
