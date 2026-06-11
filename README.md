# im_cacher 🚀


[![Go Version](img.shields.io/github/go-mod/go-version/sardor-innatov/im_cacher)](https://golang.org)
[![License](img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

**im_cacher** is a lightweight, high-performance in-memory key-value database and caching system built from scratch using **Go**. It is designed as a compact, self-contained alternative to Redis, supporting real-time data persistence to disk.

![im_cacher Banner](banner.png)

---

## ✨ Features

- **In-Memory Storage**: Lightning-fast key-value operations directly powered by RAM.
- **AOF Persistence**: Safe from data loss. All write operations are logged to an `.aof` (Append-Only File) and automatically replayed upon server restart to restore state.
- **Configuration Management**: Flexible settings for ports, memory limits, and synchronization modes using a dedicated `.conf` file.
- **Zero External Dependencies**: Pure Go implementation, ensuring a highly portable, secure, and easily buildable binary.

---

## 📂 Project Structure

```text
├── src/               # Core source code and business logic
├── .conf              # Server configuration parameters
├── .aof               # Append-only file for database recovery (auto-generated)
├── go.mod             # Go module dependencies
└── main.go            # Application entry point
```

---

## 🛠️ Quick Start

### Prerequisites
- [Go Compiler](https://golang.org) version **1.18** or higher.

### 1. Clone the Repository
```bash
git clone https://github.com
cd im_cacher
```

### 2. Configure the Server
Before running, you can modify the database parameters inside the `.conf` file:
```ini
# ==============================================================================
# IM_CACHER CONFIGURATION FILE
# ==============================================================================

# Network address and port where the cache server will listen for connections.
# Format: IP:PORT (e.g., 127.0.0.1:9990)
Addr:127.0.0.1:9990

# Time interval (in seconds) for the passive background loop.
# This background worker runs every X seconds to find and delete expired TTL keys.
EvictionLoopInterval:5

# The maximum number of key-value pairs allowed in the cache.
# Active Eviction: If a user adds a new key when this limit is reached, 
# the Least Recently Used (LRU) key is immediately deleted to make room.
MaxKeysAmount:200

# The initial capacity allocated for the internal map on startup.
# Pre-allocating memory improves performance and reduces resizing overhead.
InitialSize:50

```
### 3. Build and Run
You can execute the code directly or compile it into a production-ready binary:

```bash
# Run directly using Go CLI
go run main.go

# Compile into an executable binary
go build -o im_cacher main.go
./im_cacher
```

---

## ⚡ Architecture Blueprint

1. **Write Operations**: When a client sends a write request, the engine updates its concurrent-safe in-memory hash map and simultaneously appends the query to the `.aof` disk log.
2. **Crash Recovery**: If the server restarts unexpectedly, `im_cacher` parses the `.aof` file line by line, re-executing all recorded historical write actions to cleanly rebuild the database state.

---

## 🤝 Contributing

Contributions make the open-source community an amazing place! Any pull requests you make are highly appreciated:
1. Fork the Project.
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`).
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`).
4. Push to the Branch (`git push origin feature/AmazingFeature`).
5. Open a Pull Request.

---

## 📄 License

This project is licensed under the **MIT License**. See the [LICENSE](LICENSE) file for more details.

---
**Maintained by:** [@sardor-innatov](https://github.com/sardor-innatov)
