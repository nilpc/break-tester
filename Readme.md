💥 Break Tester

A blazing-fast, cache-busting, flexible HTTP command-line stress tester written entirely in Go.

Break Tester uses lightweight Go routines ("The Secret Shopper Army") and heavily-optimized Go Channels ("Conveyor Belts") to simulate thousands of concurrent users hitting multiple target APIs simultaneously.

Unlike standard load testers, Break Tester includes a Cache Busting Engine to dynamically append entropy to outgoing URLs, ensuring your requests punch straight through CDN edge caches (like Cloudflare) and truly stress test the Origin Server.

🚀 How it Works (The Secret Shopper Analogy)

The Mission Briefing (config.json): You define how many "Secret Shoppers" you want, how long they should shop, and what URLs they should target. You can set traffic weights (e.g., 70% go to the homepage, 30% go to checkout).

The Disguise Kit (Cache Bypass): If enabled, the engine attaches a dynamic ?disguise=182736 query string to every single request so edge caches treat every request as unique.

The Hyper-Bus (Goroutines): The engine spawns hundreds of shoppers in the background and rapidly feeds them tasks via memory-safe Channels.

🛠️ Quick Start (Local Testing)

You will need Go and Python installed on your machine.

1. Initialize the Project

# Initialize the Go module if you haven't already
go mod init break-tester


2. Start the Target Application (The Kitchen)

We included a dummy Flask application to safely test your cannon against.

# Open a terminal and run the Python target server
pip install flask
python app.py


3. Fire the Cannon!

Open a second terminal window and run the Go application. It will automatically read the config.json file.

# Execute the load test!
go run main.go

# Or, pass a custom JSON file using the CLI flag:
go run main.go --config=massive_test.json


☁️ Cloud Deployment (Docker)

This project uses a Multi-Stage Docker Build to compile the Go code down into a microscopic (~15MB) container, ready to be deployed to thousands of nodes on AWS ECS or DigitalOcean.

# Build the featherweight image
docker build -t break-tester:latest .

# Run the image locally
docker run break-tester:latest
