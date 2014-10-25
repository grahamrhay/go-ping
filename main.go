package main

import (
	"container/ring"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Samples struct {
	Last3Hours  *ring.Ring
	Last30Hours *ring.Ring
}

func main() {
	samples := &Samples{Last3Hours: ring.New(36), Last30Hours: ring.New(30)}
	startOfLast3Hours := &samples.Last3Hours
	startOfLast30Hours := &samples.Last30Hours

	startServer(startOfLast3Hours, startOfLast30Hours)

	ticker1 := time.NewTicker(5 * time.Minute)
	go func() {
		takeSample(samples)
		for _ = range ticker1.C {
			takeSample(samples)
		}
	}()

	ticker2 := time.NewTicker(1 * time.Hour)
	go func() {
		for _ = range ticker2.C {
			downSample(samples)
		}
	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Received signal: %v\n", <-ch)
	log.Println("Shutting down")
	ticker1.Stop()
	ticker2.Stop()
}

func takeSample(samples *Samples) {
	res := ping("google.com", 20)
	samples.Last3Hours.Value = res
	samples.Last3Hours = samples.Last3Hours.Next()
}

func downSample(samples *Samples) {
	log.Println("Down sampling")
	avg := 0.0
	count := 0.0
	samples.Last3Hours.Do(func(value interface{}) {
		if value == nil {
			return
		}

		res := value.(*PingResult)
		avg += res.Avg
		count++
	})
	samples.Last30Hours.Value = &PingResult{Time: time.Now(), Avg: avg / count}
	samples.Last30Hours = samples.Last30Hours.Next()
}
