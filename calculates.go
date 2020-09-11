package windows

import "math"

type WinCalculateFunc func(rw *RollingWindow) (valSum float64, total int)

func GetRollingWindowSum(rw *RollingWindow) (valSum float64, total int) {
	rw.Cal(func(b *Item) {
		valSum += b.Val
		total += b.Total
	})
	return
}

func GetRollingWindowAvg(rw *RollingWindow) (valSum float64, total int) {
	rw.Cal(func(b *Item) {
		valSum += b.Val
		total += b.Total
	})
	valSum = valSum / float64(rw.GetStaticItemNum())
	return
}

func GetRollingWindowMax(rw *RollingWindow) (float64, int) {
	var max *float64
	var _max float64
	var total int
	var cnt   int
	rw.Cal(func(b *Item) {
		if max == nil {
			_max = b.Val
			max = &_max
		} else {
			_max = math.Max(*max, b.Val)
			max = &_max
		}
		total += b.Total
		cnt += 1
	})
	if max == nil {
		return 0, total
	}
	if cnt < rw.GetStaticItemNum()  {
		_max = math.Max(*max, 0)
		max = &_max
	}
	return *max, total
}

func GetRollingWindowMin(rw *RollingWindow) (float64, int) {
	var min *float64
	var _min float64
	var total int
	var cnt   int
	rw.Cal(func(b *Item) {
		if min == nil {
			_min = b.Val
			min = &_min
		} else {
			_min = math.Min(*min, b.Val)
			min = &_min
		}
		total += b.Total
		cnt += 1
	})
	if min == nil {
		return 0, total
	}
	if cnt < rw.GetStaticItemNum()  {
		_min = math.Min(*min, 0)
		min = &_min
	}
	return *min, total
}
