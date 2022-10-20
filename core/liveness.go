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

	startTime          time.Time // 有效活跃开始时间
	endTime            time.Time // 有效活跃结束时间
	lastUpdateTime     time.Time // 上次更新时间
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
	ln.calcValidTime()

	go ln.watch()
	return ln
}

func (ln *lnClient) calcValidTime() {
	now := time.Now()
	ln.startTime = time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, now.Location())
	ln.endTime = time.Date(now.Year(), now.Month(), now.Day(), 19, 30, 0, 0, now.Location())
	fmt.Printf("今日活跃时间：%s ~ %s\n", ln.startTime.Format("2006-01-02 15:04:05"), ln.endTime.Format("2006-01-02 15:04:05"))
}

func (ln *lnClient) isContinue() bool {
	now := time.Now()
	if now.Year() != ln.startTime.Year() || now.Month() != ln.startTime.Month() || now.Day() != ln.startTime.Day() {
		// 天数变化
		ln.liveness = 0
		ln.calcValidTime()
		return true
	}
	if now.Before(ln.startTime) || now.After(ln.endTime) {
		return false
	}
	if ln.liveness >= 100 {
		return false
	}

	return true
}

func (ln *lnClient) Say() {
	if !ln.isContinue() {
		return
	}

	now := time.Now()
	if now.Sub(ln.lastUpdateTime).Seconds() < 30 {
		return
	}

	ln.lastUpdateTime = now
	ln.mu.Lock()
	defer ln.mu.Unlock()

	ln.liveness += 1.67

	go func() {
		time.Sleep(30 * time.Second)

		all := math.Ceil(100 / ln.inc)

		need := math.Ceil((100.00 - ln.liveness) / ln.inc)

		t := (time.Duration(math.Ceil(need*ln.interval)) * time.Second).Minutes()

		if ln.liveness < 100 {
			fmt.Printf("还差(%.f/%.f) 预计还需%.f分钟\n", need, all, t)
		} else {
			fmt.Println("你已经满了 快去code吧")
		}
	}()
}

func (ln *lnClient) watch() {
	ticker := time.NewTicker(time.Minute * 10)
	for {
		if ln.liveness >= 100 {
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
