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
	ring := ring.New(36)
	current := &ring

	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for _ = range ticker.C {
			res := ping("google.com", 20)
			ring.Value = res
			ring = ring.Next()
		}
	}()

	startServer(current)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Received signal: %v\n", <-ch)
	log.Println("Shutting down")
	ticker.Stop()
}
