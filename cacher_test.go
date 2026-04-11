package pgxredis_test

import (
	"context"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pgx-contrib/pgxcache"
	"github.com/pgx-contrib/pgxredis"
	"github.com/redis/go-redis/v9"
)

var _ = Describe("QueryCacher", func() {
	var (
		ctx    context.Context
		cacher *pgxredis.QueryCacher
	)

	BeforeEach(func() {
		if os.Getenv("REDIS_URL") == "" {
			Skip("REDIS_URL not set")
		}

		opts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
		Expect(err).NotTo(HaveOccurred())

		ctx = context.Background()
		cacher = &pgxredis.QueryCacher{
			Client: redis.NewClient(opts),
			Prefix: "test_pgxredis:",
		}
		Expect(cacher.Reset(ctx)).To(Succeed())
	})

	AfterEach(func() {
		if cacher != nil {
			cacher.Close() //nolint:errcheck
		}
	})

	Describe("Get", func() {
		It("returns nil, nil on cache miss", func() {
			key := &pgxcache.QueryKey{SQL: "SELECT 1"}
			item, err := cacher.Get(ctx, key)
			Expect(err).NotTo(HaveOccurred())
			Expect(item).To(BeNil())
		})

		It("returns item after Set", func() {
			key := &pgxcache.QueryKey{SQL: "SELECT 2"}
			original := &pgxcache.QueryItem{CommandTag: "SELECT 1"}
			Expect(cacher.Set(ctx, key, original, time.Minute)).To(Succeed())

			got, err := cacher.Get(ctx, key)
			Expect(err).NotTo(HaveOccurred())
			Expect(got).NotTo(BeNil())
			Expect(got.CommandTag).To(Equal(original.CommandTag))
		})

		It("returns nil after TTL expires", func() {
			key := &pgxcache.QueryKey{SQL: "SELECT 3"}
			item := &pgxcache.QueryItem{CommandTag: "SELECT 1"}
			Expect(cacher.Set(ctx, key, item, 50*time.Millisecond)).To(Succeed())

			time.Sleep(150 * time.Millisecond)

			got, err := cacher.Get(ctx, key)
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(BeNil())
		})
	})

	Describe("Set", func() {
		It("stores multiple items under different keys", func() {
			key1 := &pgxcache.QueryKey{SQL: "SELECT 1"}
			key2 := &pgxcache.QueryKey{SQL: "SELECT 2"}
			item1 := &pgxcache.QueryItem{CommandTag: "SELECT 1"}
			item2 := &pgxcache.QueryItem{CommandTag: "SELECT 2"}

			Expect(cacher.Set(ctx, key1, item1, time.Minute)).To(Succeed())
			Expect(cacher.Set(ctx, key2, item2, time.Minute)).To(Succeed())

			got1, err := cacher.Get(ctx, key1)
			Expect(err).NotTo(HaveOccurred())
			Expect(got1).NotTo(BeNil())
			Expect(got1.CommandTag).To(Equal(item1.CommandTag))

			got2, err := cacher.Get(ctx, key2)
			Expect(err).NotTo(HaveOccurred())
			Expect(got2).NotTo(BeNil())
			Expect(got2.CommandTag).To(Equal(item2.CommandTag))
		})
	})

	Describe("Reset", func() {
		It("clears all items with matching prefix", func() {
			key1 := &pgxcache.QueryKey{SQL: "SELECT 1"}
			key2 := &pgxcache.QueryKey{SQL: "SELECT 2"}
			Expect(cacher.Set(ctx, key1, &pgxcache.QueryItem{CommandTag: "SELECT 1"}, time.Minute)).To(Succeed())
			Expect(cacher.Set(ctx, key2, &pgxcache.QueryItem{CommandTag: "SELECT 2"}, time.Minute)).To(Succeed())

			Expect(cacher.Reset(ctx)).To(Succeed())

			got1, err := cacher.Get(ctx, key1)
			Expect(err).NotTo(HaveOccurred())
			Expect(got1).To(BeNil())

			got2, err := cacher.Get(ctx, key2)
			Expect(err).NotTo(HaveOccurred())
			Expect(got2).To(BeNil())
		})

		It("does not affect keys with a different prefix", func() {
			opts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
			Expect(err).NotTo(HaveOccurred())

			other := &pgxredis.QueryCacher{
				Client: redis.NewClient(opts),
				Prefix: "other_pgxredis:",
			}
			defer other.Close() //nolint:errcheck

			key := &pgxcache.QueryKey{SQL: "SELECT 1"}
			item := &pgxcache.QueryItem{CommandTag: "SELECT 1"}
			Expect(other.Set(ctx, key, item, time.Minute)).To(Succeed())

			Expect(cacher.Reset(ctx)).To(Succeed())

			got, err := other.Get(ctx, key)
			Expect(err).NotTo(HaveOccurred())
			Expect(got).NotTo(BeNil())
		})
	})
})
