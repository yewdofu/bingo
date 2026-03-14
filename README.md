# bingo

SMW (Super Mario World) bingo card generator — a Go HTTP server that creates randomized 5×5 bingo cards from a dataset of SMW game goals.

## Overview

This project is a backend service for generating Super Mario World bingo cards. Given an optional seed value, it deterministically selects 25 unique goals from a pool of 150+ SMW objectives (stored in `bingo.json`) and returns them as a JSON response.

Seeded generation ensures that the same seed always produces the same card, making it easy to share a specific card with others. Seeds can be numeric strings (parsed directly) or arbitrary strings (hashed via FNV-1a to produce a stable int64).

## Getting Started

### Prerequisites

- Go 1.25+
- (Optional) [air](https://github.com/air-verse/air) for hot-reload during development
- (Optional) Docker for containerized builds

### Running

```bash
go build -o ./tmp/main . && ./tmp/main
```

The server listens on port **8080**.

### Running with Docker

```bash
# Production image (distroless)
docker build --target builder -t bingo-builder .
docker build -t bingo .
docker run -p 8080:8080 bingo
```

## Development

```bash
# Hot-reload (requires air)
air

# Run all tests
go test ./...

# Run tests for the bingo package only
go test ./bingo

# Production build (statically linked)
CGO_ENABLED=0 GOOS=linux go build -o server
```

## API Reference

### GET /test

Health check endpoint.

**Response**

```
Hello, "/test"
```

---

### GET /create

Generates a bingo card. Returns 25 unique goals in a fixed array.

**Query Parameters**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `seed` | string | No | Seed for deterministic card generation. Numeric strings are used as-is; arbitrary strings are hashed via FNV-1a. If omitted, `time.Now().UnixNano()` is used. |

**Response**

```json
{
  "goals": [
    { "name": "6つの水中ステージをクリアする", "difficulty": 1 },
    ...
  ],
  "seed": "12345"
}
```

The `goals` array always contains exactly 25 elements. The `seed` field reflects the seed that was used (useful when no seed was provided).

**Error Response**

Returns HTTP `500` if bingo data is not initialized.

---

### GET /debug

Reinitializes the bingo data by reloading `bingo.json` from disk. Useful during development to pick up goal changes without restarting the server.

**Response**

`200 OK` (empty body)

---

## Architecture

The server is structured as two layers: `main.go` handles HTTP routing and delegates all card logic to the `bingo` package.

On startup, `bingo.InitData("./bingo.json")` loads the goal dataset into a package-level singleton using `sync.Once`, ensuring the file is read exactly once. `CreateBingoCard(seed)` then converts the seed to an `int64`, initializes a seeded `math/rand` source, and draws 25 unique goal indices without replacement.

```
startup
  └─ bingo.InitData() → reads bingo.json → package-level *BingoData (sync.Once)

GET /create?seed=X
  └─ handleCreateBingo()
       └─ bingo.CreateBingoCard(seed)
            ├─ strToInt64(seed)   → numeric parse or FNV-1a hash
            ├─ rand.NewSource()   → seeded RNG
            ├─ generateIndex()    → 25 unique indices (no replacement)
            └─ BingoCard{Goals, Seed} → JSON response
```

**Key components:**

- `main.go` — HTTP server, route registration, request/response handling
- `bingo/bingo.go` — card generation logic, data loading, seed conversion
- `bingo.json` — flat JSON object `{"goals": [...]}` with 150+ goals, each having `name` (Japanese) and `difficulty` (int)
