// Package wrr weight round robin加权轮询算法
package wrr

import "going/utils"

// Weight round robin weight item 权重项
type Weight struct {
	name      string
	weight    int
	curWeight int
}

// SWRR smooth weight round robin 平滑加权轮询，每次选择都要轮询一遍所有权重，性能较wrr稍差，但是更加平滑，应用于nginx
type SWRR struct {
	totalWeight int
	weights     []*Weight
}

// WRR weight round robin 加权轮询，以最高权重为基准，逐个对比，每轮降低一个权值，性能较高，但对于权重相差较大的情况，不够平滑, 应用于tgw lvs l5 cmlb
type WRR struct {
	weights    []*Weight
	itemSize   int
	curItem    int
	maxWeight  int
	baseWeight int
}

// NewSWRR new smooth weight round robin, nginx
func NewSWRR(weights map[string]uint32) *SWRR {
	if len(weights) < 1 {
		return nil
	}

	c := &SWRR{
		weights: make([]*Weight, 0),
	}

	for name, weight := range weights {
		c.totalWeight += int(weight)
		w := &Weight{
			name:   name,
			weight: int(weight),
		}
		c.weights = append(c.weights, w)
	}

	return c
}

// Next choose swrr next weight name, 每次循环，curWeight加上权重，选出最高curWeight，最后减去总权重，算法很简单，背后的数学原理比较高深
func (c *SWRR) Next() string {
	var max int
	var w *Weight

	for _, v := range c.weights {
		v.curWeight += v.weight

		if v.curWeight > max {
			max = v.curWeight
			w = v
		}
	}

	if w == nil {	// add by vianyan, 2018-09-11
		return ""
	}

	w.curWeight -= c.totalWeight
	return w.name
}

// NewWRR new weight round robin, tgw lvs l5 cmlb
func NewWRR(weights map[string]uint32) *WRR {
	if len(weights) < 1 {
		return nil
	}

	c := &WRR{
		weights: make([]*Weight, 0),
	}

	ws := make([]int, 0)
	for _, w := range weights {
		if w == 0 {
			return nil
		}
		ws = append(ws, int(w))
	}
	gcd := utils.GCD(ws...)

	for name, weight := range weights {
		c.itemSize++

		w := &Weight{
			name:   name,
			weight: int(weight) / gcd,
		}
		c.weights = append(c.weights, w)

		if w.weight > c.maxWeight {
			c.maxWeight = w.weight
		}
	}

	c.baseWeight = c.maxWeight

	return c
}

// Next choose wrr next weight name 每次对比的基准权重从最高权重开始， 每轮挑选出大于等于基准权重的项
func (c *WRR) Next() string {

	var w *Weight

	for {
		if c.weights[c.curItem].weight >= c.baseWeight {
			w = c.weights[c.curItem]
		}

		c.curItem++
		if c.curItem == c.itemSize { // 每一轮是一个小循环
			c.curItem = 0
			c.baseWeight--
		}
		if c.baseWeight == 0 { // 权重比较完，进入下一个大循环
			c.baseWeight = c.maxWeight
		}

		if w != nil {
			return w.name
		}
	}
}
