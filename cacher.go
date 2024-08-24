package pgxredis

import (
	"context"
	"time"

	"github.com/pgx-contrib/pgxcache"
	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v4"
)

var _ pgxcache.QueryCacher = &QueryCacher{}

// QueryCacher implements cache.QueryCacher interface to use redis as backend with
// go-redis as the redis client library.
type QueryCacher struct {
	// Client is the redis client
	Client redis.UniversalClient
	// Prefix is the prefix for the cache key
	Prefix string
}

// Get gets a cache item from redis. Returns pointer to the item, a boolean
// which represents whether key exists or not and an error.
func (r *QueryCacher) Get(ctx context.Context, key *pgxcache.QueryKey) (*pgxcache.QueryResult, error) {
	data, err := r.Client.Get(ctx, r.prefix(key)).Bytes()
	switch err {
	case nil:
		item := &pgxcache.QueryResult{}
		// unmarshal the result
		if err := msgpack.Unmarshal(data, item); err != nil {
			return nil, err
		}
		return item, nil
	case redis.Nil:
		return nil, nil
	default:
		return nil, err
	}
}

// Set sets the given item into redis with provided TTL duration.
func (r *QueryCacher) Set(ctx context.Context, key *pgxcache.QueryKey, item *pgxcache.QueryResult, ttl time.Duration) error {
	data, err := msgpack.Marshal(item)
	if err != nil {
		return err
	}

	_, err = r.Client.Set(ctx, r.prefix(key), data, ttl).Result()
	return err
}

// Close closes the redis client.
func (r *QueryCacher) Close() error {
	return r.Client.Close()
}

func (r *QueryCacher) prefix(key *pgxcache.QueryKey) string {
	return r.Prefix + key.String()
}