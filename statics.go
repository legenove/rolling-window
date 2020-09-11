package windows

import (
	"math"
)

type StaticType int

const (
	StaticBase StaticType = iota
	StaticSum
	StaticMin
	StaticMax
	StaticAvg
)

type WindowStatic interface {
	Value(i *Item) float64
	Total(i *Item) int
	Static(i *Item) (float64, int)
	Add(v float64, index int)
	Reset(i *Item, index int)
	SetIgnoreCurrent(i bool)
}

func NewStatic(staticType StaticType, size int) WindowStatic {
	switch staticType {
	case StaticBase:
		return nil
	case StaticSum:
		return newSumStatic()
	case StaticMin:
		return newMinStatic(size)
	case StaticMax:
		return newMaxStatic(size)
	case StaticAvg:
		return newAgvStatic(size)
	}
	return nil
}

type BaseStatic struct {
	TotalSum      int
	ignoreCurrent bool
}

func (s *BaseStatic) Total(i *Item) int {
	if i != nil && s.ignoreCurrent {
		return s.TotalSum - i.Total
	}
	return s.TotalSum
}

func (s *BaseStatic) SetIgnoreCurrent(i bool) {
	s.ignoreCurrent = i
}

/*
 * space o(1)
 * value o(1)
 * add [avg] o(1)
 * reset o(1)
 */

type SumStatic struct {
	BaseStatic
	ValSum float64
}

func newSumStatic() *SumStatic {
	return &SumStatic{}
}

func (s *SumStatic) Value(i *Item) float64 {
	if i != nil && s.ignoreCurrent {
		return s.ValSum - i.Val
	}
	return s.ValSum
}

func (s *SumStatic) Static(i *Item) (float64, int) {
	return s.Value(i), s.Total(i)
}

func (s *SumStatic) Add(v float64, index int) {
	s.ValSum += v
	s.TotalSum++
}

func (s *SumStatic) Reset(i *Item, index int) {
	s.ValSum -= i.Val
	s.TotalSum -= i.Total
}

/*
 * space o(k)
 * value o(1)
 * add [avg] o(1)
 * reset o(1)
 */

type MaxStatic struct {
	BaseStatic
	Size   int
	curPos int
	cap    int
	curVal float64
	Queue  *CircularQueue
}

func newMaxStatic(size int) *MaxStatic {
	return &MaxStatic{Queue: NewCircularQueue(size), Size: size}
}

func (s *MaxStatic) Value(i *Item) float64 {
	// 初始桶是0值，所以如果队列没满，桶的最大值是0值
	var max = s.Queue.First()
	if max < 0 && s.cap < s.Size-1 {
		max = 0
	}
	if s.ignoreCurrent {
		return max
	}
	return math.Max(max, i.Val)
}

func (s *MaxStatic) Static(i *Item) (float64, int) {
	return s.Value(i), s.Total(i)
}

func (s *MaxStatic) Add(v float64, index int) {
	if index != s.curPos {
		d := (index + s.Size - s.curPos) % s.Size
		s.Queue.PushEmpty(d)
		s.curPos = index
	}
	s.curVal += v
	s.TotalSum++
}

func (s *MaxStatic) getIndex(i int) int {
	return i % s.Size
}

func (s *MaxStatic) Reset(i *Item, index int) {
	if i.Val == s.Queue.First() {
		s.Queue.Shift()
	}
	if s.curVal != 0 {
		s.cap++
	}
	if index == s.getIndex(s.curPos+1) && s.curVal != 0 {
		// 说明上一个统计结束了
		for !s.Queue.IsEmpty() && s.Queue.Last() < s.curVal {
			s.Queue.Pop()
		}
		s.Queue.Push(s.curVal)
		s.curVal = 0
		s.curPos = s.getIndex(s.curPos + 1)
	}
	s.TotalSum -= i.Total
	if i.Val != 0 && s.cap > 0 {
		s.cap--
	}
}

/*
 * space o(k)
 * value o(1)
 * add [avg] o(1)
 * reset o(1)
 */

type MinStatic struct {
	BaseStatic
	Size   int
	curPos int
	curVal float64
	cap    int
	Queue  *CircularQueue
}

func newMinStatic(size int) *MinStatic {
	return &MinStatic{Queue: NewCircularQueue(size), Size: size}
}

func (s *MinStatic) Value(i *Item) float64 {
	// 初始桶是0值，所以如果队列没满，桶的最大值是0值
	var min = s.Queue.First()
	if min > 0 && s.cap < s.Size-1 {
		min = 0
	}
	if s.ignoreCurrent {
		return min
	}
	return math.Min(min, i.Val)
}

func (s *MinStatic) Static(i *Item) (float64, int) {
	return s.Value(i), s.Total(i)
}

func (s *MinStatic) Add(v float64, index int) {
	if index != s.curPos {
		d := (index + s.Size - s.curPos) % s.Size
		s.Queue.PushEmpty(d)
		s.curPos = index
	}
	s.curVal += v
	s.TotalSum++
}

func (s *MinStatic) Reset(i *Item, index int) {
	if i.Val == s.Queue.First() {
		s.Queue.Shift()
	}
	if s.curVal != 0 {
		s.cap++
	}
	if index == s.getIndex(s.curPos+1) && s.curVal != 0 {
		// 说明上一个统计结束了
		for !s.Queue.IsEmpty() && s.Queue.Last() > s.curVal {
			s.Queue.Pop()
		}
		s.Queue.Push(s.curVal)
		s.curVal = 0
		s.curPos = s.getIndex(s.curPos + 1)
	}
	s.TotalSum -= i.Total
	if i.Val != 0 && s.cap > 0 {
		s.cap--
	}
}

func (s *MinStatic) getIndex(i int) int {
	return i % s.Size
}

/*
 * space o(1)
 * value o(1)
 * add [avg] o(1)
 * reset o(1)
 */

type AgvStatic struct {
	BaseStatic
	Size   int
	ValSum float64
}

func newAgvStatic(size int) *AgvStatic {
	return &AgvStatic{Size: size}
}

func (s *AgvStatic) Value(i *Item) float64 {
	if i != nil && s.ignoreCurrent {
		return (s.ValSum - i.Val)/float64(s.Size-1)
	}
	return s.ValSum / float64(s.Size)
}

func (s *AgvStatic) Static(i *Item) (float64, int) {
	return s.Value(i), s.Total(i)
}

func (s *AgvStatic) Add(v float64, index int) {
	s.ValSum += v
	s.TotalSum++
}

func (s *AgvStatic) Reset(i *Item, index int) {
	s.ValSum -= i.Val
	s.TotalSum -= i.Total
}
