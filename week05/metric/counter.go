package metric

import (
	"sync"
	"time"
)

// CounterWin 滑动窗口
type CounterWin interface {
	Add()          // 计数+1
	Next()         // 窗口滑动
	Sample() int64 // 获取当前采样窗口内的值
}

// TimerCounterWin  时间滑动窗口
type TimerCounterWin interface {
	CounterWin
	Start()
	Close()
}

type _counter struct {
	ring           *_ring
	bucketDuration time.Duration
	winSize        int64
	sampleInterval time.Duration // 采样间隔
	done           chan struct{}
}

func (c *_counter) Add() {
	c.ring.add(1)
}

var startOnce sync.Once

func (c *_counter) Start() {
	startOnce.Do(func() {
		go func() {
			tic := time.NewTicker(c.bucketDuration)
			defer tic.Stop()
			for {
				select {
				case <-tic.C:
					c.ring.next()
				case <-c.done:
					return
				}
			}
		}()
	})
}

var closeOnce sync.Once

func (c *_counter) Close() {
	closeOnce.Do(func() {
		close(c.done)
		c.done = nil
	})
}

func (c *_counter) Sample() int64 {
	return c.ring.sum
}

func (c *_counter) Next() {
	c.ring.next()
}

func NewTimerCounter(interval time.Duration, winSize int64) TimerCounterWin {
	return &_counter{
		ring:           newRing(winSize),
		bucketDuration: interval / time.Duration(winSize),
		winSize:        winSize,
		sampleInterval: interval,
		done:           make(chan struct{}),
	}
}

type _ring struct {
	ring                    []int64
	head, tail, length, sum int64
	full                    bool
	sync.Mutex
}

func newRing(n int64) *_ring {
	return &_ring{
		ring:   make([]int64, n),
		head:   0,
		tail:   0,
		length: n,
	}
}

func (r *_ring) next() {
	r.Lock()
	defer r.Unlock()
	if r.full {
		r.sum -= r.ring[r.head]
		r.ring[r.head] = 0
		r.head++
		r.tail++
		if r.head >= r.length {
			r.head = 0
		}
		if r.tail >= r.length {
			r.tail = 0
		}
	} else {
		r.tail++
		if r.tail-r.head == r.length-1 {
			r.full = true
		}
	}
}

func (r *_ring) add(n int64) {
	r.Lock()
	defer r.Unlock()
	r.sum += n
	r.ring[r.tail] += n
}
