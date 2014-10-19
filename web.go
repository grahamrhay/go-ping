package main

import (
	"container/ring"
	"fmt"
	"log"
	"net/http"
)

func startServer(ring **ring.Ring) {
	http.HandleFunc("/", makeHandler(ring))
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
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
