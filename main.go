package main

import (
	"log"
	"os/exec"
	"strings"
)

func main() {
	target := "google.com"
	count := "10"
	out, err := exec.Command("ping", "-c", count, target).Output()
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(out), "\n")
	summary := lines[len(lines)-2]
	results := strings.Split(strings.TrimRight(strings.TrimSpace(strings.Split(summary, "=")[1]), " ms"), "/")
	min := results[0]
	avg := results[1]
	max := results[2]
	mdev := results[3]
	log.Printf("Min: %s\n", min)
	log.Printf("Avg: %s\n", avg)
	log.Printf("Max: %s\n", max)
	log.Printf("Mdev: %s\n", mdev)
}
