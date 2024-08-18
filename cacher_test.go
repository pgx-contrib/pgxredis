package pgxredis_test

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgx-contrib/pgxcache"
	"github.com/pgx-contrib/pgxredis"
	"github.com/redis/go-redis/v9"
)

var count int

func ExampleQueryCacher() {
	config, err := pgxpool.ParseConfig(os.Getenv("PGX_DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	conn, err := pgxpool.NewWithConfig(context.TODO(), config)
	if err != nil {
		panic(err)
	}

	// create a new client options
	options := &redis.UniversalOptions{}
	// Create a new cacher
	cacher := pgxredis.NewQueryCacherWithOptions(options, "PGX_PREFIX")

	// create a new querier
	querier := &pgxcache.Querier{
		// set the default query options, which can be overridden by the query
		// -- @cache-max-rows 100
		// -- @cache-ttl 30s
		Options: &pgxcache.QueryOptions{
			MaxLiftime: 30 * time.Second,
			MaxRows:    1,
		},
		Cacher:  cacher,
		Querier: conn,
	}

	row := querier.QueryRow(context.TODO(), "SELECT 1")
	if err := row.Scan(&count); err != nil {
		panic(err)
	}
}
