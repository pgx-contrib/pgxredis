# pgxredis

[![CI](https://github.com/pgx-contrib/pgxredis/actions/workflows/ci.yml/badge.svg)](https://github.com/pgx-contrib/pgxredis/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/pgx-contrib/pgxredis)](https://github.com/pgx-contrib/pgxredis/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/pgx-contrib/pgxredis.svg)](https://pkg.go.dev/github.com/pgx-contrib/pgxredis)
[![Go Version](https://img.shields.io/github/go-mod/go-version/pgx-contrib/pgxredis)](go.mod)
[![pgx](https://img.shields.io/badge/pgx-v5-blue)](https://github.com/jackc/pgx)
[![Redis](https://img.shields.io/badge/redis-go--redis%2Fv9-red)](https://github.com/redis/go-redis)

Redis cache backend for [pgxcache](https://github.com/pgx-contrib/pgxcache) — cache PostgreSQL query results in Redis using [go-redis](https://github.com/redis/go-redis).

## Features

- Implements the `pgxcache.QueryCacher` interface backed by Redis
- Key namespacing via a configurable `Prefix` field
- `Reset()` efficiently removes all cache keys matching the prefix using `SCAN` + `DEL`
- Works with any `redis.UniversalClient` (single node, sentinel, cluster)

## Installation

```sh
go get github.com/pgx-contrib/pgxredis
```

## Usage

### Basic

```go
cacher := &pgxredis.QueryCacher{
    Client: redis.NewUniversalClient(&redis.UniversalOptions{
        Addrs: []string{"localhost:6379"},
    }),
    Prefix: "pgxcache:",
}

querier := &pgxcache.Querier{
    Options: &pgxcache.QueryOptions{
        MaxLifetime: 30 * time.Second,
        MaxRows:     100,
    },
    Cacher:  cacher,
    Querier: pool,
}
```

### With a URL

```go
opts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
if err != nil {
    panic(err)
}

cacher := &pgxredis.QueryCacher{
    Client: redis.NewClient(opts),
    Prefix: "pgxcache:",
}
```

### With an existing client

```go
cacher := &pgxredis.QueryCacher{
    Client: existingRedisClient,
    Prefix: "myapp:pgxcache:",
}
```

## Development

### DevContainer

Open in VS Code with the Dev Containers extension. The environment provides Go,
Redis, and Nix automatically.

```
REDIS_URL=redis://redis:6379/0
```

### Nix

```bash
nix develop          # enter shell with Go
go tool ginkgo run -r
```

### Run tests

```bash
# Unit tests only (no Redis required)
go tool ginkgo run -r

# With integration tests
export REDIS_URL="redis://localhost:6379/0"
go tool ginkgo run -r
```

## License

[MIT](LICENSE)
