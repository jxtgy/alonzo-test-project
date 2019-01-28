// Package gb golang benchmark 服务压测工具, 类似ab, 区别于压测自身代码的标准库 testing.benchmark
package gb

/*
@version v1.0
@author nickzydeng
@copyright Copyright (c) 2018 Tencent Corporation, All Rights Reserved
@license http://opensource.org/licenses/gpl-license.php GNU Public License

You may not use this file except in compliance with the License.

Most recent version can be found at:
http://git.code.oa.com/going_proj/going_proj

Please see README.md for more information.
*/

import (
	"context"
	"flag"
	"log"
	"os"
	"sync"
	"time"
)

func init() {
	concurrent = flag.Int("c", 1, "concurrent")
	duration = flag.Duration("d", 1*time.Second, "duration")
	timeout = flag.Duration("t", 800*time.Millisecond, "request timeout")
}

var concurrent *int
var duration *time.Duration
var timeout *time.Duration
var deadline time.Time

// Result 统计成功率
type Result struct {
	failed, success                         int64
	totalTimeCost, minTimeCost, maxTimeCost int64
}

// Run 并发开启goroutine压测
func Run(req func(context.Context) error) {
	flag.Parse()
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	if req == nil {
		log.Fatal("request is nil")
	}
	if *concurrent < 1 {
		*concurrent = 1
	}
	deadline = time.Now().Add(*duration)

	results := make([]*Result, *concurrent, *concurrent)
	var wg sync.WaitGroup
	wg.Add(*concurrent)
	for i := 0; i < *concurrent; i++ {
		go func(index int) {
			defer wg.Done()
			results[index] = worker(req)
		}(i)
	}
	wg.Wait()
	var totalTimeCost, minTimeCost, maxTimeCost, averageTimeCost int64
	var success, failed int64
	minTimeCost, maxTimeCost = 800, 800
	for _, d := range results {
		totalTimeCost += d.totalTimeCost
		if d.maxTimeCost > maxTimeCost {
			maxTimeCost = d.maxTimeCost
		}
		if d.minTimeCost < minTimeCost {
			minTimeCost = d.minTimeCost
		}
		success += d.success
		failed += d.failed
	}
	if success+failed > 0 {
		averageTimeCost = totalTimeCost / (success + failed)
	}

	log.Println("----------------------------------")
	log.Printf("failed: %v", failed)
	log.Printf("success: %v", success)
	log.Printf("qps: %v", (failed+success)*1e9/duration.Nanoseconds())
	log.Printf("average time cost: %v ms", averageTimeCost/1e6)
	log.Printf("min time cost: %v ms", minTimeCost/1e6)
	log.Printf("max time cost: %v ms", maxTimeCost/1e6)
	log.Println("----------------------------------")
}

func worker(req func(context.Context) error) *Result {
	result := Result{minTimeCost: time.Hour.Nanoseconds()}
	for {
		if time.Now().After(deadline) {
			break
		}
		reqCtx, cancel := context.WithTimeout(context.Background(), *timeout)
		defer cancel()
		beginTime := time.Now()
		err := req(reqCtx)
		timeCost := time.Now().Sub(beginTime).Nanoseconds()
		result.totalTimeCost += timeCost
		if timeCost > result.maxTimeCost {
			result.maxTimeCost = timeCost
		}
		if timeCost > 1e6 && timeCost < result.minTimeCost {
			result.minTimeCost = timeCost
		}
		if err == nil {
			result.success++
		} else {
			result.failed++
		}
	}
	return &result
}
