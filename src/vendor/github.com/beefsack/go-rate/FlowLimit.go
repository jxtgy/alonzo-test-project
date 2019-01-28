package rate

import (
	"container/list"
	"sync"
	"time"
)


type FlowLimiter struct {
	limit    int64  // 限制速率
	current  int64  // 当前流量

	interval time.Duration
	mtx      sync.Mutex
	times    list.List
}


type FlowElem struct {
	time time.Time
	flow int64
}


func NewFlowLimiter(limit int64, interval time.Duration) *FlowLimiter {
	if limit < 1 || interval < 1 {
		return nil
	}

	lim := &FlowLimiter{
		limit:    limit,
		interval: interval,
	}
	lim.times.Init()
	return lim
}


func (fl *FlowLimiter) Try(flow int64) (ok bool, current int64) {
	fl.mtx.Lock()
	defer fl.mtx.Unlock()

	now := time.Now()


	// 均摊到每次，削峰!!
	for {
		front := fl.times.Front()
		var diff time.Duration = 0
		var frontFlow int64 = 0
		if front != nil {
			flowElem := front.Value.(FlowElem)
			diff = now.Sub(flowElem.time)
			frontFlow = flowElem.flow
		}

		if diff >= fl.interval {
			fl.current -= frontFlow
			fl.times.Remove(front)
		} else {
			if fl.current + flow <= fl.limit {
				elem := FlowElem{now, flow}
				fl.times.PushBack(elem)
				fl.current += flow
				return true, fl.current
			} else {
				return false, fl.current
			}
		}
	}

}

