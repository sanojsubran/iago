# iago

[![Build Status](https://github.com/sanojsubran/iago/actions/workflows/basic-test.yaml/badge.svg)](https://github.com/sanojsubran/iago/actions)  [![Go Report Card](https://goreportcard.com/badge/github.com/sanojsubran/iago)](https://goreportcard.com/report/github.com/sanojsubran/iago)

A lightweight news aggregation service written in Go. Iago periodically fetches articles from multiple sources — RSS/Atom feeds, Reddit, and Hacker News — normalizes them into a common format, and serves everything through a single JSON API endpoint.

## Features

- Aggregates news from RSS/Atom feeds, Reddit, and Hacker News
- In-memory cache refreshed every 60 minutes via a background goroutine
- Unified JSON response regardless of source type
- CORS-enabled HTTP server ready for browser clients
- Single embedded `config.json` — no external config files needed
- Static binary with no runtime dependencies

## API

### `GET /`

Returns all aggregated news grouped by topic.

```json
{
  "hackernews": [
    { "id": 12345, "title": "Article Title", "url": "https://example.com" }
  ],
  "golang": [
    { "id": 0, "title": "Another Article", "url": "https://go.dev/blog/..." }
  ]
}
```

## Configuration

Sources are defined in `config.json`, which is embedded into the binary at build time:

```json
{
  "sources": [
    {
      "url": "https://go.dev/blog/feed.atom",
      "clientType": "rss",
      "topic": "golang",
      "entryCount": 10
    },
    {
      "url": "https://www.reddit.com/r/programming/new.json?limit=10",
      "clientType": "reddit",
      "topic": "programming",
      "entryCount": 10
    },
    {
      "url": "https://hacker-news.firebaseio.com/v0/topstories.json",
      "clientType": "hackernews",
      "topic": "hackernews",
      "entryCount": 30
    }
  ]
}
```

Supported `clientType` values: `rss`, `reddit`, `hackernews`.

## Running

```bash
go run .
```

The server starts on port **8080**.

## Building

```bash
go build -v ./...
```

For cross-platform releases using [GoReleaser](https://goreleaser.com/):

```bash
goreleaser release --snapshot --clean
```

## Docker

```bash
docker build -t iago .
docker run -p 8080:8080 iago
```

## Testing

```bash
go test -v ./...
```

## License

Apache 2.0 — see [LICENSE](LICENSE).
