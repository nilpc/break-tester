package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync" 
	"time"
)
type Target struct {
	URL         string `json:"url"`
	Method      string `json:"method"`
	Weight      int    `json:"weight"`
	CacheBypass bool   `json:"cache_bypass"` // The Disguise Kit
}

type Config struct {
	Concurrency int      `json:"concurrency"`
	Duration    string   `json:"duration"`
	Targets     []Target `json:"targets"`
}

type Result struct {
	StatusCode int
	Duration   time.Duration
}

func secretShopper(shopperID int, assignmentBelt <-chan Target, reportBelt chan<- Result, wg *sync.WaitGroup) {

	defer wg.Done()
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for target := range assignmentBelt {
		url := target.URL
		if target.CacheBypass {
			fakeMustache := strconv.Itoa(rand.Intn(999999999))
			url = url + "?disguise=" + fakeMustache
		}

		start := time.Now()
		req, err := http.NewRequest(target.Method, url, nil)
		if err != nil {
			reportBelt <- Result{StatusCode: 0, Duration: time.Since(start)}
			continue
		}

		resp, err := client.Do(req)
		
		duration := time.Since(start)

		if err != nil {
			reportBelt <- Result{StatusCode: 0, Duration: duration}
		} else {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
			reportBelt <- Result{
				StatusCode: resp.StatusCode,
				Duration:   duration,
			}
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	configFile := flag.String("config", "config.json", "Path to the mission briefing JSON file")
	flag.Parse()

	fmt.Println("💥 Welcome to Break Tester!\n")
	data, err := os.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("❌ Could not read config file (%s): %v", *configFile, err)
	}

	var briefing Config
	err = json.Unmarshal(data, &briefing)
	if err != nil {
		log.Fatalf("❌ The boss wrote a bad JSON file: %v", err)
	}
	missionTime, err := time.ParseDuration(briefing.Duration)
	if err != nil {
		log.Fatalf("❌ Invalid duration format: %v", err)
	}

	fmt.Printf("✅ Mission Briefing Loaded!\n")
	fmt.Printf("👥 Army Size: %d concurrent shoppers\n", briefing.Concurrency)
	fmt.Printf("⏳ Mission Timer: %s\n\n", missionTime)

	assignmentBelt := make(chan Target, 100)
	reportBelt := make(chan Result, 100)

	var wg sync.WaitGroup

	fmt.Printf("🚀 Deploying %d secret shoppers into the background...\n", briefing.Concurrency)
	for i := 1; i <= briefing.Concurrency; i++ {
		wg.Add(1) // Add 1 shopper to the clipboard
		go secretShopper(i, assignmentBelt, reportBelt, &wg)
	}

	fmt.Printf("\n⏱️  Starting the countdown timer NOW! Firing HTTP Cannon...\n")
	
	go func() {
		timer := time.NewTimer(missionTime)
		
		for {
			roll := rand.Intn(100) + 1 
			var selectedTarget Target
			if roll <= briefing.Targets[0].Weight {
				selectedTarget = briefing.Targets[0]
			} else {
				selectedTarget = briefing.Targets[1]
			}

			select {
			case <-timer.C:
				close(assignmentBelt) 
				return
			case assignmentBelt <- selectedTarget:
			}
		}
	}()

	go func() {
		wg.Wait()
		close(reportBelt)
	}()

	successCount := 0
	totalRequests := 0
	var totalTime time.Duration

	for res := range reportBelt {
		totalRequests++
		if res.StatusCode >= 200 && res.StatusCode < 300 {
			successCount++
		}
		totalTime += res.Duration
	}

	fmt.Println("\n=======================================================")
	fmt.Printf("🎉 Mission Accomplished! All shoppers returned.\n")
	fmt.Printf("🔥 Total Requests Sent: %d\n", totalRequests)
	fmt.Printf("🎯 Successful Orders:   %d / %d\n", successCount, totalRequests)
	if totalRequests > 0 {
		fmt.Printf("⏱️  Average Speed:       %v\n", totalTime/time.Duration(totalRequests))
		rps := float64(totalRequests) / missionTime.Seconds()
		fmt.Printf("🚀 Throughput (RPS):    %.2f requests/second\n", rps)
	}
	fmt.Println("=======================================================")
}