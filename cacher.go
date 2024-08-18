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
	client redis.UniversalClient
	prefix string
}

// NewQueryCacher creates a new instance of redis backend using go-redis client.
// All keys created in redis by pgxcache will have start with prefix.
func NewQueryCacher(client redis.UniversalClient, prefix string) *QueryCacher {
	return &QueryCacher{
		client: client,
		prefix: prefix,
	}
}

// NewQueryCacherWithOptions creates a new instance of redis backend using go-redis
func NewQueryCacherWithOptions(options *redis.UniversalOptions, prefix string) *QueryCacher {
	client := redis.NewUniversalClient(options)
	// done!
	return NewQueryCacher(client, prefix)
}

// Get gets a cache item from redis. Returns pointer to the item, a boolean
// which represents whether key exists or not and an error.
func (r *QueryCacher) Get(ctx context.Context, key *pgxcache.QueryKey) (*pgxcache.QueryResult, error) {
	b, err := r.client.Get(ctx, r.prefix+key.String()).Bytes()
	switch err {
	case nil:
		var item pgxcache.QueryResult
		// unmarshal the result
		if err := msgpack.Unmarshal(b, &item); err != nil {
			return nil, err
		}
		return &item, nil
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

	_, err = r.client.Set(ctx, r.prefix+key.String(), data, ttl).Result()
	return err
}

// Close closes the redis client.
func (r *QueryCacher) Close() error {
	return r.client.Close()
}
