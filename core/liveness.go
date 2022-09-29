package core

import (
	"fmt"
	"math"
	"sync"
	"time"
)

type lnClient struct {
	liveness float64 // 当前活跃度
	inc      float64 // 增量活跃
	interval float64 // 活跃间隔

	lastUpdateTime     time.Time
	updateLivenessFunc func(float64)

	mu sync.Mutex
}

func NewLnClient(liveness float64, updateLivenessFunc func(float64)) *lnClient {
	ln := &lnClient{
		liveness: liveness,
		inc:      1.67,
		interval: 40,

		lastUpdateTime:     time.Now().Add(-30 * time.Second),
		updateLivenessFunc: updateLivenessFunc,
	}

	go ln.watch()
	return ln
}

func (ln *lnClient) Say() {
	now := time.Now()
	if now.Sub(ln.lastUpdateTime).Seconds() < 30 {
		return
	}

	ln.lastUpdateTime = now
	ln.mu.Lock()
	defer ln.mu.Unlock()

	ln.liveness += 1.67

	go func() {
		if ln.liveness == 100 {
			return
		}
		time.Sleep(30 * time.Second)

		all := math.Ceil(100 / ln.inc)

		need := math.Ceil((100.00 - ln.liveness) / ln.inc)

		t := (time.Duration(math.Ceil(need*ln.interval)) * time.Second).Minutes()

		fmt.Printf("还差(%.f/%.f) 预计还需%.f分钟\n", need, all, t)
	}()
}

func (ln *lnClient) watch() {
	ticker := time.NewTicker(time.Minute * 10)
	for {
		if ln.liveness == 100 {
			ticker.Stop()
			return
		}
		select {
		case <-ticker.C:
			ln.mu.Lock()
			ln.updateLivenessFunc(ln.liveness)
			ln.mu.Unlock()
		}
	}
}
