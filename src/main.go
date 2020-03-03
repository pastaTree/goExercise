package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type ByteSize float64
const (
	_ = iota
	KB ByteSize = 1<<(10*iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

func main() {
	Random()
    n := 20.1
    _ = IntFromFloat64(n)
}

func IntFromFloat64(x float64) int {
	if math.MinInt32 <= x && x <= math.MaxInt32 {
		whole, fraction := math.Modf(x)
		if fraction >= 0.5 {
			whole++
		}
		return int(whole)
	}
	panic(fmt.Sprintf("%g is out of int32 range", x))
}

func Random() {
	for i := 0; i < 10; i++ {
		a := rand.Int()
		fmt.Printf("%d / ", a)
	}

	for i := 0; i < 5; i++ {
		r := rand.Intn(8)
		fmt.Printf("%d / ", r)
	}

	fmt.Println()
	timens := int64(time.Now().Nanosecond())
	rand.Seed(timens)
	for i := 0; i < 10; i++ {
		fmt.Printf("%2.2f / ", 100*rand.Float32())
	}
}