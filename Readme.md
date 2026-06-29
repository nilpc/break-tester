# 💥 Break Tester

> **A blazing-fast, cache-busting HTTP stress tester — written entirely in Go.**

Ever wonder what happens to your API when 2,000 users show up at once? Break Tester finds out, so your users don't have to.

It uses Go's lightweight goroutines and optimized channels to simulate massive concurrent traffic across multiple endpoints simultaneously — and unlike most load testers, it ships with a **Cache Busting Engine** that punches straight through CDN edge caches (Cloudflare, Fastly, etc.) to hammer your actual origin server.

---

## What Makes It Different

Most load testers fire repeated requests and call it a day. The problem? CDNs are really good at their job — they'll happily serve cached responses to your "load test" without your origin server ever breaking a sweat.

Break Tester solves this by attaching a unique entropy string to every outgoing request (`?_bust=a7f3c1...`), so every single request looks fresh to the cache layer. You're testing the thing that actually matters: **your backend under real pressure**.

---

## How It Works

Think of it like running a secret shopper operation at scale.

**The Mission Briefing** (`config.json`)**:** You define the targets, traffic weights, and duration. Want 70% of traffic hitting your homepage and 30% slamming the checkout endpoint? Done.

**The Disguise Kit (Cache Bypass):** The engine appends a dynamic query string to every request so edge caches treat each one as unique — no cached hits, no shortcuts, pure origin load.

**The Shopper Army (Goroutines):** Hundreds of lightweight workers spin up in the background and receive their tasks through memory-safe Go channels — fast, efficient, and completely non-blocking.

---

## Quick Start

You'll need **Go** and **Python** installed before you begin.

### 1. Initialize the Project

```bash
go mod init break-tester
```

### 2. Start the Target Server

We've included a simple Flask app so you can test locally without pointing the cannon at anything real.

```bash
pip install flask
python app.py
```

### 3. Fire the Cannon

Open a second terminal and run:

```bash
go run main.go
```

Break Tester reads `config.json` by default. To use a custom config file, pass it as a flag:

```bash
go run main.go --config=massive_test.json
```

---

## Docker

The project uses a multi-stage Docker build to compile everything down into a ~15MB image — lean enough to deploy across hundreds of nodes on ECS, Fly.io, or DigitalOcean.

```bash
# Build the image
docker build -t break-tester:latest .

# Run it locally
docker run break-tester:latest
```

---

## CI/CD — GitHub Actions

Every push to `main` automatically builds the Docker image and ships it straight to Docker Hub. No manual steps, no stale images.

```yaml
name: Build and Deploy Load Tester
on:
  push:
    branches:
      - main

jobs:
  build-and-push:
    runs-on: ubuntu-latest

    steps:
      - name: Check out code
        uses: actions/checkout@v3

      - name: Login to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}

      - name: Build and push Docker image
        uses: docker/build-push-action@v4
        with:
          context: .
          push: true
          tags: nilpc1999/break-tester:latest
```

The workflow lives at `.github/workflows/deploy.yml` and needs two repository secrets to work:

| Secret | What it is |
|---|---|
| `DOCKER_USERNAME` | Your Docker Hub username |
| `DOCKER_PASSWORD` | A Docker Hub access token (not your account password) |

To add them: go to your repo → **Settings** → **Secrets and variables** → **Actions** → **New repository secret**.

Once set up, every merge to `main` triggers a fresh build and pushes `nilpc1999/break-tester:latest` to Docker Hub automatically. Pull it anywhere with:

```bash
docker pull nilpc1999/break-tester:latest
```

---

## Configuration

Break Tester is driven entirely by a `config.json` file. Here's what a basic setup looks like:

```json
{
  "concurrency": 500,
  "duration": "60s",
  "cache_busting": true,
  "targets": [
    { "url": "https://yoursite.com/", "weight": 70 },
    { "url": "https://yoursite.com/checkout", "weight": 30 }
  ]
}
```

| Field | Description |
|---|---|
| `concurrency` | Number of concurrent goroutines (shoppers) |
| `duration` | How long to run the test |
| `cache_busting` | Append unique query strings to bypass CDN caches |
| `targets` | List of URLs with traffic weight percentages |

---

## Output

At the end of each run, Break Tester prints a clean summary:

```
Duration:        60s
Total Requests:  84,312
Success (2xx):   83,901 (99.5%)
Errors:          411 (0.5%)
Avg Latency:     18ms
P95 Latency:     94ms
P99 Latency:     201ms
Throughput:      1,405 req/sec
```

---

## Tech Stack

- **Go** — core runtime, goroutines, channels
- **Python / Flask** — included dummy target server for local testing
- **Docker** — multi-stage build for lean, portable deployment
- **GitHub Actions** — automated build and push to Docker Hub on every merge to `main`

---

## Project Structure

```
break-tester/
├── .github/
│   └── workflows/
│       └── deploy.yml   # CI/CD — build & push to Docker Hub on push to main
├── main.go              # Entry point, flag parsing, orchestration
├── config.json          # Default test configuration
├── app.py               # Local Flask target server
├── Dockerfile           # Multi-stage production build
└── README.md
```

---

## When to Use This

- Pre-launch stress testing before a big traffic event
- Validating autoscaling policies on cloud infrastructure
- Finding the breaking point of a new endpoint before it goes to production
- Confirming your CDN is actually caching (or not caching) what you think it is

---

## A Word of Warning

Only point this at servers you own or have explicit permission to test. Firing this at someone else's infrastructure without authorization is illegal in most jurisdictions and a generally bad idea.

Stay ethical, stress responsibly. 🤝

---

## License

MIT — do whatever you want with it, just don't blame us when you find your server's limits.
