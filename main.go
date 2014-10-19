package main

import (
	"container/ring"
	"fmt"
	"log"
	"net/http"
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
				log.Printf("Time: %v\n", res.Time)
				log.Printf("Min: %f ms\n", res.Min)
				log.Printf("Avg: %f ms\n", res.Avg)
				log.Printf("Max: %f ms\n", res.Max)
				log.Printf("Mdev: %f ms\n", res.Mdev)
			}
		}
	}()

	http.HandleFunc("/", makeHandler(current))
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Received signal: %v\n", <-ch)
	log.Println("Shutting down")
	ticker.Stop()
}

func makeHandler(ring **ring.Ring) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ring.Do(func(value interface{}) {
			if value == nil {
				return
			}

			res := value.(*PingResult)
			fmt.Fprintf(w, "Time: %v\n", res.Time)
			fmt.Fprintf(w, "Min: %f ms\n", res.Min)
			fmt.Fprintf(w, "Avg: %f ms\n", res.Avg)
			fmt.Fprintf(w, "Max: %f ms\n", res.Max)
			fmt.Fprintf(w, "Mdev: %f ms\n", res.Mdev)
		})
	}
}
