package pgxredis_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgx-contrib/pgxcache"
	"github.com/pgx-contrib/pgxredis"
	"github.com/redis/go-redis/v9"
)

func ExampleQueryCacher() {
	config, err := pgxpool.ParseConfig(os.Getenv("PGX_DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	conn, err := pgxpool.NewWithConfig(context.TODO(), config)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

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
			MaxLifetime: 30 * time.Second,
			MaxRows:     1,
		},
		Cacher:  cacher,
		Querier: conn,
	}

	rows, err := querier.Query(context.TODO(), "SELECT * from customer")
	if err != nil {
		panic(err)
	}
	// close the rows
	defer rows.Close()

	// Customer struct must be defined
	type Customer struct {
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
	}

	for rows.Next() {
		customer, err := pgx.RowToStructByName[Customer](rows)
		if err != nil {
			panic(err)
		}

		fmt.Println(customer.FirstName)
	}
}
