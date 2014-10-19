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
	current := &ring

	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for _ = range ticker.C {
			log.Println("ping google.com -c 5")
			res, err := ping("google.com", 5)
			if err != nil {
				log.Printf("Error from ping: %v\n", err)
			} else {
				ring.Value = res
				ring = ring.Next()
			}
		}
	}()

	startServer(current)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Received signal: %v\n", <-ch)
	log.Println("Shutting down")
	ticker.Stop()
}
