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
	"sync" // The Manager's Clipboard
	"sync/atomic"
	"time"

	// The Prometheus Broadcaster Libraries (The Radio Tower)
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Target is a specific task we want our secret shoppers to do (e.g., ordering a burger).
type Target struct {
	URL         string `json:"url"`
	Method      string `json:"method"`
	Weight      int    `json:"weight"`
	CacheBypass bool   `json:"cache_bypass"` // The Disguise Kit
}

// Config is our "Mission Briefing" written by the boss.
type Config struct {
	Concurrency int      `json:"concurrency"` // How many shoppers to hire
	Duration    string   `json:"duration"`    // How long the shift lasts
	Targets     []Target `json:"targets"`     // The menu items to order
}

// Result is the physical receipt a single shopper brings back after a task.
type Result struct {
	StatusCode int
	Duration   time.Duration
}

// ==========================================
// 📡 THE BROADCASTER (Prometheus Metrics)
// ==========================================
var (
	// The Ledger: A giant notebook where we tally every single order placed (Success or Fail)
	orderLedger = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "break_tester_requests_total",
			Help: "Total number of HTTP requests sent by the secret shoppers.",
		},
		[]string{"status", "method"}, // We categorize the tally by the receipt status (e.g., 200 OK)
	)

	// The Stopwatch: We use this to measure exactly how fast the kitchen is cooking the food
	kitchenStopwatch = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "break_tester_request_duration_seconds",
			Help:    "Time taken for the kitchen to serve the food in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"status", "method"},
	)
)

// init() runs automatically to hand the Ledger and Stopwatch to Prometheus before the shift starts.
func init() {
	prometheus.MustRegister(orderLedger)
	prometheus.MustRegister(kitchenStopwatch)
}

// secretShopper is the actual human worker doing the heavy lifting.
func secretShopper(shopperID int, assignmentBelt <-chan Target, reportBelt chan<- Result, managerClipboard *sync.WaitGroup) {
	// CRITICAL: When the shift ends and this shopper goes home, they MUST cross their name off the clipboard.
	defer managerClipboard.Done()

	// Hand the shopper a standard web browser. If the kitchen takes >10 seconds, they walk away in disgust.
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// The shopper stands at the conveyor belt, grabbing assignments continuously as they roll down.
	for target := range assignmentBelt {
		url := target.URL

		// 🕵️ THE DISGUISE KIT (Cache Buster)
		// Put on a fake mustache (a random number) so the Bouncer (CDN) thinks we are a brand-new customer.
		if target.CacheBypass {
			fakeMustache := strconv.Itoa(rand.Intn(999999999))
			url = url + "?disguise=" + fakeMustache
		}

		// Click the stopwatch! The shopper is walking up to the cash register.
		start := time.Now()

		req, err := http.NewRequest(target.Method, url, nil)
		if err != nil {
			// Tripped on the way to the register. Report a failure.
			reportBelt <- Result{StatusCode: 0, Duration: time.Since(start)}
			continue
		}

		// 💥 Place the order! (This sends real traffic over the internet)
		resp, err := client.Do(req)

		// Stop the stopwatch the exact millisecond the food (response) arrives.
		duration := time.Since(start)

		// 📡 BROADCAST TO GRAFANA
		// The shopper grabs their walkie-talkie and radios the results back to the Prometheus tower
		receiptStatus := "error"
		if err == nil {
			receiptStatus = strconv.Itoa(resp.StatusCode)
		}

		// Tally the order in the Ledger and record the speed on the Stopwatch
		orderLedger.WithLabelValues(receiptStatus, target.Method).Inc()
		kitchenStopwatch.WithLabelValues(receiptStatus, target.Method).Observe(duration.Seconds())

		if err != nil {
			// The kitchen caught on fire, or the internet went down.
			reportBelt <- Result{StatusCode: 0, Duration: duration}
		} else {
			// Throw the food directly in the trash (io.Discard) so the shopper can instantly run back to the line.
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()

			// Drop the receipt onto the reporting conveyor belt for the boss.
			reportBelt <- Result{
				StatusCode: resp.StatusCode,
				Duration:   duration,
			}
		}
	}
}

func main() {
	// Seed the random number generator so our Target Roulette works unpredictably
	rand.Seed(time.Now().UnixNano())

	// The Boss passes the Mission Briefing file via the command line
	configFile := flag.String("config", "config.json", "Path to the mission briefing JSON file")
	flag.Parse()

	fmt.Println("💥 Welcome to Break Tester (Grafana Edition)!\n")

	// 1. Read the Mission Briefing
	data, err := os.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("❌ Could not read the Boss's briefing file (%s): %v", *configFile, err)
	}

	var briefing Config
	err = json.Unmarshal(data, &briefing)
	if err != nil {
		log.Fatalf("❌ The boss wrote a bad JSON file: %v", err)
	}

	// 2. Set the official Mission Timer
	missionTime, err := time.ParseDuration(briefing.Duration)
	if err != nil {
		log.Fatalf("❌ Invalid duration format: %v", err)
	}

	// ==========================================
	// 📡 START THE RADIO TOWER
	// ==========================================
	// We set up a live broadcast station on port 8080 so Prometheus can listen in.
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
	fmt.Println("📡 Prometheus Radio Tower live on port 8080")

	// Step 1: Fire up the physical Conveyor Belts
	assignmentBelt := make(chan Target, 100)
	reportBelt := make(chan Result, 100)

	// Step 2: Grab a blank Clipboard for the Manager
	var managerClipboard sync.WaitGroup

	fmt.Printf("🚀 Deploying %d secret shoppers into the background...\n", briefing.Concurrency)
	for i := 1; i <= briefing.Concurrency; i++ {
		managerClipboard.Add(1) // Write the shopper's name on the clipboard
		go secretShopper(i, assignmentBelt, reportBelt, &managerClipboard)
	}

	// Step 3: The Dispatcher & The Master Stopwatch
	fmt.Printf("\n⏱️  Starting the countdown timer NOW! Firing HTTP Cannon...\n")
	startTime := time.Now()

	go func() {
		// Set the alarm clock
		timer := time.NewTimer(missionTime)

		for {
			// 🎲 THE TARGET ROULETTE (Weighted routing based on the Briefing)
			roll := rand.Intn(100) + 1
			var selectedTarget Target

			if roll <= briefing.Targets[0].Weight {
				selectedTarget = briefing.Targets[0]
			} else {
				selectedTarget = briefing.Targets[1]
			}

			// The 'select' statement allows the Dispatcher to do two things at once!
			select {
			case <-timer.C:
				// DING! The alarm went off. Yell "No more orders!" and shut off the belt.
				close(assignmentBelt)
				return
			case assignmentBelt <- selectedTarget:
				// The alarm hasn't gone off yet, so keep throwing tasks onto the belt!
			}
		}
	}()

	// Step 3.5: The Floor Manager
	// This manager sits in the background staring at the clipboard. 
	// When all names are crossed off (Wait()), they close the receipt belt.
	go func() {
		managerClipboard.Wait()
		close(reportBelt)
	}()

	// Step 4: The Scorekeeper & Live Reporter
	var (
		liveTotal     int64
		liveSuccess   int64
		liveLatencyNs int64
	)
	done := make(chan struct{})

	// Live terminal reporter — prints stats every 2 seconds
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				total := atomic.LoadInt64(&liveTotal)
				success := atomic.LoadInt64(&liveSuccess)
				latencyNs := atomic.LoadInt64(&liveLatencyNs)
				var successRate, avgLatency, rps float64
				if total > 0 {
					successRate = float64(success) / float64(total) * 100
					avgLatency = float64(latencyNs) / float64(total) / 1e6
					elapsed := time.Since(startTime).Seconds()
					if elapsed > 0 {
						rps = float64(total) / elapsed
					}
				}
				fmt.Printf("\r📊 Live | Reqs: %d | Success: %.1f%% | RPS: %.0f | Avg Latency: %.0fms   ",
					total, successRate, rps, avgLatency)
			case <-done:
				fmt.Println()
				return
			}
		}
	}()

	// Read receipts off the belt until the Floor Manager closes it
	for receipt := range reportBelt {
		atomic.AddInt64(&liveTotal, 1)
		if receipt.StatusCode >= 200 && receipt.StatusCode < 300 {
			atomic.AddInt64(&liveSuccess, 1)
		}
		atomic.AddInt64(&liveLatencyNs, receipt.Duration.Nanoseconds())
	}
	close(done)

	totalOrdersPlaced := int(atomic.LoadInt64(&liveTotal))
	successfulOrders := int(atomic.LoadInt64(&liveSuccess))
	totalKitchenTime := time.Duration(atomic.LoadInt64(&liveLatencyNs))

	// Print the final scoreboard for the Boss!
	fmt.Println("\n=======================================================")
	fmt.Printf("🎉 Mission Accomplished! All shoppers clocked out and returned home.\n")
	fmt.Printf("🔥 Total Orders Placed: %d\n", totalOrdersPlaced)
	fmt.Printf("🎯 Successful Orders:   %d / %d\n", successfulOrders, totalOrdersPlaced)
	if totalOrdersPlaced > 0 {
		fmt.Printf("⏱️  Average Kitchen Speed: %v\n", totalKitchenTime/time.Duration(totalOrdersPlaced))
		rps := float64(totalOrdersPlaced) / missionTime.Seconds()
		fmt.Printf("🚀 Throughput (RPS):       %.2f requests/second\n", rps)
	}
	fmt.Println("=======================================================")
}