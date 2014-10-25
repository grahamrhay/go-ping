package main

import (
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type PingResult struct {
	Time time.Time
	Min  float64
	Avg  float64
	Max  float64
	Mdev float64
}

func ping(target string, count int) *PingResult {
	result := &PingResult{Time: time.Now()}
	log.Println("ping -c 20 google.com")
	out, err := exec.Command("ping", "-c", strconv.Itoa(count), target).Output()
	if err != nil {
		log.Printf("Error from ping: %v\n", err)
		return result
	}
	lines := strings.Split(string(out), "\n")
	summary := lines[len(lines)-2]
	results := strings.Split(strings.TrimRight(strings.TrimSpace(strings.Split(summary, "=")[1]), " ms"), "/")
	result.Min, _ = strconv.ParseFloat(results[0], 64)
	result.Avg, _ = strconv.ParseFloat(results[1], 64)
	result.Max, _ = strconv.ParseFloat(results[2], 64)
	result.Mdev, _ = strconv.ParseFloat(results[3], 64)
	return result
}
