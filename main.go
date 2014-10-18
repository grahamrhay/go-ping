package main

import (
	"container/ring"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ring := ring.New(10)

	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for _ = range ticker.C {
			log.Println("ping google.com -c 5")
			res, err := ping("google.com", 5)
			if err != nil {
				log.Printf("Error from ping: %v\n", err)
			} else {
				ring.Value = res
				ring.Next()
				log.Printf("Time: %v\n", res.Time)
				log.Printf("Min: %f ms\n", res.Min)
				log.Printf("Avg: %f ms\n", res.Avg)
				log.Printf("Max: %f ms\n", res.Max)
				log.Printf("Mdev: %f ms\n", res.Mdev)
			}
		}
	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Received signal: %v\n", <-ch)
	log.Println("Shutting down")
	ticker.Stop()
}
