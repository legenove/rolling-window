package windows

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

const size = 3
const durationStatic = time.Millisecond * 50
const testPecent = 70

//var timec = &DefaultTimec{}
var timec = newTestTimec()
var testIgnoreValue = []bool{true, false}

func elapseStatic() {
	timec.Next(durationStatic)
}

type testInput struct {
	Datas [][]float64
	Test  []bool
}

func doValied(t *testing.T, r *RollingWindow, val float64, total int, msgAndArgs ...interface{}) {
	assert.Equal(t, val, r.StaticValue(), msgAndArgs...)
	assert.Equal(t, total, r.StaticTotal(), msgAndArgs...)
}

func testgetStaticFunc(staticType StaticType) WinCalculateFunc {
	switch staticType {
	case StaticSum:
		return GetRollingWindowSum
	case StaticMax:
		return GetRollingWindowMax
	case StaticMin:
		return GetRollingWindowMin
	case StaticAvg:
		return GetRollingWindowAvg
	}
	return nil
}

func doTestStatic(t *testing.T, r *RollingWindow, fn WinCalculateFunc, count *int32, msgAndArgs ...interface{}) {
	if testRand() {
		atomic.AddInt32(count, 1)
		gv, gt := fn(r)
		doValied(t, r, gv, gt, msgAndArgs...)
	}
}

func testRand() bool {
	if rand.Int()%100 < testPecent {
		return true
	}
	return false
}

func testStatic(t *testing.T, testDatas []testInput, staticType StaticType, testTime int) {
	var count int32
	fn := testgetStaticFunc(staticType)
	for g, testData := range testDatas {
		for _, ig := range testIgnoreValue {
			for i := 0; i < testTime; i++ {
				atomic.StoreInt32(&count, 0)
				r := NewRollingWindow(size, duration, IgnoreCurrentBucket(ig), WithStatic(staticType), WithTimeC(timec))
				for i, vals := range testData.Datas {
					for _, val := range vals {
						if testData.Test[i] {
							doTestStatic(t, r, fn, &count,
								"before \nig: %v \ncv : %v, \nc: %d", ig, vals, val)
						}
						r.Add(val)
					}
					if testData.Test[i] {
						doTestStatic(t, r, fn, &count, "end \nig: %v \ncv : %v", ig, vals)
					}
					elapseStatic()
				}
				_ = g
				//testGroup := fmt.Sprintf("测试组：%d-%d\n", g, i+1)
				//testValue := fmt.Sprintf("ignore: %v, 测试值：%v\n", ig, testData)
				//testCount := fmt.Sprintf("测试量：%d\n", count)
				//t.Log(testGroup, testValue, testCount)
			}
		}

	}
}

func testStaticData() []testInput {
	return []testInput{
		{
			[][]float64{{1}, {2, 3}, {4, 5, 6}, {7}, {8}, {6}, {}, {}, {}},
			[]bool{true, true, true, true, true, true, false, true, true},
		},
		{
			[][]float64{{1}, {2, 3}, {}, {}, {}, {}, {4, 5, 6}, {7}, {8}, {6}, {}, {}, {}},
			[]bool{true, true, true, true, true, true, true, true, true, true, false, true, true},
		},
		{
			[][]float64{{-1}, {-2, -3}, {}, {}, {}, {}, {-4, -5, -6}, {-7}, {-8}, {-6}, {}, {}, {}},
			[]bool{true, true, true, true, true, true, true, true, true, true, false, true, true},
		},
		{
			[][]float64{{-1}, {-2, 3}, {}, {-1, 1}, {}, {}, {-4, 5, -6}, {-7}, {-8}, {-6}, {}, {}, {}},
			[]bool{true, true, true, true, true, true, true, true, true, true, false, true, true},
		},
	}
}

func TestStaticSum(t *testing.T) {
	fmt.Sprintf("")
	const staticType = StaticSum
	const testTime = 5
	testStatic(t, testStaticData(), staticType, testTime)
}

func TestStaticMax(t *testing.T) {
	const staticType = StaticMax
	const testTime = 5
	testStatic(t, testStaticData(), staticType, testTime)
}

func TestStaticMin(t *testing.T) {
	const staticType = StaticMin
	const testTime = 5
	testStatic(t, testStaticData(), staticType, testTime)
}


func TestStaticAvg(t *testing.T) {
	const staticType = StaticAvg
	const testTime = 5
	testStatic(t, testStaticData(), staticType, testTime)
}
