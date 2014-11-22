package main

import (
	"container/ring"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type Samples struct {
	Last3Hours  *ring.Ring
	Last30Hours *ring.Ring
}

const NLast3Hours = 36

func main() {
	store, err := openStore()
	if err != nil {
		log.Fatal("Failed to open data store: %v", err)
	}

	samples := &Samples{Last3Hours: ring.New(NLast3Hours), Last30Hours: ring.New(30)}
	startOfLast3Hours := &samples.Last3Hours
	startOfLast30Hours := &samples.Last30Hours

	startServer(startOfLast3Hours, startOfLast30Hours)

	err = loadSamples(store, samples)
	if err != nil {
		log.Fatal("Error loading samples: %v\n", err)
	}

	//ticker1 := time.NewTicker(5 * time.Minute)
	ticker1 := time.NewTicker(10 * time.Second)
	go func() {
		takeSample(samples)
		persistRing(store, startOfLast3Hours, "last3Hours")
		for _ = range ticker1.C {
			takeSample(samples)
			persistRing(store, startOfLast3Hours, "last3Hours")
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
	closeStore(store)
}

func takeSample(samples *Samples) {
	res := ping("google.com", 5)
	samples.Last3Hours.Value = res
	samples.Last3Hours = samples.Last3Hours.Next()
}

func persistRing(store *Store, ring **ring.Ring, prefix string) {
	i := 0
	var err error
	ring.Do(func(value interface{}) {
		if value == nil || err != nil {
			return
		}

		err = writeToStore(store, "google.com", value, prefix+strconv.Itoa(i))
		i++
	})
	if err != nil {
		log.Printf("Error writing result: %v\n", err)
	}
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

func loadSamples(store *Store, samples *Samples) error {
	for i := 0; i < NLast3Hours; i++ {
		val, err := getFromStore(store, "google.com", "last3Hours"+strconv.Itoa(i))
		if err != nil {
			return err
		}
		if val == nil {
			return nil
		}
		samples.Last3Hours.Value = val
		samples.Last3Hours = samples.Last3Hours.Next()
	}
	return nil
}
