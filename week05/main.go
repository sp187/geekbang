package main

import (
	"fmt"
	"time"

	"github.com/sp187/geekbang/week05/metric"
)

func main() {
	// 采样间隔1秒， bucket数量10， 每个bucket记录100ms内的请求数量
	rolling := metric.NewTimerCounter(time.Second, 10)
	done := make(chan struct{})
	rolling.Start()
	defer rolling.Close()

	// 模拟 qps = 500
	go func() {
		tic := time.NewTicker(2*time.Millisecond)
		defer tic.Stop()
		for {
			select {
			case <-tic.C:
				rolling.Add()
			case <- done:
				return
			}
		}
	}()

	// 模拟 qps = 200
	go func() {
		tic := time.NewTicker(5*time.Millisecond)
		defer tic.Stop()
		for {
			select {
			case <-tic.C:
				rolling.Add()
			case <- done:
				return
			}
		}
	}()



	// 1秒钟采样窗口内的数据 500 + 200 = 700
	go func() {
		t := time.NewTicker(1*time.Second)
		defer t.Stop()
		for i := 1; i <= 10; i++ {
			select {
			case <-t.C:
				fmt.Printf("第%d秒采样: %d\n", i, rolling.Sample())
			}
		}
	}()

	time.Sleep(10100*time.Millisecond)
	close(done)
}
