package main

import (
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

func ping(target string, count int) (*PingResult, error) {
	out, err := exec.Command("ping", "-c", strconv.Itoa(count), target).Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	summary := lines[len(lines)-2]
	results := strings.Split(strings.TrimRight(strings.TrimSpace(strings.Split(summary, "=")[1]), " ms"), "/")
	min, _ := strconv.ParseFloat(results[0], 64)
	avg, _ := strconv.ParseFloat(results[1], 64)
	max, _ := strconv.ParseFloat(results[2], 64)
	mdev, _ := strconv.ParseFloat(results[3], 64)
	result := &PingResult{Time: time.Now(), Min: min, Avg: avg, Max: max, Mdev: mdev}
	return result, nil
}
