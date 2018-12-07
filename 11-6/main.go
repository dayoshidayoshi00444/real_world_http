package main

import (
	"context"
	"golang.org/x/time/rate"
	"time"
)

func main() {
	// 1秒あたりの上限回数
	RateLimit := 10
	// トークンの最大保持数
	BucketSize := 10
	ctx := context.Background()
	e := rate.Every(time.Second / RateLimit)
	l := rate.NewLimiter(e, BucketSize)

	for _, task := range tasks {
		err := l.Wait(ctx)
		if err != nil {
			panic(err)
		}
	}
}
