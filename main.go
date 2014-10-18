package main

import (
	"log"
)

func main() {
	res, err := ping("google.com", 5)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Min: %f ms\n", res.Min)
	log.Printf("Avg: %f ms\n", res.Avg)
	log.Printf("Max: %f ms\n", res.Max)
	log.Printf("Mdev: %f ms\n", res.Mdev)
}
