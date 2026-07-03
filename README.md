# hook-shift 🪝

> A lightning-fast, zero-dependency webhook dispatcher built in Go.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**hook-shift** is a highly concurrent, fault-tolerant webhook delivery engine. It acts as a resilient middleman between your application and external APIs.

Instead of writing custom retry loops, managing Goroutines, and risking frozen threads in your main application, hand the payload to `hook-shift`. It will instantly accept the job and process it asynchronously.

## ✨ Features

- **Zero Dependencies:** No heavy message brokers required. Runs entirely in memory out of the box.
- **Concurrent Worker Pool:** Dispatches hundreds of webhooks simultaneously without blocking.
- **Exponential Backoff:** Automatically retries failed deliveries (up to 5 attempts) using a time-staggered backoff matrix to prevent DDoS-ing destination servers.
- **Graceful Shutdown:** Catches `SIGTERM/SIGINT` and safely drains the worker pool and in-flight HTTP requests before shutting down—zero data loss during container restarts or server stops.
- **Convention Over Configuration:** Configurable entirely via simple environment variables and CLI flags.

---

## 🚀 Quickstart

You can clone and run `hook-shift` locally in under 10 seconds. You just need the Go toolchain installed.

```bash
# Clone the repository
git clone [https://github.com/mithileshgupta12/hook-shift.git](https://github.com/mithileshgupta12/hook-shift.git)
cd hook-shift

# Run the server directly
go run cmd/server/main.go --port=9000 --workers=10
```

---

## 🔌 API Usage

Once the server is running, simply `POST` your payload and destination to the dispatcher route.

**Request:**

```bash
curl -X POST http://localhost:9000/v1/dispatches \
  -H "Content-Type: application/json" \
  -d '{
    "destination_url": "[https://api.example.com/webhook](https://api.example.com/webhook)",
    "payload": {
      "event": "user.created",
      "user_id": "12345"
    }
  }'
```

**Response (202 Accepted):**

```json
{
  "message": "job accepted",
  "job_id": "9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d"
}
```

---

## 🏗️ Architecture

```text
[Your App] ---> (POST /v1/dispatches) ---> [hook-shift API]
                                                 |
                                         [In-Memory Queue]
                                                 |
                                     +-----------+-----------+
                                     |           |           |
                                 [Worker]    [Worker]    [Worker]
                                     |           |           |
                                     +-----------+-----------+
                                                 |
                                       [Destination Servers]
```

---

## 🗺️ Roadmap

The core in-memory concurrency engine is stable. The following persistence drivers and features are actively in development for enterprise-scale deployments:

- [ ] **PostgreSQL Driver:** Persistent, zero-data-loss queueing utilizing `FOR UPDATE SKIP LOCKED`.
- [ ] **Redis Streams Driver:** High-throughput cluster orchestration.
- [ ] **Embedded UI Dashboard:** A `go:embed` single-page dashboard to visualize pending, processing, and dead jobs.
- [ ] **Dead-Letter Disk Flush:** Automatic local logging for permanently failed webhooks.

---

## 🤝 Contributing

Pull requests are welcome! If you want to tackle one of the roadmap items, please open an issue first to discuss the architecture.

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
