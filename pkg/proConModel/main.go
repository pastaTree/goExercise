package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/*
使用 goroutine 和 channel 实现一个计算 int64 随机数各位树和的程序
1. 开启1个 goroutine 循环生成 int64 的随机数, 发送到 jobchan
2. 开启24个 goroutine 从 jobchan 中取出随机数计算各位数的和, 将结果发送到 resultchan 中
3. 主 goroutine 从 resultchan 中取得结果, 将其打印到终端中
 */

type job struct {
	value int64
}

type result struct {
	job *job
	sum int64
}

var jobChan = make(chan *job, 100)
var resultChan = make(chan *result, 100)
var wg sync.WaitGroup

func producer(j chan<- *job) {
	defer wg.Done()
	for {
		x := rand.Int63()
		newJob := &job{
			value: x,
		}
		j <- newJob
		time.Sleep(time.Millisecond * 500)
	}
}

func consumer(j <-chan *job, r chan<- *result) {
	defer wg.Done()
	for {
		jobObj := <- j
		var sum int64
		value := jobObj.value
		for value != 0 {
			sum += value % 10
			value = value / 10
		}
		newResult := &result{
			job:    jobObj,
			sum: sum,
		}
		r <- newResult
	}
}

func main() {
	wg.Add(1)
	go producer(jobChan)
	for i := 0; i < 24; i++ {
		wg.Add(1)
		go consumer(jobChan, resultChan)
	}
	for result := range resultChan {
		fmt.Printf("random value: %d, sum: %d\n", result.job.value, result.sum)
	}
	wg.Wait()
}