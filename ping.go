package main

import (
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type PingResult struct {
	Time       time.Time
	Min        float64
	Avg        float64
	Max        float64
	Mdev       float64
	PacketLoss int64
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

	packetSummary := lines[len(lines)-3]
	result.PacketLoss, _ = strconv.ParseInt(strings.TrimRight(strings.TrimSpace(strings.Split(packetSummary, ",")[2]), "% packet loss"), 10, 64)

	rttSummary := lines[len(lines)-2]
	results := strings.Split(strings.TrimRight(strings.TrimSpace(strings.Split(rttSummary, "=")[1]), " ms"), "/")
	result.Min, _ = strconv.ParseFloat(results[0], 64)
	result.Avg, _ = strconv.ParseFloat(results[1], 64)
	result.Max, _ = strconv.ParseFloat(results[2], 64)
	result.Mdev, _ = strconv.ParseFloat(results[3], 64)
	log.Printf("Min: %vms, Avg: %vms, Max: %vms, Mdev: %vms, Packet loss: %%%v\n", result.Min, result.Avg, result.Max, result.Mdev, result.PacketLoss)
	return result
}
