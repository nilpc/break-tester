# 💥 Break Tester

> **A blazing-fast, cache-busting HTTP stress tester with real-time Grafana dashboards — written entirely in Go.**

Ever wonder what happens to your API when 2,000 users show up at once? Break Tester finds out, so your users don't have to.

It uses Go's lightweight goroutines and optimized channels to simulate massive concurrent traffic across multiple endpoints simultaneously — and unlike most load testers, it ships with a **Cache Busting Engine** that punches straight through CDN edge caches (Cloudflare, Fastly, etc.) to hammer your actual origin server.

Every request is tracked live by Prometheus and served up on a pre-built Grafana dashboard — think of it as your own personal mission control, refreshing every 5 seconds.

---

## What Makes It Different

**It's a two-headed beast.** On one side, you get the raw horsepower — hundreds of concurrent goroutines slamming your endpoints with cache-busting entropy strings so every request hits your origin server fresh. On the other side, every single result is broadcast to Prometheus in real time and rendered instantly on a beautiful Grafana dashboard with latency heatmaps, RPS counters, status breakdowns, and live percentile tracking.

**The old way:** fire a load test, wait for it to finish, read a wall of text.

**The Break Tester way:** watch the carnage unfold live on a screen full of colourful panels while your backend screams for mercy.

---

## How It Works

Think of it like running a secret shopper operation at scale, with a live TV crew broadcasting the results.

**The Mission Briefing** (`config.json`)**:** You define the targets, traffic weights, and duration. Want 70% of traffic hitting your homepage and 30% slamming the checkout endpoint? Done.

**The Disguise Kit (Cache Bypass):** The engine appends a dynamic query string to every request so edge caches treat each one as unique — no cached hits, no shortcuts, pure origin load.

**The Shopper Army (Goroutines):** Hundreds of lightweight workers spin up in the background and receive their tasks through memory-safe Go channels — fast, efficient, and completely non-blocking.

**The Radio Tower (Prometheus + Grafana):** Every shopper carries a walkie-talkie. The moment a request finishes, the status code and duration get radioed back to Prometheus. Grafana picks up the signal every 2 seconds and paints you a beautiful dashboard with every metric that matters.

---

## Project Structure

```
break-tester/
├── .github/
│   └── workflows/
│       └── deploy.yml              # CI/CD — build & push to Docker Hub
├── grafana/
│   └── provisioning/
│       ├── datasources/
│       │   └── datasource.yml      # Auto-connects to Prometheus
│       └── dashboards/
│           ├── dashboard.yml        # Dashboard provider config
│           └── break-tester-dashboard.json   # 14-panel live dashboard
├── main.go                          # Entry point, orchestration, Prometheus metrics
├── config.json                      # Default test configuration
├── app.py                           # Local Flask target server
├── prometheus.yml                   # Scrape config (every 2 seconds)
├── Dockerfile                       # Multi-stage production build
├── docker-compose.yml               # Full stack orchestration
├── go.mod / go.sum                  # Go module dependencies
└── README.md
```

---

## Quick Start

You'll need **Go**, **Docker**, and **Docker Compose** installed before you begin.

### 1. Fire Up the Full Stack

One command brings up the target app, Prometheus, Grafana, and Break Tester all at once:

```bash
docker compose up --build
```

This spins up four services:
- **target-app** — a simple Flask server on port 5000 (your test dummy)
- **prometheus** — scrapes metrics from Break Tester every 2 seconds on port 9090
- **grafana** — pre-provisioned dashboard waiting at port 3000 (login: `admin` / `secret`)
- **break-tester** — the stress cannon itself, broadcasting metrics on port 8080

### 2. Open Mission Control

Navigate to [http://localhost:3000](http://localhost:3000) and log in with **admin** / **secret**. The Break Tester dashboard loads automatically — every panel starts populating within seconds.

### 3. Watch the Show

The Grafana dashboard refreshes every 5 seconds, giving you live visibility into total requests, success rate, error rate, current RPS, average latency, p99 peak latency, throughput, and more — all updating in real time as the test runs.

---

## Running Standalone

If you just want the load tester without the monitoring stack:

```bash
# Build the image
docker build -t break-tester:latest .

# Run it (without Grafana)
docker run break-tester:latest
```

The terminal updates live every 2 seconds with real-time stats — even without Grafana:

```
📊 Live | Reqs: 4521 | Success: 99.3% | RPS: 225 | Avg Latency: 18ms
```

Break Tester reads `config.json` by default. To use a custom config file:

```bash
docker run -v /path/to/your/config.json:/app/config.json break-tester:latest
```

Or with Go directly:

```bash
go run main.go --config=massive_test.json
```

---

## Running Locally (Without Docker)

If you want to run the load tester and target app directly on your machine (no Docker), you'll need **two terminals** — one for the Flask target server and one for the Go tester.

### 1. Start the Flask Target Server

```bash
pip install flask
python app.py
```

The server starts on `http://localhost:5000`.

### 2. Update `config.json`

Change the target URLs from the Docker hostname (`target-app`) to `localhost`:

```json
{
  "concurrency": 100,
  "duration": "10s",
  "targets": [
    {
      "url": "http://localhost:5000/users",
      "method": "GET",
      "weight": 70,
      "cache_bypass": true
    },
    {
      "url": "http://localhost:5000/checkout",
      "method": "POST",
      "weight": 30,
      "cache_bypass": false
    }
  ]
}
```

The only change is `target-app:5000` → `127.0.0.1:5000` (or `localhost:5000`).

### 3. Run the Break Tester

In a second terminal, start the Go load tester:

```bash
go run main.go
```

It will start hitting your Flask server and expose Prometheus metrics at `http://localhost:8080/metrics`. The terminal shows live stats every 2 seconds:

```
📊 Live | Reqs: 4521 | Success: 99.3% | RPS: 225 | Avg Latency: 18ms
```

### 4. (Optional) View Metrics in Grafana

To see the live dashboard without Docker, run Prometheus and Grafana separately or use the Docker Compose services for just the monitoring stack. Alternatively, hit `http://localhost:8080/metrics` directly to see the raw metric output.

---

## The Grafana Dashboard

Grafana is the **visualization layer** of the monitoring stack. It queries Prometheus every 5 seconds and renders every metric from the load test on a live dashboard — no need to wait for the test to finish.

### How Grafana Is Configured

All Grafana operations are handled automatically via **provisioning files** in `grafana/provisioning/` — zero manual setup required:

| Operation | What It Does | File |
|---|---|---|
| **Datasource provisioning** | Registers Prometheus as the default data source, pointing at `http://prometheus:9090` | `grafana/provisioning/datasources/datasource.yml` |
| **Dashboard provider setup** | Tells Grafana to load dashboard JSON files from the provisioning directory on startup | `grafana/provisioning/dashboards/dashboard.yml` |
| **Dashboard loading** | Imports the 15-panel "Break Tester" dashboard with auto-refresh every 5 seconds | `grafana/provisioning/dashboards/break-tester-dashboard.json` |
| **Docker mount** | Binds `./grafana/provisioning/` into the Grafana container at `/etc/grafana/provisioning/` | `docker-compose.yml` |

On container startup, Grafana automatically reads these files, creates the Prometheus datasource, and imports the dashboard — it's ready the moment you log in at port 3000.

### Dashboard Panels

The dashboard ships with **15 panels** across 4 row groupings:

| Row | Panels |
|---|---|
| **Key Stats** | Total Requests, Success Rate, Current RPS, Avg Latency, Error Rate, Peak Latency (p99), Throughput, Active Users (via Little's Law) |
| **Traffic** | Request Rate by Status (time series), Request Rate by Method (time series) |
| **Latency** | Latency Percentiles (p50 / p90 / p99), Latency Distribution (heatmap) |
| **Breakdown** | Total Requests by Status (bar gauge), Requests by Method (pie chart), Status Breakdown by Method (table) |

All panels query Prometheus using PromQL — metrics like `break_tester_requests_total` and `break_tester_request_duration_seconds` are scraped by Prometheus every 2 seconds and visualized instantly in Grafana.

---

## Docker Compose Breakdown

The `docker-compose.yml` orchestrates everything:

```yaml
services:
  target-app:      # Your test subject (Flask, port 5000)
  prometheus:      # Metrics scraper (port 9090)
  grafana:         # Live dashboards (port 3000, admin/secret)
  break-tester:    # The cannon itself (port 8080)
```

Grafana comes with auto-provisioning configured under `grafana/provisioning/` — both the Prometheus datasource and the dashboard JSON are loaded automatically on startup. No clicking around in the UI.

---

## Configuration

Break Tester is driven entirely by a `config.json` file. Here's what a basic setup looks like:

```json
{
  "concurrency": 500,
  "duration": "60s",
  "targets": [
    { "url": "http://target-app:5000/", "weight": 70 },
    { "url": "http://target-app:5000/data", "weight": 30 }
  ]
}
```

| Field | Description |
|---|---|
| `concurrency` | Number of concurrent goroutines (shoppers) |
| `duration` | How long to run the test |
| `targets` | List of URLs with traffic weight percentages |

---

## Output

Break Tester prints live stats to the terminal every 2 seconds throughout the test:

```
📊 Live | Reqs: 4521 | Success: 99.3% | RPS: 225 | Avg Latency: 18ms
```

When the mission ends, a clean final summary is printed:

```
🎉 Mission Accomplished! All shoppers clocked out and returned home.
🔥 Total Orders Placed: 84,312
🎯 Successful Orders:   83,901 / 84,312
⏱️  Average Kitchen Speed: 18ms
🚀 Throughput (RPS):       1,405 requests/second
```

But the real magic is the Grafana dashboard — you don't have to wait for the test to finish. Every metric updates live, second by second, right in front of you.

---

## CI/CD — GitHub Actions

Every push to `main` automatically builds the Docker image and ships it straight to Docker Hub. No manual steps, no stale images.

The workflow lives at `.github/workflows/deploy.yml` and needs two repository secrets:

| Secret | What it is |
|---|---|
| `DOCKER_USERNAME` | Your Docker Hub username |
| `DOCKER_PASSWORD` | A Docker Hub access token |

Once set up, every merge to `main` triggers a fresh build and pushes `nilpc1999/break-tester:latest` to Docker Hub automatically:

```bash
docker pull nilpc1999/break-tester:latest
```

---

## Tech Stack

- **Go** — core runtime, goroutines, channels
- **Prometheus** — metrics collection and querying
- **Grafana** — pre-provisioned live dashboards
- **Python / Flask** — included dummy target server for local testing
- **Docker & Docker Compose** — full-stack orchestration
- **GitHub Actions** — CI/CD to Docker Hub

---

## When to Use This

- Pre-launch stress testing before a big traffic event
- Validating autoscaling policies on cloud infrastructure
- Finding the breaking point of a new endpoint before it goes to production
- Impressing your teammates with a live Grafana dashboard during load tests
- Confirming your CDN is actually caching (or not caching) what you think it is

---

## A Word of Warning

Only point this at servers you own or have explicit permission to test. Firing this at someone else's infrastructure without authorization is illegal in most jurisdictions and a generally bad idea.

Stay ethical, stress responsibly. 🤝

---

## License

MIT — do whatever you want with it, just don't blame us when you find your server's limits.
